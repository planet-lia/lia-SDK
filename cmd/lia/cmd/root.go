package cmd

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/liagame/lia-SDK"
	"github.com/liagame/lia-SDK/internal"
	"github.com/liagame/lia-SDK/internal/config"
	"github.com/liagame/lia-SDK/internal/settings"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var cfgFile string
var showVersion bool
var showSupportedLanguages bool

var rootCmd = &cobra.Command{
	Use:   "lia",
	Short: "The core Lia development tool",
	Long:  `lia is a CLI tool for easier development of Lia bots.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.CheckForUpdate()

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
		os.Exit(lia_SDK.Generic)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.MousetrapHelpText = `This is a command line tool.

You need to open GitBash and run commands from there (eg. "./lia help").
`

	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show tools version")

	rootCmd.Flags().BoolVarP(&showSupportedLanguages, "languages", "l", false, "show all supported languages")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetConfigType(config.SettingsFileExtension)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(lia_SDK.LiaSettingsFailure)
			return
		}

		// Search config in home directory for lia settings file
		viper.AddConfigPath(home)
		viper.SetConfigName(config.SettingsFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using settings file:", viper.ConfigFileUsed())

	} else {
		// In case we don't find a lia settings file, create one
		// using the default parameters
		fmt.Println("Creating new settings file.")
		if err := settings.Create(); err != nil {
			fmt.Printf("Failed to create lia settings file. Error: %v\n", err)
			os.Exit(lia_SDK.LiaSettingsFailure)
			return
		}

		// Successfully created new settings file
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Failed to read from newly created settings file. Error: %v\n", err)
			os.Exit(lia_SDK.LiaSettingsFailure)
			return
		}
	}

	// If analyticsAllow is set to null in the lia settings file or if
	// the analyticsAllowedVersion does not match the config version,
	// ask the user for permission to collect anonymous analytics.
	if viper.Get("analyticsAllow") == nil || viper.Get("analyticsAllowedVersion") != config.VERSION {
		analyticsOptIn := askAnalyticsOptIn()
		viper.Set("analyticsAllow", analyticsOptIn)
		viper.Set("analyticsAllowedVersion", config.VERSION)
		viper.WriteConfig()
	}

	// Get the username of currently logged in user
	usernameWrapper := viper.Get("username")
	if usernameWrapper == nil {
		config.LoggedInUser = ""
	} else {
		username, ok := usernameWrapper.(string)
		if !ok {
			config.LoggedInUser = ""
		} else {
			config.LoggedInUser = username
		}
	}

	// Get token of currently logged in user
	tokenWrapper := viper.Get("token")
	if tokenWrapper == nil {
		config.UserToken = ""
	} else {
		token, ok := tokenWrapper.(string)
		if !ok {
			config.UserToken = ""
		} else {
			config.UserToken = token
		}
	}
}

// Ask users to decide if they want to opt-in to our
// anonymous usage tracking.
func askAnalyticsOptIn() (optIn bool) {
	optIn = true

Loop:
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Opt in to anonymous usage tracking [Y/n]: ")
		scanner.Scan()
		text := scanner.Text()

		switch strings.ToUpper(text) {
		case "N", "NO", "NAH":
			optIn = false
			fallthrough
		case "Y", "YES", "WOW", "NSA":
			break Loop
		}

	}

	if optIn {
		fmt.Printf("You have successfully opt in to anonymous usage analytics. Thank you for your feedback!\n\n")
	} else {
		fmt.Print("You have successfully opt out from anonymous usage analytics. ")
		fmt.Printf("To turn it on run %s\n\n", color.GreenString("./lia settings --analytics-opt-in"))
	}

	return optIn
}
