//ff:func feature=cli type=command control=sequence level=error
//ff:what TestGeneratePayloadBackendError — generatePayload 의 backend-error 강등 경로. Complete 에러 시 backendErrorVerdict 로 합성→applyVerdict 로 래칫(Tries++)·영속화 후 handled=true(호출부가 루프 continue)·err=nil·raw="" 임을 단언. raw 가 비고 handled 가 참인 분기 시맨틱과 합성 FAIL 의 backend-error 어트리뷰션을 관측.

package cli

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestGeneratePayloadBackendError: a backend.Complete error is demoted, not
// propagated. applyVerdict ratchets the synthetic backend-error verdict (Tries++),
// generatePayload returns handled=true (so the caller continues the loop) with err=nil
// and empty raw. The demoted FAIL is observable under the reserved backend-error rule
// carrying the original error text.
func TestGeneratePayloadBackendError(t *testing.T) {
	dir := t.TempDir()
	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	backend := llm.CallFunc(func(system, user string) (string, error) {
		return "ignored-on-error", errBackend
	})

	raw, handled, err := generatePayload(stubDef{}, &LoopOptions{}, backend, "", "",
		s, it, filepath.Join(dir, "out.jsonl"), filepath.Join(dir, "session.json"), io.Discard)
	if err != nil {
		t.Fatalf("err = %v, want nil (generation error is demoted, not propagated)", err)
	}
	if !handled {
		t.Fatal("handled = false, want true (caller must continue, not gate-evaluate)")
	}
	if raw != "" {
		t.Fatalf("raw = %q, want empty on the demotion path", raw)
	}

	// The demotion ratcheted one try.
	if it.Tries != 1 {
		t.Fatalf("Tries = %d, want 1 (applyVerdict ratchets the demoted FAIL)", it.Tries)
	}

	// Observability: the synthetic FAIL surfaces under "backend-error" with the
	// original error text.
	if len(it.Log) != 1 {
		t.Fatalf("Log has %d attempts, want 1", len(it.Log))
	}
	last := it.Log[len(it.Log)-1]
	if !strings.Contains(last.Reason, backendErrorRule) {
		t.Fatalf("attempt reason = %q, want it to mention %q", last.Reason, backendErrorRule)
	}
	if !strings.Contains(last.Reason, errBackend.Error()) {
		t.Fatalf("attempt reason = %q, want it to carry the original error %q", last.Reason, errBackend.Error())
	}
}
