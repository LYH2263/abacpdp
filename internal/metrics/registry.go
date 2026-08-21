package metrics

import "sync/atomic"

type Registry struct {
        evaluations atomic.Int64
        permits     atomic.Int64
        denies      atomic.Int64
        loads       atomic.Int64
        errors      atomic.Int64
}

func NewRegistry() *Registry { return &Registry{} }

func (r *Registry) IncEvaluations(n int64) { r.evaluations.Add(n) }
func (r *Registry) IncPermits(n int64)     { r.permits.Add(n) }
func (r *Registry) IncDenies(n int64)      { r.denies.Add(n) }
func (r *Registry) IncLoads(n int64)       { r.loads.Add(n) }
func (r *Registry) IncErrors(n int64)      { r.errors.Add(n) }

type Snapshot struct {
        Evaluations int64 `json:"evaluations"`
        Permits     int64 `json:"permits"`
        Denies      int64 `json:"denies"`
        Loads       int64 `json:"loads"`
        Errors      int64 `json:"errors"`
}

func (r *Registry) Snapshot() Snapshot {
        return Snapshot{
                Evaluations: r.evaluations.Load(),
                Permits:     r.permits.Load(),
                Denies:      r.denies.Load(),
                Loads:       r.loads.Load(),
                Errors:      r.errors.Load(),
        }
}
