//ff:type feature=llm type=adapter
//ff:what execClaude — ClaudeCLI.Complete가 쓰는 서브프로세스 seam(패키지 변수). 본문은 공유 runSubprocess(run_subprocess.go)로 이관했고 여기서는 그것을 가리킨다. 패키지 변수라 테스트가 스텁을 주입해 진짜 프로세스 없이 Complete를 검증할 수 있다(httptest seam과 동형); execGrok와 분리돼 백엔드별 독립 스텁이 가능.

package llm

// execClaude is the subprocess seam used by ClaudeCLI.Complete. The body lives in
// the shared runSubprocess (run_subprocess.go); this is just the per-backend var
// so tests can inject a stub and exercise Complete without spawning a real process
// (mirrors the httptest seams), independently of execGrok.
var execClaude = runSubprocess
