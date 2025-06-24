package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var builder strings.Builder
	runes := []rune(s)
	length := len(runes)
	escaped := false

	for i := 0; i < length; i++ {
		r := runes[i]

		if isDigit(r) && i == 0 {
			return "", ErrInvalidString
		}

		if escaped {
			if err := handleEscaped(r, &builder); err != nil {
				return "", err
			}
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			continue
		}

		if isDigit(r) {
			if err := handleDigit(runes, i, &builder); err != nil {
				return "", err
			}
			continue
		}

		builder.WriteRune(r)
	}

	if escaped {
		return "", ErrInvalidString
	}

	return builder.String(), nil
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func handleEscaped(r rune, builder *strings.Builder) error {
	if !(r == '\\' || isDigit(r)) {
		return ErrInvalidString
	}
	builder.WriteRune(r)
	return nil
}

func handleDigit(runes []rune, i int, builder *strings.Builder) error {
	if i == 0 {
		return ErrInvalidString
	}

	curr := runes[i]
	prev := runes[i-1]

	// Если предыдущий символ — неэкранированная цифра
	if isDigit(prev) && !(i >= 2 && runes[i-2] == '\\') {
		return ErrInvalidString
	}

	count, err := strconv.Atoi(string(curr))
	if err != nil {
		return ErrInvalidString
	}

	if count == 0 {
		output := []rune(builder.String())
		if len(output) == 0 {
			return ErrInvalidString
		}
		builder.Reset()
		builder.WriteString(string(output[:len(output)-1]))
		return nil
	}

	builder.WriteString(strings.Repeat(string(prev), count-1))
	return nil
}
