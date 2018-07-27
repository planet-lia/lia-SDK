package internal

import (
	"github.com/liagame/lia-cli/internal/config"
	"os"
	"path/filepath"
	"time"
)

func Play(bot1Dir string, bot2Dir string, gameFlags *GameFlags, viewReplay bool) {
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
	//"2006-01-02T15:04:05Z07:00"
	fileName := time.Now().Format("2006-01-02T15-04-05") + ".lia"
	return filepath.Join(path, fileName)
}
