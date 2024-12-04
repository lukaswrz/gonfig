// Package gonfig provides utilities for reading, unmarshaling, and validating
// configuration files.
package gonfig

import (
	"errors"
	"fmt"
	"os"
)

// Validator defines an interface for validating a configuration of type T. It
// ensures that the provided configuration meets required constraints.
type Validator[T any] interface {
	// Validate checks the provided configuration and returns an error if it is
	// invalid.
	Validate(config T) []error
}

// ReadConfig reads a configuration file, unmarshals its content into the given
// configuration object, and validates it using the provided validator. If the
// primary path is empty, it searches for the configuration file in the
// fallback paths. The function returns the resolved path to the configuration
// file or a list of errors if the file cannot be located, read, unmarshaled,
// or validated.
func ReadConfig[T any](path string, paths []string, c *T, unmarshal func([]byte, *T) error, validator Validator[T]) (string, []error) {
	var err error

	path, err = findConfig(path, paths)
	if err != nil {
		return "", []error{err}
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", []error{fmt.Errorf("unable to read configuration file %s: %w", path, err)}
	}

	err = unmarshal(content, c)
	if err != nil {
		return "", []error{fmt.Errorf("unable to unmarshal configuration file %s: %w", path, err)}
	}

	return path, validator.Validate(*c)
}

// findConfig determines the path to the configuration file by using the
// provided primary path or searching through a list of fallback paths if the
// primary path is empty. If no valid file is found, or if the specified file is
// inaccessible, an error is returned.
func findConfig(path string, paths []string) (string, error) {
	var err error

	if path == "" {
		for _, p := range paths {
			_, err = os.Stat(p)
			if err != nil {
				continue
			}

			path = p
		}

		if path == "" {
			return "", errors.New("could not locate configuration file")
		}
	} else {
		_, err = os.Stat(path)
		if err != nil {
			return "", fmt.Errorf("could not stat configuration file %s: %w", path, err)
		}
	}

	return path, nil
}
