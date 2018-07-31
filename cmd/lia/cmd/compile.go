package cmd

import (
	"github.com/liagame/lia-cli/internal"
	"github.com/liagame/lia-cli/internal/analytics"
	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use:   "compile <botDir>",
	Short: "Compiles/prepares bot in specified dir",
	Long:  `Compiles or prepares (depending on the language) the bot in specified dir.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Log("command", "compile", map[string]string{
			"botDir": args[0],
		})
		internal.UpdateIfTime(true)
		internal.Compile(args[0])
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
