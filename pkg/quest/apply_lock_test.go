//ff:func feature=quest type=helper control=iteration dimension=1
//ff:what Apply가 REVIEW/SKIPPED/BLOCKED verdict를 대응 terminal 상태로 잠그고 Attempt 로그를 남기는지 검증한다.

package quest

import "testing"

func TestApplyLockVariants(t *testing.T) {
	cases := []struct {
		out  Outcome
		want State
	}{
		{OutReview, REVIEW},
		{OutSkip, SKIPPED},
		{OutBlock, BLOCKED},
	}
	for _, c := range cases {
		it := &Item{Key: "x", State: TODO}
		Apply(it, Verdict{Outcome: c.out}, "now")
		if it.State != c.want || !it.State.Terminal() {
			t.Errorf("%s: state = %s, want locked %s", c.out, it.State, c.want)
		}
		if len(it.Log) != 1 || it.Log[0].Outcome != string(c.out) {
			t.Errorf("%s: log = %+v, want one attempt with outcome %s", c.out, it.Log, c.out)
		}
	}
}
