package schema_test

import (
	"path/filepath"
	"testing"

	"github.com/Cliper27/grove/internal/schema"
)

func TestLoadSchemas(t *testing.T) {
	path := filepath.Join("..", "..", "test_data", "simple.gro")
	schemas, err := schema.LoadSchemas(path)
	if err != nil {
		t.Fatalf("LoadSchemas error: %v", err)
	}

	if len(schemas) != 2 {
		t.Fatalf("expected 2 schemas, got %d", len(schemas))
	}

	code := schemas["code"]
	if code.Version != "0.1.0" {
		t.Errorf("expected version 0.1.0, got %s", code.Version)
	}
	if len(code.Folders.Mandatory) != 1 {
		t.Errorf("expected 1 mandatory folder, got %d", len(code.Folders.Mandatory))
	}

	module := schemas["module"]
	if module.Version != "0.2.0" {
		t.Errorf("expected version 0.2.0, got %s", module.Version)
	}
	if len(module.Folders.Forbidden) != 1 || module.Folders.Forbidden[0] != "*" {
		t.Errorf("forbidden folder mismatch")
	}
}
