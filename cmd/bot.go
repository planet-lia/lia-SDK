package cmd

import (
	"github.com/liagame/lia-cli/internal"
	"github.com/spf13/cobra"
)

var botCmd = &cobra.Command{
	Use:   "bot <language> <name>",
	Short: "Create new bot",
	Long:  `Create new bot with specified language and chosen name`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		internal.UpdateIfTime(true)
		internal.FetchBotByLanguage(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(botCmd)
}