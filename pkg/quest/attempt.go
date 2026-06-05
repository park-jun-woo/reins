//ff:type feature=quest type=model
//ff:what 한 번의 게이트 평가 기록(감사 로그 1줄). Reason은 의견이 아니라 사실 그대로.

package quest

// Attempt is one gate evaluation logged against an item (audit trail).
type Attempt struct {
	Try     int    `json:"try"`
	Outcome string `json:"outcome"`
	Reason  string `json:"reason,omitempty"`
}
