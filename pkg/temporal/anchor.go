package temporal

import "github.com/park-jun-woo/reins/pkg/textmatch"

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

// anchorsContain reports whether token appears in any of the anchors.
func anchorsContain(anchors []string, token string) bool {
	for _, a := range anchors {
		if textmatch.Contains(a, token) {
			return true
		}
	}
	return false
}

// numericRuns returns the maximal runs of latin digits (0-9) found in s. Leading
// zeros are preserved so the run matches the source surface form ("01" stays "01").
func numericRuns(s string) []string {
	var runs []string
	start := -1
	for i, r := range s {
		if r >= '0' && r <= '9' {
			if start < 0 {
				start = i
			}
			continue
		}
		if start >= 0 {
			runs = append(runs, s[start:i])
			start = -1
		}
	}
	if start >= 0 {
		runs = append(runs, s[start:])
	}
	return runs
}
