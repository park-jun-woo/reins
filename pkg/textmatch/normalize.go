package textmatch

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

// Normalize applies Unicode NFC, collapses every run of whitespace to a single
// space, and trims the ends. It is applied identically to the source text and to
// each token so matching compares like-for-like surface forms — never an inferred
// mapping. NFC unifies composed/decomposed forms (é = "e"+◌́ vs "é"), which is what
// closes the diacritic/combining-mark false-negatives a whitespace-only normalize
// leaves open.
func Normalize(s string) string {
	return strings.Join(strings.Fields(norm.NFC.String(s)), " ")
}
