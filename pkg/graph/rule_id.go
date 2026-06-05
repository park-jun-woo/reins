//ff:func feature=graph type=helper control=sequence
//ff:what ruleIDFor — fn 포인터와 idSpec로부터 toulmin이 생성할 ruleID를 동일하게 재현한다(funcID + "#" + spec.String). toulmin ruleID 함수가 unexported라 Evaluate가 trace 엔트리를 노드로 매칭하려면 같은 규칙으로 ID를 미리 계산해 둬야 한다.

package graph

import (
	"fmt"
	"reflect"
	"runtime"
)

// ruleIDFor reproduces the toulmin ruleID for a function plus a single idSpec:
// funcID + "#" + fmt.Sprint(spec). toulmin's ruleID is unexported, so the node
// must precompute the same identifier to match collected trace entries.
func ruleIDFor(fn any, spec idSpec) string {
	ptr := reflect.ValueOf(fn).Pointer()
	id := fmt.Sprintf("unknown_%d", ptr)
	if f := runtime.FuncForPC(ptr); f != nil {
		id = f.Name()
	}
	return id + "#" + fmt.Sprint(spec)
}
