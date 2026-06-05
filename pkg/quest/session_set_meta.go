//ff:func feature=quest type=helper control=sequence
//ff:what 세션 Meta 슬롯에 키-값을 넣는다. nil 맵이면 지연 초기화한다(nil-safe).

package quest

// SetMeta stores v under key in the session's Meta slot, lazily allocating the
// map on first use so callers never have to nil-check before writing.
func (s *Session) SetMeta(key string, v any) {
	if s.Meta == nil {
		s.Meta = make(map[string]any)
	}
	s.Meta[key] = v
}
