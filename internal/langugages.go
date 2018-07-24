package internal

import (
	"github.com/liagame/lia-cli/config"
	"fmt"
	"github.com/fatih/color"
)

func ShowSupportedLanguages() {
	fmt.Println("Supported languages:")
	for _, langData := range config.GetCfg().Languages {
		fmt.Printf("   - %s\n", langData.Name)
	}

	fmt.Printf("\nTo create new bot check %s command.\n", color.GreenString("lia bot"))
}

