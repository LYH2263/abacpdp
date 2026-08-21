package obligation

type Plan struct{ Items []Item }

func BuildPlan(items []Item) Plan {
        out := Plan{Items: make([]Item, len(items))}
        copy(out.Items, items)
        return out
}

func (p Plan) IDs() []string {
        ids := make([]string, 0, len(p.Items))
        for _, it := range p.Items {
                ids = append(ids, it.ID)
        }
        return ids
}
