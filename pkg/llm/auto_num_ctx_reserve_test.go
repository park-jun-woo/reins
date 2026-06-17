//ff:func feature=llm type=helper control=iteration dimension=1
//ff:what TestAutoNumCtxReserve — reserve 인자가 컨텍스트 창 산정에 반영되는지 검증: reserve 0은 numPredict(2048) 폴백, 큰 reserve는 같은 프롬프트에서 num_ctx 계단을 상향시킨다(빈 프롬프트 기준 0⇒2048, 8192⇒8192, 16384⇒16384).

package llm

import "testing"

// TestAutoNumCtxReserve: reserve 0 falls back to numPredict (2048); a larger reserve
// snaps the context window up for the same prompt.
func TestAutoNumCtxReserve(t *testing.T) {
	cases := []struct {
		reserve int
		want    int
	}{
		{0, 2048},
		{8192, 8192},
		{16384, 16384},
	}
	for _, c := range cases {
		got := autoNumCtx("", c.reserve)
		if got != c.want {
			t.Fatalf("autoNumCtx(\"\", %d) = %d, want %d", c.reserve, got, c.want)
		}
	}
}
