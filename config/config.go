package config

import (
	"encoding/json"
	"fmt"
	"github.com/palantir/stacktrace"
	"io/ioutil"
	"os"
	"path/filepath"
)

const VERSION = "0.1.0"

const ReleasesUrl = "https://github.com/liagame/lia-tools/releases/latest"

var cfg *Config

// Store path to this program executables
var PathToBots string
var PathToData string

type Config struct {
	Version    string     `json:"version"`
	GamePort   int        `json:"gamePort"`
	PathToBash string     `json:"windowsPathToBash"`
	Languages  []Language `json:"languages"`
}

type Language struct {
	Name           string `json:"name"`
	BotURL         string `json:"botUrl"`
	PrepareUnix    string `json:"prepareUnix"`
	RunUnix        string `json:"runUnix"`
	PrepareWindows string `json:"prepareWindows"`
	RunWindows     string `json:"runWindows"`
}

func SetConfig(path string) error {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return stacktrace.Propagate(err, "couldn't open file. Location: %s", path)
	}

	cfg = &Config{}
	if err := json.Unmarshal(configFile, cfg); err != nil {
		return stacktrace.Propagate(err, "couldn't unmarshal config")
	}

	return nil
}

func GetCfg() *Config {
	if cfg == nil {
		if PathToBots == "" {
			// Set PathToBots to executable path
			ex, err := os.Executable()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get executable location\n %s", err)
				os.Exit(FailedToGetEnvironment)
			}
			PathToBots = filepath.Dir(ex)
		}

		PathToData = filepath.Join(PathToBots, "data")

		pathToCfg := filepath.Join(PathToData, "cli-config.json")

		if err := SetConfig(pathToCfg); err != nil {
			fmt.Fprintf(os.Stderr, "couldn't get config\n %s", err)
			os.Exit(FailedToReadConfig)
		}
	}
	return cfg
}
