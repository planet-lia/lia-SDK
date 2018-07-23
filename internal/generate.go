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
	"path/filepath"
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
		os.Exit(config.GENERIC)
	}
	uidBot2, err := generateUuid()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate uid. %s", err)
		os.Exit(config.GENERIC)
	}
	nameBot1 := args[0]
	nameBot2 := args[1]

	result := make(chan error)

	cmdBot1 := &CommandReference{}
	cmdBot2 := &CommandReference{}
	cmdGameGenerator := &CommandReference{}

	nRetries := config.GetCfg().RunBotRetries

	go func () {
		fmt.Printf("Running bot %s\n", nameBot1)
 		err := runBotWithRetries(nRetries, cmdBot1, nameBot1, uidBot1, gameFlags.Port)
		cmdBot1 = nil
		result <- err
	}()
	go func () {
		fmt.Printf("Running bot %s\n", nameBot2)
		err := runBotWithRetries(nRetries, cmdBot2, nameBot2, uidBot2, gameFlags.Port)
		cmdBot2 = nil
		result <- err
	}()
	go func () {
		fmt.Printf("Running game generator\n")
		err := runGameGenerator(cmdGameGenerator, gameFlags, nameBot1, nameBot2, uidBot1, uidBot2)
		cmdGameGenerator = nil
		result <- err
	}()

	for i := 0; i < 3; i++ {
		err = <-result
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate game\n %s\n", err)
			defer os.Exit(config.FAILED_TO_GENERATE_GAME)
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

func runBotWithRetries(nRetries int, cmdRef *CommandReference, name, uid string, port int) error {
	var err error
	var output []byte

	for i := nRetries - 1; i >= 0; i-- {
		output, err = runBot(cmdRef, name, uid, port)
		if err == nil {
			break
		} else {
			time.Sleep(time.Millisecond * time.Duration(config.GetCfg().RunBotRetriesWait))
		}
	}
	if output != nil {
		fmt.Printf("%s\n", output)
	}

	return err
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

func runBot(cmdRef *CommandReference, name, uid string, port int) ([]byte, error) {
	runScriptName := "./run.sh"
	if runtime.GOOS == "windows" {
		runScriptName = "run.bat"
	}

	botDir := filepath.Join(config.DirPath, name)

	cmd := exec.Command(runScriptName, strconv.Itoa(port), uid)
	cmdRef.cmd = cmd
	cmd.Dir = botDir

	output, err := cmd.CombinedOutput();
	if err != nil {
		return output, stacktrace.Propagate(err, "running bot %s failed", name)
	}

	return output, nil
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
		return stacktrace.Propagate(err, "game generator failed.")
	}
	time.Sleep(2000)

	return nil
}

func generateUuid() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", stacktrace.Propagate(err, "failed to get random number")
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}