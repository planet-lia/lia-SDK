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
		fmt.Println("Failed to check if name isUsed. \n", err)
		return
	} else if isUsed {
		fmt.Printf("Bot name %s already exists. Choose another name.\n", name)
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
		fmt.Println("Error while creating tmp tmpFile \n", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Download bot zip
	if err := downloadBot(repoURL, tmpFile); err != nil {
		fmt.Printf("Failed to download bot from %s.\n %s\n", repoURL, err)
		return
	}

	// Extract bot
	if err := unzipBot(tmpFile.Name(), name); err != nil {
		fmt.Printf("Failed to extract bot with name %s. %v\n", name, err)
		return
	}
}

func isNameUsed(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

/** Find repository from config file based on lang parameter */
func getRepositoryURL(lang string) (string, error) {
	e := reflect.ValueOf(*config.GetCfg().BotRepos)

	for i := 0; i < e.Type().NumField(); i++ {
		lang2 := e.Type().Field(i).Tag.Get("json")
		value := e.Field(i).Interface()

		if lang == lang2 {
			return value.(string), nil
		}
	}

	return "", stacktrace.NewError("BotRepo not found: %v", lang)
}

func downloadBot(url string, output *os.File) error {
	response, err := http.Get(url)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to download bot from %s. %s", url)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return stacktrace.NewError("Failed to download bot. %v", *response)
	}

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to store downloaded bot")
	}

	return nil
}

func unzipBot(archive string, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return stacktrace.Propagate(err, "Opening reader failed")
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return stacktrace.Propagate(err, "Creating target dir failed")
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