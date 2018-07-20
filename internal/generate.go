package internal

import (
	"os/exec"
	"fmt"
	"os"
	"github.com/liagame/lia-cli/config"
	"github.com/palantir/stacktrace"
	"crypto/rand"
	"runtime"
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

func GenerateGame(args []string, gameFlags *RunGameFlags) {
	uidBot1, err := generateUuid()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate uid. %s", err)
		return
	}
	uidBot2, err := generateUuid()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate uid. %s", err)
		return
	}
	nameBot1 := args[0]
	nameBot2 := args[1]

	// Prepare bots
/*	fmt.Printf("Preparing bot %s...\n", nameBot1)
	if err := prepareBot(nameBot1); err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare the first bot. %s", err)
		return
	}
	if nameBot1 != nameBot2 {
		fmt.Printf("Preparing bot %s...\n", nameBot2)
		if err := prepareBot(nameBot2); err != nil {
			fmt.Fprintf(os.Stderr, "failed to prepare the second bot. %s", err)
			return
		}
	}*/

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
			fmt.Fprintf(os.Stderr, "failed to generate game\n %s\n", err)
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
		fmt.Fprintf(os.Stderr, errorMsg, err)
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

	botDir := config.DirPath + "/" + name

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