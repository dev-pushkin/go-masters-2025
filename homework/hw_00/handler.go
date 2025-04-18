package cron

import "time"

func Add(task Task, t time.Time) {
	cron.add(task, t)
}
