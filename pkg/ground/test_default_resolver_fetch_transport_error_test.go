//ff:func feature=ground type=helper control=sequence
//ff:what defaultResolver.Fetch 전송 에러 경로 증명 — httptest 서버를 띄워 URL만 받고 즉시 닫아(connection refused) 전송 에러 경로를 만든다. Fetch가 에러를 내고 본문은 빈 문자열인지 결정적으로 커버한다.

package ground

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultResolverFetchTransportError(t *testing.T) {
	// Spin a server up, capture its URL, then close it → connection refused.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := srv.URL
	srv.Close()

	r := newDefaultResolver()
	body, err := r.Fetch(url)
	if err == nil {
		t.Fatalf("Fetch transport error: err=nil want error")
	}
	if body != "" {
		t.Fatalf("Fetch transport error: body=%q want empty", body)
	}
}
