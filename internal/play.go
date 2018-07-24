package internal

import (
	"github.com/liagame/lia-cli/config"
	"os"
	"path/filepath"
	"time"
)

func Play(args []string, gameFlags *GameFlags, viewReplay bool) {
	bot1Dir := args[0]
	bot2Dir := args[1]

	Compile(bot1Dir)
	if bot1Dir != bot2Dir {
		Compile(bot2Dir)
	}

	if gameFlags.ReplayPath == "" {
		gameFlags.ReplayPath = createReplayFileName()
	}

	GenerateGame(bot1Dir, bot2Dir, gameFlags)

	if viewReplay {
		ShowReplayViewer(gameFlags.ReplayPath)
	}
}

func createReplayFileName() string {
	path := filepath.Join(config.PathToBots, "replays")
	os.MkdirAll(path, os.ModePerm)
	fileName := time.Now().Format(time.RFC3339) + ".lia"
	return filepath.Join(path, fileName)
}
