package cmd

import (
	"fmt"
	"github.com/liagame/lia-SDK"
	"github.com/liagame/lia-SDK/internal"
	"github.com/liagame/lia-SDK/internal/analytics"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var debugMode bool

var playgroundCmd = &cobra.Command{
	Use:   "playground <number> <botDir>",
	Short: "Runs playground specified by number with chosen bot",
	Long:  `Runs playground specified by number with chosen bot. Number 1 represent a 1v1 battle and 
number 2 and 3 uses in house Lia bots as opponents in a normal match.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		numberStr := args[0]
		botDir := args[1]

		analytics.Log("command", "playground", map[string]string{
			"number": numberStr,
			"botDir": analytics.TrimPath(botDir),
			"debug":  analytics.ParseBoolFlagToString(cmd, "debug"),
			"width":  analytics.ParseStringFlag(cmd, "width"),
		})

		internal.UpdateIfTime(true)

		number, err := strconv.Atoi(numberStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to convert %s to number.\n %s\n", numberStr, err)
			os.Exit(lia_SDK.Generic)
		}

		internal.Playground(number, botDir, debugMode, true, replayViewerWidth)
	},
}

func init() {
	rootCmd.AddCommand(playgroundCmd)

	playgroundCmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "toggle if you want to manually run your bot (eg. "+
		"through debug mode in IDE)")
	playgroundCmd.Flags().StringVarP(&replayViewerWidth, "width", "w", "", "choose width of replay window,"+
		" height will be calculated automatically")
}
