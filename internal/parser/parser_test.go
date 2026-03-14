package parser

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestFilePath(file string) string {
	return filepath.Join("..", "..", "test_data", file)
}

func TestLoadSchema_HappyPath_WithIncludes(t *testing.T) {
	path := getTestFilePath(filepath.Join("happy", "go-project.gro"))

	loader := NewLoader()
	schema, err := loader.LoadSchema(path)
	assert.NoError(t, err)
	assert.NotNil(t, schema)
	assert.Equal(t, "go-project", schema.Name)
	assert.Equal(t, path, schema.Path)
	assert.Equal(t, "Standard Go project structure", schema.Description)

	// Loaded schemas
	for _, name := range []string{"go-project", "go-package", "go-command", "go-internal"} {
		assert.Contains(t, loader.schemaCacheByName, name, "included schema %q not loaded", name)
	}

	// Require nodes
	requireByPattern := map[string]*Node{}
	for _, n := range schema.Require {
		requireByPattern[n.Pattern] = n
	}

	assert.Len(t, requireByPattern, 6)
	assert.Equal(t, "go-command", requireByPattern["cmd"].Schema.Name)
	assert.Equal(t, "go-internal", requireByPattern["internal"].Schema.Name)
	assert.Equal(t, "Project Documentation", requireByPattern["README.md"].Description)
	assert.Equal(t, "Standard \"internal\" folder of a Go project", requireByPattern["internal"].Schema.Description)
	assert.Equal(t, "", requireByPattern[".gitignore"].Description)

	// Allow nodes
	assert.Len(t, schema.Allow, 1)
	allow := schema.Allow[0]
	assert.Equal(t, "pkg", allow.Pattern)
	assert.Equal(t, "go-package", allow.Schema.Name)

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

func TestLoadSchema_Failure_NonIncludedSchema(t *testing.T) {
	path := getTestFilePath(filepath.Join("failure", "missing-include.gro"))

	schema, err := NewLoader().LoadSchema(path)
	assert.Nil(t, schema)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ERR_MISSING_INCLUDE)
}

func TestLoadSchema_Failure_DuplicatePatter(t *testing.T) {
	path := getTestFilePath(filepath.Join("failure", "duplicate-pattern.gro"))

	schema, err := NewLoader().LoadSchema(path)
	assert.Nil(t, schema)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ERR_DUPLICATE_PATTERN)
}

func TestLoadSchema_Failure_DuplicateSchema(t *testing.T) {
	loader := NewLoader()
	// Load first schema (happy)
	path1 := getTestFilePath(filepath.Join("happy", "cmd.package", "go-package.gro"))
	s1, err := loader.LoadSchema(path1)
	assert.NoError(t, err)
	assert.NotNil(t, s1)

	// Load second schema with same name
	path2 := getTestFilePath(filepath.Join("failure", "duplicate-package.gro"))
	s2, err := loader.LoadSchema(path2)
	assert.Nil(t, s2)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ERR_DUPLICATE_SCHEMA)
}

func TestLoadSchema_Failure_IncludeCycle(t *testing.T) {
	path := getTestFilePath(filepath.Join("failure", "cycle-a.gro"))

	schema, err := NewLoader().LoadSchema(path)
	assert.Nil(t, schema)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ERR_CYCLIC_INCLUDE)
}
