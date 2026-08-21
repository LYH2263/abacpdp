package abacpdp

// Close 先刷义务审计，再清策略与缓存，最后关审计句柄。
func (p *PDP) Close() error {
        if !p.closed.CompareAndSwap(false, true) {
                return ErrClosed
        }
        if p.oblAudit != nil {
                _ = p.oblAudit.Flush(p.auditor)
        }
        p.mu.Lock()
        p.set = nil
        p.cache = nil
        p.mu.Unlock()
        if p.auditor != nil {
                _ = p.auditor.Close()
        }
        return nil
}
