package textmatch

import "strings"

// Contains reports whether token appears as a substring of source after both are
// Normalize'd. An empty/whitespace-only token returns false — never the trivially
// true strings.Contains(src, "") (the empty-anchor cheese vector; reins Phase005,
// distilled from ccnews Phase009 L0).
func Contains(source, token string) bool {
	nt := Normalize(token)
	if nt == "" {
		return false
	}
	return strings.Contains(Normalize(source), nt)
}

// MissingTokens returns the tokens that are not present in source (after
// normalization). Source is normalized once. Empty/whitespace tokens count as
// missing. Handy for a gate rule to build a Fact naming the first offender.
func MissingTokens(source string, tokens []string) []string {
	ns := Normalize(source)
	var miss []string
	for _, t := range tokens {
		nt := Normalize(t)
		if nt == "" || !strings.Contains(ns, nt) {
			miss = append(miss, t)
		}
	}
	return miss
}
