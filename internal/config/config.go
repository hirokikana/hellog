package config

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "runtime"
    "time"
)

type Config struct {
    // IntervalSeconds はログ出力の間隔（秒）
    IntervalSeconds int `json:"interval_seconds"`
    // Tag は syslog タグ / Windows イベントソース名
    Tag string `json:"tag"`
}

func (c *Config) Interval() time.Duration {
    if c.IntervalSeconds <= 0 {
        return 5 * time.Second
    }
    return time.Duration(c.IntervalSeconds) * time.Second
}

func DefaultPath() string {
    if p := os.Getenv("HELLOG_CONFIG"); p != "" {
        return p
    }
    if runtime.GOOS == "windows" {
        if pd := os.Getenv("PROGRAMDATA"); pd != "" {
            return filepath.Join(pd, "Hellog", "config.json")
        }
        return filepath.Join(`C:\\ProgramData`, "Hellog", "config.json")
    }
    if runtime.GOOS == "darwin" {
        return filepath.Join("/Library", "Application Support", "Hellog", "config.json")
    }
    return "/etc/hellog/config.json"
}

func Load(path string) (*Config, error) {
    if path == "" {
        path = DefaultPath()
    }
    b, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read config: %w", err)
    }
    var c Config
    if err := json.Unmarshal(b, &c); err != nil {
        return nil, fmt.Errorf("parse json: %w", err)
    }
    if c.Tag == "" {
        c.Tag = "hellog"
    }
    if c.IntervalSeconds < 0 {
        return nil, errors.New("interval_seconds must be >= 0")
    }
    return &c, nil
}
