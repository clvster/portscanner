package types

import "time"

type Config struct {
	Scan     ScanConfig     `yaml:"scan"`
	Database DatabaseConfig `yaml:"database"`
	Telegram TelegramConfig `yaml:"telegram"`
}

type ScanConfig struct {
	Targets  []string      `yaml:"targets"`
	Ports    string        `yaml:"ports"`
	Rate     int           `yaml:"rate"`
	SourceIP string        `yaml:"source_ip"`
	Interval time.Duration `yaml:"interval"`
	Once     bool          `yaml:"once"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

type TelegramConfig struct {
	Token  string `yaml:"token"`
	ChatID int64  `yaml:"chat_id"`
}
