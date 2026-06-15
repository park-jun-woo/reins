//ff:func feature=cli type=command control=sequence level=error
//ff:what TestRunLoopItemBackendErrorRecovers — 혼합 수렴. 1·2차 생성 오류가 Try 를 소비하지만(②(a)), MaxTries 내라면 3차 정상 페이로드가 게이트 PASS 로 수렴함을 단언(에러 소진이 회복을 막지 않음).

package cli

import (
	"io"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestRunLoopItemBackendErrorRecovers: the first two attempts error (consuming two
// tries) but the third returns a passing payload — within MaxTries=3 the item still
// converges to gate PASS. Transient generation errors burn tries yet do not bar
// recovery.
func TestRunLoopItemBackendErrorRecovers(t *testing.T) {
	dir := t.TempDir()
	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		calls++
		if calls < quest.MaxTries {
			return "", errBackend // first two attempts error
		}
		return "good", nil // the final try succeeds
	})

	err := runLoopItem(stubDef{}, &LoopOptions{}, backend, s, it,
		dir+"/out.jsonl", dir+"/session.json", io.Discard)
	if err != nil {
		t.Fatalf("runLoopItem: %v", err)
	}
	if calls != quest.MaxTries {
		t.Fatalf("backend called %d times, want %d", calls, quest.MaxTries)
	}
	if it.State != quest.PASS {
		t.Fatalf("state = %q, want PASS (recovered on the final try)", it.State)
	}
}
