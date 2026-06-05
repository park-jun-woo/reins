//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what submit 결과를 출력한다. "key -> OUTCOME (state STATE)" 한 줄에 더해, Verdict.Feedback(graph 백엔드가 렌더한 공략집)이 있으면 그걸 그대로 출력하고, 없으면 후방호환으로 Facts 루프(규칙·위치·기대·실제 심볼릭 피드백)로 폴백한다. cli는 graph/toulmin을 import하지 않고 문자열만 출력한다.

package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// printSubmit reports the outcome of one submission: the verdict line, then either
// the graph backend's pre-rendered walkthrough (Verdict.Feedback) when present, or
// the backward-compatible per-Fact loop otherwise. The CLI prints strings only — it
// does not import graph/toulmin.
func printSubmit(w io.Writer, key string, it *quest.Item, v quest.Verdict) {
	fmt.Fprintf(w, "%s -> %s (state %s)\n", key, v.Outcome, it.State)
	if v.Feedback != "" {
		for _, line := range strings.Split(strings.TrimRight(v.Feedback, "\n"), "\n") {
			fmt.Fprintf(w, "  %s\n", line)
		}
		return
	}
	for _, f := range v.Facts {
		fmt.Fprintf(w, "  - %s: %s expected=%q actual=%q\n", f.Rule, f.Where, f.Expected, f.Actual)
	}
}
