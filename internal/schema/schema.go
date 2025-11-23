package schema

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Schema struct {
	Version     string  `yaml:"version"`
	Description string  `yaml:"description"`
	Folders     Folders `yaml:"folders"`
	Files       Files   `yaml:"files"`
}

type Schemas map[string]Schema

func LoadSchemas(path string) (Schemas, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var schemas Schemas
	if err := yaml.Unmarshal(data, &schemas); err != nil {
		return nil, err
	}

	return schemas, nil
}
