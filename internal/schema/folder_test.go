package schema_test

import (
	"testing"

	"github.com/Cliper27/grove/internal/schema"
	"gopkg.in/yaml.v3"
)

func TestFolder_UnmarshalYAML_Scalar(t *testing.T) {
	yamlStr := `"lib"`

	var f schema.Folder
	if err := yaml.Unmarshal([]byte(yamlStr), &f); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if f.Name != "lib" {
		t.Errorf("expected Name=lib, got %s", f.Name)
	}
	if f.Rule != nil {
		t.Errorf("expected Rule=nil, got %+v", f.Rule)
	}
}

func TestFolder_UnmarshalYAML_Map(t *testing.T) {
	yamlStr := `
module:
  schema: module
  max_size: 10MB
`

	var f schema.Folder
	if err := yaml.Unmarshal([]byte(yamlStr), &f); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if f.Name != "module" {
		t.Errorf("expected Name=module, got %s", f.Name)
	}
	if f.Rule == nil || f.Rule.Schema != "module" || f.Rule.MaxSize != "10MB" {
		t.Errorf("unexpected Rule: %+v", f.Rule)
	}
}

func TestFolders_Sequence(t *testing.T) {
	yamlStr := `
- lib
- module:
    schema: module
    max_size: 10MB
`

	var fl []schema.Folder
	if err := yaml.Unmarshal([]byte(yamlStr), &fl); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(fl) != 2 {
		t.Fatalf("expected 2 folders, got %d", len(fl))
	}
}
