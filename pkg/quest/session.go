//ff:type feature=quest type=model
//ff:what 한 퀘스트 실행의 영속 런타임 상태. 작업 목록과 스키마 버전을 담는 SSOT.

package quest

// Session is the persisted runtime state of one quest run: the work-list plus a
// schema version. It is the single source of truth a disposable agent resumes from.
type Session struct {
	Version int     `json:"version"`
	Items   []*Item `json:"items"`
	// Meta is an optional free-form slot for consumer-specific runtime state
	// (e.g. a crawler's user-agent, cursors, or per-host robots cache) that must
	// survive the Save/Load round-trip alongside the work-list. Additive and
	// backward-compatible: omitempty keeps it absent from older sessions, and a
	// nil map round-trips cleanly.
	Meta map[string]any `json:"meta,omitempty"`
}
