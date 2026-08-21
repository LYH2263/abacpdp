package match

type Bag struct {
        Subject     map[string]any
        Resource    map[string]any
        Action      map[string]any
        Environment map[string]any
}

func (b Bag) Get(category, attr string) (any, bool) {
        var m map[string]any
        switch category {
        case "subject", "sub":
                m = b.Subject
        case "resource", "res":
                m = b.Resource
        case "action", "act":
                m = b.Action
        case "environment", "env":
                m = b.Environment
        default:
                return nil, false
        }
        if m == nil {
                return nil, false
        }
        v, ok := m[attr]
        return v, ok
}
