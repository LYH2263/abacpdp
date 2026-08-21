package combine_test

import (
        "testing"

        "github.com/LYH2263/go-abacpdp/internal/combine"
)

func TestDenyOverrides(t *testing.T) {
        r, err := combine.Merge("DenyOverrides", []combine.Result{
                {Effect: combine.Permit, Matched: []string{"a"}},
                {Effect: combine.Deny, Matched: []string{"b"}},
        })
        if err != nil {
                t.Fatal(err)
        }
        if r.Effect != combine.Deny {
                t.Fatalf("got %v", r.Effect)
        }
}

func TestPermitOverrides(t *testing.T) {
        r, _ := combine.Merge("PermitOverrides", []combine.Result{
                {Effect: combine.Deny}, {Effect: combine.Permit},
        })
        if r.Effect != combine.Permit {
                t.Fatalf("got %v", r.Effect)
        }
}

func TestUnknown(t *testing.T) {
        _, err := combine.Merge("nope", nil)
        if err == nil {
                t.Fatal("expected err")
        }
}
