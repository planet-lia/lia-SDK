package cmd

import (
	"github.com/spf13/cobra"
	"github.com/liagame/lia-cli/internal"
	"github.com/liagame/lia-cli/config"
)


var gameFlags = internal.RunGameFlags{}

var runCmd = &cobra.Command{
	Use:   "run [bot1] [bot2]",
	Short: "Run game instance",
	Long: `A longer description.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.RunGame(args, &gameFlags)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().IntVarP(&gameFlags.GameSeed,"gseed", "g", 0, "Game seed. 0 means random.")
	runCmd.Flags().IntVarP(&gameFlags.MapSeed,"mseed", "m", 0, "Map seed. 0 means random.")
	runCmd.Flags().IntVarP(&gameFlags.Port,"port", "p", config.GetCfg().GamePort, "Port on which game generator will run. Default is 8887.")
	runCmd.Flags().StringVarP(&gameFlags.MapPath, "map", "M", "", "Path to custom map settings.")
	runCmd.Flags().StringVarP(&gameFlags.ReplayPath, "replay", "r", "", "Choose custom replay name and location.")
	runCmd.Flags().StringVarP(&gameFlags.ConfigPath, "config", "c", config.GetCfg().GameConfigPath, "Choose custom replay name and location.")
	runCmd.Flags().IntSliceVarP(&gameFlags.DebugBots, "debug", "d", []int{},"Specify which bots you want to run manually." +
		"Examples: -d 1,2 -- debug bot1 and bot2, -d 2 -- debug bot2)")

}
