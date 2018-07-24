package internal

import (
	"fmt"
	"os"
	"github.com/palantir/stacktrace"
	"github.com/liagame/lia-cli/config"
	"runtime"
	"io/ioutil"
	"os/exec"
	"strings"
	"path/filepath"
)

func Compile(botName string) {
	lang := GetBotLanguage(botName)
	botDir := filepath.Join(config.PathToBots, botName)

	// Choose platform dependent preparing logic
	prepareCommands := lang.PrepareUnix
	runScriptContent := lang.RunScriptUnix
	runScriptName := "run.sh"
	if runtime.GOOS == "windows" {
		prepareCommands = lang.PrepareWindows
		runScriptContent = lang.RunScriptWindows
		runScriptName = "run.bat"
	}

	// Run prepare commands
	fmt.Println("Preparing bot...")
	for _, cmd := range prepareCommands {
		if err := runPrepareCommand(botDir, cmd); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run command for bot %s language %s\n %s", botName, lang.Name, err)
			os.Exit(config.PREPARING_BOT_FAILED)
		}
	}

	// Convert script content from []string to []byte
	var runScriptContentBytes []byte
	for i, line := range runScriptContent {
		runScriptContentBytes = append(runScriptContentBytes, []byte(line)...)
		if i < len(runScriptContent) - 1 {
			runScriptContentBytes = append(runScriptContentBytes, []byte(" && ")...)
		}
	}
	if runtime.GOOS != "windows" {
		runScriptContentBytes = append([]byte("#!/bin/bash\n"), runScriptContentBytes...)
	}

	// Create run script
	fmt.Println("Creating run script...")
	if err := createRunScript(botDir, runScriptName, runScriptContentBytes); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create run script in %s/%s\n %s", botDir, runScriptName, err)
		os.Exit(config.CREATING_RUN_SCRIPT_FAILED)
	}

	fmt.Println("Completed.")
}

func createRunScript(botDir string, name string, body []byte) error {
	runScriptPath := filepath.Join(botDir, name)

	err := ioutil.WriteFile(runScriptPath, body, 0644) // overwriting
	if err != nil {
		return stacktrace.Propagate(err, "failed to create run script in %s", runScriptPath)
	}
	// Make it executable
	if runtime.GOOS != "windows" {
		cmd := exec.Command("chmod", "+x", "run.sh")
		cmd.Dir = botDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return stacktrace.Propagate(err, "failed to make run.sh executable")
		}
	}
	return nil
}

func runPrepareCommand(botDir string, prepareCmd config.Command) error {
	cmd := exec.Command(prepareCmd.Args[0])
	cmd.Dir = botDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Append arguments
	for i := 1; i < len(prepareCmd.Args); i++ {
		cmd.Args = append(cmd.Args, prepareCmd.Args[i])
	}

	if err := cmd.Run(); err != nil {
		return stacktrace.Propagate(err, "failed to run command [%s]", strings.Join(prepareCmd.Args, " "))
	}
	return nil
}