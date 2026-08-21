package match_test

import (
        "testing"

        "github.com/LYH2263/go-abacpdp/internal/match"
        "github.com/LYH2263/go-abacpdp/internal/policy"
)

func TestOps(t *testing.T) {
        bag := match.Bag{Subject: map[string]any{"role": "admin", "level": 5.0}}
        if !match.Eval(policy.Match{Category: "subject", Attribute: "role", Op: "eq", Value: "admin"}, bag) {
                t.Fatal("eq")
        }
        if !match.Eval(policy.Match{Category: "subject", Attribute: "role", Op: "in", Value: []any{"admin", "user"}}, bag) {
                t.Fatal("in")
        }
        if !match.Eval(policy.Match{Category: "subject", Attribute: "level", Op: "gte", Value: 3}, bag) {
                t.Fatal("gte")
        }
        if !match.Eval(policy.Match{Category: "subject", Attribute: "role", Op: "prefix", Value: "ad"}, bag) {
                t.Fatal("prefix")
        }
}
