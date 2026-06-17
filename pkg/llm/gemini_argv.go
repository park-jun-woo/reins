//ff:func feature=llm type=adapter control=selection
//ff:what geminiArgv — 세션 상태(kind·curSID)에 따라 `gemini -p` argv를 3분기로 짓는다. Stateless: opt 그대로(매 호출 독립, 세션 옵션 미지정). Continue 1차(curSID==""): opt + `--session-id <newSID>` — reins가 발급한 UUID로 새 세션 생성(다음 resume이 찾도록 영속). Continue 2차+(curSID!=""): opt + `--resume latest` — 같은 cwd의 직전(=우리) 세션 이어감(loop는 순차·단일 인스턴스). gemini 봉투에 session_id가 없고 --resume이 latest|index라, reins가 UUID를 핀해 1차 세션을 만든다. yolo/--all-files는 절대 미사용(opt에 없음).

package llm

// geminiArgv builds the `gemini -p` argv in one of three shapes from the session
// state. curSID is the carried (already-issued) session UUID — empty on the first
// Continue call, non-empty thereafter. newSID is the freshly issued UUID used only on
// that first Continue call (--session-id). On later calls reins resumes its own latest
// session (--resume latest); the gemini json envelope carries no session_id and
// --resume takes latest|index, so reins pins the UUID itself.
func geminiArgv(kind sessionKind, curSID, newSID string, opt []string) []string {
	argv := append([]string{}, opt...)
	switch {
	case kind == Continue && curSID != "":
		// Continue 2nd+: resume our own latest session in this cwd.
		argv = append(argv, "--resume", "latest")
	case kind == Continue:
		// Continue first call: create the session under the reins-issued UUID.
		argv = append(argv, "--session-id", newSID)
	}
	return argv
}
