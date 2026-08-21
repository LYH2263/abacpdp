package abacpdp

import (
	"fmt"
	"strings"

	"github.com/LYH2263/go-abacpdp/internal/policy"
)

// ExplainDecision 返回人类可读的判决摘要（不含敏感参数细节）。
func ExplainDecision(d Decision) string {
	var b strings.Builder
	fmt.Fprintf(&b, "effect=%s", d.Effect)
	if d.PolicyID != "" {
		fmt.Fprintf(&b, " policy=%s", d.PolicyID)
	}
	if d.Status != "" {
		fmt.Fprintf(&b, " status=%s", d.Status)
	}
	if len(d.Matched) > 0 {
		fmt.Fprintf(&b, " matched=[%s]", strings.Join(d.Matched, ","))
	}
	if len(d.Obligations) > 0 {
		ids := make([]string, len(d.Obligations))
		for i, o := range d.Obligations {
			ids[i] = o.ID
		}
		fmt.Fprintf(&b, " obligations=[%s]", strings.Join(ids, ","))
	}
	return b.String()
}

// MustLoadSample 装载演示策略；失败则 panic（仅演示/测试辅助）。
func MustLoadSample(p *PDP) {
	if err := LoadSampleDocs(p); err != nil {
		panic(err)
	}
}

// PolicyDiffFrom 对比当前活动集与给定集。
func (p *PDP) PolicyDiffFrom(other *policy.Set) policy.DiffSummary {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return policy.Diff(p.set, other)
}
