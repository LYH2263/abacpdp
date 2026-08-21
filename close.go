package abacpdp

// Close 先把待处理的义务审计刷入审计日志，再清策略与缓存，
// 最后关审计句柄。顺序很关键：必须先 Flush 再丢弃 oblAudit、再关 auditor，
// 否则义务记录永远不会落盘，审计文件里只剩 load/evaluate。
func (p *PDP) Close() error {
	if !p.closed.CompareAndSwap(false, true) {
		return ErrClosed
	}

	// 1) 先把 pending 义务刷进审计日志；此时 auditor 仍是打开状态。
	//    必须先于 oblAudit=nil 与 auditor.Close()，否则义务记录永远落不了盘。
	var flushErr error
	if p.oblAudit != nil && p.auditor != nil {
		flushErr = p.oblAudit.Flush(p.auditor)
	}

	// 2) 清策略与缓存。
	p.mu.Lock()
	p.set = nil
	p.cache = nil
	p.mu.Unlock()

	// 3) 丢弃义务审计器（pending 已在步骤 1 排空）。
	p.oblAudit = nil

	// 4) 最后关闭审计句柄；Flush 已完成，不再会有写入。
	if p.auditor != nil {
		if err := p.auditor.Close(); err != nil && flushErr == nil {
			flushErr = err
		}
	}
	return flushErr
}
