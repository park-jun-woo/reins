//ff:type feature=ground type=model
//ff:what Snapshot — 평가 1회 동안의 ground 스냅샷·캐시. 주입된 Resolver를 통해 HTTPBody/MXResolves를 lazy 계산하되 첫 읽힘에만 resolver를 호출하고 (성공·에러 모두) 캐시한다 → 같은 URL/도메인 재읽기는 resolver 0회 추가. 이로써 평가 내 leaf가 고정(요청당 결정성, 열린결정 #8). 캐시 키는 bodyCache(url→본문) / mxCache(domain→bool), 에러는 별도 errCache로 보존.

package ground

// Snapshot is the per-evaluation ground snapshot+cache. Through the injected
// Resolver it computes HTTPBody/MXResolves lazily, calling the resolver only on the
// first read of a given url/domain and caching the result (success or error). A
// repeat read of the same url/domain adds zero resolver calls, so every leaf is
// fixed within one evaluation (per-request determinism, open decision #8). The
// caches are keyed by url (body) and domain (MX); errors are preserved alongside so
// a failed resolve is reproduced from cache rather than re-attempted.
type Snapshot struct {
	resolver Resolver

	bodyCache map[string]string
	bodyErr   map[string]error
	bodySeen  map[string]bool

	mxCache map[string]bool
	mxErr   map[string]error
	mxSeen  map[string]bool
}
