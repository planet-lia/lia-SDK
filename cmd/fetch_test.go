package cmd_test

import (
	"fmt"
	"github.com/liagame/lia-cli/config"
	"github.com/liagame/lia-cli/internal"
	"github.com/liagame/lia-cli/pkg/advancedcopy"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"strconv"
	"io"
)

func TestCmdFetch(t *testing.T) {
	cases := []struct {
		url             string
		name            string
		hasCustomBotDir bool
		exitStatus      int
		desc            string
	}{
		{
			url:             "https://github.com/liagame/java-bot/archive/master.zip",
			name:            "birko",
			hasCustomBotDir: false,
			exitStatus:      config.OK,
			desc:            "downloading bot birko and put it into working dir",
		},
		{
			url:             "https://github.com/liagame",
			name:            "mirko",
			hasCustomBotDir: false,
			exitStatus:      config.BotDownloadFailed,
			desc:            "try to download non zip file",
		},
		{
			url:             "https://github.com/liagame.zip",
			name:            "mirko",
			hasCustomBotDir: false,
			exitStatus:      config.BotDownloadFailed,
			desc:            "try to download bot from non existent file",
		},
		{
			url:             "https://github.com/liagame/java-bot/archive/master.zip",
			name:            "mirko",
			hasCustomBotDir: true,
			exitStatus:      config.OK,
			desc:            "download bot mirko and put it into custom bot Dir",
		},
	}

	// Run actual tests
	for i, c := range cases {
		if os.Getenv("RUN_FUNC") == strconv.Itoa(i) {
			setupTmpConfigPaths()
			defer cleanupTmpFiles()
			config.Setup()

			customBotDir := ""

			// Set custom bot dir
			if c.hasCustomBotDir {
				var err error
				customBotDir, err = ioutil.TempDir("", "")
				if err != nil {
					t.Error(err)
				}
			}
			defer func () {
				if err := os.RemoveAll(customBotDir); err != nil {
					t.Error(err)
				}
			}()

			// Run command
			internal.FetchBot(c.url, c.name, customBotDir)

			// Check custom bot dir
			if c.hasCustomBotDir {
				empty, err := IsEmpty(customBotDir)
				if err != nil {
					t.Error(err)
				}
				if empty {
					t.Error("hasCustomBotDir should not be empty")
				}
			}

			return
		}
	}

	// Run test and check exit status
	for i, c := range cases {
		output, exitStatus := getCmdStatus("TestCmdFetch", i)
		if exitStatus != c.exitStatus {
			t.Logf("%s", c.desc)
			t.Logf("%s", output)
			t.Errorf("exit status is %v but should be %v", exitStatus, c.exitStatus)
		}
	}
}

func getCmdStatus(funcName string, envValue int) (string, int) {
	cmd := exec.Command(os.Args[0], "-test.run="+funcName)
	cmd.Env = append(os.Environ(), "RUN_FUNC=" + strconv.Itoa(envValue))
	output, err := cmd.CombinedOutput()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		if status, ok := e.Sys().(syscall.WaitStatus); ok {
			return string(output), status.ExitStatus()
		}
	}
	return string(output), 0
}

func setupTmpConfigPaths() {
	// Set tmp path to bots
	var err error
	config.PathToBots, err = ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	// Copy data to tmp path to bots
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var index int
	if runtime.GOOS == "windows" {
		index = strings.LastIndex(wd, "\\")
	} else {
		index = strings.LastIndex(wd, "/")
	}
	pathToData := filepath.Join(wd[:index], "build", "data")
	pathToTmpData := filepath.Join(config.PathToBots, "data")

	if err := advancedcopy.Dir(pathToData, pathToTmpData); err != nil {
		msg := fmt.Sprintf("failed to advancedcopy data to executable path %s", err)
		panic(msg)
	}
}

func cleanupTmpFiles() {
	err := os.RemoveAll(config.PathToBots)
	if err != nil {
		panic(err)
	}
}


func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}