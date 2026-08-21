package obligation

import (
	"sync"

	"github.com/LYH2263/go-abacpdp/internal/audit"
)

type Item struct {
	ID        string
	Attribute string
	FulfillOn string
}

type Auditor struct {
	mu    sync.Mutex
	items []Item
}

func NewAuditor() *Auditor { return &Auditor{} }

func (a *Auditor) Record(items []Item) {
	if len(items) == 0 {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.items = append(a.items, items...)
}

func (a *Auditor) Flush(logger *audit.Logger) error {
	a.mu.Lock()
	items := append([]Item(nil), a.items...)
	a.items = nil
	a.mu.Unlock()
	if logger == nil {

		return nil
	}
	for _, it := range items {
		if err := logger.Write("obligation", map[string]any{
			"id": it.ID, "attribute": it.Attribute, "fulfill_on": it.FulfillOn,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (a *Auditor) Pending() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.items)
}
