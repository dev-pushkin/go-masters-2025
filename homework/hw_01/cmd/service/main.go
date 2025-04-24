package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go_course_master/homework/hw_01/internal/api"
	"github.com/go_course_master/homework/hw_01/internal/app"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := app.NewApp(logger)
	app.RndStrFn = func() string {
		return "12345678901234567890123456789012" // 32 байта
	}

	httpServer := &http.Server{
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		IdleTimeout:    5 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	api := api.NewApi(logger, httpServer, app)

	err := api.RegisterRouteAndServe(":8080")
	if err != nil {
		slog.Error("Error starting server", "error", err)
		return
	}
}
