package parser

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
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

	// ERR_INVALID_BYTEUNITS is returned when a byte-units string cannot
	// be parsed.
	ERR_INVALID_BYTEUNITS = errors.New("invalid byte units")

	byteUnitsRegex = regexp.MustCompile(`^(\d+)(B|KB|MB|GB|TB)?$`)
	mults          = map[string]uint64{
		"":   1,
		"B":  1,
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
	}
)

type rawSchema struct {
	Include []string   `yaml:"include,omitempty"`
	Name    string     `yaml:"name"`
	Options rawOptions `yaml:",inline"`

	Require map[string]rawNode `yaml:"require,omitempty"`
	Allow   map[string]rawNode `yaml:"allow,omitempty"`
	Deny    []string           `yaml:"deny,omitempty"`
}

type rawNode struct {
	Schema  string     `yaml:"schema,omitempty"`
	Options rawOptions `yaml:",inline"`
}

type rawOptions struct {
	MaxSize string `yaml:"maxSize,omitempty"`
}

func parsePattern(pattern string) (string, PatternEngine, NodeType) {
	engine := PatternGlob
	if strings.HasPrefix(pattern, "~") {
		engine = PatternRegex
		pattern = strings.TrimPrefix(pattern, "~")
	}

	nodeType := NodeFile
	if strings.HasSuffix(pattern, "/") {
		nodeType = NodeFolder
		pattern = strings.TrimSuffix(pattern, "/")
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

		maxSize, maxSizeErr := ParseByteUnits(raw.Options.MaxSize)
		if maxSizeErr != nil {
			return nil, maxSizeErr
		}
		n := &Node{
			Pattern: p,
			Engine:  engine,
			Type:    typ,
			Options: Options{
				MaxSize: maxSize,
			},
			Schema: schema,
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

	maxSize, maxSizeErr := ParseByteUnits(raw.Options.MaxSize)
	if maxSizeErr != nil {
		return nil, maxSizeErr
	}
	if raw.Name == "" {
		return nil, errors.New("schema name required")
	}
	schema := &Schema{
		Name: raw.Name,
		Options: Options{
			MaxSize: maxSize,
		},
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
	if !strings.HasSuffix(path, ".gro") {
		path += ".gro"
	}

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
	schema.Path = path
	l.schemaCacheByName[parsed.Name] = schema
	schema.CompilePatterns()

	return schema, err
}

// ParseByteUnits parses a human-readable byte size string.
//
// Valid formats are an unsigned integer followed by an optional unit:
// B, KB, MB, GB, or TB. Units are case-insensitive.
//
// Examples:
//
//	"10"     -> 10
//	"1KB"    -> 1024
//	"5mb"    -> 5242880
//
// An error wrapping ERR_INVALID_BYTEUNITS is returned for invalid input.
func ParseByteUnits(s string) (uint64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return 0, nil
	}

	matches := byteUnitsRegex.FindStringSubmatch(s)
	if matches == nil {
		return 0, fmt.Errorf("%w: invalid format %q (expected number + optional unit B/KB/MB/GB/TB)", ERR_INVALID_BYTEUNITS, s)
	}

	value, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: cannot parse number %q: %v", ERR_INVALID_BYTEUNITS, matches[1], err)
	}

	unit := matches[2]
	mult, ok := mults[unit]
	if !ok {
		return 0, fmt.Errorf("%w: unknown unit %q", ERR_INVALID_BYTEUNITS, unit)
	}
	if value > math.MaxUint64/mult {
		return 0, fmt.Errorf("%w: overflow", ERR_INVALID_BYTEUNITS)
	}
	return value * mult, nil
}
