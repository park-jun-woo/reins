//ff:func feature=llm type=adapter control=sequence level=error
//ff:what Gemini.Complete — Google Gemini generateContent endpoint 호출. system+user를 단일 user 턴으로 병합(Gemini 규약), temperature 0 고정, maxOutputTokens 2048. API 키는 env GEMINI_API_KEY를 x-goog-api-key 헤더로 전달(URL 쿼리 금지 — *url.Error가 전체 URL을 포함해 키가 로그로 누출되는 것을 차단). 공용 llmClient(300s 타임아웃) 사용. candidates[0].content.parts[0].text 반환.

package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Complete merges system + user into one user turn and returns the first candidate.
// The API key travels in the x-goog-api-key header, never in the URL: a transport
// failure returns a *url.Error that embeds the full URL, so a query-string key
// would leak into error strings and logs.
func (g Gemini) Complete(system, user string) (string, error) {
	apiKey, err := loadAPIKey("gemini")
	if err != nil {
		return "", err
	}
	base := g.BaseURL
	if base == "" {
		base = defaultGeminiBaseURL
	}
	endpoint := fmt.Sprintf("%s/v1beta/models/%s:generateContent", base, g.Model)

	combined := system + "\n\n" + user
	body := map[string]any{
		"contents": []map[string]any{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": combined},
				},
			},
		},
		"generationConfig": map[string]any{
			"temperature":     0,
			"maxOutputTokens": 2048,
		},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal gemini request: %w", err)
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create gemini request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", apiKey)
	resp, err := llmClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gemini request: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read gemini response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini %d: %s", resp.StatusCode, string(data))
	}
	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("parse gemini response: %w", err)
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini: empty response")
	}
	return result.Candidates[0].Content.Parts[0].Text, nil
}
