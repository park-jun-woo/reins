//ff:func feature=cli type=command control=sequence level=error
//ff:what TestGeneratePayloadRenderError — generatePayload 의 def.Render 에러 fail-fast 분기. Render 가 에러나면 backend.Complete 호출 전에 그 에러를 그대로 전파(handled=false·raw="")함을 render-에러 def 더블로 단언. 강등은 backend 인프라 에러 한정이고 프롬프트 렌더 실패는 즉시 중단임을 경계짓는다.

package cli

import (
	"io"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestGeneratePayloadRenderError: a def.Render error propagates immediately —
// generatePayload returns that error (handled=false, raw="") without ever calling
// backend.Complete. Demotion is reserved for backend infra errors; a prompt-render
// failure is a hard fail-fast.
func TestGeneratePayloadRenderError(t *testing.T) {
	dir := t.TempDir()
	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	called := false
	backend := llm.CallFunc(func(system, user string) (string, error) {
		called = true
		return "should-not-run", nil
	})

	raw, handled, err := generatePayload(errDef{renderErr: true}, &LoopOptions{}, backend, "", "",
		s, it, filepath.Join(dir, "out.jsonl"), filepath.Join(dir, "session.json"), io.Discard)
	if err == nil {
		t.Fatal("err = nil, want def.Render error to propagate")
	}
	if handled {
		t.Fatal("handled = true, want false (render error is propagated, not demoted)")
	}
	if raw != "" {
		t.Fatalf("raw = %q, want empty", raw)
	}
	if called {
		t.Fatal("backend.Complete was called, want it skipped after a Render error")
	}
}
