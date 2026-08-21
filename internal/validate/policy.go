package validate

import (
        "fmt"
        "strings"

        "github.com/LYH2263/go-abacpdp/internal/policy"
)

var knownCombine = map[string]bool{
        "denyoverrides": true, "deny-overrides": true, "deny_overrides": true,
        "permitoverrides": true, "permit-overrides": true, "permit_overrides": true,
        "firstapplicable": true, "first-applicable": true, "first_applicable": true,
        "onlyoneapplicable": true, "only-one-applicable": true, "only_one_applicable": true,
}

func PolicySet(set *policy.Set) error {
        if set == nil {
                return fmt.Errorf("nil set")
        }
        if strings.TrimSpace(set.ID) == "" {
                return fmt.Errorf("empty set id")
        }
        if set.Combine != "" && !knownCombine[strings.ToLower(strings.TrimSpace(set.Combine))] {
                return fmt.Errorf("unknown combine %q", set.Combine)
        }
        if len(set.Policies) == 0 {
                return fmt.Errorf("empty policies")
        }
        seen := map[string]bool{}
        for i := range set.Policies {
                if err := Policy(&set.Policies[i]); err != nil {
                        return fmt.Errorf("policy[%d]: %w", i, err)
                }
                if seen[set.Policies[i].ID] {
                        return fmt.Errorf("duplicate policy id %q", set.Policies[i].ID)
                }
                seen[set.Policies[i].ID] = true
        }
        return nil
}

func Policy(p *policy.Policy) error {
        if p == nil || strings.TrimSpace(p.ID) == "" {
                return fmt.Errorf("empty policy id")
        }
        if p.Combine != "" && !knownCombine[strings.ToLower(strings.TrimSpace(p.Combine))] {
                return fmt.Errorf("unknown combine %q", p.Combine)
        }
        if len(p.Rules) == 0 {
                return fmt.Errorf("no rules")
        }
        for i, r := range p.Rules {
                if strings.TrimSpace(r.ID) == "" {
                        return fmt.Errorf("rule[%d]: empty id", i)
                }
                eff := strings.ToLower(r.Effect)
                if eff != "permit" && eff != "deny" {
                        return fmt.Errorf("rule[%d]: bad effect %q", i, r.Effect)
                }
                for _, m := range append(append([]policy.Match{}, r.Target.AllOf...), r.Target.AnyOf...) {
                        if err := Match(m); err != nil {
                                return fmt.Errorf("rule[%d]: %w", i, err)
                        }
                }
        }
        for _, m := range append(append([]policy.Match{}, p.Target.AllOf...), p.Target.AnyOf...) {
                if err := Match(m); err != nil {
                        return err
                }
        }
        return nil
}

func Match(m policy.Match) error {
        switch strings.ToLower(m.Category) {
        case "subject", "sub", "resource", "res", "action", "act", "environment", "env":
        default:
                return fmt.Errorf("bad category %q", m.Category)
        }
        if strings.TrimSpace(m.Attribute) == "" {
                return fmt.Errorf("empty attribute")
        }
        switch strings.ToLower(m.Op) {
        case "eq", "equals", "neq", "not_equals", "in", "exists", "prefix", "contains", "gt", "gte", "lt", "lte":
        default:
                return fmt.Errorf("bad op %q", m.Op)
        }
        return nil
}
