package internal

import (
	"github.com/liagame/lia-SDK/internal/config"
	"path/filepath"
	"strconv"
	"os"
	"github.com/liagame/lia-SDK/x_vendor"
	"io/ioutil"
	"bytes"
	"github.com/liagame/lia-SDK/internal/analytics"
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

	// Reads the replay file
	in, err := ioutil.ReadFile(gameFlags.ReplayPath)
	if err != nil {
		t.Errorf("Failed when reading replay file: %s", err)
	}

	// Parses the replay file and outputs replayData
	replayData, err := x_vendor.GetReplayData(bytes.NewReader(in))
	if err != nil {
		t.Errorf("Failed to parse replay file: %s", err)
	}

	//Check which bot won
	var win bool = false
	if replayData.GamerWinner == x_vendor.BOT_1 {
		win = true;
	}
	//Push to analytics
	analytics.Log("playground", "win", map[string]string{
		"win": strconv.FormatBool(win),
		"playgroundNumber" : strconv.Itoa(playgroundNumber),
	})


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
