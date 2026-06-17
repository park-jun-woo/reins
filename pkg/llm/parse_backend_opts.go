//ff:func feature=llm type=helper control=iteration dimension=1 level=error
//ff:what parseBackendOpts — `--model` 쿼리(`k=v&…`)를 backendOpts로 파싱한다. '&'로 자른 각 세그먼트를 첫 '='로 key|val 분리(='없으면 에러), applyBackendOpt로 타입드 반영, present 집합에 키 기록. raw가 빈 문자열이면 빈 opts(present 비어 있음)를 반환해 현행 동치를 보존한다.

package llm

import (
	"fmt"
	"strings"
)

// parseBackendOpts parses a `k=v&…` query into typed opts plus the present-key set.
// An empty raw yields empty opts (no keys present), preserving pre-Phase017 behavior.
func parseBackendOpts(raw string) (backendOpts, error) {
	opts := backendOpts{present: map[string]bool{}}
	if raw == "" {
		return opts, nil
	}
	for _, pair := range strings.Split(raw, "&") {
		key, val, ok := strings.Cut(pair, "=")
		if !ok {
			return opts, fmt.Errorf("invalid query segment %q: expected k=v", pair)
		}
		if err := applyBackendOpt(&opts, key, val); err != nil {
			return opts, err
		}
		opts.present[key] = true
	}
	return opts, nil
}
