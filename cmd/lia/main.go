package main

import (
	"github.com/liagame/lia-SDK/cmd/lia/cmd"
	"github.com/liagame/lia-SDK/internal/config"
)

func main() {
	config.Setup()

	cmd.Execute()
}
