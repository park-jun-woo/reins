//ff:func feature=ground type=helper control=selection
//ff:what defaultResolver.Fetch — 실제 net 스택으로 url을 http.Get한다. 전송 에러나 비-2xx 상태는 에러로 반환(본문 없음), 2xx면 응답 본문 전체를 문자열로 읽어 반환한다. 테스트는 이 실네트워크 구현 대신 fake를 주입한다.

package ground

import (
	"fmt"
	"io"
)

// Fetch GETs url and returns the response body as a string. A transport error or a
// non-2xx status is returned as an error (no body).
func (r defaultResolver) Fetch(url string) (string, error) {
	resp, err := r.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("ground: fetch %q: status %d", url, resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
