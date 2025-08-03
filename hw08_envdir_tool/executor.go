package main

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) int {
	if len(cmd) == 0 {
		return 1
	}

	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = prepareEnv(env)

	if err := command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		}
		return 1
	}

	return 0
}

func prepareEnv(env Environment) []string {
	var result []string

	// Сначала добавляем существующие переменные, кроме тех, что нужно удалить
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) != 2 {
			continue
		}
		if _, shouldRemove := env[pair[0]]; !shouldRemove {
			result = append(result, e)
		}
	}

	// Затем добавляем новые переменные
	for name, value := range env {
		if !value.NeedRemove && value.Value != "" {
			result = append(result, name+"="+value.Value)
		}
	}

	return result
}
