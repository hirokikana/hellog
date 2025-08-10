package runner

import (
    "context"
    "errors"
    "testing"
    "time"

    "hellog/internal/platform"
)

type fakeLogger struct {
    calls int
    failAt int
}

func (f *fakeLogger) LogHello() error {
    f.calls++
    if f.failAt > 0 && f.calls == f.failAt {
        return errors.New("boom")
    }
    return nil
}
func (f *fakeLogger) Close() error { return nil }

func TestProcessNTicks(t *testing.T) {
    fl := &fakeLogger{}
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    tick := make(chan time.Time, 3)
    // 3回分のtickを送る
    for i := 0; i < 3; i++ {
        tick <- time.Now()
    }
    if err := processNTicks(ctx, fl, tick, 3); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if fl.calls != 3 {
        t.Fatalf("calls = %d, want 3", fl.calls)
    }
}

func TestProcessNTicks_ErrorPropagation(t *testing.T) {
    fl := &fakeLogger{failAt: 2}
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    tick := make(chan time.Time, 3)
    for i := 0; i < 3; i++ {
        tick <- time.Now()
    }
    err := processNTicks(ctx, fl, tick, 3)
    if err == nil {
        t.Fatalf("expected error, got nil")
    }
}

// コンパイル時に interface を満たしていることのチェック
var _ platform.HelloLogger = (*fakeLogger)(nil)

