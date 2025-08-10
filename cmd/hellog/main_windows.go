//go:build windows

package main

import (
    "context"
    "errors"
    "flag"
    "fmt"
    "os"
    "time"

    "golang.org/x/sys/windows/svc"

    "hellog/internal/app/runner"
    "hellog/internal/config"
    "hellog/internal/platform"
)

func main() {
    // デバッグ用: 対話セッションであれば通常実行も可能
    isWinSvc, err := svc.IsWindowsService()
    if err != nil {
        fmt.Fprintln(os.Stderr, "svc detect error:", err)
        os.Exit(1)
    }
    if !isWinSvc {
        var cfgPath string
        flag.StringVar(&cfgPath, "config", "", "設定ファイルパス (省略時は既定パス)")
        flag.Parse()
        if err := runFromConfig(cfgPath); err != nil && !errors.Is(err, context.Canceled) {
            fmt.Fprintln(os.Stderr, "run error:", err)
            os.Exit(1)
        }
        return
    }
    // サービスとして実行
    if err := svc.Run("hellog", &serviceHandler{}); err != nil {
        fmt.Fprintln(os.Stderr, "service run error:", err)
        os.Exit(1)
    }
}

type serviceHandler struct{}

func (h *serviceHandler) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
    s <- svc.Status{State: svc.StartPending}

    // 設定読み込み
    cfg, err := config.Load("")
    if err != nil {
        // サービス起動失敗コード 2 = ERROR_FILE_NOT_FOUND 等に相当させたいが、簡易に 1 を返す
        return true, 1
    }
    lg, err := platform.NewLogger(cfg.Tag)
    if err != nil {
        return true, 1
    }
    defer lg.Close()
    rnr := &runner.Runner{Logger: lg, Interval: cfg.Interval()}

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    done := make(chan error, 1)
    go func() { done <- rnr.Run(ctx) }()

    s <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

    var exitCode uint32 = 0
    loop:
    for {
        select {
        case c := <-r:
            switch c.Cmd {
            case svc.Interrogate:
                s <- c.CurrentStatus
            case svc.Stop, svc.Shutdown:
                cancel()
                break loop
            default:
                // 他は無視
            }
        case err := <-done:
            if err != nil && !errors.Is(err, context.Canceled) {
                exitCode = 1
            }
            break loop
        case <-time.After(24 * time.Hour):
            // 冗長タイムアウト（通常到達しない）
        }
    }

    s <- svc.Status{State: svc.StopPending}
    return true, exitCode
}

func runFromConfig(cfgPath string) error {
    cfg, err := config.Load(cfgPath)
    if err != nil {
        return err
    }
    lg, err := platform.NewLogger(cfg.Tag)
    if err != nil {
        return err
    }
    defer lg.Close()
    r := &runner.Runner{Logger: lg, Interval: cfg.Interval()}
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    // Ctrl+C 等で終了（デバッグ用）
    return r.Run(ctx)
}

