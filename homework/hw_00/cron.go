package cron

import (
	"sort"
	"sync"
	"time"
)

type c struct {
	// мьютекс для синхронизации доступа к задачам
	mu sync.Mutex

	// канал для входящих задач
	in chan *task

	// массив задач, которые нужно выполнить, отсортированный по времени
	// в порядке возрастания времени выполнения, более подходящий тип это связный список
	// но для простоты реализации используем массив
	tasks []*task
}

func Add(task Task, t time.Time) {
	cron.add(task, t)
}

// Логика работы планировщика:
// 1. Кейс с тикером каждую секунду проверяет, есть ли задачи в массиве tasks, которые нужно выполнить по времени
// 2. Кейс с каналом in принимает новые задачи и добавляет их в массив tasks, сортируя по времени выполнения
func (c *c) run() {
	tiker := time.NewTicker(time.Second)
	defer tiker.Stop()
	for {
		select {
		case newTask := <-c.in:
			c.mu.Lock()
			c.tasks = append(c.tasks, newTask)
			// сортируем задачи по времени выполнения
			sort.Slice(c.tasks, func(i, j int) bool {
				return c.tasks[i].scheduledTime.Before(c.tasks[j].scheduledTime)
			})
			c.mu.Unlock()
		case <-tiker.C:
			c.mu.Lock()
			// если нет задач, то выходим из цикла
			if len(c.tasks) == 0 {
				c.mu.Unlock()
				continue
			}

			firstTask := c.tasks[0]
			if firstTask.scheduledTime.After(time.Now()) {
				c.mu.Unlock()
				continue
			}

			for firstTask.scheduledTime.Before(time.Now()) {
				// запускаем задачу в отдельной горутине
				go c.execueFuncWithRecover(firstTask.taskFunc)
				// удаляем задачу из массива задач
				c.tasks = c.tasks[1:]
				if len(c.tasks) > 0 {
					firstTask = c.tasks[0]
				} else {
					break
				}

			}

			c.mu.Unlock()
		}

	}
}

func (c *c) add(taskFunc Task, t time.Time) {
	if t.Before(time.Now()) {
		return
	}
	c.in <- &task{
		scheduledTime: t,
		taskFunc:      taskFunc,
	}
}

func new() *c {
	in := make(chan *task)
	return &c{
		in: in,
	}
}
