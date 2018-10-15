package config

import (
	"encoding/json"
	"fmt"
	"github.com/liagame/lia-SDK"
	"github.com/palantir/stacktrace"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

const VERSION = "0.1.0"

const ReleasesUrl = "https://github.com/liagame/lia-SDK/releases/latest"
const SettingsFile = ".lia"
const SettingsFileExtension = "json"
const PropertyID = "UA-122844498-1" // Id of google analytics project
const TestPropertyID = "UA-122844498-2"

var OperatingSystem = runtime.GOOS

var Cfg *Config

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

	Cfg = &Config{}
	if err := json.Unmarshal(configFile, Cfg); err != nil {
		return stacktrace.Propagate(err, "couldn't unmarshal config")
	}

	return nil
}

func Setup() {
	if Cfg == nil {
		if PathToBots == "" {
			// Set PathToBots to executable path
			ex, err := os.Executable()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get executable location\n %s", err)
				os.Exit(lia_SDK.FailedToGetEnvironment)
			}
			PathToBots = filepath.Dir(ex)
		}

		PathToData = filepath.Join(PathToBots, "data")
		pathToCfg := filepath.Join(PathToData, "cli-config.json")

		if err := SetConfig(pathToCfg); err != nil {
			fmt.Fprintf(os.Stderr, "couldn't get config\n %s", err)
			os.Exit(lia_SDK.FailedToReadConfig)
		}
	}
}
