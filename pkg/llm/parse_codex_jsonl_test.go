//ff:func feature=llm type=adapter control=iteration dimension=1 level=error
//ff:what TestParseCodexJSONL — parseCodexJSONL을 테이블로 단언: (1) 다중 이벤트에서 최종 agent_message 텍스트와 thread_id sid 추출·미파싱 노이즈 무시, (2) error 이벤트는 에러, (3) agent_message 부재는 에러. wantErr 행은 err!=nil만 확인.

package llm

import "testing"

// TestParseCodexJSONL asserts parseCodexJSONL over a table: it extracts the last
// agent_message text plus the thread_id, tolerates unparseable noise lines, errors on
// an error event, and errors when no agent_message is present.
func TestParseCodexJSONL(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		wantText string
		wantSID  string
		wantErr  bool
	}{
		{
			name: "final-message",
			in: `{"type":"thread.started","thread_id":"T1"}
noise
{"type":"item.completed","item":{"type":"agent_message","text":"first"}}
{"type":"item.completed","item":{"type":"agent_message","text":"FINAL"}}
{"type":"turn.completed"}`,
			wantText: "FINAL", wantSID: "T1",
		},
		{name: "error-event", in: `{"type":"error","message":"boom"}`, wantErr: true},
		{name: "no-message", in: `{"type":"thread.started","thread_id":"T1"}`, wantErr: true},
	}
	for _, tc := range cases {
		text, sid, err := parseCodexJSONL(tc.in)
		if tc.wantErr && err == nil {
			t.Fatalf("%s: err = nil, want error", tc.name)
		}
		if tc.wantErr {
			continue
		}
		if err != nil {
			t.Fatalf("%s: unexpected err %v", tc.name, err)
		}
		if text != tc.wantText {
			t.Fatalf("%s: text = %q, want %q", tc.name, text, tc.wantText)
		}
		if sid != tc.wantSID {
			t.Fatalf("%s: sid = %q, want %q", tc.name, sid, tc.wantSID)
		}
	}
}
