package policy

// SampleDocs 返回文档访问演示策略。
func SampleDocs() *Set {
        return NewBuilder("demo").
                Version(1).
                Combine("DenyOverrides").
                Policy(Policy{
                        ID: "docs",
                        Rules: []Rule{
                                DenyRule("deny-intern-secret",
                                        Eq("subject", "role", "intern"),
                                        Eq("resource", "class", "secret"),
                                ),
                                {
                                        ID:     "allow-admin-read",
                                        Effect: "Permit",
                                        Target: Target{AllOf: []Match{
                                                Eq("subject", "role", "admin"),
                                                Eq("action", "name", "read"),
                                        }},
                                        Obligations: []Obligation{{
                                                ID: "audit-read", FulfillOn: "Permit",
                                                Parameters: map[string]string{"level": "info"},
                                        }},
                                },
                                PermitRule("allow-owner",
                                        Eq("subject", "id", "alice"),
                                        Eq("resource", "owner", "alice"),
                                ),
                        },
                }).Build()
}

// SampleMultiPolicy 多策略集，用于集合级合并测试。
func SampleMultiPolicy() *Set {
        return NewBuilder("multi").
                Combine("DenyOverrides").
                Policy(Policy{
                        ID: "hr",
                        Target: Target{AllOf: []Match{Eq("resource", "dept", "hr")}},
                        Rules: []Rule{
                                PermitRule("hr-read", Eq("action", "name", "read"), In("subject", "role", "hr", "admin")),
                                DenyRule("hr-deny-guest", Eq("subject", "role", "guest")),
                        },
                }).
                Policy(Policy{
                        ID: "eng",
                        Target: Target{AllOf: []Match{Eq("resource", "dept", "eng")}},
                        Rules: []Rule{
                                PermitRule("eng-write", Eq("action", "name", "write"), Eq("subject", "role", "engineer")),
                        },
                }).Build()
}

func NormalizeEffect(s string) string {
        switch s {
        case "permit", "Permit", "PERMIT", "allow", "Allow":
                return "Permit"
        case "deny", "Deny", "DENY", "reject", "Reject":
                return "Deny"
        default:
                return s
        }
}
