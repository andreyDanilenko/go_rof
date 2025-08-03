package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || strings.Contains(entry.Name(), "=") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		if info.Size() == 0 {
			env[entry.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var buf bytes.Buffer
		_, err = io.CopyN(&buf, file, info.Size())
		if err != nil && err != io.EOF {
			return nil, err
		}

		firstLine := bytes.SplitN(buf.Bytes(), []byte{'\n'}, 2)[0]
		value := bytes.ReplaceAll(firstLine, []byte{0}, []byte{'\n'})
		valueStr := strings.TrimRight(string(value), " \t")
		valueStr = strings.Trim(valueStr, `"`) // Удаляем кавычки если есть

		env[entry.Name()] = EnvValue{
			Value:      strings.TrimSpace(valueStr), // Удаляем все пробелы по краям
			NeedRemove: false,
		}
	}

	return env, nil
}
