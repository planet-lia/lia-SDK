package internal

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/liagame/lia-cli/config"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func Update(checkOnly bool) {
	available := isUpdateAvailable()
	if available {
		printLinkToDownloads()
	} else {
		updateLastCheckedField()
	}

	if !checkOnly {
		// TODO perform update
	}
}

func UpdateIfTime(checkOnly bool) {
	if shouldCheckForUpdate() {
		Update(checkOnly)
	}
}

func shouldCheckForUpdate() bool {
	timeToCheck, err := timeToCheckForUpdate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nCan't get local release version. Error: %s\n", err)
		printLinkToDownloads()
		os.Exit(config.FailedToReadReleaseFile)
	}

	return timeToCheck
}

func updateLastCheckedField() {
	releaseCfg, err := getLocalReleaseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get release config. %s", err)
		os.Exit(config.FailedToReadReleaseFile)
	}

	releaseCfg.LastChecked = time.Now().Format(time.RFC3339)

	path := filepath.Join(config.PathToData, "RELEASE.json")
	data, err := json.Marshal(releaseCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal release. %s", err)
		os.Exit(config.Generic)
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write release data to file. %s", err)
		os.Exit(config.Generic)
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

func isUpdateAvailable() bool {
	latestTag := getLatestReleaseTag()
	localRelease, err := getLocalReleaseConfig()

	if err != nil {
		fmt.Printf("\nCan't get local release version. Error: %s\n", err)
	} else if latestTag > localRelease.Tag {
		fmt.Printf("\nUpdate available %s -> %s.\n", color.WhiteString(localRelease.Tag), color.GreenString(latestTag))
	} else {
		return false
	}

	return true
}

type ReleaseConfig struct {
	Tag         string `json:"tag"`
	LastChecked string `json:"lastChecked"`
}

func getLocalReleaseConfig() (*ReleaseConfig, error) {
	path := filepath.Join(config.PathToData, "RELEASE.json")

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &ReleaseConfig{}
	if err := json.Unmarshal(b, cfg); err != nil {
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
		fmt.Fprintf(os.Stderr, "faied to create request. %s\n", err)
		os.Exit(config.Generic)
	}
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "faied to get latest release. %s\n", err)
		os.Exit(config.FailedToGetLatestRelease)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "faied to get latest release. Status: %s\n", res.Status)
		os.Exit(config.FailedToGetLatestRelease)
	}

	// Convert body to bytes
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read response body. %s\n", err)
		os.Exit(config.FailedToGetLatestRelease)
	}

	// Get tag
	var objmap map[string]*json.RawMessage
	if err := json.Unmarshal(body, &objmap); err != nil {
		fmt.Fprintf(os.Stderr, "failed to convert response body to json. %s\n", err)
		os.Exit(config.FailedToGetLatestRelease)
	}
	var tag string
	if err := json.Unmarshal(*objmap["tag_name"], &tag); err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshal tag_name field. %s\n", err)
		os.Exit(config.FailedToGetLatestRelease)
	}

	return tag
}

func printLinkToDownloads() {
	fmt.Printf("Visit %s to download new version. When downloaded just "+
		"replace the old content.\n\n", color.GreenString(config.ReleasesUrl))
}
