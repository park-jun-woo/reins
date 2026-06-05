//ff:func feature=cli type=helper control=sequence
//ff:what submit의 Evaluator 분기 + 후방호환 증명 — gate.Evaluator를 구현한 graphDef는 submit이 ev.Evaluate(그래프) 경로를 타 "bad"→FAIL·정상→PASS를 내고(Rules() 카탈로그 직접 평가가 아님을 그래프 산물로 확인), Evaluator 미구현 stubDef는 기존 gate.Evaluate(Rules) 경로 그대로 동작(무회귀)함을 검증한다.

package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestSubmitUsesEvaluatorWhenDefImplementsIt: a graph-backed Definition routes
// through the gate.Evaluator path (its Verdict comes from pkg/graph), while a
// Rules-only Definition keeps the legacy gate.Evaluate(Rules) path (back-compat).
func TestSubmitUsesEvaluatorWhenDefImplementsIt(t *testing.T) {
	t.Run("evaluator fail", func(t *testing.T) {
		dir := t.TempDir()
		session := filepath.Join(dir, "session.json")
		out := filepath.Join(dir, "out.jsonl")
		runCmd(t, graphDef{}, session, out, "", "scan", "a")
		got := runCmd(t, graphDef{}, session, out, "bad", "submit", "--key", "a")
		if !strings.Contains(got, "a -> FAIL") {
			t.Fatalf("evaluator FAIL not reached: %q", got)
		}
	})

	t.Run("evaluator pass", func(t *testing.T) {
		dir := t.TempDir()
		session := filepath.Join(dir, "session.json")
		out := filepath.Join(dir, "out.jsonl")
		runCmd(t, graphDef{}, session, out, "", "scan", "a")
		got := runCmd(t, graphDef{}, session, out, "good", "submit", "--key", "a")
		if !strings.Contains(got, "a -> PASS") {
			t.Fatalf("evaluator PASS not reached: %q", got)
		}
	})

	t.Run("rules-only back-compat", func(t *testing.T) {
		dir := t.TempDir()
		session := filepath.Join(dir, "session.json")
		out := filepath.Join(dir, "out.jsonl")
		runCmd(t, stubDef{}, session, out, "", "scan", "a")
		got := runCmd(t, stubDef{}, session, out, "bad", "submit", "--key", "a")
		if !strings.Contains(got, "a -> FAIL") {
			t.Fatalf("rules-only FAIL not reached: %q", got)
		}
	})
}
