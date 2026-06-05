//ff:type feature=ground type=model
//ff:what Resolver — ground 원시연산의 부작용(네트워크)을 격리하는 주입 가능한 인터페이스. Fetch(URL→본문)·LookupMX(도메인→수신가능)만 노출한다. 기본 구현은 실제 net 스택(defaultResolver), 테스트는 fake를 주입해 결정적·네트워크-free로 만든다. Snapshot이 이 인터페이스를 통해서만 외부 사실을 읽고 그 결과를 캐시한다.

package ground

// Resolver isolates the side-effecting (network) part of the ground primitives so it
// can be injected. It exposes only Fetch (url → page body) and LookupMX (domain →
// deliverable). The default implementation uses the real net stack; tests inject a
// fake to stay deterministic and network-free. A Snapshot reads external facts only
// through this interface and caches the results.
type Resolver interface {
	// Fetch returns the body of the page at url. A non-nil error means the body
	// could not be obtained (the caller reduces it to a FAIL Fact).
	Fetch(url string) (string, error)
	// LookupMX reports whether domain has at least one reachable MX record. A
	// non-nil error means the lookup could not be completed.
	LookupMX(domain string) (bool, error)
}
