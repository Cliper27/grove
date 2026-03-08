package validator

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/Cliper27/grove/internal/parser"
)

func mkdir(t *testing.T, root, name string) {
	t.Helper()
	err := os.Mkdir(filepath.Join(root, name), 0755)
	if err != nil {
		t.Fatal(err)
	}
}

func touch(t *testing.T, root, name string) {
	t.Helper()
	err := os.WriteFile(filepath.Join(root, name), []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func fileNode(pattern string) *parser.Node {
	return &parser.Node{
		Pattern: pattern,
		Type:    parser.NodeFile,
		Engine:  parser.PatternGlob,
	}
}

func dirNode(pattern string, schema *parser.Schema) *parser.Node {
	return &parser.Node{
		Pattern: pattern,
		Type:    parser.NodeFolder,
		Engine:  parser.PatternGlob,
		Schema:  schema,
	}
}

func globFile(pattern string) *parser.Node {
	return &parser.Node{
		Pattern: pattern,
		Type:    parser.NodeFile,
		Engine:  parser.PatternGlob,
	}
}

func regexFile(pattern string) *parser.Node {
	return &parser.Node{
		Pattern:         pattern,
		Type:            parser.NodeFile,
		Engine:          parser.PatternRegex,
		CompiledPattern: regexp.MustCompile(pattern),
	}
}

func findChild(node *ValidatedNode, name string) *ValidatedNode {
	for _, child := range node.Children {
		if filepath.Base(child.Path) == name {
			return child
		}
	}
	return nil
}
