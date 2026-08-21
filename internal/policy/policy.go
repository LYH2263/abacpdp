package policy

// Policy 含目标、合并算法与规则。
type Policy struct {
        ID      string `json:"id"`
        Combine string `json:"combine,omitempty"`
        Target  Target `json:"target"`
        Rules   []Rule `json:"rules"`
}

func (p Policy) Clone() Policy {
        out := p
        out.Target = p.Target.Clone()
        out.Rules = make([]Rule, len(p.Rules))
        for i := range p.Rules {
                out.Rules[i] = p.Rules[i].Clone()
        }
        return out
}
