//go:build !windows

package platform

import (
    "log/syslog"
)

type sysLogger struct {
    w *syslog.Writer
}

func newPlatformLogger(tag string) (HelloLogger, error) {
    w, err := syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, tag)
    if err != nil {
        return nil, err
    }
    return &sysLogger{w: w}, nil
}

func (s *sysLogger) LogHello() error {
    return s.w.Info("Hello")
}

func (s *sysLogger) Close() error { return s.w.Close() }

