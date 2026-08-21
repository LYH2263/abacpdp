package abacpdp

// Effect 是规则/策略效果。
type Effect string

const (
	EffectPermit        Effect = "Permit"
	EffectDeny          Effect = "Deny"
	EffectNotApplicable Effect = "NotApplicable"
	EffectIndeterminate Effect = "Indeterminate"
)

// Obligation 是判决附带义务（如审计、脱敏）。
type Obligation struct {
	ID         string            `json:"id"`
	Attribute  string            `json:"attribute,omitempty"`
	FulfillOn  Effect            `json:"fulfill_on"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

// Decision 是一次 Evaluate 的结果。
type Decision struct {
	Effect      Effect       `json:"effect"`
	Status      string       `json:"status,omitempty"`
	Obligations []Obligation `json:"obligations,omitempty"`
	Matched     []string     `json:"matched,omitempty"`
	PolicyID    string       `json:"policy_id,omitempty"`
}

// Clone 返回独立副本；Obligations 切片与 Parameters 均不共享。
func (d Decision) Clone() Decision {

	return d
}

// IsPermit 判断是否允许。
func (d Decision) IsPermit() bool { return d.Effect == EffectPermit }

// IsDeny 判断是否拒绝。
func (d Decision) IsDeny() bool { return d.Effect == EffectDeny }
