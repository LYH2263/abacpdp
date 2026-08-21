package abacpdp_test

import (
	"testing"

	abacpdp "github.com/LYH2263/go-abacpdp"
	"github.com/LYH2263/go-abacpdp/internal/policy"
)

func TestBug02_ObligationsSliceAlias(t *testing.T) {
	p := abacpdp.New()
	defer p.Close()
	set := policy.NewBuilder("obl").Policy(policy.Policy{
		ID: "p",
		Rules: []policy.Rule{{
			ID: "r", Effect: "Permit",
			Obligations: []policy.Obligation{{
				ID: "o1", FulfillOn: "Permit",
				Parameters: map[string]string{"k": "v"},
			}},
		}},
	}).Build()
	if err := p.Load(set); err != nil {
		t.Fatal(err)
	}
	d, err := p.Evaluate(abacpdp.NewAttrBag())
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Obligations) != 1 {
		t.Fatalf("obls=%v", d.Obligations)
	}
	d.Obligations[0].Parameters["k"] = "mutated"
	d.Obligations[0].ID = "hacked"
	d2, err := p.Evaluate(abacpdp.NewAttrBag())
	if err != nil {
		t.Fatal(err)
	}
	if d2.Obligations[0].Parameters["k"] != "v" || d2.Obligations[0].ID != "o1" {
		t.Fatalf("obligations aliased with cache: %+v", d2.Obligations[0])
	}
}
