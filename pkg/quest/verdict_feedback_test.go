//ff:func feature=quest type=helper control=sequence
//ff:what Verdict.Feedback가 additive임을 검증한다 — 빈 Feedback이면 Reason()/Outcome/Facts 동작 무영향(후방호환), 채워지면 그 문자열을 보존하고 Reason()은 여전히 Facts만 렌더한다.

package quest

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestVerdictFeedbackAdditive: an empty Feedback leaves Reason/Outcome/Facts intact
// (backward compatible); a set Feedback is preserved and does not bleed into Reason.
func TestVerdictFeedbackAdditive(t *testing.T) {
	v := Verdict{Outcome: OutFail, Facts: []Fact{{Rule: "R1", Where: "x"}}}
	if v.Feedback != "" {
		t.Fatalf("zero-value Feedback = %q, want empty", v.Feedback)
	}
	if got := v.Reason(); got != "R1: x" {
		t.Fatalf("Reason() = %q (Feedback must not affect it)", got)
	}

	v.Feedback = "FAIL. root cause = R1"
	if got := v.Reason(); got != "R1: x" {
		t.Fatalf("Reason() with Feedback set = %q, want unchanged", got)
	}

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(b), `"feedback":"FAIL. root cause = R1"`) {
		t.Fatalf("marshalled Feedback = %s", b)
	}

	// omitempty: empty Feedback is not serialized.
	b2, _ := json.Marshal(Verdict{Outcome: OutPass})
	if strings.Contains(string(b2), "feedback") {
		t.Fatalf("empty Feedback should be omitted, got %s", b2)
	}
}
