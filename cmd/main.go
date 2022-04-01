package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gad",
		Short: "Galaxy Aidrop Checker",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())
			return nil
		},
	}

	rootCmd.AddCommand(
		CmdImport(),
	)

	return rootCmd
}
