package cron

import (
	"fmt"
	"testing"
	"time"
)

var mockTaskFunc = func(out chan time.Time) func() {
	return func() {
		fmt.Println("Mock task executed")
		out <- time.Now()
	}
}

type MockTask struct {
	mockTask func()
}

func (m *MockTask) Exec() {
	m.mockTask()
}

func TestAdd_AddOne(t *testing.T) {
	funcExecue := make(chan time.Time, 1)

	// Создаем мок задачу
	mockTask := &MockTask{}
	mockTask.mockTask = mockTaskFunc(funcExecue)
	scheduledTime := time.Now().Add(5 * time.Second)

	// Добавляем задачу в планировщик
	Add(mockTask, scheduledTime)

	// Ждем выполнения задачи
	realTime := <-funcExecue
	// Проверяем, что задача была выполнена в нужное время
	if scheduledTime.Sub(realTime) > 1*time.Second {
		t.Errorf("Expected task to be executed at %v, but got %v", scheduledTime, realTime)
	}
}

func TestAdd_Add_many(t *testing.T) {
	funcExecue := make(chan time.Time, 4)

	// Создаем мок задачу
	mockTask1 := &MockTask{}
	mockTask1.mockTask = mockTaskFunc(funcExecue)
	scheduledTime1 := time.Now().Add(5 * time.Second)

	mockTask2 := &MockTask{}
	mockTask2.mockTask = mockTaskFunc(funcExecue)
	scheduledTime2 := time.Now().Add(10 * time.Second)

	mockTask3 := &MockTask{}
	mockTask3.mockTask = mockTaskFunc(funcExecue)
	scheduledTime3 := time.Now().Add(15 * time.Second)
	mockTask4 := &MockTask{}
	mockTask4.mockTask = mockTaskFunc(funcExecue)
	scheduledTime4 := time.Now().Add(20 * time.Second)

	// Добавляем задачи в планировщик
	Add(mockTask1, scheduledTime1)
	Add(mockTask2, scheduledTime2)
	Add(mockTask3, scheduledTime3)
	Add(mockTask4, scheduledTime4)
	// Ждем выполнения задач
	realTime1 := <-funcExecue
	realTime2 := <-funcExecue
	realTime3 := <-funcExecue
	realTime4 := <-funcExecue
	// Проверяем, что задачи были выполнены в нужное время
	if scheduledTime1.Sub(realTime1) > 1*time.Second {
		t.Errorf("Expected task 1 to be executed at %v, but got %v", scheduledTime1, realTime1)
	}

	if scheduledTime2.Sub(realTime2) > 1*time.Second {
		t.Errorf("Expected task 2 to be executed at %v, but got %v", scheduledTime2, realTime2)
	}

	if scheduledTime3.Sub(realTime3) > 1*time.Second {
		t.Errorf("Expected task 3 to be executed at %v, but got %v", scheduledTime3, realTime3)
	}

	if scheduledTime4.Sub(realTime4) > 1*time.Second {
		t.Errorf("Expected task 4 to be executed at %v, but got %v", scheduledTime4, realTime4)
	}

}

func TestAdd_panic(t *testing.T) {
	funcExecue := make(chan time.Time, 1)

	// Создаем мок задачу
	mockTask := &MockTask{}
	mockTask.mockTask = func() {
		fmt.Println("Mock task executed")
		funcExecue <- time.Now()
		panic("panic")
	}
	scheduledTime := time.Now().Add(1 * time.Second)

	// Добавляем задачу в планировщик
	Add(mockTask, scheduledTime)

	// Ждем выполнения задачи
	realTime := <-funcExecue
	// Проверяем, что задача была выполнена в нужное время
	if scheduledTime.Sub(realTime) > 1*time.Second {
		t.Errorf("Expected task 4 to be executed at %v, but got %v", scheduledTime, realTime)
	}
}
