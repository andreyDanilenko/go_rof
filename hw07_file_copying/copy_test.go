package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func checkResult(from, to string, offset, limit int64) bool {
	toFile, err := os.Open(to)
	fromFile, _ := os.Open(from)
	if err != nil {
		fmt.Println("open err:", err)
		return false
	}
	defer toFile.Close()

	toInfo, err := toFile.Stat()
	fromInfo, _ := fromFile.Stat()
	if err != nil {
		fmt.Println("info err:", err)
		return false
	}

	toSize := toInfo.Size()
	fromSize := fromInfo.Size()

	remaining := fromSize - offset
	wantSize := remaining
	if limit > 0 && limit < remaining {
		wantSize = limit
	}

	if wantSize != toSize {
		fmt.Println("file size not complete:", err)
		return false
	}

	return true
}

func TestCopy(t *testing.T) {
	// go run . -from testdata/input.txt -to testdata/expected_offset100_limit1000.txt  -offset 100  -limit 1000
	t.Run("out_offset100_limit1000", func(t *testing.T) {
		from := "testdata/input.txt"
		to := "testdata/out_offset100_limit1000.txt"
		offset := int64(100)
		limit := int64(1000)

		err := Copy(from, to, offset, limit)
		require.NoError(t, err)
		require.True(t, checkResult(from, to, offset, limit))

		expected, _ := os.ReadFile("testdata/expected_offset100_limit1000.txt")
		actual, _ := os.ReadFile(to)
		require.Equal(t, expected, actual)
	})

	// go run . -from testdata/input.txt -to testdata/expected_offset0_limit0.txt  -offset 0  -limit 0
	t.Run("out_offset0_limit0", func(t *testing.T) {
		from := "testdata/input.txt"
		to := "testdata/out_offset0_limit0.txt"
		offset := int64(0)
		limit := int64(0)

		err := Copy(from, to, offset, limit)
		require.NoError(t, err)
		require.True(t, checkResult(from, to, offset, limit))

		expected, _ := os.ReadFile("testdata/expected_offset0_limit0.txt")
		actual, _ := os.ReadFile(to)
		require.Equal(t, expected, actual)
	})
}
