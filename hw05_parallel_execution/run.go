package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	tasksChan := make(chan Task, n)

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
