//ff:func feature=cli type=helper control=sequence
//ff:what backendErrorVerdict — L0 생성 오류(backend.Complete 실패)를 재시도 가능한 FAIL로 감싼다. 예약 규칙 "backend-error" 아래 Fact 한 건에 원문 에러 텍스트를 실어, 내용 비평을 지어내지 않고도 run 로그가 '왜' 시도가 실패했는지 기록한다. RootCause=backend-error로 게이트 실패와 생성 실패를 구분 가능.

package cli

import "github.com/park-jun-woo/reins/pkg/quest"

// backendErrorRule is the reserved synthetic rule ID marking a generation-stage
// failure (an L0 backend.Complete error), distinct from any consumer rule a gate
// might fire. RootCause==backendErrorRule on a Verdict means "the generator failed",
// not "the gate rejected the content".
const backendErrorRule = "backend-error"

// backendErrorVerdict wraps an L0 generation error as a retryable FAIL: a single
// Fact under the reserved rule "backend-error" carrying the error text, so the run
// log records why the attempt failed without inventing a content critique.
func backendErrorVerdict(err error) quest.Verdict {
	return quest.Verdict{
		Outcome:   quest.OutFail,
		RootCause: backendErrorRule,
		Facts: []quest.Fact{{
			Rule:   backendErrorRule,
			Where:  "backend.Complete",
			Actual: err.Error(),
		}},
	}
}
