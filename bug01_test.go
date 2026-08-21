package abacpdp_test

import (
	"testing"

	abacpdp "github.com/LYH2263/go-abacpdp"
	"github.com/LYH2263/go-abacpdp/internal/policy"
)

func TestBug01_EvaluateSubjectAliasPollutesCache(t *testing.T) {
	p := abacpdp.New()
	defer p.Close()
	set := policy.NewBuilder("t1").Policy(policy.Policy{
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
	if err := p.Load(set); err != nil {
		t.Fatal(err)
	}
	sub := map[string]any{"role": "admin"}
	res := map[string]any{"class": "public"}
	act := map[string]any{"name": "read"}
	bag := abacpdp.AttrBag{Subject: sub, Resource: res, Action: act}
	d1, err := p.Evaluate(bag)
	if err != nil || !d1.IsPermit() {
		t.Fatalf("d1=%+v err=%v", d1, err)
	}
	// 原地改写同一批 map（不换指针），暴露按指针生成的脏缓存键
	sub["role"] = "intern"
	res["class"] = "secret"
	d2, err := p.Evaluate(bag)
	if err != nil {
		t.Fatal(err)
	}
	if !d2.IsDeny() {
		t.Fatalf("mutated Subject should deny, got %+v (cache key polluted?)", d2)
	}
}
