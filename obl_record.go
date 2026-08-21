package abacpdp

import "github.com/LYH2263/go-abacpdp/internal/obligation"

func (p *PDP) recordObligations(obls []Obligation) {
        if p.oblAudit == nil || len(obls) == 0 {
                return
        }
        items := make([]obligation.Item, len(obls))
        for i, o := range obls {
                items[i] = obligation.Item{ID: o.ID, Attribute: o.Attribute, FulfillOn: string(o.FulfillOn)}
        }
        p.oblAudit.Record(items)
}
