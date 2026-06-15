//ff:func feature=llm type=helper control=iteration dimension=1
//ff:what flagValue — argv에서 flag 다음에 오는 토큰을 반환하고, 없으면 빈 문자열을 반환하는 테스트 헬퍼.

package llm

// flagValue returns the token following flag in argv, or "" if absent.
func flagValue(argv []string, flag string) string {
	for i, a := range argv {
		if a == flag && i+1 < len(argv) {
			return argv[i+1]
		}
	}
	return ""
}
