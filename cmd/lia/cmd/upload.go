package cmd

import (
	"github.com/liagame/lia-SDK/internal"
	"github.com/liagame/lia-SDK/internal/analytics"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload <botDir>",
	Short: "Uploads the bot to Lia leaderboard",
	Long:  "Uploads the bot to Lia leaderboard",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Log("command", "upload", map[string]string{})

		botDir := args[0]

		internal.CheckForUpdate()
		internal.Upload(botDir)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
