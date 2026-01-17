package schema

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestFilePath(file string) string {
	return filepath.Join("..", "..", "test_data", file)
}

func resetSchemaCaches() {
	schemaCache = map[string]*Schema{}
	schemaCacheByName = map[string]*Schema{}
}

func TestLoadSchema_HappyPath_WithIncludes(t *testing.T) {
	resetSchemaCaches()

	path := getTestFilePath(filepath.Join("happy", "go-project.gro"))

	schema, err := LoadSchema(path)
	assert.NoError(t, err)
	assert.NotNil(t, schema)
	assert.Equal(t, "go-project", schema.Name)
	assert.Equal(t, "1GB", schema.Options.MaxSize)
	assert.Equal(t, path, schema.Path)

	// Loaded schemas
	for _, name := range []string{"go-project", "go-package", "go-command", "go-internal"} {
		assert.Contains(t, schemaCacheByName, name, "included schema %q not loaded", name)
	}

	// Require nodes
	requireByPattern := map[string]*Node{}
	for _, n := range schema.Require {
		requireByPattern[n.Pattern] = n
	}

	assert.Len(t, requireByPattern, 6)
	assert.Equal(t, "go-command", requireByPattern["cmd"].Schema.Name)
	assert.Equal(t, "go-internal", requireByPattern["internal"].Schema.Name)
	assert.Equal(t, "10MB", requireByPattern["README.md"].Options.MaxSize)

	// Allow nodes
	assert.Len(t, schema.Allow, 1)
	allow := schema.Allow[0]
	assert.Equal(t, "pkg", allow.Pattern)
	assert.Equal(t, "go-package", allow.Schema.Name)
	assert.Zero(t, allow.Options.MaxSize)

	// Deny nodes
	assert.Len(t, schema.Deny, 3)
	denyPatterns := map[string]PatternEngine{}
	for _, n := range schema.Deny {
		denyPatterns[n.Pattern] = n.Engine
	}
	assert.Equal(t, PatternGlob, denyPatterns["node_modules"])
	assert.Equal(t, PatternGlob, denyPatterns["*.exe"])
	assert.Equal(t, PatternRegex, denyPatterns["^temp_[0-9]+.bin$"])
}
