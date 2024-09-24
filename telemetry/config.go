package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Name        string            `json:"driver"`
	Config      json.RawMessage   `json:"driver_config"`
	LogLevel    LogLevel          `json:"log_level"`
	DefaultTags map[string]string `json:"default_tags"`
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

	err = validateConfig(config)
	if err != nil {
		return config, fmt.Errorf("invalid config: %v", err)
	}

	return config, nil
}

func validateConfig(config Config) error {
	var errors []string

	if config.Name == "" {
		errors = append(errors, "no driver specified")
	}

	if len(config.Config) == 0 {
		errors = append(errors, "empty config")
	}

	if config.LogLevel < 0 || config.LogLevel > 3 {
		errors = append(errors, "invalid log level")
	}

	if config.LogLevel < DebugLevel || config.LogLevel > ErrorLevel {
		errors = append(errors, fmt.Sprintf("invalid log level: %d (must be between %d and %d)", config.LogLevel, DebugLevel, ErrorLevel))
	}

	for key, value := range config.DefaultTags {
		if key == "" {
			errors = append(errors, "default tag has empty key")
		}
		if value == "" {
			errors = append(errors, "default tag has empty value")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("config validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}
