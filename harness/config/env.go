package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"

	herrors "dari/harness/errors"
)

func Require(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", herrors.New(herrors.MissingEnv, errors.New(name))
	}
	return value, nil
}

func RequireURL(name string) (*url.URL, error) {
	value, err := Require(name)
	if err != nil {
		return nil, err
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, herrors.New(herrors.InvalidEnv, fmt.Errorf("%s=%q is not an absolute URL", name, value))
	}
	return parsed, nil
}

func RequirePositiveInt(name string) (int64, error) {
	value, err := Require(name)
	if err != nil {
		return 0, err
	}
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil || number <= 0 {
		return 0, herrors.New(herrors.InvalidEnv, fmt.Errorf("%s=%q is not a positive integer", name, value))
	}
	return number, nil
}
