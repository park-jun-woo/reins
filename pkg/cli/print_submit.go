//ff:func feature=cli type=helper control=sequence
//ff:what printSubmit — 제출 1건의 결과를 보고한다(renderVerdict 래퍼).

package cli

import (
	"io"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// printSubmit reports the outcome of one submission (renderVerdict wrapper).
func printSubmit(w io.Writer, key string, it *quest.Item, v quest.Verdict) {
	renderVerdict(w, key, it, v)
}
