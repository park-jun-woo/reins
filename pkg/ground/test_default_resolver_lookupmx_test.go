//ff:func feature=ground type=helper control=sequence
//ff:what defaultResolver.LookupMX 에러 경로 증명 — 빈/불정 도메인은 net.LookupMX가 네트워크 egress 없이 즉시 에러를 내므로 (false, err) 경로를 결정적으로 커버한다. records>0 / records==0 분기는 실제 DNS가 필요해 네트워크-free로 결정화 불가(소비자가 fake Resolver 주입으로 대체) → 이 두 분기는 미커버.

package ground

import "testing"

func TestDefaultResolverLookupMXError(t *testing.T) {
	r := newDefaultResolver()
	// An empty domain name is rejected by the stdlib resolver before any network
	// I/O, so this exercises the (false, err) branch deterministically.
	ok, err := r.LookupMX("")
	if err == nil {
		t.Fatalf("LookupMX(\"\"): err=nil want error")
	}
	if ok {
		t.Fatalf("LookupMX(\"\"): ok=true want false")
	}
}
