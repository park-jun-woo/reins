//ff:func feature=quest type=helper control=iteration dimension=1 level=error
//ff:what Key로 아이템을 찾는다. 없으면 에러를 반환한다.

package quest

import "fmt"

// Find returns the item with the given Key, or an error if absent.
func (s *Session) Find(key string) (*Item, error) {
	for _, it := range s.Items {
		if it.Key == key {
			return it, nil
		}
	}
	return nil, fmt.Errorf("item not found: %s", key)
}
