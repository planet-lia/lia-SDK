package cmd

import (
	"github.com/liagame/lia-cli/config"
	"github.com/liagame/lia-cli/internal"
	"github.com/spf13/cobra"
)

var viewReplay bool

var playCmd = &cobra.Command{
	Use:   "play <bot1Dir> <bot2Dir>",
	Short: "Compiles and generates a game between bot1 and bot2.",
	Long:  `Compiles and generates a game between bot1 and bot2. If config is not specified then default config will be used and 
if at least one of the bots is set to be in debug mode, the -debug.json config will be used.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.UpdateIfTime(true)
		internal.Play(args, &gameFlags, viewReplay)
	},
}

func init() {
	rootCmd.AddCommand(playCmd)

	playCmd.Flags().BoolVarP(&viewReplay, "viewReplay", "v", true, "If set, the replay will not be opened in replay viewer.")
	playCmd.Flags().IntVarP(&gameFlags.GameSeed, "gseed", "g", 0, "Game seed. 0 means random.")
	playCmd.Flags().IntVarP(&gameFlags.MapSeed, "mseed", "m", 0, "Map seed. 0 means random.")
	playCmd.Flags().IntVarP(&gameFlags.Port, "port", "p", config.GetCfg().GamePort, "Port on which game generator will run. Default is 8887.")
	playCmd.Flags().StringVarP(&gameFlags.MapPath, "map", "M", "", "Path to custom map settings.")
	playCmd.Flags().StringVarP(&gameFlags.ReplayPath, "replay", "r", "", "Choose custom replay name and location.")
	playCmd.Flags().StringVarP(&gameFlags.ConfigPath, "config", "c", "", "Choose custom replay name and location.")
	playCmd.Flags().IntSliceVarP(&gameFlags.DebugBots, "debug", "d", []int{}, "Specify which bots you want to run manually."+
		"Examples: -d 1,2 -- debug bot1 and bot2, -d 2 -- debug bot2)")
}
