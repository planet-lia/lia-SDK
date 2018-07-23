package internal

import (
	"github.com/palantir/stacktrace"
	"io"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/liagame/lia-cli/config"
	"net/http"
	"archive/zip"
	"path/filepath"
)

func FetchBot(url string, name string)  {
	// Create temporary file
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while creating tmp tmpFile %s\n", err)
		os.Exit(config.OS_CALL_FAILED)
	}
	defer os.Remove(tmpFile.Name())

	// Download bot zip
	fmt.Printf("Downloading bot from %s...\n", url)
	if err := downloadBot(url, tmpFile); err != nil {
		fmt.Fprintf(os.Stderr, "failed to download bot from %s.\n %s\n", url, err)
		defer os.Exit(config.BOT_DOWNLOAD_FAILED)
		return // need to call like that so that other defers are called (removing files etc...)
	}

	// Extract bot
	fmt.Println("Preparing bot...")
	tmpBotParentDir, err := ioutil.TempDir("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create tmp bot dir. %s", err)
		defer os.Exit(config.OS_CALL_FAILED)
		return
	}
	defer os.RemoveAll(tmpBotParentDir)
	if err := unzipBot(tmpFile.Name(), tmpBotParentDir); err != nil {
		fmt.Fprintf(os.Stderr, "failed to extract bot with target %s. %v\n", tmpBotParentDir, err)
		defer os.Exit(config.OS_CALL_FAILED)
		return
	}

	// Get bot dir name in temporary file
	botDirName, err := getBotDirName(tmpBotParentDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get bot dir. %s\n", err)
		defer os.Exit(config.GENERIC)
		return
	}

	// Set bot name
	if name == "" {
		name = botDirName
	}

	// Check if the bot with chosen name already exists
	if isUsed, err := isNameUsed(name); err != nil {
		fmt.Fprintf(os.Stderr, "failed to check if name isUsed. %s", err)
		defer os.Exit(config.GENERIC)
		return
	} else if isUsed {
		fmt.Fprintf(os.Stderr, "bot name %s already exists. Choose another name.\n", name)
		defer os.Exit(config.BOT_EXISTS)
		return
	}

	// Move bot dir and set new name
	tmpBotDir := filepath.Join(tmpBotParentDir, botDirName)
	finalBotDir := filepath.Join(config.DirPath, name)
	if err := os.Rename(tmpBotDir, finalBotDir); err != nil {
		fmt.Fprintf(os.Stderr, "failed move bot dir from %s to %s. %s\n", botDirName, finalBotDir, err)
		defer os.Exit(config.OS_CALL_FAILED)
		return
	}

	fmt.Printf("Bot %s is ready!\n", name)
}

func isNameUsed(name string) (bool, error) {
	path := filepath.Join(config.DirPath, name)
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

func downloadBot(url string, output *os.File) error {
	response, err := http.Get(url)
	if err != nil {
		return stacktrace.Propagate(err, "failed to download bot from %s", url)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return stacktrace.NewError("failed to download bot. %v", *response)
	}

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return stacktrace.Propagate(err, "failed to store downloaded bot")
	}

	return nil
}

func unzipBot(archive string, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return stacktrace.Propagate(err, "opening reader failed")
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return stacktrace.Propagate(err, "creating target dir failed")
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return stacktrace.Propagate(err, "")
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return stacktrace.Propagate(err, "")
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return stacktrace.Propagate(err, "")
		}
	}

	return nil
}

func getBotDirName(parentDir string) (string, error) {
	files, err := ioutil.ReadDir(parentDir)
	if err != nil {
		return "", stacktrace.Propagate(err,"failed to read files from dir: %s", parentDir)
	}
	if len(files) != 1 {
		return "", stacktrace.NewError("there should be exactly 1 file in parentDir. nFiles: %v", len(files))
	}
	return files[0].Name(), nil
}