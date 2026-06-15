//ff:func feature=cli type=helper control=sequence level=error
//ff:what applyVerdict — verdict를 아이템 래칫에 적용하고 영속화하는 공용 꼬리. quest.Apply(UTC RFC3339)→s.Save→newJSONLSink→exportAndSave(Export 실패여도 Save로 Emitted 래칫 보존). verdict가 래칫을 바꾸는 단일 지점으로, 게이트 경로(evaluateAndApply)와 백엔드-에러 경로(runLoopItem)가 공유한다. PASS 잠금 권한은 여전히 게이트뿐 — 이 헬퍼는 주어진 verdict를 적용할 뿐이다.

package cli

import (
	"time"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// applyVerdict ratchets a verdict onto an item and persists it: quest.Apply(UTC) →
// Save → exportAndSave (the Emitted ratchet survives an Export failure). It is the
// single place a verdict mutates the ratchet, shared by the gate path
// (evaluateAndApply) and the backend-error path (runLoopItem). It does not lock PASS
// on its own — the caller supplies the verdict, and PASS still originates only at the
// gate.
func applyVerdict(s *quest.Session, it *quest.Item, v quest.Verdict, outPath, sessionPath string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	quest.Apply(it, v, now)
	if err := s.Save(sessionPath); err != nil {
		return err
	}
	sink, err := newJSONLSink(outPath)
	if err != nil {
		return err
	}
	if _, err := exportAndSave(s, sink, sessionPath); err != nil {
		return err
	}
	return nil
}
