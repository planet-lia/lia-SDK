package internal

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"os"
)

type LiaConfig struct {
	Language string `json:"language"`
}

func getConfig(path string) (*LiaConfig, error) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open file. Location: %s.", path)
		return nil, err
	}

	cfg := &LiaConfig{}
	if err := json.Unmarshal(configFile, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't unmarshal lia config.")
		return nil, err
	}

	return cfg, nil
}
