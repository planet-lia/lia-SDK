package internal

import (
	"encoding/json"
	"fmt"
	"github.com/liagame/lia-SDK"
	"github.com/liagame/lia-SDK/internal/config"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"github.com/fatih/color"
)

const ReleaseRequestFailed = "failed"

func CheckForUpdate() {
	if !isTimeToCheckForUpdate() {
		return
	}

	latestTag := getLatestReleaseTag()
	if latestTag == ReleaseRequestFailed {
		return
	}

	available := isUpdateAvailable(latestTag)
	if available {
		printNewUpdateAvailableNotification(latestTag)
	} else {
		updateLastCheckedField()
	}
}


func Update() {
	// TODO
}

func isTimeToCheckForUpdate() bool {
	timeToCheck, err := timeToCheckForUpdate()
	if err != nil {
		return false
	}
	return timeToCheck
}

func updateLastCheckedField() {
	releaseCfg, err := getLocalReleaseConfig()
	if err != nil {
		return
	}

	releaseCfg.LastChecked = time.Now().Format(time.RFC3339)

	path := filepath.Join(config.PathToData, "RELEASE.json")
	data, err := json.Marshal(releaseCfg)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return
	}
}

func timeToCheckForUpdate() (bool, error) {
	localRelease, err := getLocalReleaseConfig()
	if err != nil {
		return false, err
	}

	timeNow := time.Now()
	latestTime, err := time.Parse(time.RFC3339, localRelease.LastChecked)
	if err != nil {
		return false, err
	}

	return latestTime.Add(time.Hour * 24).Before(timeNow), nil
}

func isUpdateAvailable(latestTag string) bool {
	localRelease, err := getLocalReleaseConfig()
	if err != nil {
		return false
	}

	// If the major versions are different then don't look at it
	if latestTag[:2] != localRelease.Tag[:2] {
		return false
	}

	if latestTag > localRelease.Tag {
		return true
	}

	return false
}

type ReleaseConfig struct {
	Tag         string `json:"tag"`
	LastChecked string `json:"lastChecked"`
}

func getLocalReleaseConfig() (*ReleaseConfig, error) {
	path := filepath.Join(config.PathToData, "RELEASE.json")

	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Can't find data/RELEASE.json file. If you have deleted it please download Lia-SDK again" +
			" ant copy it from there.\n")
		return nil, err
	}

	cfg := &ReleaseConfig{}
	if err := json.Unmarshal(b, cfg); err != nil {
		fmt.Printf("File data/RELEASE.json is broken. If you have deleted it please download Lia-SDK again" +
			" ant replace it from there.\n")
		return nil, err
	}

	return cfg, nil
}

func getLatestReleaseTag() string {
	var client = &http.Client{
		Timeout: time.Second * 30,
	}

	url := config.ReleasesUrl
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		os.Exit(lia_SDK.Generic)
	}
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return ReleaseRequestFailed
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		os.Exit(lia_SDK.FailedToGetLatestRelease)
	}

	// Convert body to bytes
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		os.Exit(lia_SDK.FailedToGetLatestRelease)
	}

	// Get tag
	var objmap map[string]*json.RawMessage
	if err := json.Unmarshal(body, &objmap); err != nil {
		os.Exit(lia_SDK.FailedToGetLatestRelease)
	}
	var tag string
	if err := json.Unmarshal(*objmap["tag_name"], &tag); err != nil {
		os.Exit(lia_SDK.FailedToGetLatestRelease)
	}

	return tag
}

func printNewUpdateAvailableNotification(latestTag string) {
	text := fmt.Sprintf("New version %s of Lia-SDK is available.\nSee the changes made at %s.\n" +
		"Please run '%s' command to update it automatically.\n\n", latestTag, config.ReleasesUrl, "lia update")

	fmt.Printf(color.GreenString(text))
}
