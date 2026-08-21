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
	if p.store != nil {
		if err := p.store.Save(set.ID, snap); err != nil {
			return fmt.Errorf("%w: %v", ErrPersist, err)
		}
	}
	cloned := set.Clone()
	p.mu.Lock()
	p.set = cloned
	p.cache = make(map[string]Decision)
	p.mu.Unlock()
	p.loads.Add(1)
	p.metrics.IncLoads(1)
	_ = p.auditor.Write("load", map[string]any{"id": set.ID, "version": set.Version})
	return nil
}
