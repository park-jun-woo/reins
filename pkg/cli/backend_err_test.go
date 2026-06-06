//ff:type feature=cli type=model
//ff:what 테스트용 에러 타입. backendErr는 문자열을 error로 감싸 스텁 Backend가 반환할 에러(errBackend)를 만든다.

package cli

// backendErr is a string-typed error for stub backends to return.
type backendErr string

// errBackend is the canonical stub backend error.
var errBackend = backendErr("boom")
