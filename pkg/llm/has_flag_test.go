//ff:func feature=llm type=helper control=iteration dimension=1
//ff:what hasFlag — argv 슬라이스에 주어진 flag 토큰이 포함됐는지 반환하는 테스트 헬퍼.

package llm

// hasFlag reports whether argv contains the flag token.
func hasFlag(argv []string, flag string) bool {
	for _, a := range argv {
		if a == flag {
			return true
		}
	}
	return false
}
