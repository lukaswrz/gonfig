// Package gonfig provides utilities for reading, unmarshaling, and validating
// configuration files.
package gonfig

import (
	"errors"
	"fmt"
	"os"
)

// UnmarshalFunc is a function that unmarshals raw bytes into a configuration
// object of type T.
type UnmarshalFunc[T any] func([]byte, T) error

// ValidateFunc is a function that validates a configuration object of type T
// and returns an error if validation fails.
type ValidateFunc[T any] func(T) error

// ReadConfig reads a configuration file, unmarshals its content into the given
// configuration object, and validates it.
//
// If the primary path is empty, it searches for the configuration file in the
// fallback paths. Returns the resolved path or an error if the file cannot be
// located, read, unmarshaled, or validated.
func ReadConfig[T any](path string, searchPaths []string, c *T, unmarshal UnmarshalFunc[*T], validate ValidateFunc[T]) (string, error) {
	var err error

	path, err = FindConfig(path, searchPaths)
	if err != nil {
		return "", err
	}

	return path, ReadFoundConfig(path, c, unmarshal, validate)
}

// ReadFoundConfig reads and processes a configuration file from a known path.
//
// Unmarshals the file's content into the given configuration object and
// validates it. Returns an error if the file cannot be read, unmarshaled, or
// validated.
func ReadFoundConfig[T any](path string, c *T, unmarshal UnmarshalFunc[*T], validate ValidateFunc[T]) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read configuration file %s: %w", path, err)
	}

	err = unmarshal(content, c)
	if err != nil {
		return fmt.Errorf("unable to unmarshal configuration file %s: %w", path, err)
	}

	return validate(*c)
}

// FindConfig determines the path to the configuration file by using the
// provided primary path or searching through a list of fallback paths if the
// primary path is empty.
//
// Returns the resolved path or an error if no valid file is found or if the
// file is inaccessible.
func FindConfig(path string, paths []string) (string, error) {
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
