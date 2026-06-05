//ff:func feature=quest type=helper control=sequence
//ff:what 세션 Meta 슬롯에서 키 값을 읽는다. nil 맵·부재 키 모두 (nil,false)로 보고한다(nil-safe).

package quest

// GetMeta returns the value stored under key in the session's Meta slot. The
// second result is false when Meta is nil or the key is absent, mirroring the
// comma-ok map idiom.
func (s *Session) GetMeta(key string) (any, bool) {
	if s.Meta == nil {
		return nil, false
	}
	v, ok := s.Meta[key]
	return v, ok
}
