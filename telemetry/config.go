package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Name   string          `json:"driver"`
	Config json.RawMessage `json:"driver_config"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, fmt.Errorf("failed to decode config file: %v", err)
	}

	if len(config.Config) == 0 {
		return config, fmt.Errorf("empty config")
	}

	if config.Name == "" {
		return config, fmt.Errorf("no driver specified")
	}

	return config, nil
}
