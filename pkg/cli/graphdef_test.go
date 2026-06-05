//ff:type feature=cli type=model
//ff:what 테스트용 그래프형 gate.Definition. gate.Evaluator를 구현(Evaluate가 pkg/graph 그래프로 판독)해 submit 배선이 Rules() 대신 Evaluator 경로를 타는지 검증한다. Rules()도 카탈로그를 제공(rules 명령·후방호환 감사)하나 submit은 Evaluate를 쓴다. 제출물 "bad"면 FAIL 카운터 발동. 테스트 바이너리가 toulmin을 링크하는 건 go.work가 해결(소스 의존 아님).

package cli

// graphDef is a graph-backed Definition: it implements gate.Evaluator so submit
// reads its Verdict from a pkg/graph defeat graph instead of the Rules() catalog.
// Rules() still returns the catalog (for the `rules` command and back-compat audit),
// but submit uses Evaluate. The single FAIL counter fires when the submission is
// "bad".
type graphDef struct{}
