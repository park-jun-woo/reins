//ff:type feature=cli type=model
//ff:what 사소한 테스트용 gate.Definition. arg 하나당 아이템 1건, 제출물이 "bad"면 게이트 FAIL, "skip"이면 Prepare가 SKIPPED로 단락(규칙 카탈로그 우회 경로 자극).

package cli

// stubDef is a trivial Definition: one item per arg. The gate fails when the
// submission is "bad"; a submission of "skip" short-circuits to a SKIPPED verdict
// (exercising the def.Prepare short path that bypasses the rule catalog).
type stubDef struct{}
