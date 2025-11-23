package schema

import (
	"gopkg.in/yaml.v3"
)

type Folders struct {
	Mandatory []Folder `yaml:"mandatory"`
	Optional  []Folder `yaml:"optional"`
	Forbidden []string `yaml:"forbidden"`
}

type FolderRule struct {
	Schema  string `yaml:"schema"`
	MaxSize string `yaml:"max_size"`
}

type Folder struct {
	Name string
	Rule *FolderRule
}

func (f *Folder) Set(name string, rule any) {
	f.Name = name
	if r, ok := rule.(*FolderRule); ok {
		f.Rule = r
	} else {
		f.Rule = nil
	}
}

func (f *Folder) UnmarshalYAML(value *yaml.Node) error {
	return unmarshalNamedRule(value, f, &FolderRule{})
}
