package internal

import (
	"os/exec"
	"fmt"
	"github.com/liagame/lia-cli/config"
	"os"
)

func ShowReplayViewer(replayFile string)  {
	cmd := exec.Command("java", "-jar", "replay-viewer.jar")
	cmd.Dir = config.PathToLia
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if replayFile != "" {
		cmd.Args = append(cmd.Args, replayFile)
	}
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't run replay: %s\n", err)
		os.Exit(config.REPLAY_VIEWER_FAILED)
	}
}
