package abacpdp

import (
	"github.com/LYH2263/go-abacpdp/internal/audit"
	"github.com/LYH2263/go-abacpdp/internal/persist"
	"github.com/LYH2263/go-abacpdp/internal/resolve"
)

// Option 配置 PDP。
type Option func(*PDP)

// WithPersistDir 启用文件快照目录。
func WithPersistDir(dir string) Option {
	return func(p *PDP) { p.store = persist.NewFileStore(dir) }
}

// WithMemoryStore 使用内存快照。
func WithMemoryStore() Option {
	return func(p *PDP) { p.store = persist.NewMemoryStore() }
}

// WithStore 注入自定义存储。
func WithStore(s persist.Store) Option {
	return func(p *PDP) {
		if s != nil {
			p.store = s
		}
	}
}

// WithAuditPath 启用审计日志路径。
func WithAuditPath(path string) Option {
	return func(p *PDP) { p.auditor = audit.NewLogger(path) }
}

// WithCacheSize 设置判决缓存容量。
func WithCacheSize(n int) Option {
	return func(p *PDP) {
		if n > 0 {
			p.cacheSize = n
		}
	}
}

// WithRemoteResolver 注入远程属性解析器。
func WithRemoteResolver(r resolve.Resolver) Option {
	return func(p *PDP) { p.resolver = r }
}

// WithDefaultCombine 设置策略集缺省合并算法。
func WithDefaultCombine(name string) Option {
	return func(p *PDP) {
		if name != "" {
			p.defaultCombine = name
		}
	}
}
