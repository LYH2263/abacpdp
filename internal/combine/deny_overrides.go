package combine

// DenyOverrides：任一 Deny 则 Deny；否则有 Permit 则 Permit；全 NA 则 NA。
func DenyOverrides(results []Result) Result {
        var permits, denies, inds []Result
        for _, r := range results {
                switch r.Effect {
                case Deny:
                        denies = append(denies, r)
                case Permit:
                        permits = append(permits, r)
                case Indeterminate:
                        inds = append(inds, r)
                }
        }
        if len(denies) > 0 {
                return collect(Deny, denies)
        }
        if len(inds) > 0 && len(permits) == 0 {
                return collect(Indeterminate, inds)
        }
        if len(permits) > 0 {
                return collect(Permit, permits)
        }
        return Result{Effect: NotApplicable, Status: "no applicable"}
}

func collect(eff Effect, rs []Result) Result {
        out := Result{Effect: eff, PolicyID: rs[0].PolicyID}
        for _, r := range rs {
                out.Obligations = append(out.Obligations, r.Obligations...)
                out.Matched = append(out.Matched, r.Matched...)
                if out.PolicyID == "" {
                        out.PolicyID = r.PolicyID
                }
        }
        return out
}
