//ff:func feature=graph type=helper control=sequence
//ff:what 테스트 헬퍼 — tier-0 Fail 카운터(email-format) 1개와, 둘 다 "source-body"를 needs하는 tier-1 카운터 2개를 가진 그래프를 만든다(같은 ground를 두 번 읽어 스냅샷 캐시를 행사하게).

package graph

import (
	"github.com/park-jun-woo/reins/pkg/gate"
)

// stagedGraph builds a graph with one tier-0 Fail counter (email-format) and two
// tier-1 counters both needing "source-body" (so the snapshot cache is exercised).
func stagedGraph() *Graph {
	g := NewGraph("staged")
	pass := g.Warrant(alwaysTrue)
	g.Counter(fireRule("email-format", gate.LevelFail), gate.LevelFail).Attacks(pass)
	g.Counter(groundFireRule("source-lacks-email", "source-body"), gate.LevelFail).
		Attacks(pass).Needs("source-body")
	g.Counter(groundFireRule("source-stale", "source-body"), gate.LevelFail).
		Attacks(pass).Needs("source-body")
	return g
}
