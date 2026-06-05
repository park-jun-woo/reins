//ff:func feature=graph type=helper control=sequence
//ff:what ruleIDFor 증명 — 명명된 fn에서 "funcID#spec" 형태의 ruleID를 만들고(접미사가 spec.String), 같은 fn에 다른 spec이면 다른 ID·같은 spec이면 같은 ID를 내는 결정성을 검증한다.

package graph

import (
	"strings"
	"testing"
)

func TestRuleIDFor(t *testing.T) {
	id1 := ruleIDFor(alwaysTrue, idSpec{ID: "warrant#0"})
	id2 := ruleIDFor(alwaysTrue, idSpec{ID: "warrant#1"})
	id1again := ruleIDFor(alwaysTrue, idSpec{ID: "warrant#0"})

	if !strings.HasSuffix(id1, "#warrant#0") {
		t.Fatalf("id1=%q want suffix #warrant#0", id1)
	}
	if !strings.Contains(id1, "alwaysTrue") {
		t.Fatalf("id1=%q want funcID containing alwaysTrue", id1)
	}
	if id1 == id2 {
		t.Fatalf("different specs produced same id: %q", id1)
	}
	if id1 != id1again {
		t.Fatalf("same fn+spec not deterministic: %q vs %q", id1, id1again)
	}
}
