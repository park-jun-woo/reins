//ff:type feature=llm type=model
//ff:what llmClient — pkg/llm 세 어댑터(ollama/gemini/openai-compat)가 공유하는 단일 http.Client. 300s 타임아웃으로 무인 agent 루프에서 호출 1건이 행걸려 루프 전체가 영구 정지하는 것을 막는다(http.DefaultClient/http.Post 무타임아웃 대체).

package llm

import (
	"net/http"
	"time"
)

// llmTimeout bounds every adapter HTTP call so an unattended agent loop can never
// hang forever on one request. 300s leaves headroom for slow LLM generation.
const llmTimeout = 300 * time.Second

// llmClient is the single HTTP client shared by all pkg/llm adapters
// (ollama/gemini/openai-compat). It replaces the timeout-less http.DefaultClient
// and http.Post.
var llmClient = &http.Client{Timeout: llmTimeout}
