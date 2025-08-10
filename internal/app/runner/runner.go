package runner

import (
    "context"
    "time"

    "hellog/internal/platform"
)

// Runner は指定間隔で「Hello」を出力します。
type Runner struct {
    Logger   platform.HelloLogger
    Interval time.Duration
}

// Run はコンテキストがキャンセルされるまでループします。
func (r *Runner) Run(ctx context.Context) error {
    if r.Interval <= 0 {
        r.Interval = 5 * time.Second
    }
    t := time.NewTicker(r.Interval)
    defer t.Stop()
    return runWithTick(ctx, r.Logger, t.C)
}

// runWithTick はテスト用に抽出したループ本体です。
func runWithTick(ctx context.Context, l platform.HelloLogger, tick <-chan time.Time) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-tick:
            if err := l.LogHello(); err != nil {
                return err
            }
        }
    }
}

// processNTicks はテスト専用: N回処理したら終了します。
func processNTicks(ctx context.Context, l platform.HelloLogger, tick <-chan time.Time, n int) error {
    count := 0
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-tick:
            if err := l.LogHello(); err != nil {
                return err
            }
            count++
            if count >= n {
                return nil
            }
        }
    }
}

