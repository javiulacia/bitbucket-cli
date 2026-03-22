package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Save(cfg *Config) error {
	path, err := Path()
	if err != nil {
		return err
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.bitbucket.org/2.0"
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func Delete() error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
