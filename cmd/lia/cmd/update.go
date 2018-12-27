package cmd

import (
	"github.com/liagame/lia-SDK/internal"
	"github.com/liagame/lia-SDK/internal/analytics"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates Lia-SDK",
	Long:  "Updates Lia-SDK.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Log("command", "update", map[string]string{
			"check": analytics.ParseBoolFlagToString(cmd, "check"),
		})

		internal.Update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
