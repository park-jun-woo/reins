//ff:func feature=llm type=adapter control=sequence level=error
//ff:what OpenAICompat.Complete — OpenAI 호환 chat completions endpoint(xai 등) 호출. URL·Backend(키 조회용 이름)·Model을 갖고 Authorization: Bearer <env key>로 공용 llmClient(300s 타임아웃) POST. temperature 0 고정, max_tokens 2048. choices[0].message.content 반환.

package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Complete posts an OpenAI-compatible chat completion and returns the first choice.
func (o OpenAICompat) Complete(system, user string) (string, error) {
	apiKey, err := loadAPIKey(o.Backend)
	if err != nil {
		return "", err
	}
	body := map[string]any{
		"model": o.Model,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
		"max_tokens":  2048,
		"temperature": 0,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal %s request: %w", o.Backend, err)
	}
	req, err := http.NewRequest("POST", o.URL, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create %s request: %w", o.Backend, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := llmClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s request: %w", o.Backend, err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read %s response: %w", o.Backend, err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s %d: %s", o.Backend, resp.StatusCode, string(data))
	}
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("parse %s response: %w", o.Backend, err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("%s: empty choices", o.Backend)
	}
	return result.Choices[0].Message.Content, nil
}
