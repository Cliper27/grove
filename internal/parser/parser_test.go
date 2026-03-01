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
	assert.Equal(t, uint64(1*1024*1024*1024), schema.Options.MaxSize)
	assert.Equal(t, path, schema.Path)

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
	assert.Equal(t, uint64(10*1024*1024), requireByPattern["README.md"].Options.MaxSize)

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

func TestLoadSchema_Failure_NonIncludedSchema(t *testing.T) {
	path := getTestFilePath(filepath.Join("failure", "missing-include.gro"))

	schema, err := NewLoader().LoadSchema(path)
	assert.Nil(t, schema)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ERR_MISSING_INCLUDE)
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

func TestParseByteUnits(t *testing.T) {
	tests := []struct {
		input       string
		wantBytes   uint64
		expectError bool
	}{
		// Valid inputs
		{"1B", 1, false},
		{"1024B", 1024, false},
		{"1KB", 1024, false},
		{"2KB", 2048, false},
		{"1MB", 1024 * 1024, false},
		{"1GB", 1024 * 1024 * 1024, false},
		{"1TB", 1024 * 1024 * 1024 * 1024, false},
		{"10", 10, false},            // No unit defaults to bytes
		{"  5kb  ", 5 * 1024, false}, // spaces and lowercase
		{"0B", 0, false},

		// Invalid inputs
		{"", 0, false},     // empty
		{"abc", 0, true},   // invalid format
		{"123XB", 0, true}, // unknown unit
		{"1.5MB", 0, true}, // decimal not allowed
		{"-10B", 0, true},  // negative not allowed
	}

	for _, tt := range tests {
		got, err := ParseByteUnits(tt.input)
		if tt.expectError {
			assert.Error(t, err, "input: %q", tt.input)
			assert.Equal(t, uint64(0), got, "input: %q", tt.input)
		} else {
			assert.NoError(t, err, "input: %q", tt.input)
			assert.Equal(t, tt.wantBytes, got, "input: %q", tt.input)
		}
	}
}
