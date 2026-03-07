package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "check <dir> <schema>",
		Short:             "Validate a directory against a schema",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: completeCheckArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			dir := args[0]
			schema := args[1]

			fmt.Println("Directory:", dir)
			fmt.Println("Schema:", schema)

			return nil
		},
	}
}

func completeCheckArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

	// First argument → directories
	if len(args) == 0 {
		return nil, cobra.ShellCompDirectiveFilterDirs
	}

	// Second argument → .gro files
	if len(args) == 1 {
		files, _ := filepath.Glob("*.gro")
		return files, cobra.ShellCompDirectiveDefault
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}
