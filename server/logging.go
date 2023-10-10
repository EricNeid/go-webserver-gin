package server

import (
	"strings"
	"sync"
)

// LogService is an io.WriteCloser that collects log entries.
type LogService struct {
	mu               sync.Mutex
	Logs             []string
	Max              int
	MaxMessageLength int
}

// Write implements io.Writer.
func (l *LogService) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := string(p)
	if l.MaxMessageLength > 0 {
		limit := min(len(msg), l.MaxMessageLength)
		msg = msg[:limit]
	}
	l.Logs = append(l.Logs, strings.TrimSuffix(msg, "\n"))

	if l.Max > 0 && len(l.Logs) > l.Max {
		first := len(l.Logs) - l.Max
		l.Logs = l.Logs[first:]
	}

	return len(p), err
}
