package envvar

import (
	"fmt"
	"net/mail"
	"net/url"
	"os"
	"strconv"
)

func Get[T string | bool | int | uint | float64](name string) T {
	val, _ := getEnvVar[T](name)
	return val
}

func GetWithDefault[T string | bool | int | uint | float64](name string, defaultValue T) T {
	val, err := getEnvVar[T](name)
	if err != nil {
		return defaultValue
	}
	return val
}

type ValidationFunc[T string | bool | int | uint | float64] func(T) error

func GetAndValidate[T string | bool | int | uint | float64](name string, validate ValidationFunc[T]) (T, error) {
	val, err := getEnvVar[T](name)
	if err != nil {
		return val, fmt.Errorf("could not parse %q: %v", name, err)
	}
	err = validate(val)
	if err != nil {
		return *new(T), fmt.Errorf("%q failed validation: %v", name, err)
	}
	return val, nil
}

func GreaterThan[T int | uint | float64](value T) ValidationFunc[T] {
	return func(t T) error {
		if t <= value {
			return fmt.Errorf("must be greater than %v", value)
		}
		return nil
	}
}

func LessThan[T int | uint | float64](value T) ValidationFunc[T] {
	return func(t T) error {
		if t >= value {
			return fmt.Errorf("must be less than %v", value)
		}
		return nil
	}
}

func Between[T int | uint | float64](min, max T) ValidationFunc[T] {
	return func(t T) error {
		if t < min || t > max {
			return fmt.Errorf("must be between %v and %v", min, max)
		}
		return nil
	}
}

func Within[T string | bool | int | uint | float64](values []T) ValidationFunc[T] {
	return func(t T) error {
		var found bool
		for _, val := range values {
			if t == val {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%v not within %v", t, values)
		}
		return nil
	}
}

func ValidEmail() ValidationFunc[string] {
	return func(s string) error {
		_, err := mail.ParseAddress(s)
		if err != nil {
			return fmt.Errorf("%s is not a valid email", s)
		}
		return nil
	}
}

func ValidUrl() ValidationFunc[string] {
	return func(s string) error {
		_, err := url.Parse(s)
		if err != nil {
			return fmt.Errorf("%s is not a valid url", s)
		}
		return nil
	}
}

func getEnvVar[T string | bool | int | uint | float64](name string) (T, error) {
	var val T
	envVar, exists := os.LookupEnv(name)
	if !exists {
		return val, fmt.Errorf("%q is not set", name)
	}

	var ret any
	var err error
	switch any(val).(type) {
	case string:
		ret = envVar
	case bool:
		ret, err = strconv.ParseBool(envVar)
	case int:
		ret, err = strconv.Atoi(envVar)
	case uint:
		res, convErr := strconv.ParseUint(envVar, 10, 64)
		ret = uint(res)
		err = convErr
	case float64:
		ret, err = strconv.ParseFloat(envVar, 64)
	}

	return ret.(T), err
}
