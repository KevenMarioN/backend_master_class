package util

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Config
type Config struct {
	DBDRIVER      string `json:"db_driver"`
	DBSource      string `json:"db_source"`
	ServerAddress string `json:"server_address"`
}

func LoadConfig(configBasePath string) (config *Config, err error) {
	if configBasePath == "" {
		configBasePath = ""
	}

	path := configBasePath + "/config.json"

	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return
	}
	if err = json.Unmarshal(raw, &config); err != nil {
		return
	}

	if err = config.validConfig(); err != nil {
		return
	}

	return config, nil
}

func (c Config) validConfig() error {
	if c.ServerAddress == "" {
		return errors.New("empty server address")
	}
	if c.DBDRIVER == "" {
		return errors.New("empty database driver")
	}

	if c.DBSource == "" {
		return errors.New("empty database address")
	}
	return nil
}
