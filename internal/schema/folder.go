package schema

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Folders struct {
	Mandatory []Folder `yaml:"mandatory"`
	Optional  []Folder `yaml:"optional"`
	Forbidden []string `yaml:"forbidden"`
}

type Folder struct {
	Name string
	Rule *FolderRule
}

type FolderRule struct {
	Schema  string `yaml:"schema"`
	MaxSize string `yaml:"max_size"`
}

func (nf *Folder) UnmarshalYAML(value *yaml.Node) error {
	// If node is a string, just assign Name
	if value.Kind == yaml.ScalarNode {
		nf.Name = value.Value
		nf.Rule = nil
		return nil
	}

	// Otherwise it should be a map with one key
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("invalid folder node")
	}

	nf.Name = value.Content[0].Value
	var rule FolderRule
	if err := value.Content[1].Decode(&rule); err != nil {
		return err
	}
	nf.Rule = &rule
	return nil
}
