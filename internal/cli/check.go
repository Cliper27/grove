package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Cliper27/grove/internal/parser"
	"github.com/Cliper27/grove/internal/validator"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func NewCheckCmd() *cobra.Command {
	checkCmd := &cobra.Command{
		Use:               "check <dir> <schema>",
		Short:             "Validate a directory against a schema",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: completeCheckArgs,
		SilenceUsage:      false,

		RunE: func(cmd *cobra.Command, args []string) error {
			dir := args[0]
			schemaPath := args[1]

			outFormat, err := cmd.Flags().GetString("format")
			if err != nil {
				return err
			}
			switch outFormat {
			case "", "tree", "json":
			default:
				return fmt.Errorf("invalid format: %s (expected 'json' or 'tree')", outFormat)
			}

			outFile, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}

			quiet, err := cmd.Flags().GetBool("quiet")
			if err != nil {
				return err
			}

			noColor, err := cmd.Flags().GetBool("no-color")
			if err != nil {
				return err
			}

			if !isatty.IsTerminal(os.Stdout.Fd()) {
				noColor = true
			}

			cmd.SilenceUsage = true
			schema, err := parser.NewLoader().LoadSchema(schemaPath)
			if err != nil {
				return err
			}

			rootNode := validator.Validate(dir, schema)

			var outString string
			switch outFormat {
			case "tree":
				outString = rootNode.TreeDumps(!noColor)
			case "json":
				jsonStr, err := rootNode.JsonDumps()
				if err != nil {
					return err
				}
				outString = jsonStr
			default:
				if rootNode.Valid {
					outString = "✔ Valid"
				} else {
					var b strings.Builder
					b.WriteString("✗ Invalid\n")
					for _, reason := range rootNode.Reasons {
						fmt.Fprintf(&b, "  - %s\n", reason)
					}
					outString = strings.TrimSuffix(b.String(), "\n")
				}
			}

			if !quiet {
				fmt.Println(outString)
			}

			if outFile != "" {
				outString = ansiRegex.ReplaceAllString(outString, "")
				err := os.WriteFile(filepath.Clean(outFile), []byte(outString), 0644)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
	checkCmd.Flags().String("format", "", "Output format. Options are 'json' or 'tree'")
	checkCmd.RegisterFlagCompletionFunc("format", completeFormat)

	checkCmd.Flags().StringP("output", "o", "", "Output to specified file")
	checkCmd.Flags().BoolP("quiet", "q", false, "Suppress stdout")
	checkCmd.Flags().BoolP("no-color", "n", false, "Suppress cmd colors. Colors are automatically suppressed when not using a terminal that supports them.")
	return checkCmd
}

func completeCheckArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

	// First argument --> directories
	if len(args) == 0 {
		return nil, cobra.ShellCompDirectiveFilterDirs
	}

	// Second argument --> .gro files
	if len(args) == 1 {
		files, _ := filepath.Glob("*.gro")
		return files, cobra.ShellCompDirectiveDefault
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func completeFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"json", "tree"}, cobra.ShellCompDirectiveNoFileComp
}
