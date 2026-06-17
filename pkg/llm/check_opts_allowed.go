//ff:func feature=llm type=helper control=iteration dimension=1 level=error
//ff:what checkOptsAllowed — present(실제로 주어진 쿼리 키 집합)의 모든 키가 백엔드 허용 키(allowed 가변인자) 안에 드는지 검사한다. 하나라도 벗어나면 fmt.Errorf로 거부(loud, no silent caps). allowed가 비면(subprocess 백엔드) present에 키가 있는 즉시 에러.

package llm

import "fmt"

// checkOptsAllowed reports an error if present contains any key not in allowed.
// An empty allowed set (subprocess backends) rejects any present key.
func checkOptsAllowed(present map[string]bool, allowed ...string) error {
	allow := make(map[string]bool, len(allowed))
	for _, k := range allowed {
		allow[k] = true
	}
	for k := range present {
		if !allow[k] {
			return fmt.Errorf("option %q is not supported by this backend", k)
		}
	}
	return nil
}
