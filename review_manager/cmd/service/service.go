package main

import (
	"context"
	"log/slog"

	"github.com/go_course_master/review_manager/config"
	"github.com/go_course_master/review_manager/internal/app"
	"github.com/go_course_master/review_manager/internal/db/migration"
	"github.com/go_course_master/review_manager/internal/db/pg"
	"github.com/go_course_master/review_manager/internal/rate/memory"
	"github.com/go_course_master/review_manager/internal/server"
	"github.com/go_course_master/review_manager/pkg/logger"
	"github.com/go_course_master/review_manager/pkg/pgx"
	"github.com/go_course_master/review_manager/pkg/reviewrate/ollama"
	"github.com/jackc/pgx/v4/pgxpool"
)

type serviceDI struct {
	log        *slog.Logger
	config     config.Config
	pgxpool    *pgxpool.Pool  // заменить на конкретную реализацию
	db         *pg.PostgresDB // заменить на конкретную реализацию
	app        *app.App
	server     *server.Server
	rateQueue  *memory.Queue
	reviewRate *ollama.ReviewRate

	telemetryStop func(context.Context) error
}

func (s *serviceDI) initialize() error {
	var err error
	s.config, err = config.Read()
	if err != nil {
		s.log.Error("Failed to read configuration", "error", err)
	}

	s.log = logger.Init(
		logger.WithLevel(s.config.Logger.Level),
		logger.WithSerciveName("review_manager"),
		logger.WithFormat(s.config.Logger.Format),
		logger.WithGelfUdpUrl(s.config.Logger.GelfUdpUrl))

	s.log.Info("Initializing service dependencies")

	err = migration.NewGooseMigration(s.log).Migrate(s.config.DB.PgURL(), s.config.DB.MigrationsPath)
	if err != nil {
		s.log.Error("Failed to run database migrations", "error", err)
	}

	// Initialize the database connection
	s.pgxpool, err = pgx.New(s.config.PgxConfig())
	if err != nil {
		s.log.Error("Failed to initialize database connection", "error", err)
		return err
	}
	s.log.Debug("Database connection initialized")

	// Initialize the database layer
	s.log.Debug("Initializing database layer")
	s.db = pg.NewPostgresDB(s.log, s.pgxpool)
	if err != nil {
		s.log.Error("Failed to initialize database layer", "error", err)
		return err
	}

	// Initialize the review rate service
	s.log.Debug("Initializing review rate service")
	s.reviewRate = ollama.New(ollama.WithSystemPrompt("Ты оценщик отзывов, который должен по тексту выставить бал от 1 до 5, в ответ давай только число оценки"))

	// Initialize the rate queue
	s.rateQueue = memory.NewQueue(s.log, s.reviewRate, s.db)

	// Initialize the application layer
	s.app = app.New(s.log, s.db, s.config, s.reviewRate, s.rateQueue)

	// Initialize the server
	s.server = server.NewServer(s.log, s.app, s.config)

	return nil
}

func (s *serviceDI) start() error {
	s.log.Info("Starting service")

	// Start the application
	if err := s.app.Start(); err != nil {
		s.log.Error("Failed to start application", "error", err)
	}

	// Start the rate queue
	s.log.Debug("Starting rate queue")
	s.rateQueue.Run()

	// Start the server
	s.log.Debug("Starting server")
	if err := s.server.Start(); err != nil {
		return err
	}

	return nil
}

func (s *serviceDI) stop(ctx context.Context) error {
	s.log.Info("Stopping service")

	// Stop the server
	s.log.Debug("Stopping server")
	if err := s.server.Stop(ctx); err != nil {
		return err
	}

	// Stop review rate service
	s.log.Debug("Stopping review rate service")
	s.reviewRate.Stop(ctx)

	// Stop the rate queue
	s.log.Debug("Stopping rate queue")
	s.rateQueue.Stop(ctx)

	s.pgxpool.Close()
	s.log.Debug("Database connection closed")

	return nil
}
