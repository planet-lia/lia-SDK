package internal

import (
	"os"
	"github.com/liagame/lia-cli/config"
	"time"
	"path/filepath"
)

func Play(args []string, gameFlags *RunGameFlags, viewReplay bool) {
	bot1Name := args[0]
	bot2Name := args[1]

	Compile(bot1Name)
	if bot1Name != bot2Name {
		Compile(bot2Name)
	}

	// Set up replay file
	if gameFlags.ReplayPath == "" {
		path := filepath.Join(config.DirPath, "replays")
		os.MkdirAll(path, os.ModePerm)
		fileName := time.Now().Format(time.RFC3339) + ".lia"
		gameFlags.ReplayPath = filepath.Join(path, fileName)
	}

	GenerateGame(args, gameFlags)

	if viewReplay {
		ShowReplayViewer(gameFlags.ReplayPath)
	}
}
