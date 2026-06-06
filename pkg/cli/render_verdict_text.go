//ff:func feature=cli type=helper control=sequence
//ff:what renderVerdictText — renderVerdict가 출력할 텍스트를 string으로 반환한다. agent의 재시도 피드백이 사람이 보는 submit 출력과 바이트 동일하도록 한다(사람이 보는 피드백=모델이 받는 피드백).

package cli

import (
	"strings"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// renderVerdictText returns the same text renderVerdict would print, so the agent's
// retry feedback to the model is byte-identical to what submit shows a human.
func renderVerdictText(key string, it *quest.Item, v quest.Verdict) string {
	var b strings.Builder
	renderVerdict(&b, key, it, v)
	return b.String()
}
