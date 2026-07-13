package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"

	"portscanner/types"
)

func Load(path string) (*types.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg types.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	overrideFromEnv(&cfg)

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func overrideFromEnv(cfg *types.Config) {
	if v := os.Getenv("DATABASE_DSN"); v != "" {
		cfg.Database.DSN = v
	}
	if v := os.Getenv("TELEGRAM_TOKEN"); v != "" {
		cfg.Telegram.Token = v
	}
	if v := os.Getenv("TELEGRAM_CHAT_ID"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			cfg.Telegram.ChatID = id
		}
	}
}

func validate(cfg *types.Config) error {
	if len(cfg.Scan.Targets) == 0 {
		return fmt.Errorf("scan.targets is empty")
	}
	if cfg.Scan.Ports == "" {
		return fmt.Errorf("scan.ports is empty")
	}
	if cfg.Scan.Rate <= 0 {
		return fmt.Errorf("scan.rate must be positive")
	}
	if cfg.Database.DSN == "" {
		return fmt.Errorf("database.dsn is empty")
	}
	if cfg.Telegram.Token == "" || cfg.Telegram.ChatID == 0 {
		return fmt.Errorf("telegram token/chat_id required")
	}
	return nil
}
