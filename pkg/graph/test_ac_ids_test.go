//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what acIDs 테스트헬퍼 — activeCounter 슬라이스를 노드 id 슬라이스로 등록순 그대로 투영한다(applySupersession 결과 비교용).

package graph

func acIDs(acs []activeCounter) []string {
	out := make([]string, 0, len(acs))
	for _, ac := range acs {
		out = append(out, ac.node.id)
	}
	return out
}
