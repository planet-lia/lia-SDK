package internal

func Play(bot1Dir string, bot2Dir string, gameFlags *GameFlags, viewReplay bool, replayViewerWidth string) {
	Compile(bot1Dir)
	if bot1Dir != bot2Dir {
		Compile(bot2Dir)
	}

	GenerateGame(bot1Dir, bot2Dir, gameFlags)

	if viewReplay {
		ShowReplayViewer(gameFlags.ReplayPath, replayViewerWidth)
	}
}

