package cmd

import (
	"github.com/spf13/cobra"
	"github.com/liagame/lia-cli/internal"
)

var newbotCmd = &cobra.Command{
	Use:   "newbot [language] [name]",
	Short: "Create new bot.",
	Long: `Create new bot with specified language and chosen name.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		internal.FetchNewBot(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(newbotCmd)
}
