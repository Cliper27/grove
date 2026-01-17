package parser

var schemaCache = map[string]*Schema{}
var schemaCacheByName = map[string]*Schema{}

type Options struct {
	MaxSize string
}

type PatternEngine string

const (
	PatternGlob  PatternEngine = "glob"
	PatternRegex PatternEngine = "regex"
)

type NodeType string

const (
	NodeFile   NodeType = "file"
	NodeFolder NodeType = "folder"
)

type Node struct {
	Pattern string
	Engine  PatternEngine
	Type    NodeType

	Schema  *Schema
	Options Options // overrides schema options
}

// func (n *Node) Search(root string) (string, error)

type Schema struct {
	Name    string
	Path    string
	Options Options

	Require []*Node
	Allow   []*Node
	Deny    []*Node
}
