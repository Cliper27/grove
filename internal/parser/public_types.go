// Package parser implements loading, parsing, and validation of schema files.
//
// A schema defines filesystem rules using required, allowed, and denied nodes.
// Schemas may include other schemas explicitly via the `include` directive and
// reference them by name.
//
// Included schemas form a dependency graph which must be acyclic (a DAG).
// Cyclic includes and duplicate schema names are rejected at load time.
//
// Schema loading is path-based and cached globally within the package.
// Errors returned by this package are wrapped and may be inspected using errors.Is.
package parser

import (
	"fmt"
	"regexp"
)

var RegexPrefix string = "~"
var FolderSuffix string = "/"

// PatternEngine specifies how a node pattern is interpreted.
type PatternEngine string

const (
	// PatternGlob indicates that the pattern uses glob-style matching (default).
	PatternGlob PatternEngine = "glob"

	// PatternRegex indicates that the pattern uses regular expression matching.
	PatternRegex PatternEngine = "regex"
)

// NodeType defines whether a node applies to files or directories.
type NodeType string

const (
	// NodeFile represents a file node.
	NodeFile NodeType = "file"

	// NodeFolder represents a directory node.
	NodeFolder NodeType = "folder"

	// NodeSymlink represents a symbolic link node.
	NodeSymlink NodeType = "symlink"
)

// Node represents a single rule within a schema.
//
// A node matches filesystem entries using a pattern and pattern engine,
// applies optional constraints, and may reference another schema.
type Node struct {
	// Pattern is the raw pattern used to match paths.
	Pattern string

	// Compiled Regex pattern
	CompiledPattern *regexp.Regexp

	// Engine determines how the pattern is interpreted (glob or regex).
	Engine PatternEngine

	// Type indicates whether the node applies to files or folders.
	Type NodeType

	// Schema is an optional referenced schema.
	Schema *Schema

	// Brief optional description. Overrides the Schema description.
	Description string
}

// Schema represents a parsed and fully resolved schema definition.
//
// Schemas may include other schemas, forming a directed acyclic graph (DAG).
// All included schemas are loaded and validated before the schema is returned.
type Schema struct {
	// Name is the unique identifier of the schema.
	Name string

	// Path is the filesystem path from which the schema was loaded.
	Path string

	// Brief optional description for the schema.
	Description string

	// Require defines nodes that must be present.
	Require []*Node

	// Allow defines nodes that are permitted but not required.
	Allow []*Node

	// Deny defines nodes that are explicitly disallowed.
	Deny []*Node
}

// CompilePatterns precompiles regex patterns for performance
func (s *Schema) CompilePatterns() error {
	for _, group := range [][]*Node{s.Require, s.Allow, s.Deny} {
		for _, n := range group {
			if n.Engine == PatternRegex {
				re, err := regexp.Compile(n.Pattern)
				if err != nil {
					return fmt.Errorf("invalid regex pattern %q: %w", n.Pattern, err)
				}
				n.CompiledPattern = re
			}
		}
	}
	return nil
}
