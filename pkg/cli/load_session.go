//ff:func feature=cli type=helper control=sequence level=error
//ff:what sessionPath의 세션을 로드한다. 파일이 없으면 빈 세션(quest.New())을 새로 만들어 반환한다 — 첫 scan을 위해 부재를 에러로 보지 않는다. 로드 직후 MetaLoop를 무조건 삭제 — in-process 전용 신호라 loop 프로세스가 kill돼 플래그가 박힌 채 남아도 다음 프로세스가 자가 치유한다(Render의 실패 로그 tail 영구 억제 방지).

package cli

import (
	"os"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// loadSession returns the session at path, creating a fresh one if absent. It
// unconditionally clears the MetaLoop flag right after loading: the flag is an
// in-process-only signal, so any residue left by a killed loop process is
// self-healed by the next process (otherwise Render would suppress its failure
// log-tail forever).
func loadSession(path string) (*quest.Session, error) {
	s, err := quest.Load(path)
	if os.IsNotExist(err) {
		return quest.New(), nil
	}
	if err != nil {
		return nil, err
	}
	delete(s.Meta, quest.MetaLoop)
	return s, nil
}
