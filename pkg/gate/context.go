//ff:type feature=gate type=model
//ff:what 규칙이 검사하는 제출 1건의 사실 운반체. Item=대상 퀘스트, Submission=디코드된 도메인 산출물, Source=캐시된 원천(치즈방어 규칙이 재확인), Grounds=staged 평가에서 tier1 직전 lazy resolve된 네트워크 ground 값(키=ground 이름, 예 "source-body"). toulmin-free·generic — 규칙 Check가 ctx.Grounds["source-body"]로 읽는다.

package gate

import "github.com/park-jun-woo/reins/pkg/quest"

// Context carries the per-submission facts a rule inspects. Source is the cached
// ground truth (re-confirmed by cheese-defense rules); Submission is the decoded
// domain artifact. Grounds holds the network ground values lazily resolved just
// before tier-1 evaluation (keyed by ground name, e.g. "source-body"); it is empty
// for tier-0 rules (which declare no ground deps) and for the flat gate.Evaluate
// path. It is generic and toulmin-free — a rule's Check reads ctx.Grounds["..."].
type Context struct {
	Item       *quest.Item
	Submission any
	Source     string
	Grounds    map[string]string
}
