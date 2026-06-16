//ff:func feature=llm type=adapter control=iteration dimension=1 level=error
//ff:what parseCodexJSONL — `codex exec --json`의 JSONL 이벤트 스트림(줄당 1 JSON)을 훑어 최종 assistant 메시지 텍스트와 세션 id를 뽑는다. 실측 스키마(codex-cli 0.139.0): `{"type":"thread.started","thread_id":...}`(세션/conversation id), `{"type":"item.completed","item":{"type":"agent_message","text":...}}`(어시스턴트 메시지 — 여러 개면 마지막 것이 최종), `{"type":"error",...}`(에러 이벤트). 파싱 불가 줄은 건너뛰고(스트림 노이즈 허용), error 이벤트는 즉시 에러, agent_message가 하나도 없으면 에러. thread_id는 Continue resume에 쓰이도록 sid로 운반.

package llm

import (
	"encoding/json"
	"fmt"
	"strings"
)

// parseCodexJSONL scans the `codex exec --json` JSONL event stream (one JSON object
// per line) for the final assistant message text and the session id. Measured schema
// (codex-cli 0.139.0): a "thread.started" event carries the thread_id; an
// "item.completed" event whose item.type is "agent_message" carries the text (the
// last such event wins); an "error" event signals failure. Unparseable lines are
// skipped (stream noise tolerated); the absence of any agent_message is an error.
func parseCodexJSONL(stdout string) (text, sid string, err error) {
	var found bool
	for _, line := range strings.Split(stdout, "\n") {
		var ev struct {
			Type     string `json:"type"`
			ThreadID string `json:"thread_id"`
			Item     struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"item"`
		}
		if json.Unmarshal([]byte(line), &ev) != nil {
			continue
		}
		if ev.ThreadID != "" {
			sid = ev.ThreadID
		}
		if ev.Type == "error" {
			return "", "", fmt.Errorf("codex error event: %s", strings.TrimSpace(line))
		}
		if ev.Type == "item.completed" && ev.Item.Type == "agent_message" {
			text, found = ev.Item.Text, true
		}
	}
	if !found {
		return "", "", fmt.Errorf("codex: no assistant message in JSONL output")
	}
	return text, sid, nil
}
