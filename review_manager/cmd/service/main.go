package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var s *serviceDI

func main() {
	s = &serviceDI{}

	if err := s.initialize(); err != nil {
		s.log.Error("Failed to initialize service", "error", err)
		return
	}

	if err := s.start(); err != nil {
		s.log.Error("Failed to start service", "error", err)
		return
	}

	s.log.Info("Service started successfully", "host", s.config.API.Host, "port", s.config.API.Port)
	// Initialize metrics and telemetry
	if err := s.initMetricsAndTelemetry(); err != nil {
		s.log.Error("Failed to initialize metrics and telemetry", "error", err)
		return
	}

	// Wait for termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	// Handle graceful shutdown
	s.log.Info("Received shutdown signal, stopping service")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.stop(ctx); err != nil {
		s.log.Error("Failed to stop service", "error", err)
	} else {
		s.log.Info("Service stopped gracefully")
	}
}
