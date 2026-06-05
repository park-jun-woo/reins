//ff:type feature=graph type=model
//ff:what stagedFakeResolver — staged-평가 테스트용 네트워크-free ground.Resolver. Fetch 호출 횟수를 세어(스냅샷이 최대 1회 resolve함을 단언) 어떤 URL에도 미리 정한 본문을 돌려준다.

package graph

// stagedFakeResolver is a network-free ground.Resolver for the staged-evaluation
// tests. It counts Fetch calls so a test can assert the snapshot resolved at most
// once, and serves a canned body for any URL.
type stagedFakeResolver struct {
	body       string
	fetchCalls int
}
