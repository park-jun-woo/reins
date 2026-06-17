//ff:func feature=llm type=helper control=selection level=error
//ff:what applyBackendOpt — 단일 canonical 쿼리 키=값을 backendOpts에 타입드로 반영한다. max_output_tokens·num_ctx는 int(strconv.Atoi), temperature는 float64(ParseFloat), think는 bool(ParseBool). 파싱 실패와 미지 키는 에러로 표면화(loud, no silent). parseBackendOpts의 루프 본문에서 분리(루프 본문 10라인 규칙).

package llm

import (
	"fmt"
	"strconv"
)

// applyBackendOpt parses one canonical key=value pair into opts. A parse failure or
// an unrecognized key is surfaced as an error (no silent drop).
func applyBackendOpt(opts *backendOpts, key, val string) error {
	switch key {
	case "max_output_tokens":
		n, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("max_output_tokens %q: %w", val, err)
		}
		opts.MaxOutputTokens = n
	case "num_ctx":
		n, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("num_ctx %q: %w", val, err)
		}
		opts.NumCtx = n
	case "temperature":
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("temperature %q: %w", val, err)
		}
		opts.Temperature = &f
	case "think":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("think %q: %w", val, err)
		}
		opts.Think = &b
	default:
		return fmt.Errorf("unknown option key %q", key)
	}
	return nil
}
