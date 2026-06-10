//ff:func feature=cli type=helper control=sequence
//ff:what stubDef.Render. "render:<key>"를 돌려준다(테스트 더블).

package cli

import "github.com/park-jun-woo/reins/pkg/quest"

func (stubDef) Render(_ *quest.Session, it *quest.Item) (string, error) {
	return "render:" + it.Key, nil
}
