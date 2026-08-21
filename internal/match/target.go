package match

import "github.com/LYH2263/go-abacpdp/internal/policy"

func Target(t policy.Target, bag Bag) bool {
        if t.Empty() {
                return true
        }
        for _, m := range t.AllOf {
                if !Eval(m, bag) {
                        return false
                }
        }
        if len(t.AnyOf) == 0 {
                return true
        }
        for _, m := range t.AnyOf {
                if Eval(m, bag) {
                        return true
                }
        }
        return false
}
