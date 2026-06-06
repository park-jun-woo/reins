//ff:func feature=cli type=helper control=sequence
//ff:what graphDef.Prepare. raw 바이트를 문자열 Submission으로 담은 gate.Context를 만든다(테스트 더블).

package cli

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func (graphDef) Prepare(_ *quest.Session, it *quest.Item, raw []byte) (gate.Context, *quest.Verdict, error) {
	return gate.Context{Item: it, Submission: string(raw)}, nil, nil
}
