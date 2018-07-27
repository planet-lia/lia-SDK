package main

import (
	"github.com/liagame/lia-cli/cmd"
	"github.com/liagame/lia-cli/config"
)

func main() {
	config.Setup()

	cmd.Execute()
}
