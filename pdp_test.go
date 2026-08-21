package abacpdp_test

import (
        "context"
        "encoding/json"
        "errors"
        "testing"
        "time"

        abacpdp "github.com/LYH2263/go-abacpdp"
        "github.com/LYH2263/go-abacpdp/internal/persist"
        "github.com/LYH2263/go-abacpdp/internal/policy"
        "github.com/LYH2263/go-abacpdp/internal/resolve"
)

func sampleSet() *policy.Set {
        return policy.NewBuilder("t1").
                Policy(policy.Policy{
                        ID: "p1",
                        Rules: []policy.Rule{
                                policy.DenyRule("d1",
                                        policy.Eq("subject", "role", "intern"),
                                        policy.Eq("resource", "class", "secret"),
                                ),
                                policy.PermitRule("a1",
                                        policy.Eq("subject", "role", "admin"),
                                        policy.Eq("action", "name", "read"),
                                ),
                        },
                }).Build()
}

func TestDenyOverrides(t *testing.T) {
        p := abacpdp.New()
        defer p.Close()
        if err := p.Load(sampleSet()); err != nil {
                t.Fatal(err)
        }
        d, err := p.Evaluate(abacpdp.AttrBag{
                Subject: map[string]any{"role": "intern"}, Resource: map[string]any{"class": "secret"},
                Action: map[string]any{"name": "read"},
        })
        if err != nil {
                t.Fatal(err)
        }
        if !d.IsDeny() {
                t.Fatalf("want Deny, got %+v", d)
        }
}

func TestPermitAdmin(t *testing.T) {
        p := abacpdp.New()
        defer p.Close()
        _ = p.Load(sampleSet())
        d, err := p.Evaluate(abacpdp.AttrBag{
                Subject: map[string]any{"role": "admin"}, Resource: map[string]any{"class": "public"},
                Action: map[string]any{"name": "read"},
        })
        if err != nil || !d.IsPermit() {
                t.Fatalf("got %+v err=%v", d, err)
        }
}

func TestAttrBagCloneIsolatesCache(t *testing.T) {
        p := abacpdp.New()
        defer p.Close()
        _ = p.Load(sampleSet())
        bag := abacpdp.AttrBag{
                Subject: map[string]any{"role": "admin"}, Resource: map[string]any{"class": "public"},
                Action: map[string]any{"name": "read"},
        }
        d1, err := p.Evaluate(bag)
        if err != nil || !d1.IsPermit() {
                t.Fatalf("d1=%v err=%v", d1, err)
        }
        bag.Subject["role"] = "intern"
        bag.Resource["class"] = "secret"
        d2, err := p.Evaluate(bag)
        if err != nil {
                t.Fatal(err)
        }
        if !d2.IsDeny() {
                t.Fatalf("mutated bag should deny, got %+v", d2)
        }
}

func TestObligationsNotShared(t *testing.T) {
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
        _ = p.Load(set)
        d, err := p.Evaluate(abacpdp.NewAttrBag())
        if err != nil {
                t.Fatal(err)
        }
        if len(d.Obligations) != 1 {
                t.Fatalf("obls=%v", d.Obligations)
        }
        d.Obligations[0].Parameters["k"] = "mutated"
        d2, _ := p.Evaluate(abacpdp.NewAttrBag())
        if d2.Obligations[0].Parameters["k"] != "v" {
                t.Fatalf("shared params: %v", d2.Obligations[0].Parameters)
        }
}

func TestCloseThenEvaluate(t *testing.T) {
        p := abacpdp.New()
        _ = p.Load(sampleSet())
        _ = p.Close()
        _, err := p.Evaluate(abacpdp.NewAttrBag())
        if !errors.Is(err, abacpdp.ErrClosed) {
                t.Fatalf("want ErrClosed, got %v", err)
        }
}

func TestNoPolicy(t *testing.T) {
        p := abacpdp.New()
        defer p.Close()
        _, err := p.Evaluate(abacpdp.NewAttrBag())
        if !errors.Is(err, abacpdp.ErrNoPolicy) {
                t.Fatalf("want ErrNoPolicy, got %v", err)
        }
}

func TestUnknownCombineRejected(t *testing.T) {
        p := abacpdp.New()
        defer p.Close()
        raw := []byte(`{"id":"x","combine":"TotallyUnknown","policies":[{"id":"p","rules":[{"id":"r","effect":"Permit"}]}]}`)
        err := p.LoadJSON(raw)
        if err == nil || !errors.Is(err, abacpdp.ErrInvalidPolicy) {
                t.Fatalf("want ErrInvalidPolicy, got %v", err)
        }
}

func TestContextCancel(t *testing.T) {
        p := abacpdp.New(abacpdp.WithRemoteResolver(resolve.Delayed{Delay: 2 * time.Second, Key: "x", Value: 1}))
        defer p.Close()
        _ = p.Load(sampleSet())
        ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
        defer cancel()
        _, err := p.EvaluateContext(ctx, abacpdp.AttrBag{
                Subject: map[string]any{"role": "admin"}, Action: map[string]any{"name": "read"},
        })
        if err == nil {
                t.Fatal("expected cancel")
        }
}

func TestSnapshot(t *testing.T) {
        p := abacpdp.New()
        defer p.Close()
        _ = p.Load(sampleSet())
        b, err := p.SnapshotJSON()
        if err != nil {
                t.Fatal(err)
        }
        var set policy.Set
        if err := json.Unmarshal(b, &set); err != nil {
                t.Fatal(err)
        }
        if set.ID != "t1" {
                t.Fatalf("id=%s", set.ID)
        }
}

func TestPersistFailureDoesNotActivate(t *testing.T) {
        p := abacpdp.New(abacpdp.WithStore(&persist.FailingStore{}))
        defer p.Close()
        err := p.Load(sampleSet())
        if !errors.Is(err, abacpdp.ErrPersist) {
                t.Fatalf("want ErrPersist, got %v", err)
        }
        if p.PolicyView().ID != "" {
                t.Fatalf("policy should not activate: %+v", p.PolicyView())
        }
}
