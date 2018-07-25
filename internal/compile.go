package internal

import (
	"fmt"
	"github.com/liagame/lia-cli/config"
	"github.com/palantir/stacktrace"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"io"
)

func Compile(botDir string) {
	botDirAbsPath := botDir
	if !filepath.IsAbs(botDir) {
		botDirAbsPath = filepath.Join(config.PathToBots, botDir)
	}

	lang := GetBotLanguage(botDirAbsPath)

	// Run prepare commands
	fmt.Println("Preparing bot...")
	if err := prepareBot(botDirAbsPath,  lang); err != nil {
		fmt.Fprintf(os.Stderr, "failed run prepare bot script for bot %s and lang %s. %s\n", botDirAbsPath, lang.Name, err)
		os.Exit(config.PreparingBotFailed)
	}

	if err := createRunScript(botDirAbsPath, lang); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create run script for bot %s. %s\n", botDirAbsPath, err)
		os.Exit(config.CreatingRunScriptFailed)
	}

	fmt.Println("Completed.")
}

func prepareBot(botDir string, lang *config.Language) error {
	prepareScript := lang.PrepareUnix
	if runtime.GOOS == "windows" {
		prepareScript = lang.PrepareWindows
	}

	if runtime.GOOS == "windows" {
		// Execute
		cmd := exec.Command("./" + prepareScript)
		cmd.Dir = botDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return stacktrace.Propagate(err, "failed to make prepare the bot %s\n", botDir)
		}
	} else {
		pathToLanguages := filepath.Join(config.PathToData, "languages")
		if err := makeFileExecutable(pathToLanguages, prepareScript); err != nil {
			return stacktrace.Propagate(err, "")
		}

		// Execute
		cmd := exec.Command("./" + prepareScript, botDir)
		cmd.Dir = pathToLanguages
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return stacktrace.Propagate(err, "failed to make prepare the bot %s\n", botDir)
		}
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

func createRunScript(botDir string, lang *config.Language) error {
	runScript := lang.RunUnix
	if runtime.GOOS == "windows" {
		runScript = lang.RunWindows
	}
	globalRunScriptPath := filepath.Join(config.PathToData, "languages", runScript)
	botRunScriptPath := filepath.Join(botDir, "run.sh")

	// Copy run script to bot
	if err := Copy(globalRunScriptPath, botRunScriptPath); err != nil {
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

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}