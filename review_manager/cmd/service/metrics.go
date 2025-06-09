package main

import (
	"context"

	"github.com/go_course_master/review_manager/internal/metric"
	"github.com/go_course_master/review_manager/internal/telemetry"
)

func (s *serviceDI) initMetricsAndTelemetry() error {
	// Initialize metrics
	if err := metric.InitAndServe(); err != nil {
		s.log.Error("Failed to initialize metrics", "error", err)
		return err
	}
	s.log.Info("Metrics initialized successfully")

	// Initialize telemetry
	s.telemetryStop = telemetry.Init(s.config.Metric.Namespace)
	s.log.Info("Telemetry initialized successfully")

	return nil
}

func (s *serviceDI) stopMetricsAndTelemetry(ctx context.Context) error {
	// Stop telemetry
	if s.telemetryStop != nil {
		if err := s.telemetryStop(ctx); err != nil {
			s.log.Error("Failed to stop telemetry", "error", err)
			return err
		}
		s.log.Info("Telemetry stopped successfully")
	}

	return nil
}
