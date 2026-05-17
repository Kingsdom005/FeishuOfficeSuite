package data

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func NewConfig(path string) (*BootstrapConfig, error) {
	data, err := os.ReadFile(fmt.Sprintf("%s/config.yaml", path))
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg BootstrapConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}