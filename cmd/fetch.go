package cmd

import (
	"github.com/liagame/lia-cli/internal"
	"github.com/spf13/cobra"
)

var customBotDir string

var fetchCmd = &cobra.Command{
	Use:   "fetch <url> <name>",
	Short: "Fetches a bot from url and sets a new name",
	Long: `Fetches a bot from specified url, unzips it into current folder and renames it if 
the argument is provided`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		internal.UpdateIfTime(true)

		url := args[0]
		name := ""
		if len(args) == 2 {
			name = args[1]
		}
		internal.FetchBot(url, name, customBotDir)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	fetchCmd.Flags().StringVarP(&customBotDir, "dir", "d", "",
		"Specify the dir where the bot will be located")
}
