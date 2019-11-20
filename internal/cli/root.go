package cli

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:   "indexer-cli",
		Short: "CLI for the navcoin explorer indexer",
	}
)

func Execute() error {
	return rootCmd.Execute()
}
