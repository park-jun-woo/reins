//ff:func feature=llm type=adapter control=sequence level=error
//ff:what doOllamaRequest — body를 JSON 직렬화해 POST하고, 파싱한 message.content를 반환한다. 상태코드 비-200·마샬/언마샬/IO 오류를 명확한 에러로 표면화. Ollama.Complete의 전송·파싱 단계.

package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// doOllamaRequest marshals body, POSTs it, and returns the parsed message.content.
func doOllamaRequest(url string, body any) (string, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal ollama request: %w", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("ollama request: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read ollama response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama %d: %s", resp.StatusCode, string(data))
	}
	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("parse ollama response: %w", err)
	}
	return result.Message.Content, nil
}
