package rotten

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Dependency struct {
	Package  string
	Operator string
	Current  string
	Latest   string
}

func (d *Dependency) Outdated() bool {
	if d.Latest == "" {
		return false
	}
	return d.Current != d.Latest
}

type PythonPackage struct {
	Info struct {
		Version string
	}
}

func (d *Dependency) Check() error {
	url := "http://pypi.python.org/pypi/" + d.Package + "/json"

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	p := &PythonPackage{}
	err = json.Unmarshal(body, &p)

	if err != nil {
		return err
	}

	d.Latest = p.Info.Version

	return nil
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

		dep := Dependency{Package: parts[0], Operator: "==", Current: parts[1]}
		dependencies = append(dependencies, dep)
	}

	return dependencies, nil
}

func CheckPythonPackageIndex(dependencies []Dependency) ([]Dependency, error) {
	var wg sync.WaitGroup

	for i, dep := range dependencies {
		wg.Add(1)
		go func(index int, dep Dependency) {
			defer wg.Done()

			err := dep.Check()

			if err != nil {
				log.Print(err)
			}

			// Super terrible
			dependencies[index].Latest = dep.Latest
		}(i, dep)
	}

	wg.Wait()

	return dependencies, nil
}
