package main

import (
	"github.com/liagame/lia-cli/cmd"
	"fmt"
	"github.com/liagame/lia-cli/config"
	"os"
	"path/filepath"
)


func main() {
	// Set DirPath to executable path
	ex, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("Failed to get executable location. %s", err))
	}
	config.DirPath = filepath.Dir(ex)

	cmd.Execute()
}
