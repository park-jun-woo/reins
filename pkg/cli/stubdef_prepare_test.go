//ff:func feature=cli type=helper control=sequence
//ff:what stubDef.Prepare. 제출물이 "skip"이면 SKIPPED verdict로 단락하고(규칙 우회), 아니면 raw를 Submission으로 담은 Context를 돌려준다.

package cli

import (
	"strings"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func (stubDef) Prepare(it *quest.Item, raw []byte) (gate.Context, *quest.Verdict, error) {
	if strings.TrimSpace(string(raw)) == "skip" {
		return gate.Context{}, &quest.Verdict{Outcome: quest.OutSkip}, nil
	}
	return gate.Context{Item: it, Submission: string(raw)}, nil, nil
}
