//ff:type feature=cli type=model
//ff:what NewQuestCmd 튜닝 옵션. Out(기본 export 경로, 비면 "<name>-results.jsonl")와 Version(루트 명령 버전 문자열). 둘 다 선택.

package cli

// Options tunes the generated CLI. All fields are optional; zero values fall back to
// reins defaults derived from the quest name.
type Options struct {
	// Out is the default export path (JSONL). Empty ⇒ "<name>-results.jsonl".
	Out string
	// Version is shown as the root command's version string.
	Version string
}
