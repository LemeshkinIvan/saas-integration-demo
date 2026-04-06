package config

import (
	"daos_core/internal/constants"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	AppEnv = os.Getenv("DAOS_APP_ENV")
)

type Container struct {
	Amo      *constants.AmoConfig      `yaml:"amo"`
	Telegram *constants.TelegramConfig `yaml:"telegram"`
	Common   *constants.CommonConfig   `yaml:"common"`
}

func OpenCfg(path string) (*Container, error) {
	// config.yaml
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadFile, err)
	}

	var cfg Container
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParseYAML, err)
	}
	return &cfg, nil
}
