//ff:func feature=ground type=helper control=selection
//ff:what fakeResolver.Fetch — 호출 횟수를 1 올리고, 해당 url에 준비된 에러가 있으면 그걸, 없으면 준비된 본문을 돌려준다(네트워크 없음). 스냅샷 캐시·에러 전파 테스트의 결정적 백엔드.

package ground

func (f *fakeResolver) Fetch(url string) (string, error) {
	f.fetchCalls++
	if e := f.fetchErr[url]; e != nil {
		return "", e
	}
	return f.bodies[url], nil
}
