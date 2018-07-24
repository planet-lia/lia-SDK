package cmd

import (
	"github.com/spf13/cobra"
	"github.com/liagame/lia-cli/internal"
	"github.com/liagame/lia-cli/config"
)


var gameFlags = internal.GameFlags{}

var generateCmd = &cobra.Command{
	Use:   "generate <bot1> <id1> <bot2> <id2>",
	Short: "Generates a game.",
	Long: `Generates a game. This is a low level command.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		internal.GenerateGame(args[0], args[1], &gameFlags)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().IntVarP(&gameFlags.GameSeed,"gseed", "g", 0, "Game seed. 0 means random.")
	generateCmd.Flags().IntVarP(&gameFlags.MapSeed,"mseed", "m", 0, "Map seed. 0 means random.")
	generateCmd.Flags().IntVarP(&gameFlags.Port,"port", "p", config.GetCfg().GamePort, "Port on which game generator will run. Default is 8887.")
	generateCmd.Flags().StringVarP(&gameFlags.MapPath, "map", "M", "", "Path to custom map settings.")
	generateCmd.Flags().StringVarP(&gameFlags.ReplayPath, "replay", "r", "", "Choose custom replay name and location.")
	generateCmd.Flags().StringVarP(&gameFlags.ConfigPath, "config", "c", "", "Choose custom replay name and location.")
	generateCmd.Flags().IntSliceVarP(&gameFlags.DebugBots, "debug", "d", []int{},"Specify which bots you want to run manually." +
		"Examples: -d 1,2 -- debug bot1 and bot2, -d 2 -- debug bot2)")
}
