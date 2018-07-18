package internal

import (
	"os/exec"
	"fmt"
	"github.com/liagame/lia-cli/config"
	"os"
)

func ShowReplayViewer(replayFile string)  {
	cmd := exec.Command("java", "-jar", "replay-viewer.jar")
	cmd.Dir = config.DirPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if replayFile != "" {
		cmd.Args = append(cmd.Args, replayFile)
	}
	if err := cmd.Run(); err != nil {
		fmt.Printf("Couldn't run replay: %s\n", err)
		return
	}
}
