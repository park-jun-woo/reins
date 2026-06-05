//ff:func feature=graph type=helper control=sequence
//ff:what shortNameFor — toulmin collectTrace의 short name 규칙을 재현한다(마지막 "."까지 잘라낸 base + "#spec" 유지). Evaluate가 trace 엔트리(Name=short)와 노드를 매칭하는 키.

package graph

import "strings"

// shortNameFor reproduces toulmin's collectTrace shortName: drop everything up to
// the last "." in the base, keeping any "#spec" suffix. This is the key Evaluate
// uses to match a collected TraceEntry (whose Name is shortened) to a Node.
func shortNameFor(full string) string {
	base, spec, hasSpec := strings.Cut(full, "#")
	if idx := strings.LastIndex(base, "."); idx >= 0 {
		base = base[idx+1:]
	}
	if hasSpec {
		return base + "#" + spec
	}
	return base
}
