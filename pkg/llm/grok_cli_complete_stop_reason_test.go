//ff:func feature=llm type=adapter control=sequence level=error
//ff:what TestGrokCLICompleteStopReason — stopReason!="EndTurn"(예: "Cancelled") 응답이 그 stopReason을 동봉한 에러를 표면화하는지 검증(grok의 claude is_error 대응).

package llm

import (
	"context"
	"strings"
	"testing"
)

// TestGrokCLICompleteStopReason: a non-EndTurn stopReason surfaces an error carrying
// the stop reason (grok's analog of claude is_error).
func TestGrokCLICompleteStopReason(t *testing.T) {
	orig := execGrok
	defer func() { execGrok = orig }()

	execGrok = func(ctx context.Context, bin string, argv []string, stdin string) (string, string, error) {
		return `{"text":"","sessionId":"S1","stopReason":"Cancelled"}`, "", nil
	}

	g := &GrokCLI{}
	_, err := g.Complete("SYS", "USR")
	if err == nil {
		t.Fatal("Complete = nil error, want stop-reason error")
	}
	if !strings.Contains(err.Error(), "Cancelled") {
		t.Fatalf("error = %q, want it to contain stopReason", err.Error())
	}
}
