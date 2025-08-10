package config

import (
    "os"
    "path/filepath"
    "runtime"
    "testing"
)

func TestLoadAndDefaults(t *testing.T) {
    dir := t.TempDir()
    p := filepath.Join(dir, "conf.json")
    if err := os.WriteFile(p, []byte(`{"interval_seconds":2}`), 0644); err != nil {
        t.Fatal(err)
    }
    t.Setenv("HELLOG_CONFIG", p)
    c, err := Load("")
    if err != nil {
        t.Fatalf("load error: %v", err)
    }
    if c.Tag != "hellog" {
        t.Fatalf("default tag = %q", c.Tag)
    }
    if got := c.Interval().Seconds(); got != 2 {
        t.Fatalf("interval = %v", got)
    }
}

func TestDefaultPath(t *testing.T) {
    t.Setenv("HELLOG_CONFIG", "")
    dp := DefaultPath()
    if runtime.GOOS == "windows" {
        if dp == "" || filepath.Ext(dp) != ".json" {
            t.Fatalf("unexpected default path: %q", dp)
        }
    } else {
        if dp != "/etc/hellog/config.json" {
            t.Fatalf("unexpected default path: %q", dp)
        }
    }
}

