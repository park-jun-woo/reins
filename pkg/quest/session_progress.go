//ff:func feature=quest type=helper control=iteration dimension=1
//ff:what 아이템을 상태별로 집계한다(status 명령용).

package quest

// Progress tallies items by state (for the status command).
func (s *Session) Progress() map[State]int {
	out := map[State]int{}
	for _, it := range s.Items {
		out[it.State]++
	}
	return out
}
