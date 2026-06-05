//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what submit 결과를 출력한다. "key -> OUTCOME (state STATE)" 한 줄과, verdict.Facts가 있으면(주로 FAIL) 각 Fact의 규칙·위치·기대·실제를 심볼릭 피드백으로 덧붙인다.

package cli

import (
	"fmt"
	"io"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// printSubmit reports the outcome of one submission: the verdict line plus any Facts
// (the self-correction feedback shown on FAIL).
func printSubmit(w io.Writer, key string, it *quest.Item, v quest.Verdict) {
	fmt.Fprintf(w, "%s -> %s (state %s)\n", key, v.Outcome, it.State)
	for _, f := range v.Facts {
		fmt.Fprintf(w, "  - %s: %s expected=%q actual=%q\n", f.Rule, f.Where, f.Expected, f.Actual)
	}
}
