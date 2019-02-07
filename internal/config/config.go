package config

import (
	"encoding/json"
	"fmt"
	"github.com/liagame/lia-SDK"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

const VERSION = "1.0.0"
const LiaHomePage = "https://liagame.com"

const SettingsFile = ".lia"
const SettingsFileExtension = "json"
const PropertyID = "UA-122844498-1" // Id of google analytics project
const TestPropertyID = "UA-122844498-2"

const defaultReleasesBaseUrl = "https://github.com/Svigelj/sdk-test/releases/"

var ReleasesUrl string
var ReleasesZipUrlBase string

const defaultLiaBackendUrl = "https://prod.cloud1.liagame.com"

var AuthUrl string
var AuthVerifyUrl string
var BotUploadUrl string

func init() {
	// Lia releases URLs
	releasesBaseUrl := os.Getenv("RELEASES_BASE_URL")
	if releasesBaseUrl == "" {
		releasesBaseUrl = defaultReleasesBaseUrl
	} else {
		fmt.Printf("Releases base URL set to %s\n", releasesBaseUrl)
	}
	// Url returns a json with "tag_name" key containing latest version (eg. v1.0.0)
	ReleasesUrl = releasesBaseUrl + "latest"
	// Base from where the releases can be downloaded
	ReleasesZipUrlBase = releasesBaseUrl + "download"

	// Lia backend URLs
	liaBackendUrl := os.Getenv("LIA_BACKEND_URL")
	if liaBackendUrl == "" {
		liaBackendUrl = defaultLiaBackendUrl
	} else {
		fmt.Printf("Lia backend URL set to %s\n", liaBackendUrl)
	}
	AuthUrl = liaBackendUrl + "/auth/"
	AuthVerifyUrl = liaBackendUrl + "/auth/verify/"
	BotUploadUrl = liaBackendUrl + "/game/bot/upload/"
}

var OperatingSystem = runtime.GOOS

var Cfg *Config

// Store path to this program executables
var PathToBots string
var PathToData string

type Config struct {
	Version   string     `json:"version"`
	GamePort  int        `json:"gamePort"`
	Languages []Language `json:"languages"`
}

type Language struct {
	Name           string `json:"name"`
	BotURL         string `json:"botUrl"`
	PrepareUnix    string `json:"prepareUnix"`
	RunUnix        string `json:"runUnix"`
	PrepareWindows string `json:"prepareWindows"`
	RunWindows     string `json:"runWindows"`
	CleanupUnix    string `json:"cleanupUnix"`
	CleanupWindows string `json:"cleanupWindows"`
}

var LoggedInUser string
var UserToken string

func SetConfig(path string) error {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open file. Location: %s", path)
		return err
	}

	Cfg = &Config{}
	if err := json.Unmarshal(configFile, Cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't unmarshal config")
		return err
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
