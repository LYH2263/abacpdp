package abacpdp_test

import (
	"errors"
	"testing"

	abacpdp "github.com/LYH2263/go-abacpdp"
)

func TestBug05_UnknownCombineWraps(t *testing.T) {
	p := abacpdp.New()
	defer p.Close()
	raw := []byte("{\"id\":\"x\",\"combine\":\"TotallyUnknown\",\"policies\":[{\"id\":\"p\",\"rules\":[{\"id\":\"r\",\"effect\":\"Permit\"}]}]}")
	err := p.LoadJSON(raw)
	if err != nil {
		// 干净/修复：装载期应拒绝未知 combine
		if !errors.Is(err, abacpdp.ErrInvalidPolicy) {
			t.Fatalf("want ErrInvalidPolicy on Load, got %v", err)
		}
		return
	}
	// 问题版：校验被绕过，Evaluate 合并失败须挂上 ErrUnknownCombine
	_, err = p.Evaluate(abacpdp.NewAttrBag())
	if err == nil {
		t.Fatal("expected unknown combine error")
	}
	if !errors.Is(err, abacpdp.ErrUnknownCombine) {
		t.Fatalf("want errors.Is ErrUnknownCombine, got %v", err)
	}
}
