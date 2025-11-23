package schema

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Files struct {
	Mandatory []File   `yaml:"mandatory"`
	Optional  []File   `yaml:"optional"`
	Forbidden []string `yaml:"forbidden"`
}

type File struct {
	Name string
	Rule *FileRule
}

type FileRule struct {
	MaxSize string `yaml:"max_size"`
}

func (nf *File) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		nf.Name = value.Value
		nf.Rule = nil
		return nil
	}

	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("invalid file node")
	}

	nf.Name = value.Content[0].Value
	var rule FileRule
	if err := value.Content[1].Decode(&rule); err != nil {
		return err
	}
	nf.Rule = &rule
	return nil
}
