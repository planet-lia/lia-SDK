package internal

import (
	"fmt"
	"github.com/liagame/lia-cli"
	"github.com/liagame/lia-cli/internal/config"
	"os"
	"os/exec"
)

func ShowReplayViewer(replayFile string) {
	cmd := exec.Command("java", "-jar", "replay-viewer.jar")
	cmd.Dir = config.PathToData
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if replayFile != "" {
		cmd.Args = append(cmd.Args, replayFile)
	}
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't run replay: %s\n", err)
		os.Exit(lia_cli.ReplayViewerFailed)
	}
}
