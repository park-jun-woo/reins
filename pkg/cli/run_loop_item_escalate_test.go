//ff:func feature=cli type=command control=sequence level=error
//ff:what TestRunLoopItemEscalate — runLoopItem 의 능력-한계 승격 검증. primary 가 EscalateOn 의 RootCause 로 FAIL 하면 다음 시도부터 Escalate backend 로 래치 승격(그게 PASS 로 수렴); RootCause 가 EscalateOn 밖이면 승격 안 하고 primary 만 계속 호출. Escalate 백엔드는 잔여에만 불린다(비용).

package cli

import (
	"io"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// stubDef's FAIL rule has ID "not-bad" (gate.Evaluate sets RootCause to it), so a
// "bad" submission FAILs with RootCause "not-bad" — the trigger used below.
func TestRunLoopItemEscalate(t *testing.T) {
	t.Run("capability-bound FAIL promotes to Escalate backend", func(t *testing.T) {
		dir := t.TempDir()
		s := quest.New()
		it := &quest.Item{Key: "a", State: quest.TODO}
		s.Items = append(s.Items, it)

		primaryCalls, escalateCalls := 0, 0
		primary := llm.CallFunc(func(system, user string) (string, error) {
			primaryCalls++
			return "bad", nil // FAIL with RootCause "not-bad"
		})
		escalate := llm.CallFunc(func(system, user string) (string, error) {
			escalateCalls++
			return "good", nil // the stronger backend converges
		})
		opts := &LoopOptions{Escalate: escalate, EscalateOn: []string{"not-bad"}}

		if err := runLoopItem(stubDef{}, opts, primary, s, it,
			dir+"/out.jsonl", dir+"/session.json", io.Discard); err != nil {
			t.Fatalf("runLoopItem: %v", err)
		}
		if it.State != quest.PASS {
			t.Fatalf("state = %q, want PASS after escalation", it.State)
		}
		if primaryCalls != 1 {
			t.Fatalf("primary called %d times, want 1 (one FAIL then escalate)", primaryCalls)
		}
		if escalateCalls != 1 {
			t.Fatalf("escalate called %d times, want 1 (the residual only)", escalateCalls)
		}
	})

	t.Run("RootCause outside EscalateOn never escalates", func(t *testing.T) {
		dir := t.TempDir()
		s := quest.New()
		it := &quest.Item{Key: "a", State: quest.TODO}
		s.Items = append(s.Items, it)

		primaryCalls, escalateCalls := 0, 0
		primary := llm.CallFunc(func(system, user string) (string, error) {
			primaryCalls++
			return "bad", nil // always FAIL on the primary
		})
		escalate := llm.CallFunc(func(system, user string) (string, error) {
			escalateCalls++
			return "good", nil
		})
		// EscalateOn does not include "not-bad", so a format/other failure stays on
		// the cheap primary and exhausts MaxTries → DONE.
		opts := &LoopOptions{Escalate: escalate, EscalateOn: []string{"some-other-rule"}}

		if err := runLoopItem(stubDef{}, opts, primary, s, it,
			dir+"/out.jsonl", dir+"/session.json", io.Discard); err != nil {
			t.Fatalf("runLoopItem: %v", err)
		}
		if escalateCalls != 0 {
			t.Fatalf("escalate called %d times, want 0 (RootCause not in EscalateOn)", escalateCalls)
		}
		if primaryCalls != quest.MaxTries {
			t.Fatalf("primary called %d times, want MaxTries=%d", primaryCalls, quest.MaxTries)
		}
		if it.State != quest.DONE {
			t.Fatalf("state = %q, want DONE (exhausted on primary)", it.State)
		}
	})
}
