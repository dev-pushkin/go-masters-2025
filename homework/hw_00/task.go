package cron

import "time"

type Task interface {
	Exec()
}

type task struct {
	scheduledTime time.Time
	taskFunc      Task
}
