package main

import (
	exec "github.com/fixme_my_friend/hw05_parallel_execution"
)

func main() {
	tasksTest := exec.CreateTestTasks()
	n := 3
	m := 3
	err := exec.Run(tasksTest, n, m)
	_ = err

}
