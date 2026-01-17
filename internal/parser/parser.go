package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ERR_CYCLIC_INCLUDE    = errors.New("cyclic include")
	ERR_DUPLICATE_SCHEMA  = errors.New("duplicate schema")
	ERR_MISSING_INCLUDE   = errors.New("missing include")
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

		n := &Node{
			Pattern: p,
			Engine:  engine,
			Type:    typ,
			Options: Options{
				MaxSize: raw.Options.MaxSize,
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

func ParseSchema(path string, data []byte, loading map[string]bool) (*Schema, error) {
	var raw rawSchema
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	allowed := map[string]*Schema{}

	base := filepath.Dir(path)
	for _, inc := range raw.Include {
		includePath := filepath.Join(base, filepath.FromSlash(inc))

		s, err := loadSchema(includePath, loading)
		if err != nil {
			return nil, err
		}

		allowed[s.Name] = s
	}

	schema := &Schema{
		Name: raw.Name,
		Options: Options{
			MaxSize: raw.Options.MaxSize,
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

func LoadSchema(path string) (*Schema, error) {
	return loadSchema(path, map[string]bool{})
}

func loadSchema(path string, loading map[string]bool) (*Schema, error) {
	if !strings.HasSuffix(path, ".gro") {
		path += ".gro"
	}

	if loading[path] {
		return nil, fmt.Errorf("%w: %q", ERR_CYCLIC_INCLUDE, path)
	}
	loading[path] = true
	defer func() { loading[path] = false }()

	if s, ok := schemaCache[path]; ok {
		return s, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	schema := &Schema{Path: path}
	schemaCache[path] = schema

	parsed, err := ParseSchema(path, data, loading)
	if err != nil {
		delete(schemaCache, path)
		return nil, err
	}

	if existing, ok := schemaCacheByName[parsed.Name]; ok {
		return nil, fmt.Errorf("%w: %q (%s and %s)", ERR_DUPLICATE_SCHEMA, parsed.Name, existing.Path, path)
	}

	*schema = *parsed
	schema.Path = path
	schemaCacheByName[parsed.Name] = schema

	return schema, err
}

func ParseByteUnits(s string) (uint64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return 0, fmt.Errorf("%w: empty string (expected number + optional unit B/KB/MB/GB/TB)", ERR_INVALID_BYTEUNITS)
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
	return value * mult, nil
}
