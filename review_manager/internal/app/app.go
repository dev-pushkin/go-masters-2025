package app

import (
	"log/slog"

	"github.com/go_course_master/review_manager/config"
	"github.com/go_course_master/review_manager/internal/db"
	"github.com/go_course_master/review_manager/internal/rate"
	"github.com/go_course_master/review_manager/pkg/reviewrate"
	"github.com/robfig/cron"
)

type App struct {
	log       *slog.Logger
	db        db.DB
	review    reviewrate.ReviewRate
	config    config.Config
	rateQueue rate.Queue

	cron *cron.Cron
}

func New(log *slog.Logger, db db.DB, config config.Config, reviewRate reviewrate.ReviewRate, rateQueue rate.Queue) *App {
	return &App{
		log:       log,
		db:        db,
		config:    config,
		rateQueue: rateQueue,
		review:    reviewRate,
		cron:      cron.New(),
	}
}

func (a *App) Start() error {
	a.log.Info("Starting application")
	// Start the job scheduler
	a.startJobs()

	return nil
}
