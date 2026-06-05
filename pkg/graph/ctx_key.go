//ff:type feature=graph type=model
//ff:what gctxKey — reins gate.Context를 toulmin ctx에 싣고 꺼내는 약속된 키. Evaluate가 ctx.Set(gctxKey, gctx)로 넣고 어댑터가 ctx.Get(gctxKey)로 읽는다.

package graph

// gctxKey is the agreed toulmin-context key under which Evaluate stores the reins
// gate.Context. Rule adapters read it back via ctx.Get(gctxKey).
const gctxKey = "__gctx__"
