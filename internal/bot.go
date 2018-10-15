package internal

import (
	"fmt"
	"github.com/liagame/lia-SDK"
	"github.com/liagame/lia-SDK/internal/config"
	"github.com/palantir/stacktrace"
	"os"
)

func FetchBotByLanguage(lang string, name string) {
	// Check if the bot with name already exists
	if isUsed, err := isNameUsed(name); err != nil {
		fmt.Fprintf(os.Stderr, "failed to check if name isUsed. %s", err)
		os.Exit(lia_SDK.Generic)
	} else if isUsed {
		fmt.Fprintf(os.Stderr, "bot name %s already exists. Choose another name.\n", name)
		os.Exit(lia_SDK.BotExists)
	}

	// Fetch repository url for specified language
	url, err := getRepositoryURL(lang)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(lia_SDK.FailedToFindRepo)
	}
	url += "/archive/master.zip"

	FetchBot(url, name, "")
}

// Find repository from config file based on lang parameter
func getRepositoryURL(lang string) (string, error) {
	for _, langData := range config.Cfg.Languages {
		if lang == langData.Name {
			return langData.BotURL, nil
		}
	}

	return "", stacktrace.NewError("BotRepo not found: %v", lang)
}
