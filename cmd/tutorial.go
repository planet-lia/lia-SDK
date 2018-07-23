package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var debugMode bool

var tutorialCmd = &cobra.Command{
	Use:   "tutorial <number> <bot>",
	Short: "Runs tutorial specified by number with chosen bot",
	Long: `Runs tutorial specified by number with chosen bot`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tutorial called")
	},
}

func init() {
	rootCmd.AddCommand(tutorialCmd)

	tutorialCmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "Toggle if you want to manually run your bot (eg. " +
		"through debug mode in IDE).")
}
