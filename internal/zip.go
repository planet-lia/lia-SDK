package internal

import (
	"fmt"
	"github.com/liagame/lia-cli/config"
	"github.com/mholt/archiver"
	"os"
	"path/filepath"
)

func Zip(botDir string) {
	Compile(botDir)

	// Create zip
	botDirAbsPath := botDir
	if !filepath.IsAbs(botDir) {
		botDirAbsPath = filepath.Join(config.PathToBots, botDir)
	}

	if err := archiver.Zip.Make(botDirAbsPath+".zip", []string{botDirAbsPath}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to zip bot %s\n %s", botDirAbsPath, err)
		os.Exit(config.ZippingBotFailed)
	}
}
