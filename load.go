package abacpdp

import (
	"encoding/json"
	"fmt"

	"github.com/LYH2263/go-abacpdp/internal/policy"
	"github.com/LYH2263/go-abacpdp/internal/validate"
)

// LoadJSON 从 JSON 装载策略集。
func (p *PDP) LoadJSON(raw []byte) error {
	if p.closed.Load() {
		return ErrClosed
	}
	var set policy.Set
	if err := json.Unmarshal(raw, &set); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidPolicy, err)
	}
	return p.Load(&set)
}

// Load 先校验再持久化，成功后才激活新 PolicySet。
// 持久化失败时不得切换已激活的 PolicySet，否则会留下半成功状态：
// PolicyView 已见新 ID 而 Evaluate 会用到未持久化的策略。
func (p *PDP) Load(set *policy.Set) error {
	if p.closed.Load() {
		return ErrClosed
	}
	if set == nil {
		return fmt.Errorf("%w: nil set", ErrInvalidPolicy)
	}
	if set.Combine == "" {
		set.Combine = p.defaultCombine
	}
	if err := validate.PolicySet(set); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidPolicy, err)
	}
	snap, err := json.Marshal(set)
	if err != nil {
		return fmt.Errorf("%w: marshal: %v", ErrInvalidPolicy, err)
	}
	cloned := set.Clone()

	// 先持久化：失败则原样保留旧 PolicySet，不更新任何统计/审计。
	if p.store != nil {
		if err := p.store.Save(set.ID, snap); err != nil {
			return fmt.Errorf("%w: %v", ErrPersist, err)
		}
	}

	// 持久化成功后再激活，保证内存与存储一致。
	p.mu.Lock()
	p.set = cloned
	p.cache = make(map[string]Decision)
	p.mu.Unlock()
	p.loads.Add(1)
	p.metrics.IncLoads(1)
	_ = p.auditor.Write("load", map[string]any{"id": set.ID, "version": set.Version})
	return nil
}
