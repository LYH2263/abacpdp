package policy

// Builder 便于测试与示例构造策略集。
type Builder struct{ set Set }

func NewBuilder(id string) *Builder {
        return &Builder{set: Set{ID: id, Version: 1, Combine: "DenyOverrides"}}
}

func (b *Builder) Version(v int) *Builder  { b.set.Version = v; return b }
func (b *Builder) Combine(c string) *Builder { b.set.Combine = c; return b }
func (b *Builder) Policy(p Policy) *Builder {
        b.set.Policies = append(b.set.Policies, p)
        return b
}
func (b *Builder) Build() *Set { return b.set.Clone() }

func PermitRule(id string, allOf ...Match) Rule {
        return Rule{ID: id, Effect: "Permit", Target: Target{AllOf: allOf}}
}
func DenyRule(id string, allOf ...Match) Rule {
        return Rule{ID: id, Effect: "Deny", Target: Target{AllOf: allOf}}
}
func Eq(cat, attr string, value any) Match {
        return Match{Category: cat, Attribute: attr, Op: "eq", Value: value}
}
func In(cat, attr string, values ...any) Match {
        list := make([]any, len(values))
        copy(list, values)
        return Match{Category: cat, Attribute: attr, Op: "in", Value: list}
}
