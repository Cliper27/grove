package cli

import (
	"os"

	"github.com/Cliper27/grove/internal/version"
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "grove",
		Short: "Validate project directory structure using schemas",
	}

	rootCmd.Version = version.Version
	rootCmd.SetVersionTemplate("{{.Version}}\n")

	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewCheckCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
