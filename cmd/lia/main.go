package main

import (
	"github.com/liagame/lia-cli/cmd/lia/cmd"
	"github.com/liagame/lia-cli/internal/config"
)

func main() {
	config.Setup()

	cmd.Execute()
}
