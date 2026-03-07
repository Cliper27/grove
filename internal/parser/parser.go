package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	// ERR_CYCLIC_INCLUDE is returned when schemas include each other
	// directly or indirectly, forming a cycle.
	ERR_CYCLIC_INCLUDE = errors.New("cyclic include")

	// ERR_DUPLICATE_SCHEMA is returned when two schemas with the same
	// name are loaded.
	ERR_DUPLICATE_SCHEMA = errors.New("duplicate schema")

	// ERR_MISSING_INCLUDE is returned when a schema references another
	// schema that was not included.
	ERR_MISSING_INCLUDE = errors.New("missing include")
)

type Loader struct {
	// schemaCache stores loaded schemas keyed by their file path.
	//
	// It is used to avoid reloading and reparsing the same schema multiple times
	// during a single execution.
	schemaCache map[string]*Schema

	// schemaCacheByName stores loaded schemas keyed by schema name.
	//
	// It is used to detect duplicate schema names across included schemas.
	schemaCacheByName map[string]*Schema
}

func NewLoader() *Loader {
	return &Loader{
		schemaCache:       make(map[string]*Schema),
		schemaCacheByName: make(map[string]*Schema),
	}
}

type rawSchema struct {
	Include     []string `yaml:"include,omitempty"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description,omitempty"`

	Require map[string]rawNode `yaml:"require,omitempty"`
	Allow   map[string]rawNode `yaml:"allow,omitempty"`
	Deny    []string           `yaml:"deny,omitempty"`
}

type rawNode struct {
	Schema      string `yaml:"schema,omitempty"`
	Description string `yaml:"description,omitempty"`
}

func parsePattern(pattern string) (string, PatternEngine, NodeType) {
	engine := PatternGlob
	if strings.HasPrefix(pattern, RegexPrefix) {
		engine = PatternRegex
		pattern = strings.TrimPrefix(pattern, RegexPrefix)
	}

	nodeType := NodeFile
	if strings.HasSuffix(pattern, FolderSuffix) {
		nodeType = NodeFolder
		pattern = strings.TrimSuffix(pattern, FolderSuffix)
	}

	return pattern, engine, nodeType
}

func buildNodes(nodes map[string]rawNode, allowed map[string]*Schema) ([]*Node, error) {
	if len(nodes) == 0 {
		return nil, nil
	}

	var result []*Node

	for pattern, raw := range nodes {
		p, engine, typ := parsePattern(pattern)

		var schema *Schema
		if raw.Schema != "" {
			var ok bool
			schema, ok = allowed[raw.Schema]
			if !ok {
				return nil, fmt.Errorf("%w: %q", ERR_MISSING_INCLUDE, raw.Schema)
			}
		}

		n := &Node{
			Pattern:     p,
			Engine:      engine,
			Type:        typ,
			Schema:      schema,
			Description: raw.Description,
		}

		result = append(result, n)
	}

	return result, nil
}

func buildDenyNodes(patterns []string) []*Node {
	nodes := make([]*Node, 0, len(patterns))

	for _, raw := range patterns {
		pattern, engine, typ := parsePattern(raw)

		nodes = append(nodes, &Node{
			Pattern: pattern,
			Engine:  engine,
			Type:    typ,
		})
	}

	return nodes
}

func (l *Loader) parseSchema(path string, data []byte, loading map[string]bool) (*Schema, error) {
	var raw rawSchema
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	if raw.Name == "" {
		return nil, errors.New("schema name required")
	}

	allowed := map[string]*Schema{}

	base := filepath.Dir(path)
	for _, inc := range raw.Include {
		includePath := filepath.Join(base, filepath.Clean(inc))

		s, err := l.loadSchema(includePath, loading)
		if err != nil {
			return nil, err
		}

		allowed[s.Name] = s
	}

	schema := &Schema{
		Name:        raw.Name,
		Description: raw.Description,
		Path:        path,
	}

	allowed[schema.Name] = schema

	var err error
	schema.Require, err = buildNodes(raw.Require, allowed)
	if err != nil {
		return schema, err
	}

	schema.Allow, err = buildNodes(raw.Allow, allowed)
	if err != nil {
		return schema, err
	}

	schema.Deny = buildDenyNodes(raw.Deny)

	return schema, nil
}

// LoadSchema loads, parses, and resolves a schema from disk.
//
// The provided path may omit the ".gro" extension.
// Included schemas are resolved relative to the schema's directory.
//
// LoadSchema caches loaded schemas, detects include cycles, and ensures
// that schema names are globally unique.
func (l *Loader) LoadSchema(path string) (*Schema, error) {
	return l.loadSchema(path, map[string]bool{})
}

func (l *Loader) loadSchema(path string, loading map[string]bool) (*Schema, error) {
	path = filepath.Clean(path)

	if loading[path] {
		return nil, fmt.Errorf("%w: %q", ERR_CYCLIC_INCLUDE, path)
	}
	loading[path] = true
	defer func() { loading[path] = false }()

	if s, ok := l.schemaCache[path]; ok {
		return s, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	schema := &Schema{Path: path}
	l.schemaCache[path] = schema

	parsed, err := l.parseSchema(path, data, loading)
	if err != nil {
		delete(l.schemaCache, path)
		return nil, err
	}

	if existing, ok := l.schemaCacheByName[parsed.Name]; ok {
		return nil, fmt.Errorf("%w: %q (%s and %s)", ERR_DUPLICATE_SCHEMA, parsed.Name, existing.Path, path)
	}

	*schema = *parsed
	l.schemaCacheByName[parsed.Name] = schema
	schema.CompilePatterns()

	return schema, err
}
