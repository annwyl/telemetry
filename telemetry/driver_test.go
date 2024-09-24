package telemetry

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestRegisterDriver(t *testing.T) {
	err := RegisterDriver("testRegisterDriver", func(config json.RawMessage) (Driver, error) {
		return nil, nil
	})
	if err != nil {
		t.Fatalf("registerdriver gave error: %v", err)
	}
	if _, ok := registeredDrivers["testRegisterDriver"]; !ok {
		t.Errorf("wanted testRegisterDriver driver registered")
	}
}

func TestGetRegisteredDrivers(t *testing.T) {
	err := RegisterDriver("testGetRegisteredDrivers", func(config json.RawMessage) (Driver, error) {
		return nil, nil
	})
	if err != nil {
		t.Fatalf("registerdriver gave error: %v", err)
	}
	drivers := GetRegisteredDrivers()
	if _, ok := drivers["testGetRegisteredDrivers"]; !ok {
		t.Errorf("wanted testGetRegisteredDriver driver registered")
	}
}

func TestGetDriver(t *testing.T) {
	err := RegisterDriver("testGetDriver", func(config json.RawMessage) (Driver, error) {
		return nil, nil
	})
	if err != nil {
		t.Fatalf("registerdriver gave error: %v", err)
	}
	_, err = getDriver(Config{
		Name: "testGetDriver",
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

	err := RegisterDriver("error_driver", errorFactory)
	if err != nil {
		t.Fatalf("registerdriver gave error: %v", err)
	}

	config := Config{
		Name: "error_driver",
	}

	_, err = NewLogger(config)
	if err == nil {
		t.Fatal("wanted error, got nil")
	}
}
