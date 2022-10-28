package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	swissknife "github.com/Sagleft/swiss-knife"
)

func Parse(configJSONPath string) (Config, error) {
	if !swissknife.IsFileExists(configJSONPath) {
		return Config{}, errors.New("failed to find config file")
	}

	jsonBytes, err := ioutil.ReadFile(configJSONPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := Config{}
	err = json.Unmarshal(jsonBytes, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to decode config: %w", err)
	}
	return cfg, nil
}
