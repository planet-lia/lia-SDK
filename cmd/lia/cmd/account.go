package cmd

import (
	"github.com/liagame/lia-SDK/internal"
	"github.com/liagame/lia-SDK/internal/analytics"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Check which user is currently logged into Lia-SDK",
	Long:  "Check which user is currently logged into Lia-SDK.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Log("command", "account", map[string]string{})

		internal.CheckForUpdate()
		internal.CheckAccount()
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
