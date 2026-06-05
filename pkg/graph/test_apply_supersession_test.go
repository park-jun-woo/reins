//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what applySupersession 단위증명 — 활성 상류가 흡수하는 하류 ID를 제외하고 잔존만(등록순) 반환하는지, 비활성(active 목록에 없는) 상류는 흡수 안 함을, 빈 입력·흡수 없음을 테이블로 커버한다.

package graph

import "testing"

func TestApplySupersession(t *testing.T) {
	cases := []struct {
		name   string
		active []string // ids present in active list, in order
		want   []string
	}{
		{"empty", nil, []string{}},
		{"no absorption (holder alone)", []string{"holder"}, []string{"holder"}},
		{"free absorbs holder", []string{"holder", "free"}, []string{"free"}},
		{"fmt absorbs both", []string{"fmt", "holder", "free"}, []string{"fmt"}},
		{"order preserved", []string{"holder", "fmt"}, []string{"fmt"}},
		{"inactive upstream absorbs nothing", []string{"holder"}, []string{"holder"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			runApplySupersessionCase(t, c.active, c.want)
		})
	}
}
