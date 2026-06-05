//ff:func feature=graph type=helper control=sequence
//ff:what String — 노드 ID를 반환한다; toulmin ruleID는 fmt.Sprint(spec)로 "#" 접미사를 만든다.

package graph

// String returns the node ID; ruleID uses fmt.Sprint(spec) for the "#" suffix.
func (s idSpec) String() string { return s.ID }
