package cron

import "fmt"

func (c *c) execueFuncWithRecover(task Task) {
	defer func() {
		if r := recover(); r != nil {
			// обработка паники
			// например, логирование ошибки
			fmt.Println("recovered in execueFuncWithRecover", r)

		}
	}()
	task.Exec()
}
