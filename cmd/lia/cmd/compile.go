package cmd

import (
	"github.com/liagame/lia-SDK/internal"
	"github.com/liagame/lia-SDK/internal/analytics"
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"github.com/liagame/lia-SDK"
)

var compileCmd = &cobra.Command{
	Use:   "compile <botDir>",
	Short: "Compiles/prepares bot in specified dir",
	Long:  `Compiles or prepares (depending on the language) the bot in specified dir.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		botDir := args[0]

		analytics.Log("command", "compile", map[string]string{})

		internal.CheckForUpdate()

		if err := internal.Compile(botDir); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(lia_SDK.PreparingBotFailed)
		}
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
