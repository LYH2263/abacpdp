package match

import (
        "fmt"
        "strings"

        "github.com/LYH2263/go-abacpdp/internal/policy"
)

func Eval(m policy.Match, bag Bag) bool {
        v, ok := bag.Get(m.Category, m.Attribute)
        switch strings.ToLower(m.Op) {
        case "exists":
                return ok && v != nil
        case "eq", "equals":
                return ok && equal(v, m.Value)
        case "neq", "not_equals":
                if !ok {
                        return true
                }
                return !equal(v, m.Value)
        case "in":
                return ok && inList(v, m.Value)
        case "prefix":
                return ok && strings.HasPrefix(fmt.Sprint(v), fmt.Sprint(m.Value))
        case "contains":
                return ok && strings.Contains(fmt.Sprint(v), fmt.Sprint(m.Value))
        case "gt", "gte", "lt", "lte":
                return ok && compareNum(strings.ToLower(m.Op), v, m.Value)
        default:
                return false
        }
}

func equal(a, b any) bool {
        if a == nil && b == nil {
                return true
        }
        switch av := a.(type) {
        case string:
                return av == fmt.Sprint(b)
        case bool:
                if bv, ok := b.(bool); ok {
                        return av == bv
                }
                return fmt.Sprint(a) == fmt.Sprint(b)
        case float64, float32, int, int64, int32:
                return toFloat(a) == toFloat(b)
        default:
                return fmt.Sprint(a) == fmt.Sprint(b)
        }
}

func inList(v, list any) bool {
        switch xs := list.(type) {
        case []any:
                for _, x := range xs {
                        if equal(v, x) {
                                return true
                        }
                }
        case []string:
                s := fmt.Sprint(v)
                for _, x := range xs {
                        if s == x {
                                return true
                        }
                }
        }
        return false
}

func toFloat(v any) float64 {
        switch t := v.(type) {
        case float64:
                return t
        case float32:
                return float64(t)
        case int:
                return float64(t)
        case int64:
                return float64(t)
        case int32:
                return float64(t)
        default:
                var f float64
                _, _ = fmt.Sscanf(fmt.Sprint(v), "%f", &f)
                return f
        }
}

func compareNum(op string, a, b any) bool {
        x, y := toFloat(a), toFloat(b)
        switch op {
        case "gt":
                return x > y
        case "gte":
                return x >= y
        case "lt":
                return x < y
        case "lte":
                return x <= y
        }
        return false
}
