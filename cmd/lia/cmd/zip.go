package cmd

import (
	"github.com/liagame/lia-cli/internal"
	"github.com/liagame/lia-cli/internal/analytics"
	"github.com/spf13/cobra"
)

var zipCmd = &cobra.Command{
	Use:   "zip <botDir>",
	Short: "Verifies, compiles and zips the bot in botDir",
	Long:  `Verifies, compiles and zips the bot in botDir. Final zip can be uploaded to the website.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Log("command", "zip", map[string]string{
			"botDir": args[0],
		})

		internal.UpdateIfTime(true)
		internal.Zip(args[0])
	},
}

func init() {
	rootCmd.AddCommand(zipCmd)
}
