//ff:type feature=cli type=model
//ff:what NewQuestCmd 튜닝 옵션. Out(기본 export 경로, 비면 "<name>-results.jsonl")·Version(루트 버전 문자열)·ExtraCommands(루트에 부착할 소비자 서브명령)·Loop(generate→gate→retry 루프 옵트인). 모두 선택.

package cli

import (
	"github.com/spf13/cobra"
)

// Options tunes the generated CLI. All fields are optional; zero values fall back to
// reins defaults derived from the quest name.
type Options struct {
	// Out is the default export path (JSONL). Empty ⇒ "<name>-results.jsonl".
	Out string
	// Version is shown as the root command's version string.
	Version string
	// ExtraCommands are consumer-specific subcommands attached to the root
	// command after the canonical reins subcommands (e.g. ccnews "run" or
	// comail "export"). Nil or empty has no effect, so it is backward-compatible.
	ExtraCommands []*cobra.Command
	// Loop, when non-nil, opts the quest into the `loop` generate→gate→retry
	// loop. Nil ⇒ no loop command is attached.
	Loop *LoopOptions
}
