package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return errors.New("workers count must be > 0")
	}

	type taskWithID struct {
		id   int
		task Task
	}

	tasksChan := make(chan taskWithID)
	errChan := make(chan error)
	done := make(chan struct{})

	var wg sync.WaitGroup
	var errCount int

	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case t, ok := <-tasksChan:
				if !ok {
					return
				}
				select {
				case errChan <- t.task():
				case <-done:
					return
				}
			}
		}
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker()
	}

	go func() {
		for i, task := range tasks {
			select {
			case <-done:
				return
			case tasksChan <- taskWithID{i, task}:
			}
		}
		close(tasksChan)
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			if m > 0 {
				errCount++
				if errCount >= m {
					close(done)
					wg.Wait()
					return ErrErrorsLimitExceeded
				}
			}
		}
	}

	return nil
}
