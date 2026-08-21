package resolve

import (
        "context"
        "fmt"
        "sync"
        "time"
)

type Bag struct {
        Subject     map[string]any
        Resource    map[string]any
        Action      map[string]any
        Environment map[string]any
}

type Resolver interface {
        Enrich(ctx context.Context, bag *Bag) error
}

type Static struct{ Env map[string]any }

func (s Static) Enrich(ctx context.Context, bag *Bag) error {
        if err := ctx.Err(); err != nil {
                return err
        }
        if bag.Environment == nil {
                bag.Environment = map[string]any{}
        }
        for k, v := range s.Env {
                bag.Environment[k] = v
        }
        return nil
}

// Delayed 模拟远程 Wait；必须尊重 ctx。
type Delayed struct {
        Delay time.Duration
        Key   string
        Value any
}

func (d Delayed) Enrich(ctx context.Context, bag *Bag) error {
        timer := time.NewTimer(d.Delay)
        defer timer.Stop()
        select {
        case <-ctx.Done():
                return fmt.Errorf("resolve: %w", ctx.Err())
        case <-timer.C:
                if bag.Environment == nil {
                        bag.Environment = map[string]any{}
                }
                key := d.Key
                if key == "" {
                        key = "remote"
                }
                bag.Environment[key] = d.Value
                return nil
        }
}

type Cache struct {
        mu    sync.RWMutex
        inner Resolver
        data  map[string]any
}

func NewCache(inner Resolver) *Cache {
        return &Cache{inner: inner, data: make(map[string]any)}
}

func (c *Cache) Enrich(ctx context.Context, bag *Bag) error {
        if c.inner != nil {
                return c.inner.Enrich(ctx, bag)
        }
        return ctx.Err()
}
