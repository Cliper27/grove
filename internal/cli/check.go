package cli

import (
	"fmt"
	"path/filepath"

	"github.com/Cliper27/grove/internal/parser"
	"github.com/Cliper27/grove/internal/validator"
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
			schemaPath := args[1]

			schema, err := parser.NewLoader().LoadSchema(schemaPath)
			if err != nil {
				return err
			}

			rootNode := validator.Validate(dir, schema)
			fmt.Println("Is valid:", rootNode.Valid)

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
