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

func RegisterDriver(name string, factory DriverFactory) {
	registeredDrivers[name] = factory
}

func GetRegisteredDrivers() map[string]DriverFactory {
	return registeredDrivers
}

func getDriver(config Config) (Driver, error) {
	factory, ok := registeredDrivers[config.Name]
	if !ok {
		return nil, fmt.Errorf("unknown driver: %s", config.Config)
	}

	return factory(config.Config)
}
