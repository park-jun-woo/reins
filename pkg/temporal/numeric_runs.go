//ff:func feature=temporal type=helper control=iteration dimension=1
//ff:what 문자열에서 라틴숫자(0-9)의 최대 연속 런들을 추출한다. 선행 0을 보존해 원문 표면형("01"→"01")과 매칭되게 한다. 타 숫자체계(페르시아·아랍·데바나가리·한자)는 매핑표가 필요한 v2.

package temporal

// numericRuns returns the maximal runs of latin digits (0-9) found in s. Leading
// zeros are preserved so the run matches the source surface form ("01" stays "01").
// Other numeral systems need a numeral-mapping table and are v2 (deferred).
func numericRuns(s string) []string {
	var runs []string
	start := -1
	for i, r := range s {
		isDigit := r >= '0' && r <= '9'
		if isDigit && start < 0 {
			start = i
		}
		if !isDigit && start >= 0 {
			runs = append(runs, s[start:i])
			start = -1
		}
	}
	if start >= 0 {
		runs = append(runs, s[start:])
	}
	return runs
}
