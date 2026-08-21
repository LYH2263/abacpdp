package abacpdp

import (
	"testing"

	"github.com/LYH2263/go-abacpdp/internal/policy"
)

// admin-read 命中带义务 audit-read（Parameters level=info）的规则。
func adminReadBag() AttrBag {
	return AttrBag{
		Subject: map[string]any{"role": "admin"},
		Action:  map[string]any{"name": "read"},
	}
}

// TestDecisionCacheNotMutatedByCaller 复现：调用方改写返回 Decision 的义务
// Parameters / ID 后，同进程再 Evaluate 应命中缓存并返回未被改脏的判决。
func TestDecisionCacheNotMutatedByCaller(t *testing.T) {
	p := New()
	if err := p.Load(policy.SampleDocs()); err != nil {
		t.Fatalf("load: %v", err)
	}

	d1, err := p.Evaluate(adminReadBag())
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	if !d1.IsPermit() || len(d1.Obligations) == 0 {
		t.Fatalf("expected permit with obligations, got %+v", d1)
	}
	if got := d1.Obligations[0].Parameters["level"]; got != "info" {
		t.Fatalf("expected level=info, got %q", got)
	}

	// 调用方改写返回判决的义务 Parameters 与 ID（“高亮 ID”）。
	d1.Obligations[0].Parameters["level"] = "tampered"
	d1.Obligations[0].ID = "tampered-id"

	// 同一属性袋再 Evaluate，命中缓存，应返回未被改脏的判决。
	d2, err := p.Evaluate(adminReadBag())
	if err != nil {
		t.Fatalf("eval2: %v", err)
	}
	if got := d2.Obligations[0].Parameters["level"]; got != "info" {
		t.Errorf("cache polluted: obligation level = %q, want %q", got, "info")
	}
	if got := d2.Obligations[0].ID; got != "audit-read" {
		t.Errorf("cache polluted: obligation id = %q, want %q", got, "audit-read")
	}
}

// TestPolicySetNotMutatedByCaller 复现：调用方改写返回判决的义务 Parameters
// 时，已加载策略集中的义务参数不应被改脏（CopyParams 曾为 no-op）。
func TestPolicySetNotMutatedByCaller(t *testing.T) {
	p := New()
	if err := p.Load(policy.SampleDocs()); err != nil {
		t.Fatalf("load: %v", err)
	}
	snap1, err := p.SnapshotJSON()
	if err != nil {
		t.Fatalf("snap1: %v", err)
	}

	d1, err := p.Evaluate(adminReadBag())
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	d1.Obligations[0].Parameters["level"] = "tampered"

	snap2, err := p.SnapshotJSON()
	if err != nil {
		t.Fatalf("snap2: %v", err)
	}
	if string(snap1) != string(snap2) {
		t.Errorf("policy set polluted by caller mutation\n--- before ---\n%s\n--- after ---\n%s", snap1, snap2)
	}
}
