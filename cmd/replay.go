package cmd

import (
	"github.com/spf13/cobra"
	"github.com/liagame/lia-cli/internal"
)

var replayCmd = &cobra.Command{
	Use:   "replay [pathToReplay]",
	Short: "Runs a replay viewer.",
	Long: `Runs a replay viewer. If path to the replay file is set as an
argument then that replay is played, else replay chooser is opened.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		replayFile := ""
		if len(args) == 1 {
			replayFile = args[0]
		}
		internal.ShowReplayViewer(replayFile)
	},
}

func init() {
	rootCmd.AddCommand(replayCmd)
}
