package internal

import (
	"github.com/liagame/lia-cli/internal/config"
	"path/filepath"
	"strconv"
	"os"
)

func Playground(playgroundNumber int, botDir string, debug bool, viewReplay bool, replayViewerWidth string) {
	gameFlags := &GameFlags{
		GameSeed:   0,
		MapSeed:    0,
		Port:       config.Cfg.GamePort,
		MapPath:    getPlaygroundMap(playgroundNumber),
		ReplayPath: createReplayFileName(),
		ConfigPath: "",
		DebugBots:  setDebugSlice(debug),
	}

	Compile(botDir)

	playgroundBotDir := filepath.Join("data", "playgrounds", strconv.Itoa(playgroundNumber), "bot")
	GenerateGame(botDir, playgroundBotDir, gameFlags)

	if viewReplay {
		ShowReplayViewer(gameFlags.ReplayPath, replayViewerWidth)
	}
}

func getPlaygroundMap(playgroundNumber int) string {
	playgroundMap := filepath.Join(config.PathToData, "playgrounds", strconv.Itoa(playgroundNumber), "map.json")
	if _, err := os.Stat(playgroundMap); os.IsNotExist(err) {
		return ""
	}
	return playgroundMap
}

func setDebugSlice(debug bool) []int {
	slice := []int{}
	if debug {
		slice = append(slice, 1)
	}
	return slice
}
