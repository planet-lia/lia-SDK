package cmd

import (
	"github.com/spf13/cobra"
	"github.com/liagame/lia-cli/internal"
)

var verifyCmd = &cobra.Command{
	Use:   "verify <bot-dir>",
	Short: "Verifies if the content in bot-dir is valid.",
	Long: `Verifies if the content in bot-dir is valid.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		internal.GetBotLanguage(args[0])
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
