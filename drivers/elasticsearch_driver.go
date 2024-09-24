package drivers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/annwyl/telemetry/telemetry"
)

type ElasticsearchDriver struct {
	client *http.Client
	url    string
	index  string
	auth   string
}

func init() {
	err := telemetry.RegisterDriver("elasticsearch", func(config json.RawMessage) (telemetry.Driver, error) {
		var cfg struct {
			Host     string `json:"host"`
			Index    string `json:"index"`
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.Unmarshal(config, &cfg); err != nil {
			return nil, err
		}

		if cfg.Host == "" || cfg.Index == "" {
			return nil, fmt.Errorf("elasticsearch host and index required")
		}

		driver := &ElasticsearchDriver{
			client: &http.Client{},
			url:    fmt.Sprintf("%s/%s/_doc", cfg.Host, cfg.Index),
			index:  cfg.Index,
		}

		if cfg.Username != "" && cfg.Password != "" {
			driver.auth = fmt.Sprintf("%s:%s", cfg.Username, cfg.Password)
		}

		return driver, nil
	})
	if err != nil {
		panic(err)
	}
}

func (e *ElasticsearchDriver) Log(log telemetry.Log) error {
	logData := map[string]interface{}{
		"timestamp": log.Timestamp.Format(time.RFC3339),
		"level":     log.Level,
		"message":   log.Message,
		"tags":      log.Tags,
	}
	if log.TransactionID != "" {
		logData["transaction_id"] = log.TransactionID
	}

	payload, err := json.Marshal(logData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", e.url, strings.NewReader(string(payload)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if e.auth != "" {
		req.SetBasicAuth(strings.Split(e.auth, ":")[0], strings.Split(e.auth, ":")[1])
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("elasticsearch gave non-2xx status: %d", resp.StatusCode)
	}

	return nil
}

func (e *ElasticsearchDriver) Close() error {
	return nil
}
