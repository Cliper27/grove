package validator

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"runtime"
	"sync"

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
	sem := make(chan struct{}, runtime.NumCPU())
	return validateDir(fsys, ".", schema, sem)
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

func validateFile(entry fs.DirEntry, filePath string, schema *parser.Schema) *NodeValidation {
	result := &NodeValidation{
		Path:  filePath,
		Type:  parser.NodeFile,
		Valid: true,
	}
	if schema != nil && schema.Options.MaxSize > 0 && !entry.IsDir() {
		ok, err := checkMaxSize(entry, schema.Options.MaxSize)
		if err != nil || !ok {
			result.Valid = false
			result.Reasons = append(result.Reasons,
				fmt.Sprintf("File exceeds maxSize: %s", filePath),
			)
		}
	}
	return result
}

func validateDir(fsys fs.FS, dir string, schema *parser.Schema, sem chan struct{}) *NodeValidation {
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

	children := make([]*NodeValidation, len(filteredEntries))
	var wg sync.WaitGroup
	for i, entry := range filteredEntries {
		sem <- struct{}{}
		wg.Add(1)
		go func(i int, entry fs.DirEntry) {
			defer wg.Done()
			defer func() { <-sem }()

			childPath := path.Join(dir, entry.Name())
			matched := getNextNode(entry, schema)

			var result *NodeValidation

			if matched == nil {
				var nodeType parser.NodeType
				if entry.IsDir() {
					nodeType = parser.NodeFolder
				} else {
					nodeType = parser.NodeFile
				}
				result = &NodeValidation{
					Path:  childPath,
					Type:  nodeType,
					Valid: true,
				}
			} else {
				nextSchema := matched.Schema
				if nextSchema == nil {
					isDenied := findMatchingNode(entry, schema.Deny) != nil
					result = &NodeValidation{
						Path:        childPath,
						Type:        matched.Type,
						Valid:       !isDenied,
						MatchedNode: matched,
					}
					if isDenied {
						result.Reasons = append(result.Reasons,
							fmt.Sprintf("Denied %s found: %q", matched.Type, matched.Pattern),
						)
					}
					if matched.Options.MaxSize > 0 && !entry.IsDir() {
						ok, err := checkMaxSize(entry, matched.Options.MaxSize)
						if err != nil {
							result.Valid = false
							result.Reasons = append(result.Reasons, err.Error())
						} else if !ok {
							result.Valid = false
							result.Reasons = append(result.Reasons,
								fmt.Sprintf("File exceeds maxSize (%d bytes)", matched.Options.MaxSize),
							)
						}
					}
				} else {
					if entry.IsDir() {
						result = validateDir(fsys, childPath, nextSchema, sem)
					} else {
						result = validateFile(entry, childPath, nextSchema)
					}
					result.MatchedNode = matched
				}
			}

			children[i] = result
		}(i, entry)
	}

	wg.Wait()

	for _, child := range children {
		if !child.Valid {
			node.Valid = false
		}
		node.Children = append(node.Children, child)
	}

	if schema.Options.MaxSize > 0 {
		size, err := dirSize(fsys, dir)
		if err != nil {
			node.Valid = false
			node.Reasons = append(node.Reasons, err.Error())
		} else if size > schema.Options.MaxSize {
			node.Valid = false
			node.Reasons = append(node.Reasons,
				fmt.Sprintf("Directory exceeds maxSize (%d bytes)", schema.Options.MaxSize),
			)
		}
	}

	return node
}
