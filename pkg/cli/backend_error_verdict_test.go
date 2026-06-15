//ff:func feature=cli type=helper control=sequence
//ff:what TestBackendErrorVerdict — L0 생성 오류를 감싼 verdict가 재시도 가능한 FAIL(RootCause=backend-error)이고, 원문 에러를 backend.Complete 위치의 단일 Fact로 싣는지 단언. 게이트 실패와 생성 실패의 분류 경계를 고정한다.

package cli

import (
	"errors"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestBackendErrorVerdict: the verdict wrapping an L0 generation error is a retryable
// FAIL whose RootCause is the reserved "backend-error" rule, carrying the original
// error text in a single Fact located at backend.Complete.
func TestBackendErrorVerdict(t *testing.T) {
	err := errors.New("boom")
	v := backendErrorVerdict(err)

	if v.Outcome != quest.OutFail {
		t.Fatalf("outcome = %v, want %v", v.Outcome, quest.OutFail)
	}
	if v.RootCause != "backend-error" {
		t.Fatalf("root cause = %q, want %q", v.RootCause, "backend-error")
	}
	if len(v.Facts) != 1 {
		t.Fatalf("facts len = %d, want 1", len(v.Facts))
	}
	f := v.Facts[0]
	if f.Rule != "backend-error" {
		t.Fatalf("fact rule = %q, want %q", f.Rule, "backend-error")
	}
	if f.Where != "backend.Complete" {
		t.Fatalf("fact where = %q, want %q", f.Where, "backend.Complete")
	}
	if f.Actual != err.Error() {
		t.Fatalf("fact actual = %q, want %q", f.Actual, err.Error())
	}
}
