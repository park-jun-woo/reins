//ff:func feature=ground type=helper control=sequence
//ff:what defaultResolver.Fetch 비-2xx 상태 에러 경로 증명 — httptest 루프백 서버(외부 네트워크 없음)가 500을 주면 Fetch가 에러를 내고 본문은 빈 문자열인지 결정적으로 커버한다.

package ground

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultResolverFetchNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusInternalServerError)
	}))
	defer srv.Close()

	r := newDefaultResolver()
	body, err := r.Fetch(srv.URL)
	if err == nil {
		t.Fatalf("Fetch non-2xx: err=nil want error")
	}
	if body != "" {
		t.Fatalf("Fetch non-2xx: body=%q want empty", body)
	}
}
