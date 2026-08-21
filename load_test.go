package abacpdp

import (
	"errors"
	"testing"

	"github.com/LYH2263/go-abacpdp/internal/persist"
	"github.com/LYH2263/go-abacpdp/internal/policy"
)

// newPDPWithStore 构造带自定义存储的 PDP。
func newPDPWithStore(s persist.Store) *PDP {
	return New(WithStore(s))
}

// TestLoadFailingStoreKeepsOldSet 验证：持久化失败时不得激活新 PolicySet。
// 即 PolicyView 仍见旧 ID，Evaluate 仍用旧策略，Stats 的 Loads 不递增。
func TestLoadFailingStoreKeepsOldSet(t *testing.T) {
	p := newPDPWithStore(persist.NewMemoryStore())

	// 先装载一个会通过的基线策略。
	base := policy.NewBuilder("base").Version(1).Policy(policy.Policy{
		ID: "p0",
		Rules: []policy.Rule{
			policy.PermitRule("r0", policy.Eq("subject", "role", "admin")),
		},
	}).Build()
	if err := p.Load(base); err != nil {
		t.Fatalf("load base: %v", err)
	}
	if got := p.PolicyView().ID; got != "base" {
		t.Fatalf("base view id = %q, want base", got)
	}

	// 注入一个总会失败的存储，再装载新策略。
	p.store = &persist.FailingStore{}
	fresh := policy.NewBuilder("fresh").Version(2).Policy(policy.Policy{
		ID: "p1",
		Rules: []policy.Rule{
			policy.DenyRule("r1", policy.Eq("subject", "role", "admin")),
		},
	}).Build()

	loadsBefore := p.loads.Load()
	err := p.Load(fresh)
	if !errors.Is(err, ErrPersist) {
		t.Fatalf("load fresh err = %v, want ErrPersist", err)
	}

	// 旧 PolicySet 必须仍在位。
	if got := p.PolicyView().ID; got != "base" {
		t.Fatalf("view id after failed persist = %q, want base (old set must remain active)", got)
	}
	if got := p.PolicyView().Version; got != 1 {
		t.Fatalf("view version after failed persist = %d, want 1", got)
	}

	// Evaluate 仍用旧策略：admin 应 Permit，而非被新策略 Deny。
	dec, err := p.Evaluate(AttrBag{Subject: map[string]any{"role": "admin"}})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !dec.IsPermit() {
		t.Fatalf("evaluate effect = %s, want Permit (old set must drive evaluation)", dec.Effect)
	}

	// Loads 计数不得因失败而递增。
	if got := p.loads.Load(); got != loadsBefore {
		t.Fatalf("loads after failed persist = %d, want %d (no partial activation)", got, loadsBefore)
	}
}
