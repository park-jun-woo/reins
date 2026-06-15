//ff:func feature=cli type=command control=sequence level=error
//ff:what TestLoopBackendErrorContinues — BUG-002 정면 회귀. 1번 아이템 backend 가 항상 에러여도 그 실패가 2번 아이템을 막지 않음(1번은 MaxTries 후 DONE 잠금, 2번은 게이트 PASS)을 loop 명령 끝까지로 단언.

package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestLoopBackendErrorContinues: with two TODO items, item "a"'s backend always
// errors while item "b"'s succeeds. The run must not abort on "a" — it locks "a"
// DONE (residual) after MaxTries and still drives "b" through the gate to PASS. The
// backend keys off the rendered prompt (stubDef renders "render:<key>").
func TestLoopBackendErrorContinues(t *testing.T) {
	dir := t.TempDir()
	session := dir + "/session.json"
	out := dir + "/out.jsonl"

	backend := llm.CallFunc(func(system, user string) (string, error) {
		if strings.Contains(user, "render:a") {
			return "", errBackend // item "a" generation always fails
		}
		return "good", nil // item "b" converges
	})
	opts := Options{Loop: &LoopOptions{LLM: backend}}

	if _, err := newLoopRoot(t, opts, session, out, "scan", "a", "b"); err != nil {
		t.Fatalf("scan: %v", err)
	}
	stdout, err := newLoopRoot(t, opts, session, out, "loop")
	if err != nil {
		t.Fatalf("loop = %v, want nil (item 'a' failure must not abort the run)", err)
	}
	if !strings.Contains(stdout, "processed 2 item(s)") {
		t.Fatalf("loop output = %q, want both items processed", stdout)
	}

	s, err := loadSession(session)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	a, err := s.Find("a")
	if err != nil {
		t.Fatalf("find a: %v", err)
	}
	if a.State != quest.DONE || a.Tries != quest.MaxTries {
		t.Fatalf("item a = %+v, want DONE residual at MaxTries", a)
	}
	b, err := s.Find("b")
	if err != nil {
		t.Fatalf("find b: %v", err)
	}
	if b.State != quest.PASS {
		t.Fatalf("item b = %+v, want PASS (unblocked by a's failure)", b)
	}
}
