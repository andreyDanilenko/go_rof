package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidatorRule struct {
	name string
	arg  string
}

type ValidationErrors []ValidationError
type ValidatorRules []ValidatorRule

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("should struct")
	}

	typeVal := val.Type()
	var errs []ValidationError

	for i := 0; i < typeVal.NumField(); i++ {
		fieldType := typeVal.Field(i)

		fmt.Println("##", fieldType, errs)

		// `validate:"min:18|max:50"` => min:18|max:50
		validateTag := fieldType.Tag.Get("validate")

		if validateTag == "" {
			continue
		}

		fieldValue := val.Field(i)
		if err := validateField(fieldValue, validateTag); err != nil {
			fmt.Println("fieldValue", err)

			errs = append(errs, ValidationError{
				Field: fieldType.Name,
				Err:   err,
			})
		}
	}

	if len(errs) > 0 {
		return ValidationErrors(errs)
	}

	fmt.Println("#", typeVal)

	return nil
}

// rules
// int => len, max, min
// string => regexp, in, max, min, len
// slice =>

func validateField(fieldValue reflect.Value, validateTag string) error {
	fmt.Printf("Валидирую поле типа %v с тегом: %s\n",
		fieldValue.Kind(), validateTag)

	rules, _ := parseValidateTag(validateTag)

	fmt.Printf("Валидирую %v с тегом: %s\n",
		rules, validateTag)

	switch fieldValue.Kind() {
	case reflect.Int:
		value := fieldValue.Int()
		fmt.Printf("Это %v: %v\n", fieldValue.Kind(), value)
	case reflect.String:
		fmt.Printf("Это %v: %v\n", fieldValue.Kind(), fieldValue.String())
		return validateString(fieldValue.String(), rules)

	case reflect.Slice:
		fmt.Printf("Это слайс, длина: %d\n", fieldValue.Len())

	default:
		return fmt.Errorf("неподдерживаемый тип: %v", fieldValue.Kind())
	}

	return nil
}

func parseValidateTag(validateTag string) (ValidatorRules, error) {
	if validateTag == "" {
		return nil, nil
	}

	var rules ValidatorRules
	// min:18|max:50 => [min:18 max:50]
	parts := strings.Split(validateTag, "|")

	for _, part := range parts {
		// на всякий случай нужно убрать пробелы по краям
		part = strings.TrimSpace(part)
		// [min:18 max:50] => min:18 => max:50
		if part == "" {
			continue
		}

		// max:50 => name: max => arg: 50
		name, arg, found := strings.Cut(part, ":")
		if !found {
			return nil, fmt.Errorf("invalid rule: %s", part)
		}

		rules = append(rules, ValidatorRule{
			name: strings.TrimSpace(name),
			arg:  strings.TrimSpace(arg),
		})
	}

	// rules = [{min 18} {max 50}]
	return rules, nil
}

func validateString(value string, rules ValidatorRules) error {

	for _, rule := range rules {
		fmt.Println("fieldValue", rule.arg)

		switch rule.name {
		// min, max, len, regex, pattern, match
		case "len":
			// string => int
			expected, err := strconv.Atoi(rule.arg)
			if err != nil {
				return fmt.Errorf("invalid len: %s", rule.arg)
			}

			if len(value) != expected {
				return fmt.Errorf("length must be %d", expected)
			}
		case "regex", "pattern", "match":
			// Компилируем регулярное выражение
			re, err := regexp.Compile(rule.arg)
			if err != nil {
				return fmt.Errorf("invalid regex pattern: %s", rule.arg)
			}
			// Проверяем соответствие
			if !re.MatchString(value) {
				return fmt.Errorf("value must match pattern: %s", rule.arg)
			}
		case "min":
			// string => int
			min, err := strconv.Atoi(rule.arg)

			fmt.Println("fieldValue", min)

			if err != nil {
				return fmt.Errorf("invalid rule min: %s", rule.arg)
			}

			if len(value) < min {
				return fmt.Errorf("length must be at least %d", min)
			}

		case "max":
			max, err := strconv.Atoi(rule.arg)
			if err != nil {
				return fmt.Errorf("invalid rule max: %s", rule.arg)
			}

			if len(value) < max {
				return fmt.Errorf("length must be at least %d", max)
			}
		case "in":

		}
	}

	return nil
}
