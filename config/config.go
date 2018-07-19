package config

import (
	"encoding/json"
	"fmt"
	"github.com/palantir/stacktrace"
	"io/ioutil"
)

var cfg *Config

// Store path to this program executables
var DirPath string

type Config struct {
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
		if err := SetConfig("config/config.json"); err != nil {
			panic(fmt.Sprintf("couldn't get config. %s", err))
		}
	}
	return cfg
}
