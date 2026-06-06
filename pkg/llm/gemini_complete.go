//ff:func feature=llm type=adapter control=sequence level=error
//ff:what Gemini.Complete — Google Gemini generateContent endpoint 호출. system+user를 단일 user 턴으로 병합(Gemini 규약), temperature 0 고정, maxOutputTokens 2048. API 키는 env GEMINI_API_KEY. candidates[0].content.parts[0].text 반환.

package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Complete merges system + user into one user turn and returns the first candidate.
func (g Gemini) Complete(system, user string) (string, error) {
	apiKey, err := loadAPIKey("gemini")
	if err != nil {
		return "", err
	}
	base := g.BaseURL
	if base == "" {
		base = defaultGeminiBaseURL
	}
	endpoint := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", base, g.Model, apiKey)

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
	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(payload))
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
