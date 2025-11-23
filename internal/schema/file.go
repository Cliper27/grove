package schema

import (
	"gopkg.in/yaml.v3"
)

type Files struct {
	Mandatory []File   `yaml:"mandatory"`
	Optional  []File   `yaml:"optional"`
	Forbidden []string `yaml:"forbidden"`
}

type FileRule struct {
	MaxSize string `yaml:"max_size"`
}

type File struct {
	Name string
	Rule *FileRule
}

func (f *File) Set(name string, rule any) {
	f.Name = name
	if r, ok := rule.(*FileRule); ok {
		f.Rule = r
	} else {
		f.Rule = nil
	}
}

func (f *File) UnmarshalYAML(value *yaml.Node) error {
	return unmarshalNamedRule(value, f, &FileRule{})
}
