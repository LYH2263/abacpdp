package clone

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// MapAny 深拷贝 map[string]any：返回新 map，内部嵌套 map/slice 均独立，
// 不与入参共享任何底层容器。nil 入参原样返回。
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

// cloneValue 深拷贝单值，覆盖 map[string]any、[]any、[]string 与标量。
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

// fingerprintNil 是 nil map 的占位标记，与非 nil 空 map 区分开。
const fingerprintNil = "<nil>"

// FingerprintBag 为若干属性 map 生成内容相关的稳定指纹。
//
// 键必须只取决于各 map 的内容，而与 map 的指针身份、map 遍历顺序无关：
// 同一 AttrBag 在内容不变时键恒定；内容改写后键必然变化。这样判决缓存
// 不会因入参 Subject map 被就地改写（指针未变、内容已变）而脏命中。
//
// 实现上对每个 map 按 key 排序后序列化为 "<k>:<vjson>;<k2>:<v2json>;..."，
// 再对分段拼接后取 SHA-256。用哈希而非明文，避免键随属性体积膨胀撑爆
// map 内存，同时规避 nested map 内部遍历序对明文拼接的污染。
func FingerprintBag(parts ...map[string]any) string {
	var b strings.Builder
	for i, m := range parts {
		if i > 0 {
			b.WriteByte('|')
		}
		b.WriteString(fingerprintMap(m))
	}
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

// fingerprintMap 产出单个 map 的内容指纹：键排序后 "<k>:<v>;" 拼接。
// nil 与非 nil 空 map 分别给不同标记，二者内容语义不同，不应合并。
func fingerprintMap(m map[string]any) string {
	if m == nil {
		return fingerprintNil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte(':')
		// json.Marshal 给出稳定、确定性的值表示；对无法编码的值回退到 fmt。
		raw, err := json.Marshal(m[k])
		if err != nil {
			b.WriteString(fmt.Sprintf("%v", m[k]))
		} else {
			b.Write(raw)
		}
		b.WriteByte(';')
	}
	return b.String()
}
