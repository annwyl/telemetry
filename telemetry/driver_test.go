package telemetry

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestRegisterDriver(t *testing.T) {
	RegisterDriver("test", func(config json.RawMessage) (Driver, error) {
		return nil, nil
	})
	if _, ok := registeredDrivers["test"]; !ok {
		t.Errorf("wanted test driver registered")
	}
}

func TestGetRegisteredDrivers(t *testing.T) {
	RegisterDriver("test", func(config json.RawMessage) (Driver, error) {
		return nil, nil
	})
	drivers := GetRegisteredDrivers()
	if _, ok := drivers["test"]; !ok {
		t.Errorf("wanted test driver registered")
	}
}

func TestGetDriver(t *testing.T) {
	RegisterDriver("test", func(config json.RawMessage) (Driver, error) {
		return nil, nil
	})
	_, err := getDriver(Config{
		Name: "test",
	})
	if err != nil {
		t.Errorf("wanted no error, got %v", err)
	}
}

func TestGetMissingDriver(t *testing.T) {
	_, err := getDriver(Config{
		Name: "missing",
	})
	if err == nil {
		t.Errorf("wanted error, got nil")
	}
}

func TestDriverFactoryError(t *testing.T) {
	errorFactory := func(config json.RawMessage) (Driver, error) {
		return nil, errors.New("mock factory error")
	}

	RegisterDriver("error_driver", errorFactory)

	config := Config{
		Name: "error_driver",
	}

	logger := NewLogger(config)
	if logger != nil {
		t.Error("wanted nil logger, didnt get nil")
	}
}
