package drivers

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/annwyl/telemetry/telemetry"
)

type FileDriver struct {
	file *os.File
}

func init() {
	telemetry.RegisterDriver("file", func(config json.RawMessage) (telemetry.Driver, error) {
		var filename string
		if err := json.Unmarshal(config, &filename); err != nil {
			return nil, err
		}

		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		return &FileDriver{file: file}, nil
	})
}

func (f *FileDriver) Log(log telemetry.Log) error {
	_, err := fmt.Fprintf(f.file, "%s %d %s %s\n", log.Timestamp.Format(time.RFC3339), log.Level, log.Message, log.Tags)
	return err
}

func (f *FileDriver) Close() error {
	return f.file.Close()
}
