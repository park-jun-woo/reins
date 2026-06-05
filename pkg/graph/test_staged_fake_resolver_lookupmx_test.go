//ff:func feature=graph type=helper control=sequence
//ff:what stagedFakeResolver.LookupMX — 어떤 도메인에도 (true, nil)을 돌려준다(네트워크 없음). ground.Resolver 인터페이스 충족용 스텁.

package graph

func (r *stagedFakeResolver) LookupMX(domain string) (bool, error) {
	return true, nil
}
