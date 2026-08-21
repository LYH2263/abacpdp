package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LYH2263/go-abacpdp/internal/audit"
)

func TestBug09_AuditRotateClosesHandle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")
	l := audit.NewLogger(path)
	l.SetMaxBytes(40)
	payload := map[string]any{"k": "0123456789abcdef0123456789abcdef"}
	if err := l.Write("evaluate", payload); err != nil {
		t.Fatal(err)
	}
	if err := l.Write("evaluate", payload); err != nil {
		t.Fatalf("rotate/write failed (handle still held?): %v", err)
	}
	ents, _ := os.ReadDir(dir)
	rotated := 0
	for _, e := range ents {
		if e.Name() != "audit.jsonl" {
			rotated++
		}
	}
	if rotated < 1 {
		t.Fatal("expected rotated audit file")
	}
	_ = l.Close()
}
