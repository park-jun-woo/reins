//ff:type feature=quest type=model
//ff:what 작업 단위 하나(퀘스트). State는 래칫 위치, Payload는 도메인 산출물, 나머지는 감사 로그.

package quest

// Item is one quest (work unit). Payload carries the domain artifact (e.g. an
// extracted event6). State is the ratchet position; the rest is the audit trail.
type Item struct {
	Key         string    `json:"key"`
	State       State     `json:"state"`
	Tries       int       `json:"tries"`
	Verdict     string    `json:"verdict,omitempty"`
	Reason      string    `json:"reason,omitempty"`
	CollectedAt string    `json:"collected_at,omitempty"`
	Log         []Attempt `json:"log,omitempty"`
	Emitted     bool      `json:"emitted,omitempty"`
	Payload     any       `json:"payload,omitempty"`
}
