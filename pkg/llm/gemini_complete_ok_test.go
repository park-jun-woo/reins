//ff:func feature=llm type=adapter control=sequence
//ff:what TestGeminiCompleteOK — 200 응답이 candidates[0].parts[0].text를 반환하고, 요청이 system+user를 단일 user 턴으로 병합하며 키가 query에 실리는지 httptest로 검증(무네트워크).

package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestGeminiCompleteOK: a 200 response returns the first candidate's text; the
// request merges system + user into one user turn and the key rides on the query.
func TestGeminiCompleteOK(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "gk-1")
	var gotBody map[string]any
	var gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.Query().Get("key")
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &gotBody)
		io.WriteString(w, `{"candidates":[{"content":{"parts":[{"text":"reply"}]}}]}`)
	}))
	defer srv.Close()

	g := Gemini{Model: "gemini-1.5-pro", BaseURL: srv.URL}
	got, err := g.Complete("SYS", "USR")
	if err != nil {
		t.Fatalf("Complete error: %v", err)
	}
	if got != "reply" {
		t.Fatalf("text = %q, want reply", got)
	}
	if gotQuery != "gk-1" {
		t.Fatalf("key query = %q, want gk-1", gotQuery)
	}
	if !strings.Contains(gotPath, "gemini-1.5-pro:generateContent") {
		t.Fatalf("path = %q", gotPath)
	}
	// system + user merged into one user turn.
	contents := gotBody["contents"].([]any)
	first := contents[0].(map[string]any)
	parts := first["parts"].([]any)
	text := parts[0].(map[string]any)["text"].(string)
	if !strings.Contains(text, "SYS") || !strings.Contains(text, "USR") {
		t.Fatalf("merged text = %q, want SYS and USR", text)
	}
}
