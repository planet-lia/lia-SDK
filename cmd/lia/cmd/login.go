package cmd

import (
	"github.com/liagame/lia-SDK/internal"
	"github.com/liagame/lia-SDK/internal/analytics"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Lia with your account",
	Long:  "Login to Lia with your account.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Log("command", "login", map[string]string{})

		internal.CheckForUpdate()
		internal.Login()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
