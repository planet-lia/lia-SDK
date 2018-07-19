package internal

import (
	"io/ioutil"
	"github.com/palantir/stacktrace"
	"encoding/json"
)


type LiaConfig struct {
	Language string `json:"language"`
}

func getConfig(path string) (*LiaConfig, error) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, stacktrace.Propagate(err, "couldn't open file. Location: %s", path)
	}

	cfg := &LiaConfig{}
	if err := json.Unmarshal(configFile, cfg); err != nil {
		return nil, stacktrace.Propagate(err, "couldn't unmarshal lia config")
	}

	return cfg, nil
}