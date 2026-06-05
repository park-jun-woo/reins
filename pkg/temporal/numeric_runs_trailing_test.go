//ff:func feature=temporal type=helper control=sequence
//ff:what numericRuns가 문자열 끝의 숫자 런을 닫아 포함하는지 검증한다.

package temporal

import (
	"reflect"
	"testing"
)

func TestNumericRunsTrailing(t *testing.T) {
	got := numericRuns("abc42")
	want := []string{"42"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("numericRuns = %v, want %v", got, want)
	}
}
