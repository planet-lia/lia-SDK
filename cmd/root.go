package cmd

import (
	"fmt"
	"os"

	"github.com/liagame/lia-cli/internal"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var showVersion bool
var showSupportedLanguages bool

var rootCmd = &cobra.Command{
	Use:   "lia-cli",
	Short: "The core Lia development tool",
	Long:  `lia-cli is a CLI tool for easier development of Lia bots.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.UpdateIfTime(true)

		if showVersion {
			internal.ShowVersions()
		} else if showSupportedLanguages {
			internal.ShowSupportedLanguages()
		} else {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show tools version")

	rootCmd.Flags().BoolVarP(&showSupportedLanguages, "languages", "l", false, "show all supported languages")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".lia-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".lia-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
