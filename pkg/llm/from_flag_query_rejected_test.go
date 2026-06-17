//ff:func feature=llm type=helper control=iteration dimension=1 level=error
//ff:what TestFromFlagQueryRejected — 미적용/미지 쿼리 키가 loud 거부되는지 테이블로 검증: xai에 num_ctx·think, gemini에 num_ctx, subprocess(claude/grok/codex/geminicli)에 임의 쿼리, 알 수 없는 키, 파싱 실패 → 모두 에러 + nil backend.

package llm

import "testing"

// TestFromFlagQueryRejected: options not supported by the chosen backend (and parse
// failures) are rejected with an error and a nil backend.
func TestFromFlagQueryRejected(t *testing.T) {
	cases := []string{
		"xai:grok?num_ctx=4096",                 // xai does not allow num_ctx
		"xai:grok?think=false",                  // xai does not allow think
		"gemini:g?num_ctx=4096",                 // gemini does not allow num_ctx
		"claude:default?max_output_tokens=8192", // subprocess takes no query
		"grok:default?temperature=0.5",          // subprocess takes no query
		"codex:default?max_output_tokens=1",     // subprocess takes no query
		"geminicli:default?temperature=0",       // subprocess takes no query
		"ollama:m?bogus=1",                      // unknown key
		"ollama:m?max_output_tokens=abc",        // parse failure
	}
	for _, flag := range cases {
		t.Run(flag, func(t *testing.T) {
			b, err := FromFlag(flag)
			if err == nil {
				t.Fatalf("FromFlag(%q) = %v, want error", flag, b)
			}
			if b != nil {
				t.Fatalf("FromFlag(%q) backend = %v, want nil on error", flag, b)
			}
		})
	}
}
