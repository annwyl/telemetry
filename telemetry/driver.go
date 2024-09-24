package telemetry

import (
	"encoding/json"
	"fmt"
)

type Driver interface {
	Log(log Log) error
	Close() error
}

type DriverFactory func(config json.RawMessage) (Driver, error)

var registeredDrivers = make(map[string]DriverFactory)

func RegisterDriver(name string, factory DriverFactory) error {
	if _, ok := registeredDrivers[name]; ok {
		return fmt.Errorf("driver already registered: %s", name)
	}
	registeredDrivers[name] = factory
	return nil
}

func GetRegisteredDrivers() map[string]DriverFactory {
	return registeredDrivers
}

func getDriver(config Config) (Driver, error) {
	factory, ok := registeredDrivers[config.Name]
	if !ok {
		return nil, fmt.Errorf("unknown driver: %s", config.Name)
	}

	return factory(config.Config)
}
