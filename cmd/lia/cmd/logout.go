package cmd

import (
	"github.com/liagame/lia-SDK/internal"
	"github.com/liagame/lia-SDK/internal/analytics"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Lia",
	Long:  "Logout from Lia.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Log("command", "logout", map[string]string{})

		internal.CheckForUpdate()
		internal.Logout()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
