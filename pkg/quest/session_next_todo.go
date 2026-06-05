//ff:func feature=quest type=helper control=iteration dimension=1
//ff:what 아직 TODO인 첫 아이템을 반환한다. 잠긴 상태는 절대 다시 집지 않는다(래칫).

package quest

// NextTODO returns the first item still in TODO, or nil when none remain.
func (s *Session) NextTODO() *Item {
	for _, it := range s.Items {
		if it.State == TODO {
			return it
		}
	}
	return nil
}
