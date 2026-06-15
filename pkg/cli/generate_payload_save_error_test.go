//ff:func feature=cli type=command control=sequence level=error
//ff:what TestGeneratePayloadSaveError — generatePayload 의 강등 경로 fail-fast 경계. backend-error 강등 verdict 의 영속화(applyVerdict 의 Save) 실패는 강등에 묻히지 않고 err 로 전파됨(정합성 위협 침묵 금지)을 쓰기 불가 세션 경로로 단언. raw="" 이고 handled=false(전파는 handled 신호가 아님).

package cli

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestGeneratePayloadSaveError: on the backend-error demotion path, if persisting the
// demoted verdict fails (unwritable session path), generatePayload returns that error
// rather than swallowing it — persistence failures stay fatal even when the generation
// error itself is demoted. raw is empty and handled is false (propagation is not a
// continue signal).
func TestGeneratePayloadSaveError(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.jsonl")
	// A session path nested under a regular file (not a dir) makes Save fail.
	badParent := filepath.Join(dir, "afile")
	if err := os.WriteFile(badParent, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	session := filepath.Join(badParent, "nested", "session.json")

	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	backend := llm.CallFunc(func(system, user string) (string, error) { return "", errBackend })
	raw, handled, err := generatePayload(stubDef{}, &LoopOptions{}, backend, "", "",
		s, it, out, session, io.Discard)
	if err == nil {
		t.Fatal("err = nil, want save error to propagate on the backend-error path")
	}
	if handled {
		t.Fatal("handled = true, want false (a propagated error is not a continue signal)")
	}
	if raw != "" {
		t.Fatalf("raw = %q, want empty", raw)
	}
}
