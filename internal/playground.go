package internal

import (
	"bytes"
	"github.com/liagame/lia-SDK/internal/analytics"
	"github.com/liagame/lia-SDK/internal/config"
	"github.com/liagame/lia-SDK/x_vendor"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func Playground(playgroundNumber int, botDir string, debug bool, viewReplay bool, replayViewerWidth string) {
	gameFlags := &GameFlags{
		GameSeed:   0,
		MapSeed:    0,
		Port:       config.Cfg.GamePort,
		MapPath:    getPlaygroundMap(playgroundNumber),
		ReplayPath: "",
		ConfigPath: "",
		DebugBots:  setDebugSlice(debug),
	}

	Compile(botDir)

	playgroundBotDir := filepath.Join("data", "playgrounds", strconv.Itoa(playgroundNumber), "bot")
	GenerateGame(botDir, playgroundBotDir, gameFlags)

	sendPlaygroundWinnerToAnalytics(gameFlags, playgroundNumber)

	if viewReplay {
		ShowReplayViewer(gameFlags.ReplayPath, replayViewerWidth)
	}
}

func sendPlaygroundWinnerToAnalytics(gameFlags *GameFlags, playgroundNumber int) {
	userWon := "undefined"

	// Reads the replay file
	in, err := ioutil.ReadFile(gameFlags.ReplayPath)
	if err == nil {
		// Parses the replay file and outputs replayData
		replayData, err := x_vendor.GetReplayData(bytes.NewReader(in))
		if err == nil {
			// User's bot played as BOT_1
			userWon = strconv.FormatBool(replayData.GamerWinner == x_vendor.BOT_1)
		}
	}

	//Push to analytics
	analytics.Log("playground", "userWon", map[string]string{
		"userWon":          userWon,
		"playgroundNumber": strconv.Itoa(playgroundNumber),
	})
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
