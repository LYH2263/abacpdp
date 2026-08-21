package combine

func FirstApplicable(results []Result) Result {
        for _, r := range results {
                if r.Effect != NotApplicable {
                        return r
                }
        }
        return Result{Effect: NotApplicable, Status: "no applicable"}
}

func OnlyOneApplicable(results []Result) Result {
        var app []Result
        for _, r := range results {
                if r.Effect != NotApplicable {
                        app = append(app, r)
                }
        }
        switch len(app) {
        case 0:
                return Result{Effect: NotApplicable, Status: "no applicable"}
        case 1:
                return app[0]
        default:
                return Result{Effect: Indeterminate, Status: "multiple applicable"}
        }
}
