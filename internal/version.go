package internal

import (
	"encoding/json"
	"fmt"
	"github.com/liagame/lia-SDK"
	"github.com/liagame/lia-SDK/internal/config"
	"github.com/palantir/stacktrace"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func ShowVersions() {
	// Get lia cfg version
	liaCfgName := "cli-config.json"
	liaCfgVersion, err := getConfigVersion(liaCfgName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read from file %s\n %s", liaCfgName, err)
		os.Exit(lia_SDK.FailedToReadConfig)
	}

	// Get game cfg version
	gameCfgName := "game-config.json"
	gameCfgVersion, err := getConfigVersion(gameCfgName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read from file %s\n %s", gameCfgName, err)
		os.Exit(lia_SDK.FailedToReadConfig)
	}

	// Get game generator version
	cmd := exec.Command("java", "-jar", "game-engine.jar", "--version")
	cmd.Dir = config.PathToData
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get game-generator version\n %s", err)
		os.Exit(lia_SDK.GameGeneratorFailed)
	}
	gameGeneratorVersion := string(out)

	// Get replay viewer version
	cmd = exec.Command("java", "-jar", "replay-viewer.jar", "--version")
	cmd.Dir = config.PathToData
	out, err = cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get replay-viewer version\n %s", err)
		os.Exit(lia_SDK.ReplayViewerFailed)
	}
	replayViewerVersion := string(out)

	liaCliVersion := "lia-SDK version: " + config.VERSION

	fmt.Printf("%s\n%s\n%s\n%s%s",
		liaCliVersion,
		liaCfgVersion,
		gameCfgVersion,
		gameGeneratorVersion,
		replayViewerVersion)
}

func getConfigVersion(fileName string) (string, error) {
	path := filepath.Join(config.PathToData, fileName)

	// Read config file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", stacktrace.Propagate(err, "")
	}

	// Get version
	var objmap map[string]*json.RawMessage
	if err := json.Unmarshal(b, &objmap); err != nil {
		return "", stacktrace.Propagate(err, "")
	}
	var version string
	if err := json.Unmarshal(*objmap["version"], &version); err != nil {
		return "", stacktrace.Propagate(err, "")
	}

	version = fmt.Sprintf("%s version: %s", fileName, version)

	return version, nil
}
