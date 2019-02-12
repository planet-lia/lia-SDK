package internal

import (
	"fmt"
)

import (
	"bytes"
	"encoding/json"
	"github.com/liagame/lia-SDK"
	"github.com/liagame/lia-SDK/internal/config"
	"github.com/liagame/lia-SDK/pkg/advancedcopy"
	"github.com/mholt/archiver"
	"github.com/pkg/browser"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Upload(botDir string) {
	// Get the valid token
	if config.UserToken == "" || config.LoggedInUser == "" {
		Login()
	} else {
		// Check if token is valid else login
		ok, err := verifyToken(config.UserToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to verify token: %s\n", err)
			os.Exit(lia_SDK.Generic)
		} else if !ok {
			Login()
		}
	}

	// Upload the bot
	trackProgressUrl := ""
	var err error
	if trackProgressUrl, err = prepareAndUploadBot(botDir); err != nil {
		fmt.Fprintf(os.Stderr, "Bot upload failed: %s\n", err)
		os.Exit(lia_SDK.BotUploadFailed)
	}
	fmt.Printf("Bot %s was successfully uploaded\n", botDir)

	// Open user profile in browser
	fmt.Println("Opening your profile page in your default browser...")
	if err := browser.OpenURL(trackProgressUrl); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open URL %s in browser: %s\n", trackProgressUrl, err)
		os.Exit(lia_SDK.Generic)
	}
}
func prepareAndUploadBot(botDir string) (string, error) {
	// Copy bot to tmp dir
	fmt.Println("Copying bot to temporary directory.")

	tmpBotDir, err := ioutil.TempDir("", "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create tmp dir")
		return "", err
	}
	defer os.RemoveAll(tmpBotDir)

	if err := advancedcopy.Dir(botDir, tmpBotDir); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to copy bot to tmp directory")
		return "", err
	}

	// Check if bot compiles
	fmt.Println("Preparing bot to check if everything works.")
	if err := Compile(tmpBotDir); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to compile bot")
		return "", err
	}

	fmt.Println("Cleaning up bot directory.")
	if err := removeRedundantFiles(tmpBotDir); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to remove redundant files")
		return "", err
	}

	// Zip the bot
	botZip, err := zip(tmpBotDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Zipping bot failed")
		return "", err
	}
	defer os.RemoveAll(botZip)

	// Upload it to the backend
	trackProgressUrl := ""
	if trackProgressUrl, err = uploadBot(botZip); err != nil {
		fmt.Fprintln(os.Stderr, "Uploading bot failed")
		return "", err
	}

	return trackProgressUrl, nil
}

func uploadBot(botZip string) (string, error) {
	var client = &http.Client{
		Timeout: time.Second * time.Duration(300),
	}

	botZipFile, err := os.Open(botZip)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open zipped bot file")
		return "", err
	}

	values := map[string]io.Reader{
		"file": botZipFile,
	}

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return "", err
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return "", err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return "", err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", config.BotUploadUrl, &b)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+config.UserToken)

	// Submit the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// Check the response
	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("bad status: %s", resp.Status)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var objmap map[string]*json.RawMessage
	if err := json.Unmarshal(body, &objmap); err != nil {
		return "", err
	}
	var trackProgressUrl string
	if err := json.Unmarshal(*objmap["trackProgressUrl"], &trackProgressUrl); err != nil {
		return "", err
	}


	return trackProgressUrl, nil
}

func removeRedundantFiles(botDir string) error {
	botDirAbsPath := botDir
	if !filepath.IsAbs(botDir) {
		botDirAbsPath = filepath.Join(config.PathToBots, botDir)
	}

	lang, err := GetBotLanguage(botDirAbsPath)
	if err != nil {
		return err
	}

	cleanupScript := lang.CleanupUnix
	if config.OperatingSystem == "windows" {
		cleanupScript = lang.CleanupWindows
	}

	pathToLanguages := filepath.Join(config.PathToData, "languages")

	var cmd *exec.Cmd
	if config.OperatingSystem == "windows" {
		cmd = exec.Command(".\\"+cleanupScript, botDirAbsPath)
	} else {
		cmd = exec.Command("/bin/bash", cleanupScript, botDirAbsPath)
	}
	cmd.Dir = pathToLanguages
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "CleanupUnix script failed for bot %s\n", botDir)
		return err
	}

	return nil
}

func zip(botDir string) (string, error) {
	botDirAbsPath := botDir
	if !filepath.IsAbs(botDir) {
		botDirAbsPath = filepath.Join(config.PathToBots, botDir)
	}

	zipFile := botDirAbsPath + ".zip"
	if err := archiver.NewZip().Archive([]string{botDirAbsPath}, zipFile); err != nil {
		fmt.Fprintf(os.Stderr, "failed to zip bot %s\n", botDirAbsPath)
		return "", err
	}

	return zipFile, nil
}

func verifyToken(token string) (bool, error) {
	url := config.AuthVerifyUrl

	var jsonStr = []byte(fmt.Sprintf(`{"token":"%s"}`, token))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, nil
	}

	return true, nil
}
