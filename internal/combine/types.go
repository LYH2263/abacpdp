package combine

type Effect string

const (
        Permit        Effect = "Permit"
        Deny          Effect = "Deny"
        NotApplicable Effect = "NotApplicable"
        Indeterminate Effect = "Indeterminate"
)

type Obligation struct {
        ID         string
        Attribute  string
        FulfillOn  string
        Parameters map[string]string
}

type Result struct {
        Effect      Effect
        Status      string
        Obligations []Obligation
        Matched     []string
        PolicyID    string
}
