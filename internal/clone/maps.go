package clone

import (
        "fmt"
        "sort"
        "strings"
)

// MapAny 深拷贝 map[string]any。
func MapAny(in map[string]any) map[string]any {
        if in == nil {
                return nil
        }
        out := make(map[string]any, len(in))
        for k, v := range in {
                out[k] = cloneValue(v)
        }
        return out
}

func cloneValue(v any) any {
        switch t := v.(type) {
        case map[string]any:
                return MapAny(t)
        case []any:
                out := make([]any, len(t))
                for i, x := range t {
                        out[i] = cloneValue(x)
                }
                return out
        case []string:
                return Strings(t)
        default:
                return t
        }
}

func Strings(in []string) []string {
        if in == nil {
                return nil
        }
        out := make([]string, len(in))
        copy(out, in)
        return out
}

func StringMap(in map[string]string) map[string]string {
        if in == nil {
                return nil
        }
        out := make(map[string]string, len(in))
        for k, v := range in {
                out[k] = v
        }
        return out
}

func FingerprintBag(parts ...map[string]any) string {
        var b strings.Builder
        for i, m := range parts {
                if i > 0 {
                        b.WriteByte('|')
                }
                keys := make([]string, 0, len(m))
                for k := range m {
                        keys = append(keys, k)
                }
                sort.Strings(keys)
                for _, k := range keys {
                        fmt.Fprintf(&b, "%s=%v;", k, m[k])
                }
        }
        return b.String()
}
