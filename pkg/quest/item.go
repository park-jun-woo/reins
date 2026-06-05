package quest

// Attempt is one gate evaluation logged against an item (audit trail).
type Attempt struct {
	Try     int    `json:"try"`
	Outcome string `json:"outcome"`
	Reason  string `json:"reason,omitempty"`
}

// Fact is fact-based feedback a gate rule emits when it fires: located and
// quantified, with "no room to flatter" (how-make-quest). On FAIL it is handed back
// to the agent so a sycophantic model converges on a correction instead of arguing.
type Fact struct {
	Rule     string `json:"rule,omitempty"`
	Where    string `json:"where,omitempty"`
	Expected string `json:"expected,omitempty"`
	Actual   string `json:"actual,omitempty"`
}

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
