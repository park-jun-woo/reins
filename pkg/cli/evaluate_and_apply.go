//ff:func feature=cli type=helper control=sequence level=error
//ff:what evaluateAndApply — submit·agent 공용 헬퍼. def.Prepare(s,it,raw)→(short verdict, def가 gate.Evaluator면 ev.Evaluate(그래프), 아니면 gate.Evaluate(Rules))→quest.Apply(UTC RFC3339)→Save→quest.Export→Save 후 verdict 반환. 게이트가 PASS를 잠그는 단일 지점. submit과 agent가 같은 판정·래칫·export 경로를 공유한다(DRY).

package cli

import (
	"time"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// evaluateAndApply runs the gate over one submission, applies the ratchet
// transition, persists the session, and exports terminal items. It is the single
// place that locks PASS, shared by submit and the agent loop.
func evaluateAndApply(def gate.Definition, s *quest.Session, it *quest.Item, raw []byte, outPath, sessionPath string) (quest.Verdict, error) {
	ctx, short, err := def.Prepare(s, it, raw)
	if err != nil {
		return quest.Verdict{}, err
	}
	var verdict quest.Verdict
	if short != nil {
		verdict = *short
	} else if ev, ok := def.(gate.Evaluator); ok {
		verdict = ev.Evaluate(ctx)
	} else {
		verdict = gate.Evaluate(def.Rules(), ctx)
	}
	now := time.Now().UTC().Format(time.RFC3339)
	quest.Apply(it, verdict, now)
	if err := s.Save(sessionPath); err != nil {
		return verdict, err
	}
	sink, err := newJSONLSink(outPath)
	if err != nil {
		return verdict, err
	}
	if _, err := quest.Export(s, sink); err != nil {
		return verdict, err
	}
	if err := s.Save(sessionPath); err != nil {
		return verdict, err
	}
	return verdict, nil
}
