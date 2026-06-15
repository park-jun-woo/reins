//ff:func feature=cli type=command control=sequence level=error
//ff:what TestRunLoopItem — runLoopItem 직접 호출 검증. 통과 페이로드는 1회로 PASS 잠금, "bad"→"good" 시퀀스는 FAIL 피드백을 되먹여 재시도 후 PASS로 수렴(게이트만 잠금).

package cli

import (
	"io"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestRunLoopItem drives the per-item generate→gate→retry loop directly: a passing
// payload locks PASS in one shot, while a "bad"→"good" backend converges to PASS on
// retry (only the gate locks PASS).
func TestRunLoopItem(t *testing.T) {
	run := func(t *testing.T, replies []string) (*quest.Item, int) {
		t.Helper()
		dir := t.TempDir()
		s := quest.New()
		it := &quest.Item{Key: "a", State: quest.TODO}
		s.Items = append(s.Items, it)

		calls := 0
		backend := llm.CallFunc(func(system, user string) (string, error) {
			r := replies[calls]
			calls++
			return r, nil
		})
		err := runLoopItem(stubDef{}, &LoopOptions{}, backend, s, it,
			dir+"/out.jsonl", dir+"/session.json", io.Discard)
		if err != nil {
			t.Fatalf("runLoopItem: %v", err)
		}
		return it, calls
	}

	t.Run("pass in one shot", func(t *testing.T) {
		it, calls := run(t, []string{"good"})
		if it.State != quest.PASS {
			t.Fatalf("state = %q, want PASS", it.State)
		}
		if calls != 1 {
			t.Fatalf("backend called %d times, want 1", calls)
		}
	})

	t.Run("fail then converge", func(t *testing.T) {
		it, calls := run(t, []string{"bad", "good"})
		if it.State != quest.PASS {
			t.Fatalf("state = %q, want PASS after retry", it.State)
		}
		if calls != 2 {
			t.Fatalf("backend called %d times, want 2 (FAIL then PASS)", calls)
		}
	})

	t.Run("backend error demotes to item FAIL and does not abort", func(t *testing.T) {
		dir := t.TempDir()
		s := quest.New()
		it := &quest.Item{Key: "a", State: quest.TODO}
		s.Items = append(s.Items, it)
		backend := llm.CallFunc(func(system, user string) (string, error) { return "", errBackend })
		err := runLoopItem(stubDef{}, &LoopOptions{}, backend, s, it,
			dir+"/out.jsonl", dir+"/session.json", io.Discard)
		if err != nil {
			t.Fatalf("err = %v, want nil (generation error is a demoted item FAIL, not a run abort)", err)
		}
		if it.Tries != quest.MaxTries {
			t.Fatalf("Tries = %d, want %d (each generation error consumes a try)", it.Tries, quest.MaxTries)
		}
		if it.State != quest.DONE {
			t.Fatalf("state = %q, want DONE (MaxTries exhausted locks the residual)", it.State)
		}
	})

	t.Run("render error propagates", func(t *testing.T) {
		dir := t.TempDir()
		s := quest.New()
		it := &quest.Item{Key: "a", State: quest.TODO}
		s.Items = append(s.Items, it)
		backend := llm.CallFunc(func(system, user string) (string, error) { return "good", nil })
		err := runLoopItem(errDef{renderErr: true}, &LoopOptions{}, backend, s, it,
			dir+"/out.jsonl", dir+"/session.json", io.Discard)
		if err == nil {
			t.Fatalf("err = nil, want render error")
		}
	})

	t.Run("evaluate error propagates", func(t *testing.T) {
		dir := t.TempDir()
		s := quest.New()
		it := &quest.Item{Key: "a", State: quest.TODO}
		s.Items = append(s.Items, it)
		backend := llm.CallFunc(func(system, user string) (string, error) { return "good", nil })
		err := runLoopItem(errDef{prepareErr: true}, &LoopOptions{}, backend, s, it,
			dir+"/out.jsonl", dir+"/session.json", io.Discard)
		if err == nil {
			t.Fatalf("err = nil, want evaluateAndApply (prepare) error")
		}
	})
}
