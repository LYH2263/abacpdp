package abacpdp

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/LYH2263/go-abacpdp/internal/policy"
)

// nullParamsPermitRule 返回一条 Permit 规则，其义务 Parameters 在 JSON 中为 null。
// 复现：策略里有一条 Permit 义务 Parameters 压根没配（json 里是 null），
// Evaluate 第一次直接盖 rule 戳时 panic。
func nullParamsPermitRule(t *testing.T) *PDP {
	t.Helper()
	set := &policy.Set{
		ID:      "t-null-params",
		Version: 1,
		Combine: "DenyOverrides",
		Policies: []policy.Policy{{
			ID: "p",
			Rules: []policy.Rule{
				{
					ID:     "allow-with-null-params",
					Effect: "Permit",
					Target: policy.Target{AllOf: []policy.Match{
						policy.Eq("subject", "role", "admin"),
					}},
					Obligations: []policy.Obligation{{
						ID:        "audit-read",
						FulfillOn: "Permit",
						// Parameters 故意留空，对应 JSON 中的 null
					}},
				},
			},
		}},
	}
	// 经 JSON 往返，确保 Parameters 字段确实为 JSON null（而非 nil-but-omitted）。
	raw, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(raw), `"parameters":null`) {
		// 字段带 omitempty + nil 时不会被序列化为 null；强制注入 null 以贴近报告。
		raw = []byte(strings.Replace(string(raw),
			`"fulfill_on":"Permit","attribute":""}`,
			`"fulfill_on":"Permit","parameters":null}`, 1))
	}
	p := New()
	if err := p.LoadJSON(raw); err != nil {
		t.Fatalf("load: %v", err)
	}
	return p
}

// TestEvaluate_NullObligationParameters_NoPanic 验证义务 Parameters 为 null 时
// Evaluate 不再 panic，且盖上的 rule 戳正确写入。
func TestEvaluate_NullObligationParameters_NoPanic(t *testing.T) {
	p := nullParamsPermitRule(t)
	defer p.Close()

	bag := NewAttrBag()
	bag.Set("subject", "role", "admin")

	dec, err := p.Evaluate(bag)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if dec.Effect != EffectPermit {
		t.Fatalf("effect = %q, want Permit", dec.Effect)
	}
	if len(dec.Obligations) != 1 {
		t.Fatalf("obligations = %d, want 1", len(dec.Obligations))
	}
	o := dec.Obligations[0]
	if o.ID != "audit-read" {
		t.Errorf("obligation id = %q, want audit-read", o.ID)
	}
	if o.Parameters["rule"] != "allow-with-null-params" {
		t.Errorf("rule stamp = %q, want allow-with-null-params", o.Parameters["rule"])
	}
	if _, ok := o.Parameters["level"]; ok {
		t.Errorf("unexpected level param present")
	}
}

// TestEvaluate_FailureDoesNotPoisonCache 验证：求值失败（含 panic 后恢复路径）
// 绝不污染 decision cache，同一属性袋再次 Evaluate 不会命中半写入的空义务 Permit。
func TestEvaluate_FailureDoesNotPoisonCache(t *testing.T) {
	p := nullParamsPermitRule(t)
	defer p.Close()

	bag := NewAttrBag()
	bag.Set("subject", "role", "admin")

	// 第一次：历史上会 panic（Bug2）且缓存被半写入（Bug1）；修复后应正常返回 Permit。
	dec1, err := p.Evaluate(bag)
	if err != nil {
		t.Fatalf("first evaluate: %v", err)
	}
	if dec1.Effect != EffectPermit || len(dec1.Obligations) != 1 {
		t.Fatalf("first: effect=%q obligations=%d", dec1.Effect, len(dec1.Obligations))
	}

	// 第二次：命中缓存，必须返回完整义务，而非空义务的 Permit。
	dec2, err := p.Evaluate(bag)
	if err != nil {
		t.Fatalf("second evaluate: %v", err)
	}
	if dec2.Effect != EffectPermit {
		t.Fatalf("second: effect=%q, want Permit", dec2.Effect)
	}
	if len(dec2.Obligations) != 1 {
		t.Fatalf("second: obligations=%d, want 1 — cache returned half-written empty-obligation Permit", len(dec2.Obligations))
	}
	if dec2.Obligations[0].Parameters["rule"] != "allow-with-null-params" {
		t.Errorf("second: rule stamp lost from cached decision: %q", dec2.Obligations[0].Parameters["rule"])
	}
}

// TestEvaluate_NoPolicy_ErrorPathLeavesCacheClean 验证错误路径不落缓存：
// 无策略时返回 ErrNoPolicy，随后即使装载了策略，旧键也不应残留 Permit。
func TestEvaluate_NoPolicy_ErrorPathLeavesCacheClean(t *testing.T) {
	p := New()
	defer p.Close()

	bag := NewAttrBag()
	bag.Set("subject", "role", "admin")

	if _, err := p.Evaluate(bag); err != ErrNoPolicy {
		t.Fatalf("want ErrNoPolicy, got %v", err)
	}

	// 装载一条真实 Permit 策略。
	if err := p.LoadJSON([]byte(`{
		"id":"late","version":1,"combine":"DenyOverrides",
		"policies":[{"id":"p","rules":[
			{"id":"r","effect":"Permit","target":{"all_of":[{"category":"subject","attribute":"role","op":"eq","value":"admin"}]},"obligations":[{"id":"ob","fulfill_on":"Permit","parameters":{"level":"info"}}]}
		]}]
	}`)); err != nil {
		t.Fatalf("load: %v", err)
	}

	dec, err := p.Evaluate(bag)
	if err != nil {
		t.Fatalf("evaluate after load: %v", err)
	}
	if dec.Effect != EffectPermit || len(dec.Obligations) != 1 {
		t.Fatalf("after load: effect=%q obligations=%d", dec.Effect, len(dec.Obligations))
	}
}
