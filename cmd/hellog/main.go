package main

import (
    "context"
    "errors"
    "flag"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "hellog/internal/app/runner"
    "hellog/internal/config"
    "hellog/internal/platform"
)

// 非Windows: systemd から起動される前提。手動実行時は -config でパス指定可能。
func main() {
    var cfgPath string
    flag.StringVar(&cfgPath, "config", "", "設定ファイルパス (省略時は既定パス)")
    flag.Parse()

    cfg, err := config.Load(cfgPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "config error: %v\n", err)
        os.Exit(2)
    }

    lg, err := platform.NewLogger(cfg.Tag)
    if err != nil {
        fmt.Fprintf(os.Stderr, "logger init error: %v\n", err)
        os.Exit(1)
    }
    defer lg.Close()

    r := &runner.Runner{Logger: lg, Interval: cfg.Interval()}

    ctx, cancel := signalContext()
    defer cancel()

    if err := r.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
        fmt.Fprintf(os.Stderr, "run error: %v\n", err)
        os.Exit(1)
    }
}

func signalContext() (context.Context, context.CancelFunc) {
    ctx, cancel := context.WithCancel(context.Background())
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sig
        cancel()
    }()
    return ctx, cancel
}
