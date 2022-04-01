package cmd

import (
	"fmt"

	"github.com/galaxies-labs/airdrop/pkg/impt"
	"github.com/spf13/cobra"
)

func CmdImport() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "import [json_file_name] [type]",
		Aliases:                    []string{"i"},
		Short:                      "Import records json to sql(mysql)",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE: func(cmd *cobra.Command, args []string) error {
			file := args[0]
			fileType := args[1]
			if len(file) == 0 || len(fileType) == 0 {
				return fmt.Errorf("invalid args args[0] : %v, args[1] : %v", file, fileType)
			}
			return impt.ImportSnapshot(file, fileType)
		},
	}
	return cmd
}
