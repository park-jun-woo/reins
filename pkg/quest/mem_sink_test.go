//ff:type feature=quest type=model
//ff:what 테스트용 quest.Sink. 방출된 아이템의 Key를 순서대로 모아 emit-once 동작을 검증하게 한다.

package quest

type memSink struct{ keys []string }
