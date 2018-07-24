package cmd

import (
	"github.com/liagame/lia-cli/internal"
	"github.com/spf13/cobra"
)

var checkForUpdate bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates Lia development tools",
	Long:  "Updates Lia development tools",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		internal.Update(checkForUpdate)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVarP(&checkForUpdate, "check", "c", false,
		"If set the command will only check for updates but will not proceed with the update")
}
