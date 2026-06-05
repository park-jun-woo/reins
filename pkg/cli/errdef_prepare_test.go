//ff:func feature=cli type=helper control=sequence level=error
//ff:what errDef.Prepare. prepareErr면 에러를, 아니면 raw를 Submission으로 담은 Context를 돌려준다(에러 분기 테스트 더블).

package cli

import (
	"errors"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

func (d errDef) Prepare(it *quest.Item, raw []byte) (gate.Context, *quest.Verdict, error) {
	if d.prepareErr {
		return gate.Context{}, nil, errors.New("prepare boom")
	}
	return gate.Context{Item: it, Submission: string(raw)}, nil, nil
}
