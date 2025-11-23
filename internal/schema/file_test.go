package schema_test

import (
	"testing"

	"github.com/Cliper27/grove/internal/schema"
	"gopkg.in/yaml.v3"
)

func TestFile_UnmarshalYAML_Scalar(t *testing.T) {
	yamlStr := `"*.go"`

	var f schema.File
	if err := yaml.Unmarshal([]byte(yamlStr), &f); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if f.Name != "*.go" {
		t.Errorf("expected Name=*.go, got %s", f.Name)
	}
	if f.Rule != nil {
		t.Errorf("expected Rule=nil, got %+v", f.Rule)
	}
}

func TestFile_UnmarshalYAML_Map(t *testing.T) {
	yamlStr := `
"*.py":
  max_size: 10MB
`

	var f schema.File
	if err := yaml.Unmarshal([]byte(yamlStr), &f); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if f.Name != "*.py" {
		t.Errorf("expected Name=*.py, got %s", f.Name)
	}
	if f.Rule == nil || f.Rule.MaxSize != "10MB" {
		t.Errorf("unexpected Rule: %+v", f.Rule)
	}
}

func TestFiles_Sequence(t *testing.T) {
	yamlStr := `
- "*.py":
    max_size: 10MB
- "*.go"
`

	var fl []schema.File
	if err := yaml.Unmarshal([]byte(yamlStr), &fl); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(fl) != 2 {
		t.Fatalf("expected 2 files, got %d", len(fl))
	}
}
