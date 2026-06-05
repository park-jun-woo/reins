//ff:type feature=graph type=model
//ff:what supersessionEvent — 평가 중 흡수된 곁가지 1건의 기록. absorbed는 superseded되어 집계에서 제외된 활성 카운터, by는 그를 흡수한 활성 상류 카운터. 공략집 렌더의 "X에 의해 superseded → 곁가지" 한 줄 입력.

package graph

// supersessionEvent records one side-branch absorbed during evaluation: absorbed is
// the active counter that was superseded (excluded from aggregation), by is the
// active upstream counter that absorbed it. Both carry their Fact so the walkthrough
// can label them by gate rule ID. It feeds the "superseded by X → side-branch" line.
type supersessionEvent struct {
	absorbed activeCounter
	by       activeCounter
}
