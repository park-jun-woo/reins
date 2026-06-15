//ff:func feature=llm type=adapter control=sequence level=error
//ff:what ClaudeCLI.Complete — `claude -p` argv를 빌드(plain claude -p, --bare 미사용)하고 execClaude seam으로 돌려 단일턴 생성을 수행한다. --output-format json·--append-system-prompt system(전체 교체 아님)·--max-turns(MaxTurns||1, 도구 루프 차단)·--permission-mode dontAsk·--tools ""(전 도구 비활성 — 도구 시도→max_turns 에러 방지, L0는 순수 텍스트); Model 있으면 --model, Stateless→--no-session-persistence·Continue→sid 있으면 --resume sid. user는 stdin(긴 본문 arg-length 회피), ctx는 llmTimeout 300s. JSON(result/session_id/is_error/subtype) 파싱: 비정상 종료(stderr 동봉)·파싱 실패·is_error 각각 에러; Continue면 sid 갱신; result 반환.

package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// Complete builds the `claude -p` argv, runs it via execClaude, and parses the
// JSON envelope. In Continue mode the returned session_id is carried into the
// next call's --resume. The user prompt travels on stdin to dodge the OS
// arg-length limit on long bodies.
func (c *ClaudeCLI) Complete(system, user string) (string, error) {
	bin := c.Bin
	if bin == "" {
		bin = "claude"
	}
	maxTurns := c.MaxTurns
	if maxTurns == 0 {
		maxTurns = 1
	}
	argv := []string{
		"-p",
		"--output-format", "json",
		"--append-system-prompt", system,
		"--max-turns", strconv.Itoa(maxTurns),
		"--permission-mode", "dontAsk",
		// Disable all tools: L0 generation is pure text. Without this the model may
		// attempt a tool call, which with --max-turns 1 burns the only turn and exits
		// error_max_turns (stop_reason tool_use) before producing a result — turning a
		// generation into a hard backend error that aborts the loop.
		"--tools", "",
	}
	if c.Model != "" {
		argv = append(argv, "--model", c.Model)
	}
	if c.Session.Kind == Continue {
		if c.sid != "" {
			argv = append(argv, "--resume", c.sid)
		}
	} else {
		argv = append(argv, "--no-session-persistence")
	}

	ctx, cancel := context.WithTimeout(context.Background(), llmTimeout)
	defer cancel()

	stdout, stderr, err := execClaude(ctx, bin, argv, user)
	if err != nil {
		return "", fmt.Errorf("claude exec: %w: %s", err, stderr)
	}
	var r struct {
		Result    string `json:"result"`
		SessionID string `json:"session_id"`
		IsError   bool   `json:"is_error"`
		Subtype   string `json:"subtype"`
	}
	if err := json.Unmarshal([]byte(stdout), &r); err != nil {
		return "", fmt.Errorf("parse claude response: %w", err)
	}
	if r.IsError {
		return "", fmt.Errorf("claude: %s", r.Subtype)
	}
	if c.Session.Kind == Continue {
		c.sid = r.SessionID
	}
	return r.Result, nil
}
