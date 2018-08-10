package internal

import (
	"github.com/liagame/lia-cli/internal/config"
	"path/filepath"
	"strconv"
)

func Tutorial(tutorialNumber int, botDir string, debug bool, replayViewerWidth string) {

	gameFlags := &GameFlags{
		GameSeed:   0,
		MapSeed:    0,
		Port:       config.Cfg.GamePort,
		MapPath:    getTutorialMap(tutorialNumber),
		ReplayPath: createReplayFileName(),
		ConfigPath: "",
		DebugBots:  setDebugSlice(debug),
	}

	Compile(botDir)

	tutorialBotDir := filepath.Join("data", "tutorials", strconv.Itoa(tutorialNumber), "bot")
	GenerateGame(botDir, tutorialBotDir, gameFlags)

	ShowReplayViewer(gameFlags.ReplayPath, replayViewerWidth)
}

func getTutorialMap(tutorialNumber int) string {
	return filepath.Join(config.PathToData, "tutorials", strconv.Itoa(tutorialNumber), "map.json")
}

func setDebugSlice(debug bool) []int {
	slice := []int{}
	if debug {
		slice = append(slice, 1)
	}
	return slice
}
