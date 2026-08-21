package abacpdp_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	abacpdp "github.com/LYH2263/go-abacpdp"
	"github.com/LYH2263/go-abacpdp/internal/policy"
)

func TestBug10_CloseFlushesObligationsFirst(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")
	p := abacpdp.New(abacpdp.WithAuditPath(path))
	set := policy.NewBuilder("obl").Policy(policy.Policy{
		ID: "p",
		Rules: []policy.Rule{{
			ID: "r", Effect: "Permit",
			Obligations: []policy.Obligation{{
				ID: "o-flush", FulfillOn: "Permit",
				Parameters: map[string]string{"k": "v"},
			}},
		}},
	}).Build()
	if err := p.Load(set); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Evaluate(abacpdp.NewAttrBag()); err != nil {
		t.Fatal(err)
	}
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(raw)
	if !strings.Contains(s, "o-flush") && !strings.Contains(s, `"kind":"obligation"`) {
		t.Fatalf("Close dropped obligation audit before flush; file=%q", s)
	}
}
