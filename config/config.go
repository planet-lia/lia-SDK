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

var cfg *Config

// Store path to this program executables
var PathToBots string
var PathToLia  string

type Config struct {
	Version           string     `json:"version"`
	GamePort          int        `json:"gamePort"`
	Languages         []Language `json:"languages"`
}

type Language struct {
	Name             string    `json:"name"`
	BotURL           string    `json:"botUrl"`
	PrepareUnix      []Command `json:"prepareUnix"`
	RunScriptUnix    []string  `json:"runScriptUnix"`
	PrepareWindows   []Command `json:"prepareWindows"`
	RunScriptWindows []string  `json:"runScriptWindows"`
}

type Command struct {
	Args []string `json:"cmdArgs"`
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
		// Set PathToBots to executable path
		ex, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get executable location\n %s", err)
			os.Exit(FAILED_TO_GET_ENVIRONMENT)
		}
		PathToBots = filepath.Dir(ex)
		PathToLia = filepath.Join(PathToBots, "lia")

		pathToCfg := filepath.Join(PathToLia, "cli-config.json")
		if err := SetConfig(pathToCfg); err != nil {
			fmt.Fprintf(os.Stderr, "couldn't get config\n %s", err)
			os.Exit(FAILED_TO_READ_CONFIG)
		}
	}
	return cfg
}
