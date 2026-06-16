//ff:func feature=llm type=adapter control=selection
//ff:what codexArgv — 세션 상태(kind·sid)에 따라 `codex exec` argv를 3분기로 짓는다. Stateless: `exec --ephemeral`+flagsExec+flagsCommon(매 호출 독립, 세션 미영속). Continue 1차(sid==""): `exec`+flagsExec+flagsCommon(--ephemeral 없음 — 다음 resume이 찾을 수 있게 영속). Continue 2차+(sid!=""): `exec`+flagsExec+`resume <sid>`+flagsCommon — flagsExec(-s read-only)를 resume "앞"(exec 레벨)에 둔다(실측: codex exec resume는 -s 미수용, --help 확인). 모든 분기 끝에 "-"(stdin 프롬프트). flagsExec=exec 전용(-s read-only), flagsCommon=exec·resume 공통(-m?/--json/--skip-git-repo-check/--ignore-user-config/--ignore-rules).

package llm

// codexArgv builds the `codex exec` argv in one of three shapes from the session
// state. The read-only sandbox flag (flagsExec, "-s read-only") is exec-level only —
// `codex exec resume` does not accept -s (verified via --help) — so on the resume
// branch it is placed BEFORE the resume subcommand to keep the L0 side-effect block
// on the Continue path. (TODO: if exec-level -s turns out not to propagate to a
// resumed session's sandbox, fall back to `-c sandbox_mode="read-only"`, which resume
// does accept; "just drop -s" is not an option.)
func codexArgv(kind sessionKind, sid string, flagsExec, flagsCommon []string) []string {
	argv := []string{"exec"}
	switch {
	case kind == Continue && sid != "":
		// resume subcommand: -s lives at exec level (resume rejects it); read-only kept.
		argv = append(argv, flagsExec...)
		argv = append(argv, "resume", sid)
		argv = append(argv, flagsCommon...)
	case kind == Continue:
		// Continue first call: must persist so a later resume finds it → no --ephemeral.
		argv = append(argv, flagsExec...)
		argv = append(argv, flagsCommon...)
	default:
		// Stateless (default): ephemeral, no session files written.
		argv = append(argv, "--ephemeral")
		argv = append(argv, flagsExec...)
		argv = append(argv, flagsCommon...)
	}
	return append(argv, "-")
}
