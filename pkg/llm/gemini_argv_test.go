//ff:func feature=llm type=adapter control=iteration dimension=1
//ff:what TestGeminiArgv — geminiArgv의 3분기 argv 형태를 테이블로 단언: Stateless(opt 그대로, 세션 옵션 없음), Continue 1차(curSID=="" → opt + `--session-id <newSID>`), Continue 2차+(curSID!="" → opt + `--resume latest`, --session-id 없음). reflect.DeepEqual로 전체 슬라이스 비교.

package llm

import (
	"reflect"
	"testing"
)

// TestGeminiArgv asserts the three argv shapes geminiArgv builds from the session
// state, comparing the full slice. The first Continue call pins --session-id with the
// reins-issued UUID; later calls switch to --resume latest.
func TestGeminiArgv(t *testing.T) {
	opt := []string{"-p", "", "-o", "json", "--approval-mode", "plan"}
	cases := []struct {
		name   string
		kind   sessionKind
		curSID string
		newSID string
		want   []string
	}{
		{"stateless", Stateless, "", "U1", []string{"-p", "", "-o", "json", "--approval-mode", "plan"}},
		{"continue-first", Continue, "", "U1", []string{"-p", "", "-o", "json", "--approval-mode", "plan", "--session-id", "U1"}},
		{"continue-resume", Continue, "U1", "U1", []string{"-p", "", "-o", "json", "--approval-mode", "plan", "--resume", "latest"}},
	}
	for _, tc := range cases {
		got := geminiArgv(tc.kind, tc.curSID, tc.newSID, opt)
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("%s: geminiArgv = %v, want %v", tc.name, got, tc.want)
		}
	}
}
