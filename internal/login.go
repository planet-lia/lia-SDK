package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/liagame/lia-SDK"
	"github.com/liagame/lia-SDK/internal/config"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"
)

func Login() {
	fmt.Printf("Login with your Lia credentials. If you haven't registered yet visit %s\n", config.LiaHomePage)

	// Ask for username
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Printf("USERNAME 1: %v\n", username)

	// Ask for password
	fmt.Print("Password: ")

	passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		scanner.Scan()
		passwordBytes = scanner.Bytes()
	}
	password := string(passwordBytes)

	fmt.Println("\nRetreiving token.")

	// Get the token
	token, err := getToken(username, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get token. %s\n", err)
		os.Exit(lia_SDK.FailedToGetToken)
	}

	// Store token and username in config
	viper.Set("user", username)
	viper.Set("token", token)
	viper.WriteConfig()

	config.LoggedInUser = username
	config.UserToken = token
	fmt.Println("Login successful!")
}

func getToken(username string, password string) (string, error) {
	url := config.AuthUrl

	var data = []byte(fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 201 {
		return "", errors.New(string(body))
	}

	var objmap map[string]*json.RawMessage
	if err := json.Unmarshal(body, &objmap); err != nil {
		return "", err
	}
	var token string
	if err := json.Unmarshal(*objmap["token"], &token); err != nil {
		return "", err
	}

	return token, nil
}

func CheckAccount() {
	username := config.LoggedInUser

	if username == "" {
		fmt.Println("No user is currently logged in.")
	} else {
		fmt.Printf("Logged in user: %s\n", username)
	}
}
