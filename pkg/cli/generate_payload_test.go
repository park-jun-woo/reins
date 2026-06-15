//ff:func feature=cli type=command control=sequence level=error
//ff:what TestGeneratePayload — generatePayload 정상 생성 경로. backend.Complete 가 정상 payload 를 돌려주면 그 출력을 raw 로 그대로 반환하고 handled=false(호출부가 계속 게이트 평가하라는 신호)·err=nil 임을 단언. composeSystem(System,ruleCoach)·def.Render→Complete 가 호출돼 system/prompt 가 backend 로 전달됨도 관측.

package cli

import (
	"io"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestGeneratePayload: when backend.Complete succeeds, generatePayload returns the
// backend's output verbatim as raw, with handled=false (the caller should proceed to
// gate evaluation) and err=nil. The composed system + rendered prompt+feedback reach
// the backend.
func TestGeneratePayload(t *testing.T) {
	dir := t.TempDir()
	s := quest.New()
	it := &quest.Item{Key: "a", State: quest.TODO}
	s.Items = append(s.Items, it)

	var gotSystem, gotUser string
	backend := llm.CallFunc(func(system, user string) (string, error) {
		gotSystem, gotUser = system, user
		return "PAYLOAD-OK", nil
	})

	opts := &LoopOptions{System: "SYS"}
	raw, handled, err := generatePayload(stubDef{}, opts, backend, "COACH", "FEEDBACK",
		s, it, dir+"/out.jsonl", dir+"/session.json", io.Discard)
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if handled {
		t.Fatal("handled = true, want false (success defers to caller's gate evaluation)")
	}
	if raw != "PAYLOAD-OK" {
		t.Fatalf("raw = %q, want backend output %q", raw, "PAYLOAD-OK")
	}

	// The composed system carries opts.System and the rule coach.
	if !strings.Contains(gotSystem, "SYS") || !strings.Contains(gotSystem, "COACH") {
		t.Fatalf("system = %q, want it to compose opts.System and ruleCoach", gotSystem)
	}
	// The user prompt is def.Render output with feedback appended.
	if !strings.Contains(gotUser, "render:a") || !strings.Contains(gotUser, "FEEDBACK") {
		t.Fatalf("user = %q, want rendered prompt + feedback", gotUser)
	}

	// No demotion occurred: the item stays TODO with no try consumed.
	if it.Tries != 0 || it.State != quest.TODO {
		t.Fatalf("item mutated on success: Tries=%d State=%q, want 0/TODO", it.Tries, it.State)
	}
}
