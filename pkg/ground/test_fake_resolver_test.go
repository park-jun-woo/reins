//ff:type feature=ground type=model
//ff:what fakeResolver — 테스트용 결정적·네트워크-free Resolver. Fetch/LookupMX 호출 횟수를 세어(스냅샷이 키당 최대 1회만 resolve함을 증명) 미리 준비한 본문/MX 응답 또는 미리 준비한 에러를 돌려준다.

package ground

// fakeResolver is a deterministic, network-free Resolver for tests. It counts how
// many times Fetch/LookupMX were called (so a test can prove the snapshot resolves
// at most once per key) and serves canned bodies/MX answers or canned errors.
type fakeResolver struct {
	bodies   map[string]string
	mx       map[string]bool
	fetchErr map[string]error
	mxErr    map[string]error

	fetchCalls  int
	lookupCalls int
}
