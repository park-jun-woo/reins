//ff:func feature=ground type=helper control=sequence
//ff:what NewSnapshot — 새 ground Snapshot을 만든다. resolver가 nil이면 실네트워크 기본 resolver를 채택하고(가드), 테스트는 fake를 넘긴다. Snapshot은 평가 1회용 — 캐시가 평가 간 새지 않도록 Evaluate마다 하나씩 만든다.

package ground

// NewSnapshot builds a fresh ground Snapshot. If resolver is nil the real-network
// default resolver is used; tests pass a fake. Each Snapshot is single-evaluation:
// build one per Evaluate call so caches do not leak across evaluations.
func NewSnapshot(resolver Resolver) *Snapshot {
	if resolver == nil {
		resolver = newDefaultResolver()
	}
	return &Snapshot{
		resolver:  resolver,
		bodyCache: make(map[string]string),
		bodyErr:   make(map[string]error),
		bodySeen:  make(map[string]bool),
		mxCache:   make(map[string]bool),
		mxErr:     make(map[string]error),
		mxSeen:    make(map[string]bool),
	}
}
