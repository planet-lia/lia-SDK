package internal

import (
	"fmt"
)

import (
	"github.com/pkg/browser"
	"os"
)

func Upload(botDir string) {

	// Check if there is a token in config


	// Get the token
	// If not login


	// Check if token is valid
	// If not login
	// Login()


	// Upload the bot



	// Open user profile in browser
	fmt.Printf("Bot %s was successfully uploaded.\n" +
		"Opening your profile page in your default browser...", botDir)

	url := "https://liagame.com/games/1"
	if err := browser.OpenURL(url); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open URL %s in browser: %s\n", url, err)
	}
}