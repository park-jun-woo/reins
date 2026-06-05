//ff:func feature=cli type=command control=sequence level=error
//ff:what 한 퀘스트의 표준 reins CLI(루트 + scan/next/submit/status/export/rules)를 조립한다. persistent 플래그 --session(기본 session.json)·--out(기본 "<name>-results.jsonl")를 달고, 도메인 로직은 gate.Definition만 끼우면 된다.

package cli

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/spf13/cobra"
)

// NewQuestCmd builds the standard reins CLI for a quest: a root command named name
// with the canonical subcommands scan/next/submit/status/export plus rules. The
// quest's domain logic is supplied by def; reins drives the ratchet, the level
// aggregation, and the export.
//
// Persistent flags: --session (default "session.json") and --out (default
// "<name>-results.jsonl"). Subcommands load the session, mutate it via the pure
// quest core, and save.
func NewQuestCmd(name string, def gate.Definition, opts Options) *cobra.Command {
	defaultOut := opts.Out
	if defaultOut == "" {
		defaultOut = name + "-results.jsonl"
	}

	var (
		sessionPath string
		outPath     string
	)

	root := &cobra.Command{
		Use:           name,
		Short:         name + " — a reins quest CLI",
		Version:       opts.Version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().StringVar(&sessionPath, "session", "session.json", "session state file")
	root.PersistentFlags().StringVar(&outPath, "out", defaultOut, "export output file (JSONL)")

	load := func() (*quest.Session, error) { return loadSession(sessionPath) }

	root.AddCommand(
		newScanCmd(def, &sessionPath, load),
		newNextCmd(def, load),
		newSubmitCmd(def, &sessionPath, &outPath, load),
		newStatusCmd(load),
		newExportCmd(&sessionPath, &outPath, load),
		newRulesCmd(def),
	)
	return root
}
