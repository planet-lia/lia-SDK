package config

import (
	"encoding/json"
	"fmt"
	"github.com/palantir/stacktrace"
	"io/ioutil"
	"path/filepath"
	"os"
)

const CLI_VERISON = "0.1.0"

var cfg *Config

// Store path to this program executables
var DirPath string

type Config struct {
	Version 	   string 		  `json:"version"`
	GamePort       int        `json:"gamePort"`
	GameConfigPath string     `json:"gameConfigPath"`
	Languages      []Language `json:"languages"`
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
		// Set DirPath to executable path
		ex, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get executable location\n %s", err)
			os.Exit(FAILED_TO_GET_ENVIRONMENT)
		}
		DirPath = filepath.Dir(ex)

		pathToCfg := DirPath + "/.lia/cli-config.json"
		if err := SetConfig(pathToCfg); err != nil {
			fmt.Fprintf(os.Stderr, "couldn't get config\n %s", err)
			os.Exit(FAILED_TO_READ_CONFIG)
		}
	}
	return cfg
}
