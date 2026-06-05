//ff:func feature=temporal type=helper control=sequence
//ff:what 명세를 그레고리력 ISO Result로 정규화한다(순수, now는 ref 주입). Relative=ref+OffsetDays, Absolute+Gregorian=Start(있으면 End까지) 검증 후 단일/기간, 비그레고리·파싱실패=Determined=false. 역법 변환은 v1 미지원(v2).

package temporal

import "time"

// isoLayout is the canonical Gregorian date layout (ISO 8601 calendar date).
const isoLayout = "2006-01-02"

// Resolve normalizes a Spec to a Gregorian ISO Result. It is pure and deterministic
// (now is injected via ref):
//
//   - Relative: ref.AddDate(0,0,OffsetDays), formatted "2006-01-02"; Determined=true.
//   - Absolute + Gregorian: Start (and End, if set) are validated as "2006-01-02".
//     A single date yields Value=start; an interval yields Value="start/end",
//     IsInterval=true.
//   - Any non-Gregorian calendar, or an unparseable component, returns
//     Result{Determined:false} — v2 (deferred) would add calendar conversion here;
//     reins does not import calendar-conversion libraries in v1.
func Resolve(spec Spec, ref time.Time) Result {
	if spec.Kind == Relative {
		v := ref.AddDate(0, 0, spec.OffsetDays).Format(isoLayout)
		return Result{Value: v, Determined: true}
	}

	// Absolute. v1 only converts the Gregorian calendar.
	if spec.Calendar != Gregorian {
		return Result{Determined: false}
	}
	start, ok := parseGregorian(spec.Start)
	if !ok {
		return Result{Determined: false}
	}
	if spec.End == "" {
		return Result{Value: start, Determined: true}
	}
	end, ok := parseGregorian(spec.End)
	if !ok {
		return Result{Determined: false}
	}
	return Result{
		Value:      start + "/" + end,
		IsInterval: true,
		Determined: true,
	}
}
