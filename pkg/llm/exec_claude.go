//ff:type feature=llm type=adapter
//ff:what execClaude — ClaudeCLI.Complete가 쓰는 서브프로세스 seam(패키지 변수). bin+argv를 exec.CommandContext로 돌리고 stdin(user 프롬프트)을 먹인 뒤 stdout/stderr를 반환한다. 패키지 변수라 테스트가 스텁을 주입해 진짜 프로세스 없이 Complete를 검증할 수 있다(httptest seam과 동형).

package llm

import (
	"context"
	"os/exec"
	"strings"
)

// execClaude is the subprocess seam: it runs bin with argv, feeds stdin, and
// returns stdout/stderr. A package var so tests can inject a stub and exercise
// Complete without spawning a real process (mirrors the httptest seams).
var execClaude = func(ctx context.Context, bin string, argv []string, stdin string) (stdout, stderr string, err error) {
	cmd := exec.CommandContext(ctx, bin, argv...)
	var out, errBuf strings.Builder
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err = cmd.Run()
	return out.String(), errBuf.String(), err
}
