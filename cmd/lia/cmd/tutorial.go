package cmd

import (
	"fmt"
	"github.com/liagame/lia-cli"
	"github.com/liagame/lia-cli/internal"
	"github.com/liagame/lia-cli/internal/analytics"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var debugMode bool

var tutorialCmd = &cobra.Command{
	Use:   "tutorial <number> <botDir>",
	Short: "Runs tutorial specified by number with chosen bot",
	Long:  `Runs tutorial specified by number with chosen bot.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Log("command", "tutorial", map[string]string{
			"number": args[0],
			"botDir": args[1],
			"debug":  analytics.ParseBoolFlagToString(cmd, "debug"),
			"width":  analytics.ParseStringFlag(cmd, "width"),
		})

		internal.UpdateIfTime(true)

		tutorialNumber, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to convert %s to number.\n %s\n", args[0], err)
			os.Exit(lia_cli.Generic)
		}
		botDir := args[1]

		internal.Tutorial(tutorialNumber, botDir, debugMode, replayViewerWidth)
	},
}

func init() {
	rootCmd.AddCommand(tutorialCmd)

	tutorialCmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "toggle if you want to manually run your bot (eg. "+
		"through debug mode in IDE)")
	tutorialCmd.Flags().StringVarP(&replayViewerWidth, "width", "w", "", "choose width of replay window,"+
		" height will be calcualted automatically")
}
