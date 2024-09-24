package drivers

import (
	"encoding/json"
	"fmt"

	"github.com/annwyl/telemetry/telemetry"
)

type ConsoleDriver struct{}

func (c *ConsoleDriver) Log(log telemetry.Log) error {
	fmt.Println(log)
	return nil
}

func (c *ConsoleDriver) Close() error {
	return nil
}

func init() {
	err := telemetry.RegisterDriver("console", func(_ json.RawMessage) (telemetry.Driver, error) {
		return &ConsoleDriver{}, nil
	})
	if err != nil {
		panic(err)
	}
}
