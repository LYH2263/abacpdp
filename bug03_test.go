package abacpdp_test

import (
	"errors"
	"testing"

	abacpdp "github.com/LYH2263/go-abacpdp"
	"github.com/LYH2263/go-abacpdp/internal/policy"
)

func TestBug03_EvaluateAfterCloseNoPanic(t *testing.T) {
	p := abacpdp.New()
	set := policy.NewBuilder("t1").Policy(policy.Policy{
		ID: "p1",
		Rules: []policy.Rule{
			policy.PermitRule("a1", policy.Eq("subject", "role", "admin")),
		},
	}).Build()
	if err := p.Load(set); err != nil {
		t.Fatal(err)
	}
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("Evaluate after Close panicked: %v", rec)
		}
	}()
	_, err := p.Evaluate(abacpdp.AttrBag{
		Subject: map[string]any{"role": "admin"},
		Action:  map[string]any{"name": "read"},
	})
	if !errors.Is(err, abacpdp.ErrClosed) {
		t.Fatalf("want ErrClosed, got %v", err)
	}
}
