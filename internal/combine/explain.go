package combine

import "fmt"

type Trace struct {
        Algorithm string   `json:"algorithm"`
        Inputs    []string `json:"inputs"`
        Output    string   `json:"output"`
        Reason    string   `json:"reason"`
}

func MergeTrace(algo string, results []Result) (Result, Trace, error) {
        inputs := make([]string, len(results))
        for i, r := range results {
                inputs[i] = fmt.Sprintf("%s/%s", r.PolicyID, r.Effect)
        }
        out, err := Merge(algo, results)
        tr := Trace{Algorithm: algo, Inputs: inputs}
        if err != nil {
                tr.Reason = err.Error()
                return out, tr, err
        }
        tr.Output = string(out.Effect)
        switch out.Effect {
        case Deny:
                tr.Reason = "deny selected by algorithm priority"
        case Permit:
                tr.Reason = "permit selected; no higher-priority deny"
        case NotApplicable:
                tr.Reason = "no applicable results"
        case Indeterminate:
                tr.Reason = out.Status
        }
        return out, tr, nil
}
