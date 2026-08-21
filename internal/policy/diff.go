package policy

import "strings"

// DiffSummary 比较两个策略集的粗粒度差异，便于热更新日志。
type DiffSummary struct {
	IDChanged       bool     `json:"id_changed"`
	VersionDelta    int      `json:"version_delta"`
	CombineChanged  bool     `json:"combine_changed"`
	AddedPolicies   []string `json:"added_policies,omitempty"`
	RemovedPolicies []string `json:"removed_policies,omitempty"`
	ChangedPolicies []string `json:"changed_policies,omitempty"`
	RuleCountBefore int      `json:"rule_count_before"`
	RuleCountAfter  int      `json:"rule_count_after"`
}

// Diff 对比旧集与新集。
func Diff(old, neu *Set) DiffSummary {
	var d DiffSummary
	if old == nil && neu == nil {
		return d
	}
	if old == nil {
		d.AddedPolicies = policyIDs(neu)
		d.RuleCountAfter = CountRules(neu)
		if neu != nil {
			d.IDChanged = true
			d.CombineChanged = neu.Combine != ""
			d.VersionDelta = neu.Version
		}
		return d
	}
	if neu == nil {
		d.RemovedPolicies = policyIDs(old)
		d.RuleCountBefore = CountRules(old)
		d.IDChanged = true
		return d
	}
	d.IDChanged = old.ID != neu.ID
	d.VersionDelta = neu.Version - old.Version
	d.CombineChanged = !strings.EqualFold(old.Combine, neu.Combine)
	d.RuleCountBefore = CountRules(old)
	d.RuleCountAfter = CountRules(neu)

	oldMap := map[string]Policy{}
	for _, p := range old.Policies {
		oldMap[p.ID] = p
	}
	newMap := map[string]Policy{}
	for _, p := range neu.Policies {
		newMap[p.ID] = p
	}
	for id := range newMap {
		if _, ok := oldMap[id]; !ok {
			d.AddedPolicies = append(d.AddedPolicies, id)
		}
	}
	for id := range oldMap {
		if _, ok := newMap[id]; !ok {
			d.RemovedPolicies = append(d.RemovedPolicies, id)
		}
	}
	for id, np := range newMap {
		op, ok := oldMap[id]
		if !ok {
			continue
		}
		if len(op.Rules) != len(np.Rules) || !strings.EqualFold(op.Combine, np.Combine) {
			d.ChangedPolicies = append(d.ChangedPolicies, id)
		}
	}
	return d
}

func policyIDs(set *Set) []string {
	if set == nil {
		return nil
	}
	out := make([]string, 0, len(set.Policies))
	for _, p := range set.Policies {
		out = append(out, p.ID)
	}
	return out
}

// HasDenyRule 是否包含 Deny 效果规则。
func HasDenyRule(set *Set) bool {
	found := false
	WalkRules(set, func(_ string, r Rule) bool {
		if strings.EqualFold(r.Effect, "Deny") {
			found = true
			return false
		}
		return true
	})
	return found
}
