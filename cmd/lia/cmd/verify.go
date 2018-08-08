package cmd

import (
	"github.com/liagame/lia-cli/internal"
	"github.com/liagame/lia-cli/internal/analytics"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify <botDir>",
	Short: "Verifies if the content in bot-dir is valid",
	Long:  `Verifies if the content in bot-dir is valid.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Log("command", "verify", map[string]string{
			"botDir": args[0],
		})

		internal.UpdateIfTime(true)
		internal.GetBotLanguage(args[0])
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
