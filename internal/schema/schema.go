package schema

import (
	"os"

	"gopkg.in/yaml.v3"
)

type rawSchema struct {
	Version string `yaml:"version"`
}

// LoadVersion reads a .gro YAML file and returns its version.
func LoadVersion(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var rs rawSchema
	if err := yaml.Unmarshal(data, &rs); err != nil {
		return "", err
	}

	return rs.Version, nil
}
