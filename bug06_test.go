package abacpdp_test

import (
	"errors"
	"testing"

	abacpdp "github.com/LYH2263/go-abacpdp"
	"github.com/LYH2263/go-abacpdp/internal/persist"
	"github.com/LYH2263/go-abacpdp/internal/policy"
)

func TestBug06_PersistFailureDoesNotActivate(t *testing.T) {
	p := abacpdp.New(abacpdp.WithStore(&persist.FailingStore{}))
	defer p.Close()
	set := policy.NewBuilder("t1").Policy(policy.Policy{
		ID: "p1",
		Rules: []policy.Rule{
			policy.PermitRule("a1", policy.Eq("subject", "role", "admin")),
		},
	}).Build()
	err := p.Load(set)
	if !errors.Is(err, abacpdp.ErrPersist) {
		t.Fatalf("want ErrPersist, got %v", err)
	}
	if p.PolicyView().ID != "" {
		t.Fatalf("policy should not activate after persist failure: %+v", p.PolicyView())
	}
}
