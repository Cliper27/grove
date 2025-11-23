package schema_test

import (
	"path/filepath"
	"testing"

	"github.com/Cliper27/grove/internal/schema"
)

func getTestFilePath() string {
	return filepath.Join("..", "..", "test_data", "simple.gro")
}

// --- Test: YAML file can be loaded ---
func TestLoadSchemas_FileExists(t *testing.T) {
	path := getTestFilePath()

	_, err := schema.LoadSchemas(path)
	if err != nil {
		t.Fatalf("LoadSchemas error: %v", err)
	}
}

// --- Test: correct number of schemas ---
func TestLoadSchemas_SchemaCount(t *testing.T) {
	path := getTestFilePath()
	schemas, _ := schema.LoadSchemas(path)

	if len(schemas) != 2 {
		t.Fatalf("expected 2 schemas, got %d", len(schemas))
	}
}

// --- Test: code schema fields ---
func TestLoadSchemas_CodeSchema(t *testing.T) {
	path := getTestFilePath()
	schemas, _ := schema.LoadSchemas(path)

	code := schemas["code"]

	if code.Version != "0.1.0" {
		t.Errorf("expected version 0.1.0, got %s", code.Version)
	}

	if code.Description != "Source code layout" {
		t.Errorf("unexpected description: %s", code.Description)
	}

	// Folders
	if len(code.Folders.Mandatory) != 1 {
		t.Errorf("expected 1 mandatory folder, got %d", len(code.Folders.Mandatory))
	}
	if code.Folders.Mandatory[0].Name != "module" {
		t.Errorf("expected mandatory folder 'module', got '%s'", code.Folders.Mandatory[0].Name)
	}
	if code.Folders.Mandatory[0].Rule == nil || code.Folders.Mandatory[0].Rule.Schema != "module" {
		t.Errorf("mandatory folder rule incorrect")
	}

	// Optional folders
	if len(code.Folders.Optional) != 1 {
		t.Errorf("expected 1 optional folder, got %d", len(code.Folders.Optional))
	}
	if code.Folders.Optional[0].Name != "lib" {
		t.Errorf("expected optional folder 'lib', got '%s'", code.Folders.Optional[0].Name)
	}
	if code.Folders.Optional[0].Rule != nil {
		t.Errorf("optional folder 'lib' should have no rule")
	}

	if len(code.Folders.Forbidden) != 1 || code.Folders.Forbidden[0] != "tmp" {
		t.Errorf("forbidden folders incorrect")
	}

	// Files
	if len(code.Files.Mandatory) != 0 {
		t.Errorf("expected 0 mandatory files, got %d", len(code.Files.Mandatory))
	}

	// Optional files
	if len(code.Files.Optional) != 2 {
		t.Errorf("expected 2 optional files, got %d", len(code.Files.Optional))
	}

	py := code.Files.Optional[0]
	if py.Name != "*.py" || py.Rule == nil || py.Rule.MaxSize != "10MB" {
		t.Errorf("optional file '*.py' incorrect")
	}

	goFile := code.Files.Optional[1]
	if goFile.Name != "*.go" || goFile.Rule != nil {
		t.Errorf("optional file '*.go' should have no rule")
	}

	if len(code.Files.Forbidden) != 1 || code.Files.Forbidden[0] != "*.exe" {
		t.Errorf("forbidden files incorrect")
	}
}

// --- Test: module schema fields ---
func TestLoadSchemas_ModuleSchema(t *testing.T) {
	path := getTestFilePath()
	schemas, _ := schema.LoadSchemas(path)

	module := schemas["module"]

	if module.Version != "0.2.0" {
		t.Errorf("expected version 0.2.0, got %s", module.Version)
	}

	if module.Description != "Module layout" {
		t.Errorf("unexpected description: %s", module.Description)
	}

	if len(module.Folders.Mandatory) != 0 {
		t.Errorf("expected 0 mandatory folders, got %d", len(module.Folders.Mandatory))
	}

	if len(module.Folders.Forbidden) != 1 || module.Folders.Forbidden[0] != "*" {
		t.Errorf("forbidden folders incorrect")
	}

	if len(module.Folders.Optional) != 0 {
		t.Errorf("expected 0 optional folders, got %d", len(module.Folders.Optional))
	}
}
