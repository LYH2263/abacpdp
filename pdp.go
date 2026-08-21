package abacpdp

import (
	"sync"
	"sync/atomic"

	"github.com/LYH2263/go-abacpdp/internal/audit"
	"github.com/LYH2263/go-abacpdp/internal/metrics"
	"github.com/LYH2263/go-abacpdp/internal/obligation"
	"github.com/LYH2263/go-abacpdp/internal/persist"
	"github.com/LYH2263/go-abacpdp/internal/policy"
	"github.com/LYH2263/go-abacpdp/internal/resolve"
)

// PDP 是 ABAC 策略判决点。
type PDP struct {
	mu             sync.RWMutex
	set            *policy.Set
	store          persist.Store
	auditor        *audit.Logger
	oblAudit       *obligation.Auditor
	resolver       resolve.Resolver
	metrics        *metrics.Registry
	cache          map[string]Decision
	cacheSize      int
	defaultCombine string
	closed         atomic.Bool
	evaluations    atomic.Int64
	permits        atomic.Int64
	denies         atomic.Int64
	cacheHits      atomic.Int64
	loads          atomic.Int64
}

// New 创建 PDP。
func New(opts ...Option) *PDP {
	p := &PDP{
		store:          persist.NewMemoryStore(),
		metrics:        metrics.NewRegistry(),
		cache:          make(map[string]Decision),
		cacheSize:      256,
		defaultCombine: "DenyOverrides",
		oblAudit:       obligation.NewAuditor(),
	}
	for _, o := range opts {
		o(p)
	}
	if p.auditor == nil {
		p.auditor = audit.NewLogger("")
	}
	return p
}

// IsClosed 报告是否已关闭。
func (p *PDP) IsClosed() bool { return p.closed.Load() }

// PolicyView 返回当前策略集只读视图。
func (p *PDP) PolicyView() PolicySetView {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return viewFromSet(p.set)
}

// Stats 返回运行统计。
func (p *PDP) Stats() Stats {
	p.mu.RLock()
	defer p.mu.RUnlock()
	st := Stats{
		Evaluations: p.evaluations.Load(),
		Permits:     p.permits.Load(),
		Denies:      p.denies.Load(),
		CacheHits:   p.cacheHits.Load(),
		Loads:       p.loads.Load(),
		Closed:      p.closed.Load(),
	}
	if p.set != nil {
		st.PolicyID = p.set.ID
		st.PolicyVer = p.set.Version
	}
	return st
}
