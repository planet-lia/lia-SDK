package internal

import (
	"fmt"
	"github.com/liagame/lia-SDK"
	"github.com/liagame/lia-SDK/internal/config"
	"os"
	"os/exec"
)

func ShowReplayViewer(replayFile string, replayViewerWidth string) {
	var args []string
	if config.OperatingSystem == "darwin" {
		args = append(args, "-XstartOnFirstThread", "-Dorg.lwjgl.system.allocator=system")
	}
	args = append(args, "-jar", "replay-viewer.jar")
	if replayFile != "" {
		args = append(args, replayFile)
	}
	if replayViewerWidth != "" {
		args = append(args, "-w", replayViewerWidth)
	}
	cmd := exec.Command("java", args...)
	cmd.Dir = config.PathToData
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't run replay: %s\n", err)
		os.Exit(lia_SDK.ReplayViewerFailed)
	}
}
