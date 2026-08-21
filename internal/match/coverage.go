package match

import (
	"fmt"
	"sort"

	"github.com/LYH2263/go-abacpdp/internal/policy"
)

// Coverage 统计属性袋相对目标条件的覆盖情况。
type Coverage struct {
	Required int      `json:"required"`
	Hit      int      `json:"hit"`
	Missing  []string `json:"missing,omitempty"`
}

// CoverMatches 统计匹配条件所需属性是否齐全。
func CoverMatches(matches []policy.Match, bag Bag) Coverage {
	var cov Coverage
	for _, m := range matches {
		cov.Required++
		key := m.Category + "." + m.Attribute
		if _, ok := bag.Get(m.Category, m.Attribute); ok {
			cov.Hit++
		} else {
			cov.Missing = append(cov.Missing, key)
		}
	}
	sort.Strings(cov.Missing)
	return cov
}

// CoverTarget 覆盖 Target.AllOf。
func CoverTarget(t policy.Target, bag Bag) Coverage {
	return CoverMatches(t.AllOf, bag)
}

// CollectKeys 收集属性袋全部键路径。
func CollectKeys(bag Bag) []string {
	var out []string
	add := func(cat string, m map[string]any) {
		for k := range m {
			out = append(out, cat+"."+k)
		}
	}
	add("subject", bag.Subject)
	add("resource", bag.Resource)
	add("action", bag.Action)
	add("environment", bag.Environment)
	sort.Strings(out)
	return out
}

// FormatMatch 格式化单条匹配描述。
func FormatMatch(m policy.Match) string {
	return fmt.Sprintf("%s.%s %s %v", m.Category, m.Attribute, m.Op, m.Value)
}

// Ratio 命中比例，Required 为 0 时返回 1。
func (c Coverage) Ratio() float64 {
	if c.Required == 0 {
		return 1
	}
	return float64(c.Hit) / float64(c.Required)
}
