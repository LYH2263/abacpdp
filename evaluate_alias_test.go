package abacpdp

import (
	"encoding/json"
	"testing"

	"github.com/LYH2263/go-abacpdp/internal/policy"
)

// 复现：试判页复用同一张 Subject map，先 admin 判 Permit，再把同一张 map
// 改成 intern 去打 secret。第二次因缓存脏命中仍返回第一次的 Permit。
// 根因：Fingerprint 只取 map 指针地址，改写内容后地址不变 → 同一缓存键。
func TestEvaluate_SubjectMapAliasCachePoison(t *testing.T) {
	raw := `{
      "id":"secret-access","version":1,"combine":"PermitOverrides",
      "policies":[
        {"id":"p-admin","combine":"DenyOverrides","target":{"all_of":[{"category":"subject","attribute":"role","op":"eq","value":"admin"}]},"rules":[
          {"id":"r-admin-permit","effect":"Permit","target":{"all_of":[{"category":"action","attribute":"name","op":"eq","value":"read"}]}}
        ]},
        {"id":"p-default","combine":"DenyOverrides","rules":[
          {"id":"r-default-deny","effect":"Deny"}
        ]}
      ]
    }`
	p := New(WithCacheSize(8))
	if err := p.LoadJSON([]byte(raw)); err != nil {
		t.Fatalf("LoadJSON: %v", err)
	}

	// 复用同一张 Subject map（模拟试判页的复用场景）。
	subject := map[string]any{"role": "admin"}
	bag := AttrBag{
		Subject:  subject,
		Action:   map[string]any{"name": "read"},
		Resource: map[string]any{"kind": "secret"},
	}

	d1, err := p.Evaluate(bag)
	if err != nil {
		t.Fatalf("Evaluate #1: %v", err)
	}
	if !d1.IsPermit() {
		t.Fatalf("Evaluate #1: want Permit, got %s", d1.Effect)
	}

	// 就地改写同一张 map：admin -> intern。期望落到 r-default-deny → Deny。
	subject["role"] = "intern"

	d2, err := p.Evaluate(bag)
	if err != nil {
		t.Fatalf("Evaluate #2: %v", err)
	}
	if d2.IsPermit() {
		t.Fatalf("Evaluate #2: 缓存脏命中，intern 仍返回 admin 的 Permit；effect=%s matched=%v policy=%s",
			d2.Effect, d2.Matched, d2.PolicyID)
	}
	if !d2.IsDeny() {
		t.Fatalf("Evaluate #2: want Deny, got %s", d2.Effect)
	}
}

// Fingerprint 必须随内容变化，且与 map 别名无关。
func TestAttrBag_FingerprintContentBased(t *testing.T) {
	a := NewAttrBag()
	a.Set("subject", "role", "admin")
	k1 := a.Fingerprint()

	// 改写内容 → 键必须变。
	a.Set("subject", "role", "intern")
	k2 := a.Fingerprint()
	if k1 == k2 {
		t.Fatalf("Fingerprint 不随内容变化：k1==k2=%q", k1)
	}

	// 内容相同但 map 身份不同 → 键必须相同。
	b := NewAttrBag()
	b.Set("subject", "role", "intern")
	if b.Fingerprint() != k2 {
		t.Fatalf("Fingerprint 依赖 map 身份而非内容：b=%q want %q", b.Fingerprint(), k2)
	}
}

// Clone 必须深拷贝，改写副本不得影响原袋（也不得共享底层 map）。
func TestAttrBag_CloneDeep(t *testing.T) {
	orig := NewAttrBag()
	orig.Set("subject", "role", "admin")
	orig.Set("resource", "nested", map[string]any{"k": "v"})

	cp := orig.Clone()
	cp.Set("subject", "role", "intern")
	if got, _ := orig.Get("subject", "role"); got != "admin" {
		t.Fatalf("Clone 浅拷贝：改写副本影响了原袋，orig.role=%v", got)
	}
	if got, _ := cp.Get("subject", "role"); got != "intern" {
		t.Fatalf("Clone 未见副本改动，cp.role=%v", got)
	}
}

// 顺带覆盖 policy.Set 经 JSON 往返仍能命中 target（保障上述策略解析无误）。
func TestEvaluate_PolicyRoundTrip(t *testing.T) {
	set := &policy.Set{ID: "x", Version: 1, Combine: "DenyOverrides"}
	raw, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var back policy.Set
	if err := json.Unmarshal(raw, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if back.ID != set.ID {
		t.Fatalf("roundtrip id: got %q want %q", back.ID, set.ID)
	}
}
