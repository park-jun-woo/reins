//ff:func feature=llm type=helper control=sequence
//ff:what TestApplyBackendOpt — 단일 키=값이 backendOpts 타입드 필드에 반영되는지(int·float64·bool 각각) 직접 검증. think=true→*Think true, temperature=0.25→*Temperature 0.25.

package llm

import "testing"

// TestApplyBackendOpt: each canonical key writes its typed field.
func TestApplyBackendOpt(t *testing.T) {
	opts := backendOpts{present: map[string]bool{}}
	if err := applyBackendOpt(&opts, "max_output_tokens", "4096"); err != nil {
		t.Fatalf("max_output_tokens: %v", err)
	}
	if opts.MaxOutputTokens != 4096 {
		t.Fatalf("MaxOutputTokens = %d, want 4096", opts.MaxOutputTokens)
	}
	if err := applyBackendOpt(&opts, "temperature", "0.25"); err != nil {
		t.Fatalf("temperature: %v", err)
	}
	if opts.Temperature == nil || *opts.Temperature != 0.25 {
		t.Fatalf("Temperature = %v, want 0.25", opts.Temperature)
	}
	if err := applyBackendOpt(&opts, "think", "true"); err != nil {
		t.Fatalf("think: %v", err)
	}
	if opts.Think == nil || *opts.Think != true {
		t.Fatalf("Think = %v, want true", opts.Think)
	}
}
