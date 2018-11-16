package settings

import (
	"encoding/json"
	"fmt"
	"github.com/liagame/lia-SDK/internal/config"
	"github.com/mitchellh/go-homedir"
	"github.com/satori/go.uuid"
	"os"
	"path/filepath"
)

var defaultSettings = struct {
	TrackingId              string `json:"trackingId"`
	AnalyticsAllow          *bool  `json:"analyticsAllow"`
	AnalyticsAllowedVersion string `json:"analyticsAllowedVersion"`
}{
	GenerateTrackingId(),
	nil,
	config.VERSION,
}

// Create a new settings file in the user's default home directory using the default settings
// file contents.
func Create() error {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find homedir, could not generate .lia.json file")
		return err
	}

	newConfigPath := filepath.Join(home, fmt.Sprintf("%s.%s", config.SettingsFile, config.SettingsFileExtension))

	f, err := os.Create(newConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create config file in: %s\n", newConfigPath)
		return err
	}
	defer f.Close()

	defaultFileContents, _ := json.Marshal(defaultSettings)
	_, err = f.Write(defaultFileContents)
	return err
}

func GenerateTrackingId() string {
	return uuid.Must(uuid.NewV4()).String()
}
