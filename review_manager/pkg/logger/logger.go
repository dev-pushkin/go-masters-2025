package logger

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/Graylog2/go-gelf.v1/gelf"
)

var defaultConfig = &config{
	Level:      "error",
	Format:     "text",
	GelfUdpUrl: "",
}

type Option func(*config)

type config struct {
	Level       string
	Format      string
	GelfUdpUrl  string
	ServiceName string
}

func WithLevel(level string) Option {
	return func(c *config) {
		c.Level = level
	}
}

func WithFormat(format string) Option {
	return func(c *config) {
		c.Format = format
	}
}

func WithGelfUdpUrl(gelfUdpUrl string) Option {
	return func(c *config) {
		c.GelfUdpUrl = gelfUdpUrl
	}
}

func WithSerciveName(serviceName string) Option {
	return func(c *config) {
		c.ServiceName = serviceName
	}
}

func Init(options ...Option) *slog.Logger {
	cfg := defaultConfig
	for _, opt := range options {
		opt(cfg)
	}

	slogOptions := &slog.HandlerOptions{
		Level:     toSlogLevel(cfg.Level),
		AddSource: false,
	}

	var outputWriter io.Writer = os.Stdout // Default output writer
	if cfg.GelfUdpUrl != "" {
		gelfWriter, err := gelf.NewWriter(cfg.GelfUdpUrl)
		if err != nil {
			slog.Default().Error("create GELF writer, use default os Stdout", "error", err)
		} else {
			outputWriter = gelfWriter
		}
	}

	return slog.New(slog.NewTextHandler(outputWriter, slogOptions)).WithGroup(cfg.ServiceName)
}

func toSlogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelError // Default to error if unknown level
	}
}
