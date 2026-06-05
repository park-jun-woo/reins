package quest

// Sink receives terminal item records during export. Implementations choose the
// format (JSONL, CSV, …); reins ships a JSONL file sink in package cli.
type Sink interface {
	Emit(it *Item) error
}

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
