//ff:func feature=cli type=helper control=sequence level=error
//ff:what errDef.Render. renderErr면 에러를, 아니면 "render:<key>"를 돌려준다(에러 분기 테스트 더블).

package cli

import (
	"errors"

	"github.com/park-jun-woo/reins/pkg/quest"
)

func (d errDef) Render(_ *quest.Session, it *quest.Item) (string, error) {
	if d.renderErr {
		return "", errors.New("render boom")
	}
	return "render:" + it.Key, nil
}
