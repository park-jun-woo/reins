//ff:type feature=graph type=model
//ff:what activeCounter — trace에서 활성으로 판독된 카운터 노드 1건과 그 Fact(어댑터 evidence)를 묶는다. 집계·supersession 입력.

package graph

import "github.com/park-jun-woo/reins/pkg/quest"

// activeCounter pairs an activated counter node read from the trace with its Fact
// (the adapter's evidence). It is the input to aggregation and supersession.
type activeCounter struct {
	node *Node
	fact quest.Fact
}
