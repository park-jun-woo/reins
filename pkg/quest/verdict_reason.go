//ff:func feature=quest type=helper control=iteration dimension=1
//ff:what Verdict의 Facts를 감사 로그용 한 줄 문자열로 렌더한다. Facts가 없으면 Outcome만 반환.

package quest

import "strings"

// Reason renders the Facts into a single human-readable line for the audit log.
func (v Verdict) Reason() string {
	if len(v.Facts) == 0 {
		return string(v.Outcome)
	}
	parts := make([]string, 0, len(v.Facts))
	for _, f := range v.Facts {
		seg := f.Where
		if f.Expected != "" || f.Actual != "" {
			seg += " (expected " + f.Expected + ", got " + f.Actual + ")"
		}
		if f.Rule != "" {
			seg = f.Rule + ": " + seg
		}
		parts = append(parts, seg)
	}
	return strings.Join(parts, "; ")
}
