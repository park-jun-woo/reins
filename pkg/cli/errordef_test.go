//ff:type feature=cli type=model
//ff:what 테스트용 gate.Definition. seed/render/prepare 플래그가 켜지면 해당 단계가 실패해 명령의 에러 분기(level=error)를 자극한다. 플래그가 꺼지면 stubDef처럼 동작해 happy path도 배선된다.

package cli

// errDef is a Definition whose Seed/Render/Prepare each fail when the corresponding
// flag is set, so the command error branches (level=error) can be exercised. When a
// flag is unset it behaves like stubDef so the happy path still wires up.
type errDef struct {
	seedErr    bool
	renderErr  bool
	prepareErr bool
}
