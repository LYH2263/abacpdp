package abacpdp_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	abacpdp "github.com/LYH2263/go-abacpdp"
	"github.com/LYH2263/go-abacpdp/internal/policy"
)

func TestBug07_EvaluateContextHonorsCancel(t *testing.T) {
	p := abacpdp.New()
	defer p.Close()
	set := &policy.Set{ID: "big", Combine: "DenyOverrides"}
	for i := 0; i < 30; i++ {
		set.Policies = append(set.Policies, policy.Policy{
			ID: fmt.Sprintf("pol%d", i),
			Rules: []policy.Rule{
				{ID: "r1", Effect: "Permit", Target: policy.Target{AllOf: []policy.Match{policy.Eq("action", "name", "read")}}},
			},
		})
	}
	if err := p.Load(set); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := p.EvaluateContext(ctx, abacpdp.AttrBag{
		Subject: map[string]any{"role": "admin"},
		Action:  map[string]any{"name": "read"},
	})
	if err == nil {
		t.Fatal("expected cancel error")
	}
	if !errors.Is(err, abacpdp.ErrCanceled) && !errors.Is(err, context.Canceled) {
		t.Fatalf("want canceled, got %v", err)
	}
}
