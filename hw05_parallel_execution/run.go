package hw05parallelexecution

import (
	"errors"
	"log"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {

	log.Printf("Starting: %d tasks, %d workers, %d max errors", len(tasks), n, m)

	taskTests := []Task{
		func() error {
			return nil
		},
		func() error {
			return nil
		},
		func() error {
			return nil
		},
	}

	_ = tasks
	_ = taskTests

	// Place your code here.
	return nil
}

func СreateTestTasks() []Task {
	return []Task{
		func() error {
			return nil
		},
		func() error {
			return nil
		},
		func() error {
			return nil
		},
	}
}
