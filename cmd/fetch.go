package cmd

import (
	"fmt"

	"github.com/liagame/lia-cli/internal"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch <url> [name]",
	Short: "Fetches a bot from url and sets a new name.",
	Long: `Fetches a bot from specified url, unzips it into current folder and renames it if 
the argument is provided`,
	Args: func(cmd *cobra.Command, args []string) error {
		if !(len(args) == 1 || len(args) == 2) {
			return fmt.Errorf("there should be 1 or 2 arguments and not %d", len(args))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		internal.UpdateIfTime(true)

		url := args[0]
		name := ""
		if len(args) == 2 {
			name = args[1]
		}
		internal.FetchBot(url, name)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
