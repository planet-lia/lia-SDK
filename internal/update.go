package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/inconshreveable/go-update"
	"github.com/liagame/lia-SDK"
	"github.com/liagame/lia-SDK/internal/config"
	"github.com/liagame/lia-SDK/pkg/advancedcopy"
	"github.com/mholt/archiver"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"strings"
	"strconv"
)

const ReleaseRequestFailed = "failed"

func CheckForUpdate() {
	if !isTimeToCheckForUpdate() {
		return
	}

	latestTag, available := isUpdateAvailable()
	if available {
		printNewUpdateAvailableNotification(latestTag)
	} else {
		updateLastCheckedField()
	}
}

func Update() {
	// Allows running deferred functions before exiting
	osExitStatus := -1
	defer func() {
		if osExitStatus != -1 {
			os.Exit(osExitStatus)
		}
	}()

	// Don't upgrade if the version is the same
	latestTag, available := isUpdateAvailable()
	if !available {
		fmt.Println("Lia-SDK is up to date.")
		return
	}

	fmt.Printf("Upgrading Lia-SDK to version %s\n", latestTag)

	// Create a releaseUrl to the correct release
	releaseUrl := getReleaseZipUrl(latestTag)

	// Create temporary file
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating tmp file %s\n", err)
		osExitStatus = lia_SDK.OsCallFailed
		return
	}
	defer os.Remove(tmpFile.Name())

	// Download update zip
	fmt.Printf("Downloading latest release from %s.\n", releaseUrl)
	if err := downloadZip(releaseUrl, tmpFile, 500); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to download update from %s.\n %s\n", releaseUrl, err)
		osExitStatus = lia_SDK.UpdateDownloadFailed
		return
	}

	// Extract update zip to tmpReleaseParentDir
	tmpReleaseParentDir, err := ioutil.TempDir("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create tmp update dir. %s", err)
		osExitStatus = lia_SDK.OsCallFailed
		return
	}
	defer os.RemoveAll(tmpReleaseParentDir)

	if err := archiver.NewZip().Unarchive(tmpFile.Name(), tmpReleaseParentDir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to extract update with target %s\n%v\n", tmpReleaseParentDir, err)
		osExitStatus = lia_SDK.OsCallFailed
		return
	}

	// Get update dir name in temporary directory
	releaseDirName, err := getDirName(tmpReleaseParentDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get update dir name. %s\n", err)
		osExitStatus = lia_SDK.Generic
		return
	}

	// Get the path to the update dir
	releaseDirPath := filepath.Join(tmpReleaseParentDir, releaseDirName)

	// Check if data directory exists
	if _, err := os.Stat(config.PathToData); !os.IsNotExist(err) {

		fmt.Println("Removing current data/ directory.")
		if err := os.RemoveAll(config.PathToData); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to delete current data/ directory. "+
				"If nothing else helps please download and replace it manualy from %s. Error: %s\n", releaseUrl, err)
			osExitStatus = lia_SDK.Generic
			return
		}
	}

	fmt.Println("Replacing old data/ directory with a new one.")
	pathToNewDataDir := filepath.Join(releaseDirPath, "/data")

	if err := advancedcopy.Dir(pathToNewDataDir, config.PathToData); err != nil {
		fmt.Fprintf(os.Stderr, "Failed move new data dir from %s to %s. %s\n",
			pathToNewDataDir, config.PathToData, err)
		osExitStatus = lia_SDK.OsCallFailed
		return
	}

	// Remove tmp directory
	if err := os.RemoveAll(tmpReleaseParentDir); err != nil {
		fmt.Fprintf(os.Stderr, "failed to remove tmp dir %s, error: %s\n", tmpReleaseParentDir, err)
	}

	fmt.Println("Replacing lia executable.")

	pathToNewLiaExecutable := releaseDirPath + "/lia"
	if runtime.GOOS == "windows" {
		pathToNewLiaExecutable += ".exe"
	}

	liaExecutableBytes, err := ioutil.ReadFile(pathToNewLiaExecutable)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read lia executable to buffer. %s\n", err)
		osExitStatus = lia_SDK.Generic
		return
	}

	// Replace lia executable
	if err := update.Apply(bytes.NewReader(liaExecutableBytes), update.Options{}); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update lia executable. %s\n", err)
		osExitStatus = lia_SDK.Generic
		return
	}

	fmt.Println("Lia was updated sucessfully!")
}

func getReleaseZipUrl(latestTag string) string {
	releaseUrl := config.ReleasesZipUrlBase + "/" + latestTag + "/lia-sdk-"
	switch config.OperatingSystem {
	case "windows":
		releaseUrl += "windows.zip"
	case "darwin":
		releaseUrl += "macos.zip"
	default:
		releaseUrl += "linux.zip"
	}
	return releaseUrl
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

	return latestTime.Add(time.Hour * 3).Before(timeNow), nil
}

func isUpdateAvailable() (string, bool) {
	latestTag := getLatestReleaseTag()
	if latestTag == ReleaseRequestFailed {
		return "", false
	}

	localRelease, err := getLocalReleaseConfig()
	if err != nil {
		return "", false
	}

	larger := isNewUpdateLargerTanCurrent(latestTag, localRelease.Tag)

	return latestTag, larger
}

func isNewUpdateLargerTanCurrent(new, current string) bool {
	new = strings.TrimPrefix(new, "v")
	current = strings.TrimPrefix(current, "v")

	newVersions := strings.Split(new, ".")
	currentVersions := strings.Split(current, ".")

	if len(newVersions) != len(currentVersions) || len(newVersions) != 3 {
		fmt.Printf("Error: New or current release tag is of wrong format (new len=%d, current len=%d)\n",
			len(newVersions), len(currentVersions))
		return false
	}

	// Only check if the major versions are the same
	if newVersions[0] != currentVersions[0] {
		return false
	}

	// Compare minor versions
	if result, err := compareVersions(newVersions[1], currentVersions[1]); err != nil {
		fmt.Printf("Failed to compare minor versions (new minor=%s, current minor=%s, err=%v)\n",
			newVersions[1], currentVersions[1], err)
		return false
	} else if result > 0 {
		return true
	} else if result < 0 {
		return false
	}

	// Compare bug fix versions
	// TODO after the bug fix version there must be no spaces or "- some text" messages, add support
	if result, err := compareVersions(newVersions[2], currentVersions[2]); err != nil {
		fmt.Printf("Failed to compare bug fix versions (new bug fix=%s, current bug fix=%s, err=%v)\n",
			newVersions[1], currentVersions[1], err)
		return false
	} else if result > 0 {
		return true
	} else if result < 0 {
		return false
	}

	return false
}

// Returns 0 if same, less than 0 if v1 < v2 and more than 0 if v1 > v2
func compareVersions(v1, v2 string) (int64, error) {
	v1Number, err := strconv.ParseInt(v1, 10, 0)
	if err != nil {
		return -1, err
	}

	v2Number, err := strconv.ParseInt(v2, 10, 0)
	if err != nil {
		return -1, err
	}

	return v1Number - v2Number, nil
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
	text := fmt.Sprintf("New version %s of Lia-SDK is available.\nSee the changes made at %s.\n"+
		"Please run '%s' command to update it automatically.\n\n", latestTag, config.ReleasesUrl, "lia update")

	fmt.Printf(color.GreenString(text))
}
