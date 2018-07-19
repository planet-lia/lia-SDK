package internal

import (
	"reflect"
	"github.com/liagame/lia-cli/config"
	"github.com/palantir/stacktrace"
	"fmt"
	"os"
	"net/http"
	"io"
	"archive/zip"
	"path/filepath"
	"io/ioutil"
)

func FetchNewBot(lang string, name string)  {
	// Check if the bot with name already exists
	if isUsed, err := isNameUsed(name); err != nil {
		fmt.Printf("failed to check if name isUsed. %s", err)
		return
	} else if isUsed {
		fmt.Printf("bot name %s already exists. Choose another name.\n", name)
		return
	}

	// Fetch repository url for specified language
	repoURL, err := getRepositoryURL(lang)
	if err != nil {
		fmt.Println(err)
		return
	}
	repoURL += "/archive/master.zip"

	// Create temporary file
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		fmt.Printf("error while creating tmp tmpFile \n", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Download bot zip
	fmt.Printf("Downloading bot from %s...\n", repoURL)
	if err := downloadBot(repoURL, tmpFile); err != nil {
		fmt.Printf("failed to download bot from %s.\n %s\n", repoURL, err)
		return
	}

	// Extract bot
	fmt.Println("Preparing bot...")
	tmpBotParentDir, err := ioutil.TempDir("", "")
	if err != nil {
		fmt.Printf("failed to create tmp bot dir. %s", err)
		return
	}
	defer os.RemoveAll(tmpBotParentDir)
	if err := unzipBot(tmpFile.Name(), tmpBotParentDir); err != nil {
		fmt.Printf("failed to extract bot with target %s. %v\n", tmpBotParentDir, err)
		return
	}

	// Rename bot
	tmpBotDir, err := getBotDir(tmpBotParentDir)
	if err != nil {
		fmt.Printf("failed to get bot dir. %s\n", err)
		return
	}

	// Move bot dir
	finalBotDir := fmt.Sprintf("%s/%s", config.DirPath, name)
	if err := os.Rename(tmpBotDir, finalBotDir); err != nil {
		fmt.Printf("failed move bot dir from %s to %s. %s\n", tmpBotDir, finalBotDir, err)
		return
	}

	fmt.Printf("Bot %s is ready!\n", name)
}

func isNameUsed(name string) (bool, error) {
	path := fmt.Sprintf("%s/%s", config.DirPath, name)
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

/** Find repository from config file based on lang parameter */
func getRepositoryURL(lang string) (string, error) {
	for _, langData := range config.GetCfg().Languages {
		e := reflect.ValueOf(langData)

		lang2 := e.Field(0).Interface()
		value := e.Field(1).Interface()

		if lang == lang2 {
			return value.(string), nil
		}
	}

	return "", stacktrace.NewError("BotRepo not found: %v", lang)
}

func downloadBot(url string, output *os.File) error {
	response, err := http.Get(url)
	if err != nil {
		return stacktrace.Propagate(err, "failed to download bot from %s. %s", url)
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

func getBotDir(parentDir string) (string, error) {
	files, err := ioutil.ReadDir(parentDir)
	if err != nil {
		return "", stacktrace.Propagate(err,"failed to read files from dir: %s", parentDir)
	}
	if len(files) != 1 {
		return "", stacktrace.NewError("there should be exactly 1 file in parentDir. nFiles: %v", len(files))
	}
	botDir := fmt.Sprintf("%s/%s", parentDir, files[0].Name())
	return botDir, nil
}