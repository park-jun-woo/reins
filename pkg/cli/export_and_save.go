//ff:func feature=cli type=helper control=sequence level=error
//ff:what exportAndSave — quest.Export 후 결과와 무관하게 세션을 먼저 Save해 Emit 성공분의 Emitted 래칫을 영속화한다(중간 Emit 실패 시 Save 없이 리턴하면 기방출 아이템이 다음 실행에서 중복 방출됨 — emit-once 보존). Export 에러가 Save 에러보다 우선 전파. evaluateAndApply와 export 명령이 공유.

package cli

import "github.com/park-jun-woo/reins/pkg/quest"

// exportAndSave runs quest.Export and then always saves the session, so the
// Emitted ratchet of any successfully emitted item is persisted even when a later
// Emit fails (otherwise the next run would re-emit it, breaking emit-once).
// Partial progress is a ratchet; preserving it is correct. The export error takes
// precedence over a save error. Shared by evaluateAndApply and the export command.
func exportAndSave(s *quest.Session, sink quest.Sink, sessionPath string) (int, error) {
	n, exportErr := quest.Export(s, sink)
	saveErr := s.Save(sessionPath)
	if exportErr != nil {
		return n, exportErr
	}
	return n, saveErr
}
