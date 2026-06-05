//ff:func feature=ground type=helper control=sequence
//ff:what defaultResolver.Fetch 2xx 성공 경로 증명 — httptest 루프백 서버(외부 네트워크·DNS 없음)가 200+본문을 주면 Fetch가 그 본문 문자열을 err=nil로 반환하는지 결정적으로 커버한다. 실제 외부 호출 없음.

package ground

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultResolverFetchSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello-body"))
	}))
	defer srv.Close()

	r := newDefaultResolver()
	body, err := r.Fetch(srv.URL)
	if err != nil {
		t.Fatalf("Fetch success: err=%v want nil", err)
	}
	if body != "hello-body" {
		t.Fatalf("Fetch success: body=%q want %q", body, "hello-body")
	}
}
