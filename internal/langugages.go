package internal

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/liagame/lia-cli/config"
)

func ShowSupportedLanguages() {
	fmt.Println("Supported languages:")
	for _, langData := range config.Cfg.Languages {
		fmt.Printf("   - %s\n", langData.Name)
	}

	fmt.Printf("\nTo create new bot check %s command.\n", color.GreenString("lia bot"))
}
