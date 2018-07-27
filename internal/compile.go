package internal

import (
	"fmt"
	"github.com/liagame/lia-cli/config"
	"github.com/liagame/lia-cli/pkg/advancedcopy"
	"github.com/palantir/stacktrace"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func Compile(botDir string) {
	botDirAbsPath := botDir
	if !filepath.IsAbs(botDir) {
		botDirAbsPath = filepath.Join(config.PathToBots, botDir)
	}

	lang := GetBotLanguage(botDirAbsPath)

	// Prepare bot
	fmt.Println("Preparing bot...")
	if err := prepareBot(botDirAbsPath, lang); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run prepare bot script for bot %s and lang %s. %s\n", botDirAbsPath, lang.Name, err)
		os.Exit(config.PreparingBotFailed)
	}

	// Copy run script into bot dir
	if err := copyRunScript(botDirAbsPath, lang); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create run script for bot %s. %s\n", botDirAbsPath, err)
		os.Exit(config.CopyingRunScriptFailed)
	}

	fmt.Println("Completed.")
}

func prepareBot(botDir string, lang *config.Language) error {
	prepareScript := lang.PrepareUnix
	if runtime.GOOS == "windows" {
		prepareScript = lang.PrepareWindows
	}

	pathToLanguages := filepath.Join(config.PathToData, "languages")

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(config.Cfg.PathToBash, prepareScript, botDir)
	} else {
		if err := makeFileExecutable(pathToLanguages, prepareScript); err != nil {
			return stacktrace.Propagate(err, "")
		}
		cmd = exec.Command("./"+prepareScript, botDir)
	}
	cmd.Dir = pathToLanguages
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return stacktrace.Propagate(err, "prepare script failed %s\n", botDir)
	}

	return nil
}

func makeFileExecutable(dir string, file string) error {
	cmd := exec.Command("chmod", "+x", file)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return stacktrace.Propagate(err, "failed to make %s executable", file)
	}
	return nil
}

func copyRunScript(botDir string, lang *config.Language) error {
	runScript := lang.RunUnix
	if runtime.GOOS == "windows" {
		runScript = lang.RunWindows
	}
	globalRunScriptPath := filepath.Join(config.PathToData, "languages", runScript)
	botRunScriptPath := filepath.Join(botDir, "run.sh")

	// Copy run script to bot
	if err := advancedcopy.File(globalRunScriptPath, botRunScriptPath); err != nil {
		return stacktrace.Propagate(err, "failed to copy run script from %s to %s\n", globalRunScriptPath, botRunScriptPath)
	}

	// Make it executable
	if runtime.GOOS != "windows" {
		if err := makeFileExecutable(botDir, "run.sh"); err != nil {
			return stacktrace.Propagate(err, "")
		}
	}
	return nil
}
