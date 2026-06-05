//ff:func feature=graph type=helper control=sequence
//ff:what SpecName — 노드 ID를 반환한다(toulmin ruleID 구성에 쓰임).

package graph

// SpecName returns the node ID (used to build the toulmin ruleID).
func (s idSpec) SpecName() string { return s.ID }
