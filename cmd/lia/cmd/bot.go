package cmd

import (
	"github.com/liagame/lia-SDK/internal"
	"github.com/liagame/lia-SDK/internal/analytics"
	"github.com/spf13/cobra"
)

var botCmd = &cobra.Command{
	Use:   "bot <language> <name>",
	Short: "Create new bot",
	Long:  `Create new bot with specified language and chosen name.
Run "./lia --languages" to view all supported languages.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Log("command", "bot", map[string]string{
			"language": args[0],
		})

		internal.CheckForUpdate()
		internal.FetchBotByLanguage(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(botCmd)
}
