package textmatch

import "testing"

func TestContainsNFC(t *testing.T) {
	// "é" decomposed (e + combining acute) vs composed should match after NFC.
	source := "café society"
	if !Contains(source, "café") {
		t.Fatal("NFC normalization should make composed/decomposed forms equal")
	}
}

func TestContainsEmptyTokenFalse(t *testing.T) {
	if Contains("anything", "   ") {
		t.Fatal("empty/whitespace token must not match (empty-anchor cheese)")
	}
}

func TestMissingTokens(t *testing.T) {
	miss := MissingTokens("the quick brown fox", []string{"quick", "lazy", ""})
	if len(miss) != 2 || miss[0] != "lazy" || miss[1] != "" {
		t.Fatalf("missing = %v", miss)
	}
}
