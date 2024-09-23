package drivers

import (
	"encoding/json"
	"os"

	"github.com/annwyl/telemetry/telemetry"
)

type JSONDriver struct {
	file    *os.File
	encoder *json.Encoder
}

func init() {
	telemetry.RegisterDriver("json", func(config json.RawMessage) (telemetry.Driver, error) {
		var filename string
		if err := json.Unmarshal(config, &filename); err != nil {
			return nil, err
		}

		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		return &JSONDriver{
			file:    file,
			encoder: json.NewEncoder(file),
		}, nil
	})
}

func (j *JSONDriver) Log(log telemetry.Log) error {
	return j.encoder.Encode(log)
}

func (j *JSONDriver) Close() error {
	return j.file.Close()
}
