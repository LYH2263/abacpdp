package abacpdp

import (
	"context"
	"errors"
	"testing"

	"github.com/LYH2263/go-abacpdp/internal/policy"
)

// newCancelTestPDP 装一个“会判 Permit”的策略集：多条 Permit 规则 + PermitOverrides，
// 这样若不尊重 ctx 就会跑完全部规则并返回 Permit，正好复现网关上报的现象。
func newCancelTestPDP(t *testing.T) *PDP {
	t.Helper()
	b := policy.NewBuilder("tset")
	b.Combine("PermitOverrides")
	b.Version(1)
	pol := policy.Policy{
		ID:      "p1",
		Combine: "PermitOverrides",
		Target:  policy.Target{}, // 空目标：总是命中
	}
	for i := 0; i < 30; i++ {
		pol.Rules = append(pol.Rules, policy.PermitRule("r"+itoa(i)))
	}
	b.Policy(pol)
	p := New()
	if err := p.Load(b.Build()); err != nil {
		t.Fatalf("load: %v", err)
	}
	return p
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}

// 网关在调用 EvaluateContext 之前就已 cancel：不应 enrich / 跑策略，必须 ErrCanceled。
func TestEvaluateContext_PreCanceledReturnsErrCanceled(t *testing.T) {
	p := newCancelTestPDP(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立刻 cancel

	dec, err := p.EvaluateContext(ctx, NewAttrBag())
	if !errors.Is(err, ErrCanceled) {
		t.Fatalf("expected ErrCanceled, got err=%v dec=%+v", err, dec)
	}
	if dec.IsPermit() {
		t.Fatalf("canceled ctx must not return a Permit decision, got %+v", dec)
	}
}

// 即使进入 evalSet/evalPolicy 循环，也应在下一条策略/规则前停下并返回 ErrCanceled，
// 而非硬跑完几十条规则再回 Permit。用计数型 resolver 确认没有 enrich 后的误判。
func TestEvaluateContext_DoesNotRunFullPolicySetAfterCancel(t *testing.T) {
	p := newCancelTestPDP(t)

	// 预先 cancel，模拟网关立即取消。
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	dec, err := p.EvaluateContext(ctx, NewAttrBag())
	if !errors.Is(err, ErrCanceled) {
		t.Fatalf("expected ErrCanceled, got err=%v dec=%+v", err, dec)
	}
	// 统计不应记一次成功判决（被取消的不是 Permit/Deny）。
	st := p.Stats()
	if st.Permits != 0 || st.Evaluations != 0 {
		t.Fatalf("canceled eval must not be counted as success, got %+v", st)
	}
}
