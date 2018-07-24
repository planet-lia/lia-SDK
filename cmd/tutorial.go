package cmd

import (
	"github.com/spf13/cobra"
	"github.com/liagame/lia-cli/internal"
	"strconv"
	"fmt"
	"os"
	"github.com/liagame/lia-cli/config"
)

var debugMode bool

var tutorialCmd = &cobra.Command{
	Use:   "tutorial <number> <bot>",
	Short: "Runs tutorial specified by number with chosen bot",
	Long: `Runs tutorial specified by number with chosen bot`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tutorialNumber, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to convert %s to number.\n %s\n", args[0], err)
			os.Exit(config.GENERIC)
		}
		botDir := args[1]

		internal.Tutorial(tutorialNumber, botDir, debugMode)
	},
}

func init() {
	rootCmd.AddCommand(tutorialCmd)

	tutorialCmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "Toggle if you want to manually run your bot (eg. " +
		"through debug mode in IDE).")
}
