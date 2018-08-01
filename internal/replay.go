package internal

import (
	"fmt"
	"github.com/liagame/lia-cli"
	"github.com/liagame/lia-cli/internal/config"
	"os"
	"os/exec"
	"runtime"
)

func ShowReplayViewer(replayFile string) {
	var args []string
	if runtime.GOOS == "darwin" {
		args = append(args, "-XstartOnFirstThread", "-Dorg.lwjgl.system.allocator=system")
	}
	args = append(args, "-jar", "replay-viewer.jar")
	if replayFile != "" {
		args = append(args, replayFile)
	}

	cmd := exec.Command("java", args...)
	cmd.Dir = config.PathToData
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't run replay: %s\n", err)
		os.Exit(lia_cli.ReplayViewerFailed)
	}
}
