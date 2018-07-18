package config

import (
	"encoding/json"
	"github.com/palantir/stacktrace"
	"io/ioutil"
	"fmt"
)

var cfg *Config
// Store path to this program executables
var DirPath string

type Config struct {
	GamePort int          `json:"gamePort"`
	GameConfigPath string `json:"gameConfigPath"`
	BotRepos []BotRepo    `json:"botRepos"`
}

type BotRepo struct {
	Name string `json:"name"`
	BotURL string `json:"url"`
}

func SetConfig(path string) error {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return stacktrace.Propagate(err, "Couldn't open a file. Location: %s", path)
	}

	cfg = &Config{}
	if err := json.Unmarshal(configFile, cfg); err != nil {
		return stacktrace.Propagate(err, "Couldn't unmarshal config")
	}

	return nil
}

func GetCfg() *Config {
	if cfg == nil {
		if err := SetConfig("config/config.json"); err != nil {
			panic(fmt.Sprintf("Couldn't get config. %s", err))
		}
	}
	return cfg
}