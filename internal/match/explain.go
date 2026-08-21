package match

import (
        "fmt"
        "strings"

        "github.com/LYH2263/go-abacpdp/internal/policy"
)

type Explanation struct {
        Matched bool     `json:"matched"`
        Steps   []string `json:"steps"`
}

func ExplainTarget(t policy.Target, bag Bag) Explanation {
        var steps []string
        if t.Empty() {
                return Explanation{Matched: true, Steps: []string{"empty target => true"}}
        }
        ok := true
        for i, m := range t.AllOf {
                hit := Eval(m, bag)
                steps = append(steps, fmt.Sprintf("all_of[%d] %s.%s %s %v => %v", i, m.Category, m.Attribute, m.Op, m.Value, hit))
                if !hit {
                        ok = false
                }
        }
        if len(t.AnyOf) > 0 {
                anyHit := false
                for i, m := range t.AnyOf {
                        hit := Eval(m, bag)
                        steps = append(steps, fmt.Sprintf("any_of[%d] %s.%s %s %v => %v", i, m.Category, m.Attribute, m.Op, m.Value, hit))
                        if hit {
                                anyHit = true
                        }
                }
                ok = ok && anyHit
        }
        return Explanation{Matched: ok, Steps: steps}
}

func SummarizeBag(bag Bag) string {
        var b strings.Builder
        writeCat := func(name string, m map[string]any) {
                b.WriteString(name)
                b.WriteByte('{')
                first := true
                for k, v := range m {
                        if !first {
                                b.WriteByte(',')
                        }
                        first = false
                        fmt.Fprintf(&b, "%s=%v", k, v)
                }
                b.WriteByte('}')
        }
        writeCat("sub", bag.Subject)
        writeCat("res", bag.Resource)
        writeCat("act", bag.Action)
        writeCat("env", bag.Environment)
        return b.String()
}
