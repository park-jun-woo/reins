//ff:func feature=llm type=adapter control=sequence level=error
//ff:what GrokCLI.Complete — `grok` argv를 빌드해 execGrok seam으로 돌려 단일턴 생성을 수행한다. user는 --prompt-file 임시파일로 전송(grok -p는 프롬프트를 인자값으로 받고 -p -는 stdin을 안 먹어 "-"를 리터럴 프롬프트로 취급 — 긴 본문 arg-length 회피, defer os.Remove). system은 no-tools 프리앰블 선행 후 --rules로 append(전체 교체 아님). --output-format json·--max-turns(MaxTurns||1, 도구 루프 차단)·--permission-mode dontAsk·--tools ""(빈 허용목록=전 내장도구 비활성)·--disable-web-search·--no-subagents·--no-memory. Model 있으면 -m, Continue & sid 있으면 -r sid(Stateless는 resume 미전달). ctx는 llmTimeout 300s. JSON(text/sessionId/stopReason) 파싱: 비정상 종료(stderr 동봉)·파싱 실패·stopReason!="EndTurn" 각각 에러; Continue면 sid 갱신; text 반환.

package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Complete builds the `grok` argv, runs it via execGrok, and parses the JSON
// envelope. The user prompt travels via a --prompt-file temp file: grok's -p takes
// the prompt as an arg value (and `-p -` does NOT read stdin — it treats "-" as the
// literal prompt), so a temp file both dodges the OS arg-length limit on long
// bodies and is the only reliable carrier. In Continue mode the returned sessionId
// is carried into the next call's -r.
func (g *GrokCLI) Complete(system, user string) (string, error) {
	bin := g.Bin
	if bin == "" {
		bin = "grok"
	}
	maxTurns := g.MaxTurns
	if maxTurns == 0 {
		maxTurns = 1
	}

	pf, err := os.CreateTemp("", "reins-grok-prompt-*.txt")
	if err != nil {
		return "", fmt.Errorf("grok prompt file: %w", err)
	}
	defer os.Remove(pf.Name())
	if _, err := pf.WriteString(user); err != nil {
		pf.Close()
		return "", fmt.Errorf("grok prompt file write: %w", err)
	}
	if err := pf.Close(); err != nil {
		return "", fmt.Errorf("grok prompt file close: %w", err)
	}

	argv := []string{
		"--prompt-file", pf.Name(),
		"--output-format", "json",
		"--max-turns", strconv.Itoa(maxTurns),
		"--permission-mode", "dontAsk",
		// system is appended (not replaced) via --rules, after the no-tools preamble.
		"--rules", withNoToolsPreamble(system),
		// Disable all tools: L0 generation is pure text. --tools "" is an empty
		// allow-list (no built-in tools); web search and subagents are off too. Without
		// this the model may spend its only turn on a tool call and stop short of a
		// result, turning a generation into a hard backend error that aborts the loop.
		"--tools", "",
		"--disable-web-search",
		"--no-subagents",
		"--no-memory",
	}
	if g.Model != "" {
		argv = append(argv, "-m", g.Model)
	}
	if g.Session.Kind == Continue && g.sid != "" {
		argv = append(argv, "-r", g.sid)
	}

	ctx, cancel := context.WithTimeout(context.Background(), llmTimeout)
	defer cancel()

	stdout, stderr, err := execGrok(ctx, bin, argv, "")
	if err != nil {
		return "", fmt.Errorf("grok exec: %w: %s", err, stderr)
	}
	var r struct {
		Text       string `json:"text"`
		SessionID  string `json:"sessionId"`
		StopReason string `json:"stopReason"`
	}
	if err := json.Unmarshal([]byte(stdout), &r); err != nil {
		return "", fmt.Errorf("parse grok response: %w", err)
	}
	// grok signals a clean single-turn completion with stopReason "EndTurn"; any
	// other terminal (Cancelled, MaxTurns from a tool attempt, …) is a failed L0.
	if r.StopReason != "EndTurn" {
		return "", fmt.Errorf("grok: stop reason %s", r.StopReason)
	}
	if g.Session.Kind == Continue {
		g.sid = r.SessionID
	}
	return r.Text, nil
}
