package cmd

import (
	"github.com/liagame/lia-cli/internal"
	"github.com/liagame/lia-cli/internal/analytics"
	"github.com/spf13/cobra"
)

var checkForUpdate bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates Lia development tools",
	Long:  "Updates Lia development tools.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Log("command", "update", map[string]string{
			"check": analytics.ParseBoolFlagToString(cmd, "check"),
		})

		internal.Update(checkForUpdate)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVarP(&checkForUpdate, "check", "c", false,
		"if set the command will only check for updates but will not proceed with the update")
}
