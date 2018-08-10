package cmd

import (
	"github.com/liagame/lia-cli/internal"
	"github.com/liagame/lia-cli/internal/analytics"
	"github.com/spf13/cobra"
)

var replayViewerWidth string

var replayCmd = &cobra.Command{
	Use:   "replay [pathToReplay]",
	Short: "Runs a replay viewer",
	Long: `Runs a replay viewer. If path to the replay file is set as an
argument then that replay is played, else replay chooser is opened.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Log("command", "replay", map[string]string{
			"pathToReplay": args[0],
			"width":        analytics.ParseStringFlag(cmd, "width"),
		})

		internal.UpdateIfTime(true)

		replayFile := ""
		if len(args) == 1 {
			replayFile = args[0]
		}
		internal.ShowReplayViewer(replayFile, replayViewerWidth)
	},
}

func init() {
	rootCmd.AddCommand(replayCmd)
	replayCmd.Flags().StringVarP(&replayViewerWidth, "width", "w", "", "choose width of replay window, height will be calcualted automatically")
}
