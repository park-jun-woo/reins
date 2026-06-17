//ff:func feature=llm type=adapter control=sequence level=error
//ff:what GeminiCLI.Complete — 공통 opt(`-p ""`·-m?·`-o json`·`--approval-mode plan`)를 짓고 geminiArgv로 세션 분기 argv를 만들어 execGemini seam으로 돌린다. gemini엔 system 채널이 없어 withNoToolsPreamble(system)+"\n\n"+user를 stdin으로 선결합한다(-p ""+stdin: --help "Appended to input on stdin"). --approval-mode plan(read-only)+no-tools 프리앰블로 도구 루프 차단(--max-turns 부재 대체). Continue 1차엔 geminiSessionID로 UUID를 발급해 --session-id로 새 세션을 만들고 c.sid에 운반(2차+는 --resume latest). ctx는 llmTimeout 300s. UUID 생성 실패·비정상 종료(stderr 동봉)·봉투 파싱/error 봉투 각각 에러. parseGeminiJSON이 뽑은 response 반환.

package llm

import (
	"context"
	"fmt"
)

// Complete builds the gemini -p flags, assembles the session-branched argv via
// geminiArgv, runs it through execGemini, and extracts the response from the
// `-o json` envelope. gemini has no system channel, so the no-tools preamble and the
// consumer system prompt are prepended to the user prompt and sent on stdin (the
// empty -p value triggers headless mode and stdin is appended to it). In Continue mode
// reins issues a UUID on the first call (--session-id) and resumes its own latest
// session thereafter (--resume latest).
func (c *GeminiCLI) Complete(system, user string) (string, error) {
	bin := c.Bin
	if bin == "" {
		bin = "gemini"
	}
	// -p "" triggers headless mode (the real prompt rides on stdin); -m only when a
	// model is pinned; -o json for the {response,stats,error} envelope;
	// --approval-mode plan is read-only (side-effect block). yolo/--all-files are
	// never used.
	opt := []string{"-p", ""}
	if c.Model != "" {
		opt = append(opt, "-m", c.Model)
	}
	opt = append(opt, "-o", "json", "--approval-mode", "plan")

	// Continue: issue a UUID on the first call so a later --resume latest finds it.
	newSID := c.sid
	if c.Session.Kind == Continue && c.sid == "" {
		id, err := geminiSessionID()
		if err != nil {
			return "", err
		}
		newSID = id
	}
	argv := geminiArgv(c.Session.Kind, c.sid, newSID, opt)
	stdin := withNoToolsPreamble(system) + "\n\n" + user

	ctx, cancel := context.WithTimeout(context.Background(), llmTimeout)
	defer cancel()

	stdout, stderr, err := execGemini(ctx, bin, argv, stdin)
	if err != nil {
		return "", fmt.Errorf("gemini -p: %w: %s", err, stderr)
	}
	text, perr := parseGeminiJSON(stdout)
	if perr != nil {
		return "", fmt.Errorf("%w: %s", perr, stderr)
	}
	if c.Session.Kind == Continue {
		c.sid = newSID
	}
	return text, nil
}
