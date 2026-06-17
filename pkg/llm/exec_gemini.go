//ff:type feature=llm type=adapter
//ff:what execGemini — GeminiCLI.Complete가 쓰는 서브프로세스 seam(패키지 변수). 본문은 공유 runSubprocess(run_subprocess.go)를 가리킨다. execClaude/execGrok/execCodex와 분리된 별도 변수라 gemini 테스트가 다른 백엔드를 건드리지 않고 독립적으로 스텁을 주입해 진짜 프로세스 없이 Complete를 검증할 수 있다.

package llm

// execGemini is the subprocess seam used by GeminiCLI.Complete. The body lives in the
// shared runSubprocess (run_subprocess.go); this is the per-backend var, kept separate
// from execClaude/execGrok/execCodex so gemini tests can stub it in isolation and
// exercise Complete without spawning a real process.
var execGemini = runSubprocess
