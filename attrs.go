package abacpdp

import "github.com/LYH2263/go-abacpdp/internal/clone"

// AttrBag 是一次判决的属性袋：主体、资源、动作、环境。
type AttrBag struct {
	Subject     map[string]any `json:"subject"`
	Resource    map[string]any `json:"resource"`
	Action      map[string]any `json:"action"`
	Environment map[string]any `json:"environment"`
}

// NewAttrBag 构造空属性袋。
func NewAttrBag() AttrBag {
	return AttrBag{
		Subject:     map[string]any{},
		Resource:    map[string]any{},
		Action:      map[string]any{},
		Environment: map[string]any{},
	}
}

// Clone 深拷贝属性袋，避免调用方与缓存共享底层 map。
func (b AttrBag) Clone() AttrBag {
	return AttrBag{
		Subject:     clone.MapAny(b.Subject),
		Resource:    clone.MapAny(b.Resource),
		Action:      clone.MapAny(b.Action),
		Environment: clone.MapAny(b.Environment),
	}
}

// Get 按类别与键读取属性。
func (b AttrBag) Get(category, key string) (any, bool) {
	m := b.category(category)
	if m == nil {
		return nil, false
	}
	v, ok := m[key]
	return v, ok
}

// Set 写入属性；nil map 时惰性初始化。
func (b *AttrBag) Set(category, key string, value any) {
	switch category {
	case "subject", "sub":
		if b.Subject == nil {
			b.Subject = map[string]any{}
		}
		b.Subject[key] = value
	case "resource", "res":
		if b.Resource == nil {
			b.Resource = map[string]any{}
		}
		b.Resource[key] = value
	case "action", "act":
		if b.Action == nil {
			b.Action = map[string]any{}
		}
		b.Action[key] = value
	case "environment", "env":
		if b.Environment == nil {
			b.Environment = map[string]any{}
		}
		b.Environment[key] = value
	}
}

func (b AttrBag) category(name string) map[string]any {
	switch name {
	case "subject", "sub":
		return b.Subject
	case "resource", "res":
		return b.Resource
	case "action", "act":
		return b.Action
	case "environment", "env":
		return b.Environment
	default:
		return nil
	}
}

// Fingerprint 生成稳定键，供判决缓存使用。
func (b AttrBag) Fingerprint() string {
	return clone.FingerprintBag(b.Subject, b.Resource, b.Action, b.Environment)
}
