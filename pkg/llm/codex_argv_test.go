//ff:func feature=llm type=adapter control=iteration dimension=1
//ff:what TestCodexArgv — codexArgv의 3분기 argv 형태를 테이블로 단언: Stateless(`exec --ephemeral` + flagsExec + flagsCommon + "-"), Continue 1차(--ephemeral 없음, flagsExec 선두), Continue 2차(`-s` 가 `resume <sid>` "앞", --ephemeral 없음). reflect.DeepEqual로 전체 슬라이스 비교.

package llm

import (
	"reflect"
	"testing"
)

// TestCodexArgv asserts the three argv shapes codexArgv builds from the session
// state, comparing the full slice. The resume branch keeps -s before the resume
// subcommand and never emits --ephemeral.
func TestCodexArgv(t *testing.T) {
	exec := []string{"-s", "read-only"}
	common := []string{"--json", "-x"}
	cases := []struct {
		name string
		kind sessionKind
		sid  string
		want []string
	}{
		{"stateless", Stateless, "", []string{"exec", "--ephemeral", "-s", "read-only", "--json", "-x", "-"}},
		{"continue-first", Continue, "", []string{"exec", "-s", "read-only", "--json", "-x", "-"}},
		{"continue-resume", Continue, "T1", []string{"exec", "-s", "read-only", "resume", "T1", "--json", "-x", "-"}},
	}
	for _, tc := range cases {
		got := codexArgv(tc.kind, tc.sid, exec, common)
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("%s: codexArgv = %v, want %v", tc.name, got, tc.want)
		}
	}
}
