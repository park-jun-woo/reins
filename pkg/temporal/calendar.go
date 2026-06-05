//ff:type feature=temporal type=model
//ff:what 절대 명세 성분이 어느 역법인지. v1은 Gregorian만 그레고리력 ISO로 변환하고 나머지(Persian/Islamic/Hebrew/Chinese)는 Determined=false로 정직 반환한다. 역법 변환 라이브러리는 v2(관찰될 때만).

package temporal

// Calendar identifies the calendar an absolute spec's components are stated in. Only
// Gregorian is converted in v1; the others return Determined=false until v2 adds
// conversion tables. reins does not import calendar-conversion libraries in v1.
type Calendar string

const (
	Gregorian Calendar = "gregorian"
	Persian   Calendar = "persian"
	Islamic   Calendar = "islamic"
	Hebrew    Calendar = "hebrew"
	Chinese   Calendar = "chinese"
)
