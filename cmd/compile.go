package cmd

import (
	"github.com/spf13/cobra"
	"github.com/liagame/lia-cli/internal"
)

var compileCmd = &cobra.Command{
	Use:   "compile <dir>",
	Short: "Compiles/prepares bot in specified dir.",
	Long: `Compiles or prepares (depending on the language) the bot in specified dir.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		internal.Compile(args[0])
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
