package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
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
	}

	return pattern, engine, nodeType
}

func keys(m map[string]*Schema) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func buildNodes(
	nodes map[string]rawNode,
	allowed map[string]*Schema,
) ([]*Node, error) {
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
				return nil, fmt.Errorf(
					"schema %q is not included (allowed: %v)",
					raw.Schema,
					keys(allowed),
				)
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

func ParseSchema(path string, data []byte) (*Schema, error) {
	var raw rawSchema
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	allowed := map[string]*Schema{}

	base := filepath.Dir(path)
	for _, inc := range raw.Include {
		includePath := filepath.Join(base, inc)

		s, err := LoadSchema(includePath)
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
	if !strings.HasSuffix(path, ".gro") {
		path += ".gro"
	}

	if s, ok := schemaCache[path]; ok {
		return s, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	schema := &Schema{Path: path}
	schemaCache[path] = schema

	parsed, err := ParseSchema(path, data)
	if err != nil {
		delete(schemaCache, path)
		return nil, err
	}

	if existing, ok := schemaCacheByName[parsed.Name]; ok {
		return nil, fmt.Errorf(
			"duplicate schema name %q (%s and %s)",
			parsed.Name,
			existing.Path,
			path,
		)
	}

	*schema = *parsed
	schema.Path = path

	schemaCacheByName[parsed.Name] = schema

	return schema, err
}
