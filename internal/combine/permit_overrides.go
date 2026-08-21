package combine

func PermitOverrides(results []Result) Result {
        var permits, denies, inds []Result
        for _, r := range results {
                switch r.Effect {
                case Permit:
                        permits = append(permits, r)
                case Deny:
                        denies = append(denies, r)
                case Indeterminate:
                        inds = append(inds, r)
                }
        }
        if len(permits) > 0 {
                return collect(Permit, permits)
        }
        if len(inds) > 0 && len(denies) == 0 {
                return collect(Indeterminate, inds)
        }
        if len(denies) > 0 {
                return collect(Deny, denies)
        }
        return Result{Effect: NotApplicable, Status: "no applicable"}
}
