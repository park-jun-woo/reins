//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what graphDef.Seed. args 하나당 TODO 아이템을 시드한다(테스트 더블).

package cli

import "github.com/park-jun-woo/reins/pkg/quest"

func (d graphDef) Seed(args []string) ([]*quest.Item, error) {
	items := make([]*quest.Item, len(args))
	for i, a := range args {
		items[i] = &quest.Item{Key: a, State: quest.TODO}
	}
	return items, nil
}
