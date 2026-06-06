//ff:func feature=cli type=helper control=sequence level=error
//ff:what backendErr.Error — 래핑한 문자열을 error 메시지로 반환한다(error 인터페이스 만족).

package cli

// Error satisfies the error interface for backendErr.
func (e backendErr) Error() string { return string(e) }
