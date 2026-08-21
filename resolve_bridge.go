package abacpdp

import "github.com/LYH2263/go-abacpdp/internal/resolve"

func resolveBag(b *AttrBag) *resolve.Bag {
        return &resolve.Bag{
                Subject: b.Subject, Resource: b.Resource,
                Action: b.Action, Environment: b.Environment,
        }
}
