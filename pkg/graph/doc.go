//ff:type feature=graph type=model
//ff:what graph 패키지 개요. toulmin defeat-그래프를 reins 게이트 백엔드로 감싸는 편의 빌더. tautology PASS 워런트 1개 + 위반=공격 Counter 토폴로지. toulmin 의존을 이 패키지에 격리(pkg/gate는 graph를 import하지 않음 — 단방향). 노드별 reins Level 메타 보유(toulmin Strength 미의존). Supersedes는 reins-side 우선순위 맵(toulmin 엣지 아님). 엣지0 그래프는 현 gate.Evaluate와 동치.

// Package graph wraps the toulmin defeat-graph engine as a reins gate backend.
//
// Topology: one tautology PASS warrant node plus violation Counter nodes that
// attack it. Evaluate reads the warrant's trace (per-counter Activated) and joins
// each counter's reins Level to produce a quest.Verdict — never via the
// h-Categoriser float (which cannot drive a single decisive violation to FAIL).
//
// The toulmin dependency is isolated here: pkg/gate does NOT import pkg/graph
// (one-way), so existing consumers (comail, cli) stay green and toulmin-free.
//
// Supersession (Node.Supersedes) is a reins-side precedence map, not a toulmin
// edge: an active upstream counter excludes superseded downstream counters from
// aggregation. A graph with no edges and no supersession is equivalent to the
// current gate.Evaluate level aggregation.
package graph
