package internal

import (
	"fmt"
	"github.com/liagame/lia-SDK/internal/config"
	"os"
	"path/filepath"
)

func GetBotLanguage(botDir string) (*config.Language, error) {
	botConfigPath := filepath.Join(botDir, "lia.json")

	liaConfig, err := getConfig(botConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read %s\n", botConfigPath)
		return nil, err
	}
	for _, langData := range config.Cfg.Languages {
		if langData.Name == liaConfig.Language {
			return &langData, nil
		}
	}

	fmt.Fprintf(os.Stderr, "language %s was not found\n", liaConfig.Language)
	return nil, err
}
