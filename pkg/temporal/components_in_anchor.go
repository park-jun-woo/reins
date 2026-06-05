//ff:func feature=temporal type=helper control=iteration dimension=1
//ff:what Start/End의 라틴숫자 런 성분이 모두 어떤 anchor에 실재하는지(textmatch 재사용) 검사해 엉뚱앵커(Türkiye·Kurtulmuş류)를 차단한다. 성분이 0개면 묶을 게 없어 true. 비라틴 숫자체계는 v1 미지원(v2).

package temporal

// ComponentsInAnchor reports whether every numeric component of the spec's Start/End
// dates actually appears in some anchor (via textmatch.Contains). This blocks the
// "date pinned to the wrong token" cheese: a spec is only trustworthy if its numbers
// are literally present in the source tokens the AI claims to have read.
//
// v1 extracts latin-digit runs from Start/End; other numeral systems (Persian,
// Arabic, Devanagari, Chinese) are v2 — they need a numeral-mapping table before
// their components can be matched against an anchor. When the spec has no numeric
// components, there is nothing to tie, so it returns true.
func ComponentsInAnchor(spec Spec, anchors []string) bool {
	comps := append(numericRuns(spec.Start), numericRuns(spec.End)...)
	for _, c := range comps {
		if !anchorsContain(anchors, c) {
			return false
		}
	}
	return true
}
