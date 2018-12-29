package internal

import (
	"fmt"
	"github.com/liagame/lia-SDK"
	"github.com/liagame/lia-SDK/internal/config"
	"github.com/mholt/archiver"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/pkg/errors"
)

func FetchBot(url string, name string, customBotDir string) {
	// Create temporary file
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while creating tmp file %s\n", err)
		os.Exit(lia_SDK.OsCallFailed)
	}
	defer os.Remove(tmpFile.Name())

	// Download bot zip
	fmt.Printf("Downloading bot from %s...\n", url)
	if err := downloadZip(url, tmpFile, 30); err != nil {
		fmt.Fprintf(os.Stderr, "failed to download bot from %s.\n %s\n", url, err)
		defer os.Exit(lia_SDK.BotDownloadFailed)
		return // need to call like that so that other defers are called (removing files etc...)
	}

	// Extract bot
	fmt.Println("Preparing bot...")
	tmpBotParentDir, err := ioutil.TempDir("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create tmp bot dir. %s", err)
		defer os.Exit(lia_SDK.OsCallFailed)
		return
	}
	defer os.RemoveAll(tmpBotParentDir)

	if err := archiver.Zip.Open(tmpFile.Name(), tmpBotParentDir); err != nil {
		fmt.Fprintf(os.Stderr, "failed to extract bot with target %s\n%v\n", tmpBotParentDir, err)
		defer os.Exit(lia_SDK.OsCallFailed)
		return
	}

	// Get bot dir name in temporary file
	botDirName, err := getDirName(tmpBotParentDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get bot dir. %s\n", err)
		defer os.Exit(lia_SDK.Generic)
		return
	}

	// Set bot name
	if name == "" {
		name = botDirName
	}

	// Check if the bot with chosen name already exists
	if isUsed, err := isNameUsed(name); err != nil {
		fmt.Fprintf(os.Stderr, "failed to check if name isUsed. %s", err)
		defer os.Exit(lia_SDK.Generic)
		return
	} else if isUsed {
		fmt.Fprintf(os.Stderr, "bot name %s already exists. Choose another name.\n", name)
		defer os.Exit(lia_SDK.BotExists)
		return
	}

	// Move bot dir and set new name
	tmpBotDir := filepath.Join(tmpBotParentDir, botDirName)
	finalBotDir := customBotDir
	if finalBotDir == "" {
		finalBotDir = filepath.Join(config.PathToBots, name)
	} else {
		finalBotDir = filepath.Join(customBotDir, name)
	}
	if err := os.Rename(tmpBotDir, finalBotDir); err != nil {
		fmt.Fprintf(os.Stderr, "failed move bot dir from %s to %s. %s\n", botDirName, finalBotDir, err)
		defer os.Exit(lia_SDK.OsCallFailed)
		return
	}

	fmt.Printf("Bot %s is ready!\n", name)
}

func isNameUsed(name string) (bool, error) {
	path := filepath.Join(config.PathToBots, name)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func downloadZip(url string, output *os.File, timeoutSeconds int) error {
	if !strings.HasSuffix(url, ".zip") {
		return errors.New("wrong suffix")
	}

	var netClient = &http.Client{
		Timeout: time.Second * time.Duration(timeoutSeconds),
	}

	response, err := netClient.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to download bot from %s", url)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("failed to download bot. %v", *response)
	}

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to store downloaded bot")
		return err
	}

	return nil
}

func getDirName(parentDir string) (string, error) {
	files, err := ioutil.ReadDir(parentDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read files from dir: %s", parentDir)
		return "", err
	}

	switch len(files) {
	case 1:
		return files[0].Name(), nil
	case 2:
		switch {
		case files[0].Name() == "__MACOSX":
			return files[1].Name(), nil
		case files[1].Name() == "__MACOSX":
			return files[0].Name(), nil
		}
	}

	return "", fmt.Errorf("there should be exactly 1 directory in parentDir"+
		"(on mac osx can also be __MACOSX. nFiles: %v", len(files))
}
