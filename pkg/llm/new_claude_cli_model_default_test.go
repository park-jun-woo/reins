//ff:func feature=llm type=adapter control=sequence
//ff:what TestNewClaudeCLIModelDefault — newClaudeCLI에서 model=="default"는 빈 Model(CLI 기본)로, 그 외 토큰은 그대로 Model에 실리는지 검증. 무서브프로세스.

package llm

import (
	"testing"
)

// TestNewClaudeCLIModelDefault: the "default" sentinel maps to an empty Model
// (CLI default model); any other token is carried verbatim.
func TestNewClaudeCLIModelDefault(t *testing.T) {
	if c := newClaudeCLI("default"); c.Model != "" {
		t.Fatalf("Model = %q, want empty for default sentinel", c.Model)
	}
	if c := newClaudeCLI("opus"); c.Model != "opus" {
		t.Fatalf("Model = %q, want opus", c.Model)
	}
}
