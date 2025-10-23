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

	tasksChan := make(chan Task)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errCount int

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range tasksChan {
				if err := task(); err != nil {
					mu.Lock()
					errCount++
					mu.Unlock()
				}
			}
		}()
	}

	for _, task := range tasks {
		tasksChan <- task

		mu.Lock()
		if m > 0 && errCount >= m {
			mu.Unlock()
			break
		}
		mu.Unlock()

	}

	close(tasksChan)

	wg.Wait()

	if m > 0 && errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
