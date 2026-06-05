//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what 주어진 규칙들의 RuleMeta 목록 반환 — cli `rules`가 출력하는 자동 rulebook(막는 치즈 목록, ID·레벨).

package gate

// Catalog returns the metas of the given rules — the auto-generated rulebook that
// the cli `rules` command prints (every cheese this gate blocks, by ID and level).
func Catalog(rules []Rule) []RuleMeta {
	metas := make([]RuleMeta, 0, len(rules))
	for _, r := range rules {
		metas = append(metas, r.Meta)
	}
	return metas
}
