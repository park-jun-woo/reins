//ff:func feature=cli type=command control=sequence level=error
//ff:what TestRunLoopItemBackendError — runLoopItem 의 생성 오류 강등 검증(BUG-002 핵심). 항상 에러나는 backend 가 아이템을 abort 아닌 재시도 가능한 FAIL 로 강등해 MaxTries 후 DONE 잠금(종료 보장·nil 반환), 그 합성 FAIL 이 verdict.Facts 에 backend-error + 원문 에러로 관측됨을 단언.

package cli

import (
	"io"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestRunLoopItemBackendError: a generation error is demoted to a retryable item
// FAIL — runLoopItem returns nil (no abort), each error consumes a try, and the item
// locks DONE at MaxTries (termination guaranteed). The synthetic FAIL is observable
// in the item's Attempt log under the reserved "backend-error" rule carrying the
// original error text.
func TestRunLoopItemBackendError(t *testing.T) {
	dir := t.TempDir()
	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		calls++
		return "", errBackend
	})

	err := runLoopItem(stubDef{}, &LoopOptions{}, backend, s, it,
		dir+"/out.jsonl", dir+"/session.json", io.Discard)
	if err != nil {
		t.Fatalf("err = %v, want nil (generation error is demoted, not propagated)", err)
	}

	// Termination: each error consumes a try; MaxTries exhausts and locks DONE.
	if calls != quest.MaxTries {
		t.Fatalf("backend called %d times, want MaxTries=%d", calls, quest.MaxTries)
	}
	if it.Tries != quest.MaxTries {
		t.Fatalf("Tries = %d, want %d", it.Tries, quest.MaxTries)
	}
	if it.State != quest.DONE {
		t.Fatalf("state = %q, want DONE (MaxTries exhausted locks the residual)", it.State)
	}

	// Observability: the synthetic FAIL surfaces under "backend-error" with the
	// original error text in the Attempt reason (no invented content critique).
	if len(it.Log) != quest.MaxTries {
		t.Fatalf("Log has %d attempts, want %d", len(it.Log), quest.MaxTries)
	}
	last := it.Log[len(it.Log)-1]
	if !strings.Contains(last.Reason, backendErrorRule) {
		t.Fatalf("attempt reason = %q, want it to mention %q", last.Reason, backendErrorRule)
	}
	if !strings.Contains(last.Reason, errBackend.Error()) {
		t.Fatalf("attempt reason = %q, want it to carry the original error %q", last.Reason, errBackend.Error())
	}
}
