package abacpdp_test

import (
	"testing"

	abacpdp "github.com/LYH2263/go-abacpdp"
	"github.com/LYH2263/go-abacpdp/internal/policy"
)

func TestBug04_NilObligationParamsNoDirtyCache(t *testing.T) {
	p := abacpdp.New()
	defer p.Close()
	set := policy.NewBuilder("nil-obl").Combine("DenyOverrides").Policy(policy.Policy{
		ID: "docs",
		Rules: []policy.Rule{{
			ID:     "allow",
			Effect: "Permit",
			Target: policy.Target{AllOf: []policy.Match{
				policy.Eq("subject", "role", "admin"),
				policy.Eq("action", "name", "read"),
			}},
			Obligations: []policy.Obligation{{
				ID: "audit", FulfillOn: "Permit", Parameters: nil,
			}},
		}},
	}).Build()
	if err := p.Load(set); err != nil {
		t.Fatal(err)
	}
	bag := abacpdp.AttrBag{
		Subject: map[string]any{"role": "admin"},
		Action:  map[string]any{"name": "read"},
	}
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("Evaluate panicked on nil obligation params: %v", rec)
		}
	}()
	d1, err := p.Evaluate(bag)
	if err != nil {
		t.Fatal(err)
	}
	if !d1.IsPermit() {
		t.Fatalf("effect=%s", d1.Effect)
	}
	if len(d1.Obligations) != 1 || d1.Obligations[0].Parameters == nil {
		t.Fatalf("obligations=%+v", d1.Obligations)
	}
	if d1.Obligations[0].Parameters["rule"] != "allow" {
		t.Fatalf("missing rule stamp: %+v", d1.Obligations[0].Parameters)
	}
	d2, err := p.Evaluate(bag)
	if err != nil {
		t.Fatal(err)
	}
	if len(d2.Obligations) != 1 || d2.Obligations[0].Parameters["rule"] != "allow" {
		t.Fatalf("dirty/partial cache returned: %+v (hits=%d)", d2, p.Stats().CacheHits)
	}
}
