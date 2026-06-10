//ff:type feature=quest type=model
//ff:what 세션 Meta 예약 키 상수. MetaLoop는 loop 명령 실행 중임을 Render에 알려, Render가 자체 로그-tail 실패사유 표시를 끄도록 한다(loop가 renderVerdict 피드백을 따로 덧붙이므로 이중노출 방지). 외부 next/submit 수동 흐름에서는 미설정 → Render가 평소대로 tail을 보인다.

package quest

// MetaLoop is the reserved session Meta key the loop command sets (to true) while
// it runs, so a Definition.Render can suppress its own "last failure" log-tail (the
// loop appends renderVerdict feedback itself, avoiding double exposure). It is unset
// in the manual next/submit flow.
const MetaLoop = "reins.loop"
