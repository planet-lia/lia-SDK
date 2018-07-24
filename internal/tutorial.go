package internal

import (
	"github.com/liagame/lia-cli/config"
	"path/filepath"
	"strconv"
)

func Tutorial(tutorialNumber int, botDir string, debug bool) {

	gameFlags := &GameFlags{
		GameSeed: 0,
		MapSeed: 0,
		Port: config.GetCfg().GamePort,
		MapPath: getTutorialMap(tutorialNumber),
		ReplayPath: createReplayFileName(),
		ConfigPath: "",
		DebugBots: setDebugSlice(debug),
	}

	Compile(botDir)

	tutorialBotDir := filepath.Join("lia", "tutorials",  strconv.Itoa(tutorialNumber), "bot")
	GenerateGame(botDir, tutorialBotDir, gameFlags)

	ShowReplayViewer(gameFlags.ReplayPath)
}

func getTutorialMap(tutorialNumber int) string {
	return filepath.Join(config.PathToLia, "tutorials", strconv.Itoa(tutorialNumber), "map.json")
}

func setDebugSlice(debug bool) []int {
	slice := []int{}
	if debug {
		slice = append(slice, 1)
	}
	return slice
}