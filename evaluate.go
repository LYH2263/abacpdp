package abacpdp

import (
	"context"
	"fmt"

	"github.com/LYH2263/go-abacpdp/internal/combine"
	"github.com/LYH2263/go-abacpdp/internal/match"
	"github.com/LYH2263/go-abacpdp/internal/obligation"
	"github.com/LYH2263/go-abacpdp/internal/policy"
)

// Evaluate 判决属性袋。
func (p *PDP) Evaluate(bag AttrBag) (Decision, error) {
	return p.EvaluateContext(context.Background(), bag)
}

// EvaluateContext 支持取消；取消后尽快返回 ErrCanceled。
func (p *PDP) EvaluateContext(ctx context.Context, bag AttrBag) (Decision, error) {
	if p.closed.Load() {
		return Decision{}, ErrClosed
	}
	if err := ctx.Err(); err != nil {
		return Decision{}, fmt.Errorf("%w: %v", ErrCanceled, err)
	}

	// 克隆入参，避免调用方事后改写 Subject 等污染缓存键。
	snap := bag.Clone()
	if p.resolver != nil {
		if err := p.resolver.Enrich(ctx, resolveBag(&snap)); err != nil {
			return Decision{}, err
		}
	}

	key := snap.Fingerprint()
	p.mu.RLock()
	if d, ok := p.cache[key]; ok {
		p.mu.RUnlock()
		p.cacheHits.Add(1)
		return d.Clone(), nil
	}
	set := p.set
	p.mu.RUnlock()

	if set == nil || len(set.Policies) == 0 {
		return Decision{}, ErrNoPolicy
	}

	dec, err := p.evalSet(ctx, set, snap)
	if err != nil {
		return Decision{}, err
	}

	p.mu.Lock()
	if !p.closed.Load() {
		if len(p.cache) >= p.cacheSize {
			p.cache = make(map[string]Decision)
		}
		p.cache[key] = dec.Clone()
	}
	p.mu.Unlock()

	p.evaluations.Add(1)
	p.metrics.IncEvaluations(1)
	if dec.IsPermit() {
		p.permits.Add(1)
		p.metrics.IncPermits(1)
	} else if dec.IsDeny() {
		p.denies.Add(1)
		p.metrics.IncDenies(1)
	}
	_ = p.auditor.Write("evaluate", map[string]any{"effect": string(dec.Effect), "policy": dec.PolicyID})
	p.recordObligations(dec.Obligations)
	return dec.Clone(), nil
}

func (p *PDP) evalSet(ctx context.Context, set *policy.Set, bag AttrBag) (Decision, error) {
	results := make([]combine.Result, 0, len(set.Policies))
	for _, pol := range set.Policies {
		if err := ctx.Err(); err != nil {
			return Decision{}, fmt.Errorf("%w: %v", ErrCanceled, err)
		}
		d, err := p.evalPolicy(ctx, &pol, bag)
		if err != nil {
			return Decision{}, err
		}
		results = append(results, toCombine(d, pol.ID))
	}
	merged, err := combine.Merge(set.Combine, results)
	if err != nil {

		return Decision{}, fmt.Errorf("unknown combine: %v", err)
	}
	return fromCombine(merged), nil
}

func (p *PDP) evalPolicy(ctx context.Context, pol *policy.Policy, bag AttrBag) (Decision, error) {
	if !match.Target(pol.Target, bagToMatch(bag)) {
		return Decision{Effect: EffectNotApplicable, PolicyID: pol.ID}, nil
	}
	results := make([]combine.Result, 0, len(pol.Rules))
	for i := range pol.Rules {
		if err := ctx.Err(); err != nil {
			return Decision{}, fmt.Errorf("%w: %v", ErrCanceled, err)
		}
		d := evalRule(&pol.Rules[i], bag)
		results = append(results, toCombine(d, pol.ID))
	}
	algo := pol.Combine
	if algo == "" {
		algo = "DenyOverrides"
	}
	merged, err := combine.Merge(algo, results)
	if err != nil {

		return Decision{}, fmt.Errorf("unknown combine: %v", err)
	}
	out := fromCombine(merged)
	out.PolicyID = pol.ID
	return out, nil
}

func evalRule(rule *policy.Rule, bag AttrBag) Decision {
	if !match.Target(rule.Target, bagToMatch(bag)) {
		return Decision{Effect: EffectNotApplicable}
	}
	eff := Effect(rule.Effect)
	if eff != EffectPermit && eff != EffectDeny {
		eff = EffectIndeterminate
	}
	obls := make([]Obligation, 0, len(rule.Obligations))
	for _, o := range rule.Obligations {
		if o.FulfillOn == "" || o.FulfillOn == string(eff) {
			obls = append(obls, Obligation{
				ID:         o.ID,
				Attribute:  o.Attribute,
				FulfillOn:  Effect(o.FulfillOn),
				Parameters: obligation.CopyParams(o.Parameters),
			})
		}
	}
	return Decision{Effect: eff, Obligations: obls, Matched: []string{rule.ID}}
}

func bagToMatch(bag AttrBag) match.Bag {
	return match.Bag{
		Subject:     bag.Subject,
		Resource:    bag.Resource,
		Action:      bag.Action,
		Environment: bag.Environment,
	}
}

func toCombine(d Decision, policyID string) combine.Result {
	obls := make([]combine.Obligation, len(d.Obligations))
	for i, o := range d.Obligations {
		obls[i] = combine.Obligation{
			ID: o.ID, Attribute: o.Attribute,
			FulfillOn: string(o.FulfillOn), Parameters: o.Parameters,
		}
	}
	pid := d.PolicyID
	if pid == "" {
		pid = policyID
	}
	return combine.Result{
		Effect: combine.Effect(d.Effect), Obligations: obls,
		Matched: d.Matched, PolicyID: pid,
	}
}

func fromCombine(m combine.Result) Decision {
	obls := make([]Obligation, len(m.Obligations))
	for i, o := range m.Obligations {
		obls[i] = Obligation{
			ID: o.ID, Attribute: o.Attribute,
			FulfillOn: Effect(o.FulfillOn), Parameters: obligation.CopyParams(o.Parameters),
		}
	}
	return Decision{
		Effect: Effect(m.Effect), Status: m.Status,
		Obligations: obls, Matched: append([]string(nil), m.Matched...),
		PolicyID: m.PolicyID,
	}
}
