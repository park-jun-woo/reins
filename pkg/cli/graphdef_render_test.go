//ff:func feature=cli type=helper control=sequence
//ff:what graphDef.Render. 아이템 키 앞에 "render:"를 붙여 렌더 문자열을 돌려준다(테스트 더블).

package cli

import "github.com/park-jun-woo/reins/pkg/quest"

func (graphDef) Render(_ *quest.Session, it *quest.Item) (string, error) {
	return "render:" + it.Key, nil
}
