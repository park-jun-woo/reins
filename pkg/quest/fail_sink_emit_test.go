//ff:func feature=quest type=helper control=sequence level=error
//ff:what failSink.Emit. failAt번째 호출에서 에러를 내고 그 외엔 nil을 돌려준다(에러 경로 자극용).

package quest

import "errors"

func (f *failSink) Emit(it *Item) error {
	f.calls++
	if f.calls == f.failAt {
		return errors.New("sink boom")
	}
	return nil
}
