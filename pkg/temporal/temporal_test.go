package temporal

import (
	"testing"
	"time"
)

func TestResolveRelative(t *testing.T) {
	ref := time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC)
	r := Resolve(Spec{Kind: Relative, OffsetDays: -1}, ref)
	if !r.Determined || r.Value != "2026-06-04" {
		t.Fatalf("got %+v, want 2026-06-04 determined", r)
	}
}

func TestResolveAbsoluteSingle(t *testing.T) {
	r := Resolve(Spec{Kind: Absolute, Calendar: Gregorian, Start: "2017-01-10"}, time.Time{})
	if !r.Determined || r.IsInterval || r.Value != "2017-01-10" {
		t.Fatalf("got %+v", r)
	}
}

func TestResolveAbsoluteInterval(t *testing.T) {
	r := Resolve(Spec{Kind: Absolute, Calendar: Gregorian, Start: "2017-01-10", End: "2017-01-14"}, time.Time{})
	if !r.Determined || !r.IsInterval || r.Value != "2017-01-10/2017-01-14" {
		t.Fatalf("got %+v", r)
	}
}

func TestResolveUndetermined(t *testing.T) {
	if Resolve(Spec{Kind: Absolute, Calendar: Persian, Start: "1395-10-21"}, time.Time{}).Determined {
		t.Fatal("non-gregorian should be undetermined")
	}
	if Resolve(Spec{Kind: Absolute, Calendar: Gregorian, Start: "nope"}, time.Time{}).Determined {
		t.Fatal("unparseable should be undetermined")
	}
}

func TestComponentsInAnchor(t *testing.T) {
	spec := Spec{Kind: Absolute, Calendar: Gregorian, Start: "2017-01-10"}
	if !ComponentsInAnchor(spec, []string{"on 2017-01-10 in Davos"}) {
		t.Fatal("components present should pass")
	}
	if ComponentsInAnchor(spec, []string{"Türkiye Kurtulmuş said"}) {
		t.Fatal("components absent should fail")
	}
	if !ComponentsInAnchor(Spec{Kind: Relative}, nil) {
		t.Fatal("no components ⇒ nothing to tie ⇒ true")
	}
}
