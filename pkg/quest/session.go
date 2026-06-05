//ff:type feature=quest type=model
//ff:what 한 퀘스트 실행의 영속 런타임 상태. 작업 목록과 스키마 버전을 담는 SSOT.

package quest

// Session is the persisted runtime state of one quest run: the work-list plus a
// schema version. It is the single source of truth a disposable agent resumes from.
type Session struct {
	Version int     `json:"version"`
	Items   []*Item `json:"items"`
}
