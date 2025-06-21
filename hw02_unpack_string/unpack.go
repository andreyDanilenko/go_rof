package unpack

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidString = errors.New("invalid string")
)

func Unpack(s string) (string, error) {
	var builder strings.Builder
	runes := []rune(s)
	length := len(runes)

	escaped := false
	for i := 0; i < length; i++ {
		r := runes[i]

		if escaped {
			if !(r == '\\' || isDigit(r)) {
				return "", ErrInvalidString
			}

			builder.WriteRune(r)
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			continue
		}

		if isDigit(r) {
			// если первый символ - цифра, либо повторяется цифра
			if i == 0 {
				return "", ErrInvalidString
			}

			prev := runes[i-1]
			// проверка: если предыдущий тоже был цифрой (и не был экранирован)
			if isDigit(prev) && !(i >= 2 && runes[i-2] == '\\') {
				return "", ErrInvalidString
			}

			count, err := strconv.Atoi(string(r))
			if err != nil {
				return "", ErrInvalidString
			}

			if count == 0 {
				// удаляем последний символ
				output := []rune(builder.String())
				builder.Reset()
				builder.WriteString(string(output[:len(output)-1]))
				continue
			}

			// повторить предыдущий символ count-1 раз (он уже записан 1 раз)
			builder.WriteString(strings.Repeat(string(runes[i-1]), count-1))
		} else {
			builder.WriteRune(r)
		}
	}

	if escaped {
		return "", ErrInvalidString
	}

	return builder.String(), nil
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
