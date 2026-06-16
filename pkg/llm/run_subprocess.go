//ff:func feature=llm type=adapter control=sequence
//ff:what runSubprocess — claude/grok 서브프로세스 백엔드가 공유하는 기본 실행 구현. bin+argv를 exec.CommandContext로 돌리고 stdin을 먹인 뒤 stdout/stderr를 반환한다. execClaude·execGrok 패키지 변수가 이 함수를 가리켜 본문을 DRY化하고, 각 백엔드 변수는 테스트가 독립적으로 스텁 주입할 수 있게 분리 유지(N=2 추출, 동작 불변).

package llm

import (
	"context"
	"os/exec"
	"strings"
)

// runSubprocess is the shared default implementation behind the per-backend exec
// seams (execClaude, execGrok). It runs bin with argv, feeds stdin, and returns
// stdout/stderr. Each backend keeps its own package var pointing here so a test
// can stub one backend in isolation, while the body stays DRY (extracted at N=2).
func runSubprocess(ctx context.Context, bin string, argv []string, stdin string) (stdout, stderr string, err error) {
	cmd := exec.CommandContext(ctx, bin, argv...)
	var out, errBuf strings.Builder
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err = cmd.Run()
	return out.String(), errBuf.String(), err
}
