package schema_test

import (
	"path/filepath"
	"testing"

	"github.com/Cliper27/grove/internal/schema"
)

func TestLoadVersion(t *testing.T) {
	path := filepath.Join("..", "..", "test_data", "simple.gro")

	v, err := schema.LoadVersion(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "0.1.0"
	if v != expected {
		t.Fatalf("expected version %q, got %q", expected, v)
	}
}
