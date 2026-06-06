//ff:type feature=quest type=model
//ff:what 테스트용 도메인 페이로드. SetPayload/DecodePayload 라운드트립의 무손실성을 비교 가능한 값으로 검증하게 한다.

package quest

type payloadDoc struct {
	URL  string `json:"url"`
	Lang string `json:"lang"`
}
