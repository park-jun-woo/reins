//ff:func feature=llm type=helper control=iteration dimension=1
//ff:what argvIndex — argv 슬라이스에서 주어진 토큰의 첫 인덱스를 반환하고, 없으면 -1을 반환하는 테스트 헬퍼. 플래그 배치 순서(예: -s가 resume 앞) 단언에 쓴다.

package llm

// argvIndex returns the index of the first occurrence of tok in argv, or -1 if
// absent. Used to assert flag ordering (e.g. -s before resume).
func argvIndex(argv []string, tok string) int {
	for i, a := range argv {
		if a == tok {
			return i
		}
	}
	return -1
}
