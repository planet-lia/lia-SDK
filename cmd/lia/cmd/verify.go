package cmd

import (
	"github.com/liagame/lia-SDK/internal"
	"github.com/liagame/lia-SDK/internal/analytics"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify <botDir>",
	Short: "Verifies if the content in bot-dir is valid",
	Long:  `Verifies if the content in bot-dir is valid.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		botDir := args[0]

		analytics.Log("command", "verify", map[string]string{
			"botDir": analytics.TrimPath(botDir),
		})

		internal.UpdateIfTime(true)
		internal.GetBotLanguage(botDir)
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
