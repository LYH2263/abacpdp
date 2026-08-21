package abacpdp_test

import (
	"context"
	"errors"
	"testing"
	"time"

	abacpdp "github.com/LYH2263/go-abacpdp"
	"github.com/LYH2263/go-abacpdp/internal/policy"
	"github.com/LYH2263/go-abacpdp/internal/resolve"
)

func TestBug08_RemoteResolveHonorsContext(t *testing.T) {
	p := abacpdp.New(abacpdp.WithRemoteResolver(resolve.Delayed{Delay: 2 * time.Second, Key: "x", Value: 1}))
	defer p.Close()
	set := policy.NewBuilder("t1").Policy(policy.Policy{
		ID: "p1",
		Rules: []policy.Rule{
			policy.PermitRule("a1", policy.Eq("subject", "role", "admin")),
		},
	}).Build()
	if err := p.Load(set); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := p.EvaluateContext(ctx, abacpdp.AttrBag{
		Subject: map[string]any{"role": "admin"},
		Action:  map[string]any{"name": "read"},
	})
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("expected cancel from remote resolve")
	}
	if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, abacpdp.ErrCanceled) {
		t.Fatalf("want ctx cancel/deadline, got %v", err)
	}
	if elapsed > 500*time.Millisecond {
		t.Fatalf("resolve ignored cancel, elapsed=%v", elapsed)
	}
}
