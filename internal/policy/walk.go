package policy

func WalkRules(set *Set, fn func(policyID string, rule Rule) bool) {
        if set == nil {
                return
        }
        for _, p := range set.Policies {
                for _, r := range p.Rules {
                        if !fn(p.ID, r) {
                                return
                        }
                }
        }
}

func CountRules(set *Set) int {
        n := 0
        WalkRules(set, func(string, Rule) bool { n++; return true })
        return n
}

func FindRule(set *Set, ruleID string) (policyID string, rule Rule, ok bool) {
        WalkRules(set, func(pid string, r Rule) bool {
                if r.ID == ruleID {
                        policyID, rule, ok = pid, r, true
                        return false
                }
                return true
        })
        return
}

func ListEffects(set *Set) []string {
        var out []string
        WalkRules(set, func(_ string, r Rule) bool {
                out = append(out, r.Effect)
                return true
        })
        return out
}
