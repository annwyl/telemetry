package telemetry

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configContent := `{
		"driver": "mock",
		"driver_config": "",
		"log_level": 1,
		"default_tags": {"environment": "test"}
	}`
	tmpfile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	config, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("loadconfig returned error: %v", err)
	}

	if config.Name != "mock" {
		t.Errorf("wanted driver name 'mock', got '%s'", config.Name)
	}
	if config.LogLevel != InfoLevel {
		t.Errorf("wanted log level %v, got %v", InfoLevel, config.LogLevel)
	}
	if config.DefaultTags["environment"] != "test" {
		t.Errorf("wanted default tag 'environment: test', got '%s'", config.DefaultTags["environment"])
	}
}

func TestBrokenConfig(t *testing.T) {
	configContent := `{
		"driver": "",
		"driver_config": "",
		"log_level": 1,
		"default_tags": {"environment": "test"}
	}
	`
	tmpfile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(tmpfile.Name())
	if err == nil {
		t.Fatal("wanted error, got nil")
	}
}

func TestEmptyConfig(t *testing.T) {
	configContent := `{}`
	tmpfile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(tmpfile.Name())
	if err == nil {
		t.Fatal("wanted error, got nil")
	}
}

func TestMissingConfig(t *testing.T) {
	_, err := LoadConfig("missing.json")
	if err == nil {
		t.Fatal("wanted error, got nil")
	}
}
