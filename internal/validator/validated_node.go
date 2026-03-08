package validator

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/Cliper27/grove/internal/parser"
)

type Color string

const (
	Reset  Color = "\033[0m"
	Green  Color = "\033[32m"
	Red    Color = "\033[31m"
	Yellow Color = "\033[33m"
	Blue   Color = "\033[34m"
)

type ValidatedNode struct {
	Path string // relative, slash-separated
	Type parser.NodeType

	Valid   bool
	Reasons []string

	// The rule that determined how this node was validated
	MatchedNode *parser.Node // nil if no rule matched

	Children []*ValidatedNode // only for folders
}

// JsonDumps converts the ValidatedNode (and its children) into a JSON string.
func (v *ValidatedNode) JsonDumps() (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// TreeDumps returns a tree-like string representation of the node and its children.
func (v *ValidatedNode) TreeDumps(colors bool) string {
	return v.treeDumps("", colors, false)
}

func (v *ValidatedNode) treeDumps(prefix string, colors bool, isLast bool) string {
	// Choose color based on validity
	color := Red
	status := "[✗]"
	if v.Valid {
		color = Green
		status = "[✔]"
	}

	if colors {
		status = fmt.Sprintf("%s%s%s", color, status, Reset)
	}

	// Include first reason if invalid
	reason := ""
	if !v.Valid && len(v.Reasons) > 0 {
		if colors {
			reason = fmt.Sprintf(" (%s%s%s)", Yellow, v.Reasons[0], Reset)
		} else {
			reason = fmt.Sprintf(" (%s)", v.Reasons[0])
		}
	}

	// Include description if available
	desc := ""
	if v.MatchedNode != nil {
		desc = v.MatchedNode.GetDescription()
	}
	if desc != "" {
		desc = "# " + desc
	}
	if colors {
		desc = fmt.Sprintf("%s%s%s", Blue, desc, Reset)
	}

	branch := "├─ "
	if isLast {
		branch = "└─ "
	}

	item := path.Base(v.Path)
	if v.Type == parser.NodeFolder {
		item += "/"
	}

	line := fmt.Sprintf("%s%s%s %s%s %s", prefix, branch, item, status, reason, desc)

	if len(v.Children) == 0 {
		return line
	}

	var builder strings.Builder
	builder.WriteString(line)

	if len(v.Children) > 0 {
		builder.WriteString("\n")
		for i, child := range v.Children {
			childPrefix := prefix
			if isLast {
				childPrefix += "   "
			} else {
				childPrefix += "│  "
			}
			builder.WriteString(child.treeDumps(childPrefix, colors, i == len(v.Children)-1))
			if i < len(v.Children)-1 {
				builder.WriteString("\n")
			}
		}
	}
	return builder.String()
}
