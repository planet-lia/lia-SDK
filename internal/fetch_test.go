package internal_test

import (
	"github.com/liagame/lia-cli"
	"github.com/liagame/lia-cli/internal"
	"github.com/liagame/lia-cli/internal/config"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
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
			exitStatus:      lia_cli.OK,
			desc:            "downloading bot birko and put it into working dir",
		},
		{
			url:             "https://github.com/liagame",
			name:            "mirko",
			hasCustomBotDir: false,
			exitStatus:      lia_cli.BotDownloadFailed,
			desc:            "try to download non zip file",
		},
		{
			url:             "https://github.com/liagame.zip",
			name:            "mirko",
			hasCustomBotDir: false,
			exitStatus:      lia_cli.BotDownloadFailed,
			desc:            "try to download bot from non existent file",
		},
		{
			url:             "https://github.com/liagame/java-bot/archive/master.zip",
			name:            "mirko",
			hasCustomBotDir: true,
			exitStatus:      lia_cli.OK,
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
			defer func() {
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
