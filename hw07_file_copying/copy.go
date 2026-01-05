package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	src, err := os.Open(fromPath)
	if err != nil {
		fmt.Println("open err:", err)
		return err
	}
	defer src.Close()

	// Получаем информацию о файле
	info, err := src.Stat()
	if err != nil {
		fmt.Println("info err:", err)
		return err
	}

	// Узнаем размер всего файла из которого копируем
	size := info.Size()
	if size < offset {
		fmt.Printf("file size %x is less than offset %x\n", size, offset)
		return err
	}

	// Перемещаемся на offset
	_, err = src.Seek(offset, io.SeekStart)
	if err != nil {
		fmt.Println("you can't move to offset:", err)
		return err
	}

	// Оапеделяем размер копирования
	remaining := size - offset
	copySize := remaining
	if limit > 0 && limit < remaining {
		copySize = limit
	}

	// Создаем файл для копирования
	dst, err := os.Create(toPath)
	if err != nil {
		fmt.Println("info err:", err)
		return err
	}
	defer dst.Close()

	// Для прогресс-бара используем copySize - это реальное количество байт, которые будем копировать
	bar := pb.StartNew(int(copySize))

	reader := bar.NewProxyReader(src)

	// Использовать обернутый reader для копирования
	io.CopyN(dst, reader, copySize)

	// Завершить прогресс-бар
	bar.Finish()
	return nil
}
