//ff:func feature=temporal type=helper control=sequence
//ff:what numericRuns가 라틴숫자의 최대 연속 런을 추출하고 선행 0을 보존하는지 검증한다.

package temporal

import (
	"reflect"
	"testing"
)

func TestNumericRunsExtracts(t *testing.T) {
	// "01" preserved (leading zero); a trailing run ("10") closed at end of string.
	got := numericRuns("2017-01-10")
	want := []string{"2017", "01", "10"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("numericRuns = %v, want %v", got, want)
	}
}
