//ff:func feature=cli type=helper control=sequence level=error
//ff:what TestEvaluateAndApplyGraphFailRootCause — Evaluator(graph) 분기가 defeat 그래프에서 verdict를 채우고 FAIL이 RootCause를 담으며 아이템을 잠그지 않는지 검증.

package cli

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestEvaluateAndApplyGraphFailRootCause: the Evaluator (graph) branch fills the
// verdict from the defeat graph, and a FAIL carries a RootCause.
func TestEvaluateAndApplyGraphFailRootCause(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")

	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	v, err := evaluateAndApply(graphDef{}, s, it, []byte("bad"), out, session)
	if err != nil {
		t.Fatalf("evaluateAndApply error: %v", err)
	}
	if v.Outcome != quest.OutFail {
		t.Fatalf("outcome = %v, want FAIL", v.Outcome)
	}
	if v.RootCause == "" {
		t.Fatalf("graph FAIL verdict missing RootCause: %+v", v)
	}
	if it.State != quest.TODO {
		t.Fatalf("state = %v, want TODO (FAIL does not lock)", it.State)
	}
}
