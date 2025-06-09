package server

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go_course_master/review_manager/config"
	"github.com/go_course_master/review_manager/internal/app"
	"github.com/go_course_master/review_manager/internal/metric"
	"github.com/go_course_master/review_manager/internal/telemetry"
	"github.com/gorilla/mux"
)

type Server struct {
	log        *slog.Logger
	app        *app.App
	config     config.Config
	httpServer *http.Server
	router     *mux.Router
}

func NewServer(log *slog.Logger, app *app.App, config config.Config) *Server {
	s := &http.Server{
		Addr:         config.API.Host + ":" + strconv.Itoa(config.API.Port),
		ReadTimeout:  time.Duration(config.API.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.API.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.API.IdleTimeout) * time.Second,
	}

	return &Server{
		log:        log,
		app:        app,
		config:     config,
		httpServer: s,
		router:     mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	s.log.Info("Starting server", "host", s.config.API.Host, "port", s.config.API.Port)

	s.initRoutes()
	s.httpServer.Handler = s.router
	// Start the HTTP server in a goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("Failed to start server", "error", err)
			// Handle the error appropriately, e.g., log it or return it
			return
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("Stopping server")
	return s.httpServer.Shutdown(ctx)
}
func (s *Server) Run() error {
	s.log.Info("Running server", "host", s.config.API.Host, "port", s.config.API.Port)

	// Start the server
	if err := s.Start(); err != nil {
		s.log.Error("Failed to start server", "error", err)
		return err
	}

	// Placeholder for actual server run logic.
	// In a real application, you would block here until the server is stopped.
	return nil
}

func (s *Server) initRoutes() {
	// add middleware
	s.router.Use(
		telemetry.TelemetryMiddleware("rate-review-api"),
		recoverMiddleware(s.log),
		metric.StoreMetricsMiddleware(s.config.API.Host))

	// health check endpoint
	s.router.HandleFunc("/health", s.handleHealthCheck).Methods("GET")

	// Add routes for reviews
	s.router.HandleFunc("/review", s.handleReviewAdd).Methods("POST")
	s.router.HandleFunc("/reviews", s.handleReviewList).Methods("GET")
	s.router.HandleFunc("/reviews/rate/{objectID}", s.handleObjectRate).Methods("GET")

	s.httpServer.Handler = s.router
}
