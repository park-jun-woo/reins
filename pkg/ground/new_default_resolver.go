//ff:func feature=ground type=helper control=sequence
//ff:what newDefaultResolver — 타임아웃 한도(10s)를 가진 http.Client를 들고 defaultResolver를 만든다. NewSnapshot(nil)이 기본 Resolver로 이걸 채택한다.

package ground

import (
	"net/http"
	"time"
)

// newDefaultResolver builds a defaultResolver with a bounded-timeout http.Client.
func newDefaultResolver() defaultResolver {
	return defaultResolver{client: &http.Client{Timeout: 10 * time.Second}}
}
