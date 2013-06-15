package rotten

import (
	"errors"
	"io/ioutil"
	"strings"
)

type Dependency struct {
	Package  string
	Operator string
	Version  string
}

func ParseRequirements(path string) ([]Dependency, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.New("Failed to read " + path)
	}

	lines := strings.Split(string(content), "\n")

	dependencies := []Dependency{}

	for _, line := range lines {
		parts := strings.Split(line, "==")

		if len(parts) != 2 {
			continue
		}

		dependencies = append(dependencies, Dependency{parts[0], "==", parts[1]})
	}

	return dependencies, nil
}
