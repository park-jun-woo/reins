//ff:type feature=quest type=model
//ff:what 테스트용 quest.Sink. n번째 Emit(1-based)에서 에러를 내 Export의 에러 경로(중단·미방출 보존)를 검증하게 한다.

package quest

// failSink fails on the nth Emit (1-based) to exercise the error path.
type failSink struct {
	calls  int
	failAt int
}
