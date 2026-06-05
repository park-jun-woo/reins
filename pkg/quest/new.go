//ff:func feature=quest type=helper control=sequence
//ff:what 현재 스키마 버전(1)의 빈 세션을 새로 만든다.

package quest

// New returns an empty session at the current schema version.
func New() *Session { return &Session{Version: 1} }
