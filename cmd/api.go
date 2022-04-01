package cmd

import (
	"github.com/galaxies-labs/airdrop/pkg/api"
	"github.com/spf13/cobra"
)

func CmdApi() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "api",
		Aliases:                    []string{"i"},
		Short:                      "Start api airdrop searcher",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE: func(cmd *cobra.Command, args []string) error {
			return api.NewApiServer()
		},
	}
	return cmd
}
