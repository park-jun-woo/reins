//ff:type feature=temporal type=model
//ff:what 시간 명세 종류. Absolute=어떤 역법의 절대 날짜, Relative=기준 시각 대비 일수 오프셋. AI가 원문을 읽어 식별한 값(검증기는 이 분류대로 변환·산술만 한다).

package temporal

// Kind is whether a time spec is an absolute date (in some calendar) or relative to
// a reference instant (offset in days).
type Kind string

const (
	Absolute Kind = "absolute"
	Relative Kind = "relative"
)
