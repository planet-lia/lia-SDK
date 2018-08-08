package cmd

import (
	"github.com/liagame/lia-cli/internal"
	"github.com/liagame/lia-cli/internal/analytics"
	"github.com/spf13/cobra"
)

var gameFlags = internal.GameFlags{}

var generateCmd = &cobra.Command{
	Use:   "generate <bot1Dir> <id1> <bot2Dir> <id2>",
	Short: "Generates a game",
	Long:  `Generates a game. This is a low level command.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Log("command", "generate", map[string]string{
			"bot1Dir":    args[0],
			"id1":        args[1],
			"bot2Dir":    args[2],
			"id2":        args[3],
			"viewReplay": analytics.ParseBoolFlagToString(cmd, "viewReplay"),
			"gseed":      analytics.ParseIntFlagToString(cmd, "gseed"),
			"mseed":      analytics.ParseIntFlagToString(cmd, "mseed"),
			"port":       analytics.ParseIntFlagToString(cmd, "port"),
			"map":        analytics.ParseStringFlag(cmd, "map"),
			"replay":     analytics.ParseStringFlag(cmd, "replay"),
			"config":     analytics.ParseStringFlag(cmd, "config"),
			"debug":      analytics.ParseIntSliceFlagToString(cmd, "debug"),
		})

		internal.UpdateIfTime(true)
		internal.GenerateGame(args[0], args[1], &gameFlags)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().IntVarP(&gameFlags.GameSeed, "gseed", "g", 0, "game seed. 0 means random")
	generateCmd.Flags().IntVarP(&gameFlags.MapSeed, "mseed", "m", 0, "map seed. 0 means random")
	generateCmd.Flags().IntVarP(&gameFlags.Port, "port", "p", 0, "port on which game generator will run. Default is 8887")
	generateCmd.Flags().StringVarP(&gameFlags.MapPath, "map", "M", "", "path to custom map settings")
	generateCmd.Flags().StringVarP(&gameFlags.ReplayPath, "replay", "r", "", "choose custom replay name and location")
	generateCmd.Flags().StringVarP(&gameFlags.ConfigPath, "config", "c", "", "choose custom config")
	generateCmd.Flags().IntSliceVarP(&gameFlags.DebugBots, "debug", "d", []int{}, "specify which bots you want to run manually, "+
		"examples: -d 1,2 -- debug bot1 and bot2, -d 2 -- debug bot2)")
}
