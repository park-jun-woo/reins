//ff:func feature=temporal type=helper control=sequence
//ff:what numericRuns가 숫자가 없는 문자열에서 nil을 내는지 검증한다.

package temporal

import "testing"

func TestNumericRunsNone(t *testing.T) {
	if got := numericRuns("no digits here"); got != nil {
		t.Fatalf("numericRuns = %v, want nil", got)
	}
}
