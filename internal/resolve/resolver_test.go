package resolve

import (
	"context"
	"testing"
	"time"
)

// 取消必须尽早返回，而不是 Sleep 完整 Delay。
func TestDelayed_Enrich_CancelShortCircuits(t *testing.T) {
	d := Delayed{Delay: 2 * time.Second, Key: "remote", Value: "v"}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := d.Enrich(ctx, &Bag{})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("期望因取消返回错误，得到 nil")
	}
	if elapsed >= 500*time.Millisecond {
		t.Fatalf("取消应尽快返回，实际耗时 %v", elapsed)
	}
}

// 未取消时应正常写入属性。
func TestDelayed_Enrich_AppliesValue(t *testing.T) {
	d := Delayed{Delay: 5 * time.Millisecond, Key: "remote", Value: "v"}
	bag := &Bag{}
	if err := d.Enrich(context.Background(), bag); err != nil {
		t.Fatalf("未取消不应报错: %v", err)
	}
	if got := bag.Environment["remote"]; got != "v" {
		t.Fatalf("Environment[remote] = %v, 期望 v", got)
	}
}

// 已取消的 ctx 应立即返回，不等待 Delay。
func TestDelayed_Enrich_AlreadyCanceled(t *testing.T) {
	d := Delayed{Delay: 2 * time.Second, Key: "remote", Value: "v"}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	err := d.Enrich(ctx, &Bag{})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("已取消的 ctx 应返回错误")
	}
	if elapsed >= 100*time.Millisecond {
		t.Fatalf("已取消应立即返回，实际耗时 %v", elapsed)
	}
}
