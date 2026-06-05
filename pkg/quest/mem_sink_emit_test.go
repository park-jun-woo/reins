//ff:func feature=quest type=helper control=sequence
//ff:what memSink.Emit. 방출된 아이템의 Key를 기록한다(테스트 검증용).

package quest

func (m *memSink) Emit(it *Item) error {
	m.keys = append(m.keys, it.Key)
	return nil
}
