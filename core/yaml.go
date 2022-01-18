package core

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// QodanaYaml A standard qodana.yaml (or qodana.yml) format for Qodana configuration.
// https://github.com/JetBrains/qodana-profiles/blob/master/schemas/qodana-yaml-1.0.json
type QodanaYaml struct {
	// The qodana.yaml version of this log file.
	Version string `yaml:"version,omitempty"`

	// Linter to run.
	Linter string `yaml:"linter"`

	// Profile is the profile configuration for Qodana analysis.
	Profile Profile `yaml:"profile,omitempty"`

	// FailThreshold is a number of problems to fail the analysis (to exit from Qodana with code 255).
	FailThreshold int `yaml:"failThreshold,omitempty"`

	// Exclude property to disable the wanted checks on the wanted paths.
	Exclude []Exclude `yaml:"exclude,omitempty"`

	// Include property to enable the wanted checks.
	Include []Include `yaml:"include,omitempty"`

	// Properties property to override IDE properties
	Properties map[string]string `yaml:"properties,omitempty"`

	// Plugins property containing plugins to install
	Plugins []string `yaml:"plugins,omitempty"`
}

// Profile A profile is some template set of checks to run with Qodana analysis.
type Profile struct {
	Name string `yaml:"name,omitempty"`
	Path string `yaml:"path,omitempty"`
}

// Exclude A check id to disable.
type Exclude struct {
	// The name of check to exclude.
	Name string `yaml:"name"`

	// Relative to the project root path to disable analysis.
	Paths []string `yaml:"paths,omitempty"`
}

// Include A check id to enable.
type Include struct {
	// The name of check to exclude.
	Name string `yaml:"name"`
}

// GetQodanaYaml gets Qodana YAML from the project.
func GetQodanaYaml(project string) *QodanaYaml {
	q := &QodanaYaml{}
	qodanaYamlPath := filepath.Join(project, "qodana.yaml")
	if _, err := os.Stat(qodanaYamlPath); errors.Is(err, os.ErrNotExist) {
		qodanaYamlPath = filepath.Join(project, "qodana.yml")
	}
	if _, err := os.Stat(qodanaYamlPath); errors.Is(err, os.ErrNotExist) {
		return q
	}
	yamlFile, err := ioutil.ReadFile(qodanaYamlPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, q)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return q
}

// TODO: remove me
func (q *QodanaYaml) excludeDotQodana() {
	excluded := false
	for i, exclude := range q.Exclude {
		if exclude.Name == "All" {
			if !Contains(exclude.Paths, ".qodana/") {
				exclude.Paths = append(exclude.Paths, ".qodana/")
				q.Exclude[i] = exclude
			}
			excluded = true
			break
		}
	}
	if !excluded {
		q.Exclude = append(q.Exclude, Exclude{Name: "All", Paths: []string{".qodana/"}})
	}
}

// WriteQodanaYaml writes the qodana.yaml file to the given path.
func WriteQodanaYaml(path string, linters []string) {
	q := GetQodanaYaml(path)
	if q.Version == "" {
		q.Version = "1.0"
	}
	q.Linter = linters[0]
	q.excludeDotQodana()
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(&q)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filepath.Join(path, "qodana.yaml"), b.Bytes(), 0o644)
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}
}