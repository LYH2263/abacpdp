package policy

type Target struct {
        AllOf []Match `json:"all_of,omitempty"`
        AnyOf []Match `json:"any_of,omitempty"`
}

type Match struct {
        Category  string `json:"category"`
        Attribute string `json:"attribute"`
        Op        string `json:"op"`
        Value     any    `json:"value,omitempty"`
}

func (t Target) Clone() Target {
        out := Target{}
        if t.AllOf != nil {
                out.AllOf = append([]Match(nil), t.AllOf...)
        }
        if t.AnyOf != nil {
                out.AnyOf = append([]Match(nil), t.AnyOf...)
        }
        return out
}

func (t Target) Empty() bool {
        return len(t.AllOf) == 0 && len(t.AnyOf) == 0
}
