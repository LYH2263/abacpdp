package abacpdp

import (
        "encoding/json"
        "fmt"
)

// SnapshotJSON 导出当前活动策略集。
func (p *PDP) SnapshotJSON() ([]byte, error) {
        if p.closed.Load() {
                return nil, ErrClosed
        }
        p.mu.RLock()
        defer p.mu.RUnlock()
        if p.set == nil {
                return nil, ErrNoPolicy
        }
        b, err := json.MarshalIndent(p.set, "", "  ")
        if err != nil {
                return nil, fmt.Errorf("snapshot: %w", err)
        }
        return b, nil
}

// ClearCache 清空判决缓存。
func (p *PDP) ClearCache() {
        p.mu.Lock()
        defer p.mu.Unlock()
        p.cache = make(map[string]Decision)
}
