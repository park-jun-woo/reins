//ff:type feature=ground type=model
//ff:what defaultResolver — 실제 네트워크를 쓰는 기본 Resolver. Fetch는 http.Get으로 본문을 읽고(2xx만 성공), LookupMX는 net.LookupMX로 MX 레코드 존재를 본다. 타임아웃 한도를 가진 http.Client를 보유. 테스트는 이걸 쓰지 않고 fake를 주입한다(네트워크 금지). NewSnapshot(nil)이 이 resolver를 기본 채택.

package ground

import "net/http"

// defaultResolver is the Resolver backed by the real network: Fetch uses http.Get
// (2xx only counts as success) and LookupMX uses net.LookupMX to see whether the
// domain has any MX record. It holds an http.Client with a timeout bound. Tests do
// not use this — they inject a fake to avoid real network. NewSnapshot(nil) adopts
// this resolver by default.
type defaultResolver struct {
	client *http.Client
}
