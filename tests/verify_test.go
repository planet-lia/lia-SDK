package tests

import (
	"github.com/liagame/lia-cli/internal"
	"github.com/liagame/lia-cli/internal/config"
	"os"
	"strconv"
	"testing"

	"github.com/liagame/lia-cli"
)

func TestCmdVerify(t *testing.T) {
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
	}

	// Run actual tests
	for i, c := range cases {
		if os.Getenv("RUN_FUNC") == strconv.Itoa(i) {
			setupTmpConfigPaths()
			defer cleanupTmpFiles()
			config.Setup()

			// Run command
			internal.FetchBot(c.url, c.name, customBotDir)

			return
		}
	}

	// Run test and check exit status
	for i, c := range cases {
		output, exitStatus := getCmdStatus("TestCmdVerify", i)
		if exitStatus != c.exitStatus {
			t.Logf("%s", c.desc)
			t.Logf("%s", output)
			t.Errorf("exit status is %v but should be %v", exitStatus, c.exitStatus)
		}
	}
}
