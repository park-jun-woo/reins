//ff:type feature=gate type=model
//ff:what 퀘스트 도메인 계약(4메서드). Seed(입력→초기 TODO), Render(작성 프롬프트+검증 컨텍스트), Prepare(제출 디코드→Context; short≠nil이면 게이트 단락), Rules(위반-규칙 카탈로그). 이 인터페이스만 구현하면 reins가 래칫·명령·집계·export를 공급한다.

package gate

import "github.com/park-jun-woo/reins/pkg/quest"

// Definition is the per-quest domain contract. Implement these four and reins
// supplies the ratchet, command skeleton, aggregation, and export.
type Definition interface {
	// Seed creates the initial TODO items from CLI args (files, dirs, a stream).
	Seed(args []string) ([]*quest.Item, error)
	// Render returns the authoring prompt + verification context shown by `next`.
	Render(it *quest.Item) (string, error)
	// Prepare decodes a raw submission into an evaluation Context. A non-nil short
	// verdict short-circuits the gate (e.g. SKIPPED when the source is untrusted).
	Prepare(it *quest.Item, raw []byte) (ctx Context, short *quest.Verdict, err error)
	// Rules is the gate's violation-rule catalog.
	Rules() []Rule
}
