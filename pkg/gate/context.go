//ff:type feature=gate type=model
//ff:what 규칙이 검사하는 제출 1건의 사실 운반체. Item=대상 퀘스트, Submission=디코드된 도메인 산출물, Source=캐시된 원천(치즈방어 규칙이 재확인).

package gate

import "github.com/park-jun-woo/reins/pkg/quest"

// Context carries the per-submission facts a rule inspects. Source is the cached
// ground truth (re-confirmed by cheese-defense rules); Submission is the decoded
// domain artifact.
type Context struct {
	Item       *quest.Item
	Submission any
	Source     string
}
