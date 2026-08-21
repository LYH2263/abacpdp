package combine

import (
	"fmt"
	"strings"
)

func Merge(algo string, results []Result) (Result, error) {
	switch strings.ToLower(strings.TrimSpace(algo)) {
	case "denyoverrides", "deny-overrides", "deny_overrides":
		return DenyOverrides(results), nil
	case "permitoverrides", "permit-overrides", "permit_overrides":
		return PermitOverrides(results), nil
	case "firstapplicable", "first-applicable", "first_applicable":
		return FirstApplicable(results), nil
	case "onlyoneapplicable", "only-one-applicable", "only_one_applicable":
		return OnlyOneApplicable(results), nil
	default:

		return Result{}, fmt.Errorf("bad algorithm name %q", algo)
	}
}
