//ff:func feature=quest type=helper control=iteration dimension=1 level=error
//ff:what terminal·미방출 아이템을 sink로 방출하고 Emitted를 세팅한다. 증분, 1회 보장. 새로 방출한 건수 반환.

package quest

// Export emits every terminal, not-yet-emitted item to the sink and sets Emitted,
// so each item is exported at most once across runs (the export ratchet). It
// returns the number of newly emitted records.
func Export(s *Session, sink Sink) (int, error) {
	n := 0
	for _, it := range s.Items {
		if !it.State.Terminal() || it.Emitted {
			continue
		}
		if err := sink.Emit(it); err != nil {
			return n, err
		}
		it.Emitted = true
		n++
	}
	return n, nil
}
