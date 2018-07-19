package internal

import (
	"os/exec"
	"fmt"
	"os"
	"github.com/liagame/lia-cli/config"
	"github.com/palantir/stacktrace"
	"crypto/rand"
	"runtime"
	"strings"
	"io/ioutil"
	"strconv"
	"time"
)

type RunGameFlags struct {
	GameSeed int
	MapSeed int
	Port int
	MapPath string
	ReplayPath string
	ConfigPath string
	DebugBots []int
}

func RunGame(args []string, gameFlags *RunGameFlags) {
	uidBot1, err := generateUuid()
	if err != nil {
		fmt.Printf("failed to generate uid. %s", err)
		return
	}
	uidBot2, err := generateUuid()
	if err != nil {
		fmt.Printf("failed to generate uid. %s", err)
		return
	}
	nameBot1 := args[0]
	nameBot2 := args[1]

	// Prepare bots
	fmt.Printf("Preparing bot %s...\n", nameBot1)
	if err := prepareBot(nameBot1); err != nil {
		fmt.Printf("failed to prepare the first bot. %s", err)
		return
	}
	if nameBot1 != nameBot2 {
		fmt.Printf("Preparing bot %s...\n", nameBot2)
		if err := prepareBot(nameBot2); err != nil {
			fmt.Printf("failed to prepare the second bot. %s", err)
			return
		}
	}

	result := make(chan error)

	cmdBot1 := &CommandReference{}
	cmdBot2 := &CommandReference{}
	cmdGameGenerator := &CommandReference{}

	go func () {
		fmt.Printf("Running bot %s\n", nameBot1)
		err := runBot(cmdBot1, nameBot1, uidBot1, gameFlags.Port)
		cmdBot1 = nil
		result <- err
	}()
	go func () {
		fmt.Printf("Running bot %s\n", nameBot2)
		err := runBot(cmdBot2,nameBot2, uidBot2, gameFlags.Port)
		cmdBot2 = nil
		result <- err
	}()
	go func () {
		fmt.Printf("Running game generator\n")
		err := runGameGenerator(cmdGameGenerator, gameFlags, nameBot1, nameBot2, uidBot1, uidBot2)
		cmdBot2 = nil
		result <- err
	}()

	for err := range result {
		if err != nil {
			fmt.Printf("failed to generate game\n %s\n", err)
			break
		}
	}
	killProcess(cmdBot1, fmt.Sprintf("failed to kill bot %s\n %s\n", nameBot1, err))
	killProcess(cmdBot2, fmt.Sprintf("failed to kill bot %s\n %s\n", nameBot2, err))
	killProcess(cmdGameGenerator, fmt.Sprintf("failed to kill game generator\n %s\n", err))

	// Wait for outputs to appear on the console so that the output
	// wont come in late
	time.Sleep(time.Millisecond * 100)
}

func killProcess(cmdRef *CommandReference, errorMsg string) {
	if cmdRef == nil {
		return
	}
	if err := cmdRef.cmd.Process.Kill(); err != nil {
		fmt.Printf(errorMsg, err)
	}
}

type CommandReference struct {
	cmd *exec.Cmd
}

func runBot(cmdRef *CommandReference, name, uid string, port int) error {
	runScriptName := "./run.sh"
	if runtime.GOOS == "windows" {
		runScriptName = "run.bat"
	}

	botDir := fmt.Sprintf("%s/%s", config.DirPath, name)

	cmd := exec.Command(runScriptName, strconv.Itoa(port), uid)
	cmdRef.cmd = cmd
	cmd.Dir = botDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return stacktrace.Propagate(err, "running bot %s failed", name)
	}

	return nil
}

func prepareBot(name string) error {
	lang, err := getBotLanguage(name)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to get language for bot %s", name)
	}

	botDir := fmt.Sprintf("%s/%s", config.DirPath, name)

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
	for _, cmd := range prepareCommands {
		if err := runPrepareCommand(botDir, cmd); err != nil {
			return stacktrace.Propagate(err, "Failed to run command for bot %s language %s", name, lang.Name)
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
	if err := createRunScript(botDir, runScriptName, runScriptContentBytes); err != nil {
		return stacktrace.Propagate(err, "Failed to create run script in %s/%s", botDir, runScriptName)
	}

	return nil
}

func createRunScript(botDir string, name string, body []byte) error {
	runScriptPath := fmt.Sprintf("%s/%s", botDir, name)

	err := ioutil.WriteFile(runScriptPath, body, 0644) // overwriting
	if err != nil {
		return stacktrace.Propagate(err, "Failed to create run script in %s", runScriptPath)
	}
	// Make it executable
	if runtime.GOOS == "windows" {

	} else {
		cmd := exec.Command("chmod", "+x", "run.sh")
		cmd.Dir = botDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return stacktrace.Propagate(err, "Failed to make run.sh executable.")
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
		return stacktrace.Propagate(err, "Failed to run command [%s]", strings.Join(prepareCmd.Args, " "))
	}
	return nil
}


func getBotLanguage(name string) (*config.Language, error) {
	configPath := fmt.Sprintf("%s/%s/lia.json", config.DirPath, name)
	liaConfig, err := getConfig(configPath)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to read %s", configPath)
	}
	for _, langData := range config.GetCfg().Languages {
		if langData.Name == liaConfig.Language {
			return &langData, nil
		}
	}

	return nil, stacktrace.NewError("Language %s was not found", liaConfig.Language)
}

func runGameGenerator(cmdRef *CommandReference, gameFlags *RunGameFlags, nameBot1, nameBot2, uidBot1, uidBot2 string) error {
	cmd := exec.Command(
		"java", "-jar", "game-generator.jar",
		"-g", fmt.Sprint(gameFlags.GameSeed),
		"-m", fmt.Sprint(gameFlags.MapSeed),
		"-p", fmt.Sprint(gameFlags.Port),
	)
	cmdRef.cmd = cmd

	// Append string flags if they are not empty
	if len(gameFlags.MapPath) > 0 {cmd.Args = append(cmd.Args, "-M", gameFlags.MapPath)}
	if len(gameFlags.ReplayPath) > 0 {cmd.Args = append(cmd.Args, "-r", gameFlags.ReplayPath)}
	if len(gameFlags.ConfigPath) > 0 {cmd.Args = append(cmd.Args, "-c", gameFlags.ConfigPath)}
	// Append bot1 and his uid
	cmd.Args = append(cmd.Args, nameBot1, uidBot1)
	// Append bot2 and his uid
	cmd.Args = append(cmd.Args, nameBot2, uidBot2)

	cmd.Dir = config.DirPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return stacktrace.Propagate(err, "Game generator script failed.")
	}

	return nil
}

func generateUuid() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", stacktrace.Propagate(err, "Failed to get random number")
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}