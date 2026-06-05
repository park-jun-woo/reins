//ff:func feature=ground type=helper control=selection
//ff:what fakeResolver.LookupMX — 호출 횟수를 1 올리고, 해당 도메인에 준비된 에러가 있으면 그걸, 없으면 준비된 bool을 돌려준다(네트워크 없음). 스냅샷 MX 캐시·에러 전파 테스트의 결정적 백엔드.

package ground

func (f *fakeResolver) LookupMX(domain string) (bool, error) {
	f.lookupCalls++
	if e := f.mxErr[domain]; e != nil {
		return false, e
	}
	return f.mx[domain], nil
}
