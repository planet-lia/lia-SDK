package internal

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"os"
	"github.com/liagame/lia-cli/config"
	"github.com/palantir/stacktrace"
	"os/exec"
)

func ShowVersions() {
	// Get lia cfg version
	liaCfgName := "cli-config.json"
	liaCfgVersion, err := getConfigVersion(liaCfgName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read from file %s\n %s", liaCfgName, err)
		os.Exit(config.FAILED_TO_READ_CONFIG)
	}

	// Get game cfg version
	gameCfgName := "game-config.json"
	gameCfgVersion, err := getConfigVersion(gameCfgName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read from file %s\n %s", gameCfgName, err)
		os.Exit(config.FAILED_TO_READ_CONFIG)
	}

	// Get game generator version
	cmd := exec.Command("java", "-jar", "game-generator.jar", "--version")
	cmd.Dir = config.DirPath
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get game-generator version\n %s", err)
		os.Exit(config.GAME_GENERATOR_FAILED)
	}
	gameGeneratorVersion := string(out)

	// Get replay viewer version
	cmd = exec.Command("java", "-jar", "replay-viewer.jar", "--version")
	cmd.Dir = config.DirPath
	out, err = cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get replay-viewer version\n %s", err)
		os.Exit(config.REPLAY_VIEWER_FAILED)
	}
	replayViewerVersion := string(out)

	liaCliVersion := "lia-cli version: " + config.VERSION

	fmt.Printf("%s\n%s\n%s\n%s%s",
		liaCliVersion,
		liaCfgVersion,
		gameCfgVersion,
		gameGeneratorVersion,
		replayViewerVersion)
}

func getConfigVersion(fileName string) (string, error) {
	path := config.DirPath + "/.lia/" + fileName

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