//ff:func feature=cli type=helper control=iteration dimension=1 level=error
//ff:what errDef.Seed. seedErr면 에러를, 아니면 args 하나당 TODO 아이템을 시드한다(에러 분기 테스트 더블).

package cli

import (
	"errors"

	"github.com/park-jun-woo/reins/pkg/quest"
)

func (d errDef) Seed(args []string) ([]*quest.Item, error) {
	if d.seedErr {
		return nil, errors.New("seed boom")
	}
	items := make([]*quest.Item, len(args))
	for i, a := range args {
		items[i] = &quest.Item{Key: a, State: quest.TODO}
	}
	return items, nil
}
