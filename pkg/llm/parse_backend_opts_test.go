//ff:func feature=llm type=helper control=iteration dimension=1
//ff:what TestParseBackendOpts — 빈 raw는 빈 present, 4개 canonical 키(max_output_tokens·num_ctx·temperature·think)가 타입드 필드와 present 집합에 모두 반영되는지 검증.

package llm

import "testing"

// TestParseBackendOpts: empty raw yields no present keys; all four canonical keys
// populate both the typed fields and the present set.
func TestParseBackendOpts(t *testing.T) {
	empty, err := parseBackendOpts("")
	if err != nil {
		t.Fatalf("parseBackendOpts(\"\") error: %v", err)
	}
	if len(empty.present) != 0 {
		t.Fatalf("empty raw present = %v, want none", empty.present)
	}

	opts, err := parseBackendOpts("max_output_tokens=8192&num_ctx=16384&temperature=0.5&think=false")
	if err != nil {
		t.Fatalf("parseBackendOpts error: %v", err)
	}
	if opts.MaxOutputTokens != 8192 {
		t.Fatalf("MaxOutputTokens = %d, want 8192", opts.MaxOutputTokens)
	}
	if opts.NumCtx != 16384 {
		t.Fatalf("NumCtx = %d, want 16384", opts.NumCtx)
	}
	if opts.Temperature == nil || *opts.Temperature != 0.5 {
		t.Fatalf("Temperature = %v, want 0.5", opts.Temperature)
	}
	if opts.Think == nil || *opts.Think != false {
		t.Fatalf("Think = %v, want false", opts.Think)
	}
	for _, k := range []string{"max_output_tokens", "num_ctx", "temperature", "think"} {
		if !opts.present[k] {
			t.Fatalf("present missing %q: %v", k, opts.present)
		}
	}
}
