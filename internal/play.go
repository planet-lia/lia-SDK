package internal

import (
	"fmt"
	"os"
	"github.com/liagame/lia-SDK"
)

func Play(bot1Dir string, bot2Dir string, gameFlags *GameFlags, viewReplay bool, replayViewerWidth string) {
	if err := Compile(bot1Dir); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(lia_SDK.PreparingBotFailed)
	}

	if bot1Dir != bot2Dir {
		if err := Compile(bot2Dir); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(lia_SDK.PreparingBotFailed)
		}
	}

	GenerateGame(bot1Dir, bot2Dir, gameFlags)

	if viewReplay {
		ShowReplayViewer(gameFlags.ReplayPath, replayViewerWidth)
	}
}

