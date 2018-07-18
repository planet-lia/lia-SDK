package internal

import (
	"os/exec"
	"fmt"
	"os"
	"github.com/liagame/lia-cli/config"
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
	cmd := exec.Command(
		"java", "-jar", "game-generator.jar",
		"-g", fmt.Sprint(gameFlags.GameSeed),
		"-m", fmt.Sprint(gameFlags.MapSeed),
		"-p", fmt.Sprint(gameFlags.Port),
	)

	// Append string flags if they are not empty
	if len(gameFlags.MapPath) > 0 {cmd.Args = append(cmd.Args, "-M", gameFlags.MapPath)}
	if len(gameFlags.ReplayPath) > 0 {cmd.Args = append(cmd.Args, "-r", gameFlags.ReplayPath)}
	if len(gameFlags.ConfigPath) > 0 {cmd.Args = append(cmd.Args, "-c", gameFlags.ConfigPath)}
	// Append bot names
	for _, arg := range args {
		cmd.Args = append(cmd.Args, arg)
	}
	cmd.Dir = config.DirPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Couldn't run replay: %s\n", err)
		return
	}
}