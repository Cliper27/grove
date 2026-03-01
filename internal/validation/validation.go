package validation

import (
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/Cliper27/grove/internal/parser"
)

type NodeValidation struct {
	Path string // relative, slash-separated
	Type parser.NodeType

	Valid   bool
	Reasons []string

	// The rule that determined how this node was validated
	MatchedNode *parser.Node // nil if no rule matched

	Children []*NodeValidation // only for folders
}

func checkRequired(entries []fs.DirEntry, requiredNodes []*parser.Node, node *NodeValidation) {
	for _, require := range requiredNodes {
		if !matchesAny(entries, require) {
			node.Valid = false
			node.Reasons = append(
				node.Reasons,
				fmt.Sprintf("Missing Required %s: %q", require.Type, require.Pattern),
			)
		}
	}
}

func checkDenied(entries []fs.DirEntry, deniedNodes []*parser.Node, node *NodeValidation) {
	for _, deny := range deniedNodes {
		if matchesAny(entries, deny) {
			node.Valid = false
			node.Reasons = append(
				node.Reasons,
				fmt.Sprintf("Denied %s found: %q", deny.Type, deny.Pattern),
			)
		}
	}
}

func Validate(root string, schema *parser.Schema) *NodeValidation {
	return ValidateFS(os.DirFS(root), schema)
}

func ValidateFS(fsys fs.FS, schema *parser.Schema) *NodeValidation {
	return validateDir(fsys, ".", schema)
}

func getNextNode(entry fs.DirEntry, schema *parser.Schema) *parser.Node {
	if n := findMatchingNode(entry, schema.Deny); n != nil {
		return n
	}
	if n := findMatchingNode(entry, schema.Require); n != nil {
		return n
	}
	return findMatchingNode(entry, schema.Allow)
}

func validateFile(fsys fs.FS, dir string, schema *parser.Schema) *NodeValidation {
	// TODO: check maxSize
	return &NodeValidation{
		Path:  dir,
		Type:  parser.NodeFile,
		Valid: true,
	}
}

func validateDir(fsys fs.FS, dir string, schema *parser.Schema) *NodeValidation {
	node := &NodeValidation{
		Path:  dir,
		Type:  parser.NodeFolder,
		Valid: true,
	}

	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		node.Valid = false
		node.Reasons = append(node.Reasons, err.Error())
		return node
	}

	filteredEntries := make([]fs.DirEntry, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		childPath := path.Join(dir, name)
		if entry.Type()&fs.ModeSymlink != 0 {
			child := &NodeValidation{
				Path:    childPath,
				Type:    parser.NodeSymlink,
				Valid:   true,
				Reasons: []string{"Symlinks not checked"},
			}
			node.Children = append(node.Children, child)
			continue
		}
		filteredEntries = append(filteredEntries, entry)
	}

	checkDenied(filteredEntries, schema.Deny, node)
	checkRequired(filteredEntries, schema.Require, node)

	for _, entry := range filteredEntries {
		childPath := path.Join(dir, entry.Name())

		matched := getNextNode(entry, schema)
		if matched == nil {
			var nodeType parser.NodeType
			if entry.IsDir() {
				nodeType = parser.NodeFolder
			} else {
				nodeType = parser.NodeFile
			}
			node.Children = append(node.Children, &NodeValidation{
				Path:  childPath,
				Type:  nodeType,
				Valid: true,
			})
			continue
		}

		// TODO: check maxSize
		nextSchema := matched.Schema
		if nextSchema == nil {
			node.Children = append(node.Children, &NodeValidation{
				Path:        childPath,
				Type:        matched.Type,
				Valid:       true,
				MatchedNode: matched,
			})
			continue
		}

		var child *NodeValidation
		if entry.IsDir() {
			child = validateDir(fsys, childPath, nextSchema)
		} else {
			child = validateFile(fsys, childPath, nextSchema)
		}
		child.MatchedNode = matched
		if !child.Valid {
			node.Valid = false
		}
		node.Children = append(node.Children, child)
	}

	return node
}
