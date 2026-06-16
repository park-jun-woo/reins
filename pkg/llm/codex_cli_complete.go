//ff:func feature=llm type=adapter control=sequence level=error
//ff:what CodexCLI.Complete — flagsExec(-s read-only, exec 전용)·flagsCommon(-m?/--json/--skip-git-repo-check/--ignore-user-config/--ignore-rules)을 짓고 codexArgv로 세션 분기 argv를 만들어 execCodex seam으로 돌린다. codex엔 system 채널이 없어 withNoToolsPreamble(system)+"\n\n"+user를 stdin으로 선결합한다(no-tools 프리앰블로 도구 루프 차단; --max-turns 부재 대체). ctx는 llmTimeout 300s. 비정상 종료(stderr 동봉)·JSONL 파싱 실패(stderr 동봉) 각각 에러. parseCodexJSONL이 뽑은 thread_id는 Continue면 c.sid에 운반(다음 호출 resume에 사용). text 반환.

package llm

import (
	"context"
	"fmt"
)

// Complete builds the codex flags, assembles the session-branched argv via codexArgv,
// runs it through execCodex, and extracts the final assistant message from the JSONL
// event stream. codex has no system channel, so the no-tools preamble and the consumer
// system prompt are prepended to the user prompt and sent on stdin. In Continue mode
// the response's thread_id is carried into the next call's resume.
func (c *CodexCLI) Complete(system, user string) (string, error) {
	bin := c.Bin
	if bin == "" {
		bin = "codex"
	}
	// flagsExec is exec-level only (resume rejects -s); flagsCommon is valid on both
	// exec and resume. -m is prepended only when a model is pinned.
	flagsExec := []string{"-s", "read-only"}
	flagsCommon := []string{"--json", "--skip-git-repo-check", "--ignore-user-config", "--ignore-rules"}
	if c.Model != "" {
		flagsCommon = append([]string{"-m", c.Model}, flagsCommon...)
	}
	argv := codexArgv(c.Session.Kind, c.sid, flagsExec, flagsCommon)
	stdin := withNoToolsPreamble(system) + "\n\n" + user

	ctx, cancel := context.WithTimeout(context.Background(), llmTimeout)
	defer cancel()

	stdout, stderr, err := execCodex(ctx, bin, argv, stdin)
	if err != nil {
		return "", fmt.Errorf("codex exec: %w: %s", err, stderr)
	}
	text, sid, perr := parseCodexJSONL(stdout)
	if perr != nil {
		return "", fmt.Errorf("%w: %s", perr, stderr)
	}
	if c.Session.Kind == Continue {
		c.sid = sid
	}
	return text, nil
}
