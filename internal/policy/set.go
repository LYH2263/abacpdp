package policy

// Set 是策略集 AST 根。
type Set struct {
        ID       string   `json:"id"`
        Version  int      `json:"version"`
        Combine  string   `json:"combine"`
        Policies []Policy `json:"policies"`
}

func (s *Set) Clone() *Set {
        if s == nil {
                return nil
        }
        out := &Set{ID: s.ID, Version: s.Version, Combine: s.Combine, Policies: make([]Policy, len(s.Policies))}
        for i := range s.Policies {
                out.Policies[i] = s.Policies[i].Clone()
        }
        return out
}
