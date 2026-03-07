package validator

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/Cliper27/grove/internal/parser"
)

func findMatchingNode(entry fs.DirEntry, nodes []*parser.Node) *parser.Node {
	for _, n := range nodes {
		if match(entry, n) {
			return n
		}
	}
	return nil
}

func matchesAny(entries []fs.DirEntry, rule *parser.Node) bool {
	for _, entry := range entries {
		if match(entry, rule) {
			return true
		}
	}
	return false
}

func match(entry fs.DirEntry, rule *parser.Node) bool {
	var entryType parser.NodeType
	if entry.IsDir() {
		entryType = parser.NodeFolder
	} else {
		entryType = parser.NodeFile
	}

	if rule.Type != "" && entryType != rule.Type {
		return false
	}

	name := entry.Name()
	switch rule.Engine {
	case parser.PatternGlob, "":
		ok, err := filepath.Match(rule.Pattern, name)
		if err != nil {
			return false
		}
		return ok
	case parser.PatternRegex:
		if rule.CompiledPattern == nil {
			return false
		}
		return rule.CompiledPattern.MatchString(name)
	default:
		panic(fmt.Sprintf("unknown pattern engine: %q", rule.Engine))
	}
}
