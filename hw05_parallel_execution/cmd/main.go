package main

import (
	"log"

	exec "github.com/fixme_my_friend/hw05_parallel_execution"
)

func main() {
	tasksTest := exec.СreateTestTasks()
	n := 1
	m := 3
	log.Printf("Starting: %d tasks, %d workers, %d max errors", len(tasksTest), n, m)

	err := exec.Run(tasksTest, n, m)

	_ = err

	// Temporary fix to use the variable

}
