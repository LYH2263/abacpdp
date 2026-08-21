package abacpdp

import (
        "encoding/json"

        "github.com/LYH2263/go-abacpdp/internal/policy"
)

// LoadSampleDocs 装载文档演示策略。
func LoadSampleDocs(p *PDP) error {
        return p.Load(policy.SampleDocs())
}

// EvaluateMaps 从四张 map 便捷判决。
func (p *PDP) EvaluateMaps(sub, res, act, env map[string]any) (Decision, error) {
        return p.Evaluate(AttrBag{Subject: sub, Resource: res, Action: act, Environment: env})
}

// DecisionJSON 序列化判决。
func DecisionJSON(d Decision) ([]byte, error) { return json.Marshal(d) }

// ParseAttrBagJSON 解析属性袋。
func ParseAttrBagJSON(raw []byte) (AttrBag, error) {
        var b AttrBag
        err := json.Unmarshal(raw, &b)
        return b, err
}
