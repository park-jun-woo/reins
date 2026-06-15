//ff:func feature=cli type=command control=sequence level=error
//ff:what TestRunLoopItemBackendErrorSavePropagates — 에러 분류 경계 단언. 생성 오류는 강등하지만 그 강등 verdict 의 영속화(applyVerdict 의 Save) 실패는 여전히 fail-fast 로 전파됨(정합성 위협은 침묵 금지)을 쓰기 불가 세션 경로로 검증.

package cli

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestRunLoopItemBackendErrorSavePropagates: a generation error is demoted, but if
// persisting the demoted verdict fails (unwritable session path), runLoopItem still
// returns that error — persistence failures remain fatal even on the demotion path.
func TestRunLoopItemBackendErrorSavePropagates(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.jsonl")
	// A session path under a regular file (not a dir) makes Save fail.
	badParent := filepath.Join(dir, "afile")
	if err := os.WriteFile(badParent, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	session := filepath.Join(badParent, "nested", "session.json")

	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	backend := llm.CallFunc(func(system, user string) (string, error) { return "", errBackend })
	err := runLoopItem(stubDef{}, &LoopOptions{}, backend, s, it,
		out, session, io.Discard)
	if err == nil {
		t.Fatal("err = nil, want save error to propagate even on the backend-error path")
	}
}
