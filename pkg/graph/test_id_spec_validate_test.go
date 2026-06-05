//ff:func feature=graph type=helper control=sequence
//ff:what Validate 증명 — idSpec.Validate가 빈 ID는 거부(error)하고 비어있지 않은 ID는 통과(nil)하는지 검증한다.

package graph

import "testing"

func TestIDSpecValidate(t *testing.T) {
	if err := (idSpec{ID: ""}).Validate(); err == nil {
		t.Fatalf("empty ID: want error, got nil")
	}
	if err := (idSpec{ID: "node-1"}).Validate(); err != nil {
		t.Fatalf("non-empty ID: want nil, got %v", err)
	}
}
