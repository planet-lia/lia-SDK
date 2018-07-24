package internal

import (
	"github.com/liagame/lia-cli/config"
	"fmt"
	"os"
	"path/filepath"
)


func GetBotLanguage(botName string) *config.Language {
	botConfigPath := filepath.Join(config.PathToBots, botName, "lia.json")
	liaConfig, err := getConfig(botConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read %s\n %s", botConfigPath, err)
		os.Exit(config.FAILED_TO_GET_LIA_JSON)
	}
	for _, langData := range config.GetCfg().Languages {
		if langData.Name == liaConfig.Language {
			return &langData
		}
	}

	fmt.Fprintf(os.Stderr, "language %s was not found", liaConfig.Language)
	os.Exit(config.FAILED_GETTING_BOT_LANG)
	return nil
}