package abacpdp

import "github.com/LYH2263/go-abacpdp/internal/policy"

// PolicySetView 当前策略集只读视图。
type PolicySetView struct {
	ID       string   `json:"id"`
	Version  int      `json:"version"`
	Combine  string   `json:"combine"`
	Policies []string `json:"policies"`
}

// Stats 运行统计。
type Stats struct {
	Evaluations int64  `json:"evaluations"`
	Permits     int64  `json:"permits"`
	Denies      int64  `json:"denies"`
	CacheHits   int64  `json:"cache_hits"`
	Loads       int64  `json:"loads"`
	Closed      bool   `json:"closed"`
	PolicyID    string `json:"policy_id,omitempty"`
	PolicyVer   int    `json:"policy_version"`
}

func viewFromSet(ps *policy.Set) PolicySetView {
	if ps == nil {
		return PolicySetView{}
	}
	ids := make([]string, 0, len(ps.Policies))
	for _, p := range ps.Policies {
		ids = append(ids, p.ID)
	}
	return PolicySetView{ID: ps.ID, Version: ps.Version, Combine: ps.Combine, Policies: ids}
}
