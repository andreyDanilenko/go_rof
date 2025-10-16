package hw05parallelexecution

import (
	"errors"
	"log"
	"time"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	log.Printf("Starting: %d tasks, %d workers, %d max errors", len(tasks), n, m)

	tasksChan := make(chan Task)

	// Запускаем n воркеров
	for i := 0; i < n; i++ {
		workerID := i + 1
		go func() {
			log.Printf("Worker %d: запущен и ждет задачи...", workerID)

			for task := range tasksChan {
				log.Printf("Worker %d: получил задачу", workerID)
				err := task()
				if err != nil {
					log.Printf("Worker %d: ошибка - %v", workerID, err)
				} else {
					log.Printf("Worker %d: задача выполнена успешно", workerID)
				}
			}

			log.Printf("Worker %d: завершил работу", workerID)
		}()
	}

	// Отправляем задачи в канал с номерами
	log.Printf("Отправляю задачи в канал...")
	for i, task := range tasks {
		taskNumber := i + 1
		log.Printf("Отправляю задачу %d в канал", taskNumber)
		tasksChan <- task
	}

	log.Printf("Все задачи отправлены, закрываю канал...")
	close(tasksChan)

	time.Sleep(1 * time.Second)
	log.Printf("Run завершен")
	return nil
}

func CreateTestTasks() []Task {
	return []Task{
		func() error {
			log.Printf("Задача 1 выполняется")
			return nil
		},
		func() error {
			log.Printf("Задача 2 выполняется")
			return nil
		},
	}
}
