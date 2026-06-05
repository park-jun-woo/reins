//ff:func feature=graph type=helper control=sequence
//ff:what stagedFakeResolver.Fetch — 호출 횟수를 1 올리고 어떤 URL에도 미리 정한 본문을 err=nil로 돌려준다(네트워크 없음). staged 평가의 G5·캐시 단언을 위한 결정적 백엔드.

package graph

func (r *stagedFakeResolver) Fetch(url string) (string, error) {
	r.fetchCalls++
	return r.body, nil
}
