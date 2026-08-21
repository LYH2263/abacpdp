package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Logger 追加写审计事件；按大小轮转并关闭旧句柄。
type Logger struct {
	mu       sync.Mutex
	path     string
	f        *os.File
	maxBytes int64
	written  int64
	sink     []map[string]any
}

func NewLogger(path string) *Logger {
	return &Logger{path: path, maxBytes: 1 << 20}
}

func (l *Logger) Write(kind string, fields map[string]any) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	ev := map[string]any{"ts": time.Now().UTC().Format(time.RFC3339Nano), "kind": kind}
	for k, v := range fields {
		ev[k] = v
	}
	if l.path == "" {
		l.sink = append(l.sink, ev)
		return nil
	}
	if err := l.ensureLocked(); err != nil {
		return err
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	n, err := l.f.Write(b)
	l.written += int64(n)
	if err != nil {
		return err
	}
	if l.written >= l.maxBytes {
		return l.rotateLocked()
	}
	return nil
}

func (l *Logger) ensureLocked() error {
	if l.f != nil {
		return nil
	}
	dir := filepath.Dir(l.path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	l.f = f
	if st, err := f.Stat(); err == nil {
		l.written = st.Size()
	}
	return nil
}

func (l *Logger) rotateLocked() error {
	rotated := l.path + "." + time.Now().UTC().Format("20060102T150405")
	// Windows 下当前句柄仍占用 audit.jsonl，直接 Rename 会因文件被占用而失败
	// （归档写不出，后续 Write 也跟着失败）。改名前必须先 Close 旧句柄。
	if l.f != nil {
		_ = l.f.Sync()
		_ = l.f.Close()
		l.f = nil
	}
	if err := os.Rename(l.path, rotated); err != nil {
		return err
	}
	l.written = 0
	return l.ensureLocked()
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f != nil {
		err := l.f.Close()
		l.f = nil
		return err
	}
	return nil
}

func (l *Logger) SetMaxBytes(n int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if n > 0 {
		l.maxBytes = n
	}
}

func (l *Logger) MemoryEvents() []map[string]any {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]map[string]any, len(l.sink))
	copy(out, l.sink)
	return out
}
