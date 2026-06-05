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
	start, err := time.Parse(isoLayout, spec.Start)
	if err != nil {
		return Result{Determined: false}
	}
	if spec.End == "" {
		return Result{Value: start.Format(isoLayout), Determined: true}
	}
	end, err := time.Parse(isoLayout, spec.End)
	if err != nil {
		return Result{Determined: false}
	}
	return Result{
		Value:      start.Format(isoLayout) + "/" + end.Format(isoLayout),
		IsInterval: true,
		Determined: true,
	}
}
