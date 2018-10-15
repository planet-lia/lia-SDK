package tests

import (
	"os"
	"github.com/liagame/lia-SDK/internal/config"
	"github.com/liagame/lia-SDK/internal"
	"fmt"
	"github.com/liagame/lia-SDK"
	"testing"
	"strconv"
)

func TestPlayground(t *testing.T) {
	// Run actual test
	if os.Getenv("RUN_FUNC") != "" {
		SetupTmpConfigPaths()
		defer CleanupTmpFiles()
		config.Setup()

		botName := "birko"

		// Fetch bot
		internal.FetchBotByLanguage("java", botName)

		playgroundNumber, err := strconv.Atoi(os.Getenv("RUN_FUNC"))
		if err != nil {
			t.Fatal(err)
		}
		internal.Playground(playgroundNumber, botName, false, false, "")

		return
	}

	// Run test for all supported playgrounds and check exit status
	for i := 1; i < 3; i++ {
		fmt.Printf("testing playground %d...\n", i)
		output, exitStatus := GetCmdStatus("TestPlayground", i, true)
		if exitStatus != lia_SDK.OK {
			t.Logf("playground %d failed\n", i)
			t.Logf("%s", output)
			t.Fatalf("exit status is %v but should be %v", exitStatus, lia_SDK.OK)
		}
	}
}