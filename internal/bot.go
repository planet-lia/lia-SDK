package internal

import (
	"fmt"
	"github.com/liagame/lia-cli/config"
	"github.com/palantir/stacktrace"
	"os"
)

func FetchBotByLanguage(lang string, name string) {
	// Check if the bot with name already exists
	if isUsed, err := isNameUsed(name); err != nil {
		fmt.Fprintf(os.Stderr, "failed to check if name isUsed. %s", err)
		os.Exit(config.Generic)
	} else if isUsed {
		fmt.Fprintf(os.Stderr, "bot name %s already exists. Choose another name.\n", name)
		os.Exit(config.BotExists)
	}

	// Fetch repository url for specified language
	url, err := getRepositoryURL(lang)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(config.FailedToFindRepo)
	}
	url += "/archive/master.zip"

	FetchBot(url, name)
}

// Find repository from config file based on lang parameter
func getRepositoryURL(lang string) (string, error) {
	for _, langData := range config.GetCfg().Languages {
		if lang == langData.Name {
			return langData.BotURL, nil
		}
	}

	return "", stacktrace.NewError("BotRepo not found: %v", lang)
}
