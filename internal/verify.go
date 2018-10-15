package internal

import (
	"fmt"
	"github.com/liagame/lia-SDK"
	"github.com/liagame/lia-SDK/internal/config"
	"os"
	"path/filepath"
)

func GetBotLanguage(botDir string) *config.Language {
	botConfigPath := filepath.Join(botDir, "lia.json")

	liaConfig, err := getConfig(botConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read %s\n %s\n", botConfigPath, err)
		os.Exit(lia_SDK.FailedToGetLiaJson)
	}
	for _, langData := range config.Cfg.Languages {
		if langData.Name == liaConfig.Language {
			return &langData
		}
	}

	fmt.Fprintf(os.Stderr, "language %s was not found\n", liaConfig.Language)
	os.Exit(lia_SDK.FailedGettingBotLang)
	return nil
}
