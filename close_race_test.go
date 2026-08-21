package abacpdp

import (
	"errors"
	"sync"
	"testing"
	"time"
)

// bagForEval 构造一个会命中文档演示策略的属性袋。
func bagForEval() AttrBag {
	return AttrBag{
		Subject:  map[string]any{"role": "admin", "id": "alice"},
		Resource: map[string]any{"owner": "alice"},
		Action:   map[string]any{"name": "read"},
	}
}

// TestEvaluateAfterClose 模拟策略服务滚动重启：
// 先 Close 了 PDP，健康检查还打了一次 Evaluate。
// 修复前 combine 跑完后写 nil cache 触发 panic；修复后应直接 ErrClosed。
func TestEvaluateAfterClose(t *testing.T) {
	p := New()
	if err := LoadSampleDocs(p); err != nil {
		t.Fatalf("load: %v", err)
	}

	// 正常一次，确保有命中策略的路径。
	if _, err := p.Evaluate(bagForEval()); err != nil {
		t.Fatalf("warm-up evaluate: %v", err)
	}

	if err := p.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	// 关闭后再判决：不应 panic，应返回 ErrClosed。
	dec, err := p.Evaluate(bagForEval())
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("after close: want ErrClosed, got dec=%+v err=%v", dec, err)
	}
}

// TestEvaluateConcurrentClose 复现真正的竞态：Evaluate 跑到写缓存前，
// 另一个 goroutine Close 把 cache 拆成 nil。修复前会 nil map panic。
func TestEvaluateConcurrentClose(t *testing.T) {
	const n = 200
	for run := 0; run < n; run++ {
		p := New()
		if err := LoadSampleDocs(p); err != nil {
			t.Fatalf("load: %v", err)
		}

		var wg sync.WaitGroup
		wg.Add(2)
		evaluateErr := make(chan error, 1)
		closeErr := make(chan error, 1)

		go func() {
			defer wg.Done()
			// 持续 evaluate，直到被 close 中断或返回错误。
			for {
				_, err := p.Evaluate(bagForEval())
				if err != nil {
					evaluateErr <- err
					return
				}
				if p.IsClosed() {
					evaluateErr <- nil
					return
				}
			}
		}()

		go func() {
			defer wg.Done()
			// 错开一点，争取正好踩在写缓存的时间窗里 Close。
			time.Sleep(time.Microsecond)
			closeErr <- p.Close()
		}()

		wg.Wait()
		close(evaluateErr)
		close(closeErr)

		// Close 应成功。
		if err := <-closeErr; err != nil {
			t.Fatalf("close: %v", err)
		}
		// evaluate 的结局只能是 ErrClosed（被关闭挡下）或 nil（在 close 前正常返回）。
		// 任何其它错误、尤其 panic，会让本测试直接失败。
		// 不要求 errors.Is(ErrClosed)：评测可能先于 close 正常完成。
		select {
		case err := <-evaluateErr:
			if err != nil && !errors.Is(err, ErrClosed) {
				t.Fatalf("run %d: evaluate after concurrent close: %v", run, err)
			}
		default:
		}
	}
}
