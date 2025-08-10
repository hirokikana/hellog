//go:build windows

package platform

import (
    "errors"
    "fmt"
    "os"
    "path/filepath"

    "golang.org/x/sys/windows/svc/eventlog"
)

type winLogger struct {
    l *eventlog.Log
}

func ensureInstalled(source string) error {
    exe, err := os.Executable()
    if err != nil {
        return err
    }
    exe = filepath.Clean(exe)
    // 既にインストール済みでも Install は失敗するため、Open を先に試す
    if _, err := eventlog.Open(source); err == nil {
        return nil
    }
    // ソースが無い場合のみインストールを試みる
    // Windowsでは管理者権限が必要な場合があります。
    // 最新の x/sys では Install のシグネチャが (string, string, bool, uint32)
    // に拡張されているため、第4引数に 0 を渡す。
    if err := eventlog.Install(source, exe, false, 0); err != nil {
        return fmt.Errorf("install event source: %w", err)
    }
    return nil
}

func newPlatformLogger(tag string) (HelloLogger, error) {
    if tag == "" {
        return nil, errors.New("empty event source tag")
    }
    if err := ensureInstalled(tag); err != nil {
        // インストールに失敗しても Open を試す（既存の可能性）
        // それでも失敗したらエラーを返す
    }
    l, err := eventlog.Open(tag)
    if err != nil {
        return nil, err
    }
    return &winLogger{l: l}, nil
}

func (w *winLogger) LogHello() error {
    // 1 はイベントIDの例。適宜変更可能。
    return w.l.Info(1, "Hello")
}

func (w *winLogger) Close() error { return w.l.Close() }
