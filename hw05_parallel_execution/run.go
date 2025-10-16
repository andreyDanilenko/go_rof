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

	type taskWithID struct {
		id   int
		task Task
	}

	tasksChan := make(chan taskWithID)
	resultsChan := make(chan error, len(tasks))

	// Запускаем n воркеров
	for i := 0; i < n; i++ {
		go func(workerID int) {
			log.Printf("Worker %d: запущен и ждет задачи...", workerID)

			// Запускается бесконечный цикл ожидания задач из канала
			for taskItem := range tasksChan {
				log.Printf("Worker %d: получил задачу %d", workerID, taskItem.id)
				resultsChan <- taskItem.task()
			}

			log.Printf("Worker %d: завершил работу", workerID)
		}(i + 1)
	}

	// Отправляем задачи в канал с номерами
	go func() {
		for i, task := range tasks {
			taskNumber := i + 1
			tasksChan <- taskWithID{task: task, id: taskNumber}
		}
		close(tasksChan)
	}()

	errorCount := 0
	completedTask := 0

	// Ловим результат и ошибки
	for i := 0; i < len(tasks); i++ {
		result := <-resultsChan

		if result != nil {
			errorCount++
			if m > 0 && errorCount >= m {
				log.Printf("ДОСТИГНУТ ЛИМИТ %d ОШИБОК! ПРЕРЫВАЮ ВЫПОЛНЕНИЕ!", m)
				return ErrErrorsLimitExceeded
			}
		}

		if result == nil {
			completedTask++
		}
	}

	log.Printf("Total tasks %d %d", errorCount, completedTask)
	// time.Sleep(1 * time.Second)
	log.Printf("Run завершен")
	return nil
}
