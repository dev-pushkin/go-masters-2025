package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	// DefaultReadTimeout is the default read timeout for the servers.
	DefaultReadTimeout  = 10 * time.Second
	DefaultWriteTimeout = 10 * time.Second
	DefaultIdleTimeout  = 60 * time.Second
)

type Server struct {
	log *slog.Logger

	stdServer  *http.Server
	fastServer *fasthttp.Server
}

func New(l *slog.Logger) *Server {
	return &Server{
		log: l,
		stdServer: &http.Server{
			ReadTimeout:  DefaultReadTimeout,
			WriteTimeout: DefaultWriteTimeout,
			IdleTimeout:  DefaultIdleTimeout,
		},
		fastServer: &fasthttp.Server{
			ReadTimeout:  DefaultReadTimeout,
			WriteTimeout: DefaultWriteTimeout,
			IdleTimeout:  DefaultIdleTimeout,
			Concurrency:  10000,
		},
	}
}

func (s *Server) Start() {
	go func() {
		if err := s.startStandartHttp(":8082"); err != nil {
			s.log.Error("failed to start net/http server", slog.String("error", err.Error()))
		}
	}()
	go func() {
		if err := s.startFastHttp(":8083"); err != nil {
			s.log.Error("failed to start fasthttp server", slog.String("error", err.Error()))
		}
	}()
}

func (s *Server) startStandartHttp(addr string) error {
	s.stdServer.Handler = http.HandlerFunc(s.stdHttpHealthCheck)
	fmt.Printf("Starting net/http server on %s...\n", addr)
	s.stdServer.Addr = addr
	return s.stdServer.ListenAndServe()
}

func (s *Server) startFastHttp(addr string) error {
	s.fastServer.Handler = fasthttp.RequestHandler(s.fastHttpHealthCheck)
	fmt.Printf("Starting fasthttp server on %s...\n", addr)
	return s.fastServer.ListenAndServe(addr)
}

func (s *Server) Stop(ctx context.Context) {
	fmt.Println("Stopping servers...")
	if err := s.stdServer.Shutdown(ctx); err != nil {
		fmt.Printf("failed graysfully stopping net/http server: %v\n", err)
	}

	if err := s.fastServer.ShutdownWithContext(ctx); err != nil {
		fmt.Printf("failed graysfully stopping fasthttp server: %v\n", err)
	}

	fmt.Println("Server graysfully stopped.")
}
