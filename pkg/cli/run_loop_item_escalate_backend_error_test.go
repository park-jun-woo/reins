//ff:func feature=cli type=command control=sequence level=error
//ff:what TestRunLoopItemEscalateBackendError — 에스컬레이션 상호작용. primary 의 능력-한계 FAIL 로 승격된 뒤 강한 Escalate backend 가 에러나도, 그 에러는 abort 가 아니라 강등 FAIL 로 Tries 를 소진해 DONE 잠금(승격 상태 유지·nil 반환)됨을 단언(Phase011 에스컬레이션과 충돌 없음).

package cli

import (
	"io"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestRunLoopItemEscalateBackendError: the primary FAILs once with the
// capability-bound RootCause "not-bad", latching the item onto the stronger Escalate
// backend. That backend then errors — but the error is demoted (not propagated), so
// the remaining tries are spent on Escalate and the item locks DONE. runLoopItem
// returns nil; the next item would proceed.
func TestRunLoopItemEscalateBackendError(t *testing.T) {
	dir := t.TempDir()
	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	primaryCalls, escalateCalls := 0, 0
	primary := llm.CallFunc(func(system, user string) (string, error) {
		primaryCalls++
		return "bad", nil // FAIL with RootCause "not-bad" → escalate
	})
	escalate := llm.CallFunc(func(system, user string) (string, error) {
		escalateCalls++
		return "", errBackend // the stronger backend now errors
	})
	opts := &LoopOptions{Escalate: escalate, EscalateOn: []string{"not-bad"}}

	err := runLoopItem(stubDef{}, opts, primary, s, it,
		dir+"/out.jsonl", dir+"/session.json", io.Discard)
	if err != nil {
		t.Fatalf("err = %v, want nil (escalated backend error is demoted, not propagated)", err)
	}
	if primaryCalls != 1 {
		t.Fatalf("primary called %d times, want 1 (one FAIL then escalate)", primaryCalls)
	}
	// One try spent on the primary FAIL, the rest on the (erroring) Escalate backend.
	if escalateCalls != quest.MaxTries-1 {
		t.Fatalf("escalate called %d times, want %d", escalateCalls, quest.MaxTries-1)
	}
	if it.Tries != quest.MaxTries {
		t.Fatalf("Tries = %d, want %d", it.Tries, quest.MaxTries)
	}
	if it.State != quest.DONE {
		t.Fatalf("state = %q, want DONE (exhausted on escalated backend)", it.State)
	}
}
