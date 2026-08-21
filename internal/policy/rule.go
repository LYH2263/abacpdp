package policy

type Rule struct {
        ID          string       `json:"id"`
        Effect      string       `json:"effect"`
        Target      Target       `json:"target"`
        Obligations []Obligation `json:"obligations,omitempty"`
}

type Obligation struct {
        ID         string            `json:"id"`
        Attribute  string            `json:"attribute,omitempty"`
        FulfillOn  string            `json:"fulfill_on,omitempty"`
        Parameters map[string]string `json:"parameters,omitempty"`
}

func (r Rule) Clone() Rule {
        out := r
        out.Target = r.Target.Clone()
        if r.Obligations != nil {
                out.Obligations = make([]Obligation, len(r.Obligations))
                for i, o := range r.Obligations {
                        out.Obligations[i] = o.Clone()
                }
        }
        return out
}

func (o Obligation) Clone() Obligation {
        out := o
        if o.Parameters != nil {
                out.Parameters = make(map[string]string, len(o.Parameters))
                for k, v := range o.Parameters {
                        out.Parameters[k] = v
                }
        }
        return out
}
