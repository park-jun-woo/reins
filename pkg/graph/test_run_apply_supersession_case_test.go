//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what runApplySupersessionCase 테스트헬퍼 — 활성 id 목록을 activeCounter로 빌드해 applySupersession에 넣고, 잔존 id가 기대 순서·길이와 일치하는지 단언한다(케이스 1건 판정).

package graph

import "testing"

func runApplySupersessionCase(t *testing.T, active, want []string) {
	t.Helper()
	g, byID := buildSupersessionGraph()
	acs := make([]activeCounter, 0, len(active))
	for _, id := range active {
		acs = append(acs, activeCounter{node: byID[id]})
	}
	got := acIDs(g.applySupersession(acs))
	if len(got) != len(want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("got %v want %v", got, want)
		}
	}
}
