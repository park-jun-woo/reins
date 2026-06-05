package cli

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/spf13/cobra"
)

// Options tunes the generated CLI. All fields are optional; zero values fall back to
// reins defaults derived from the quest name.
type Options struct {
	// Out is the default export path (JSONL). Empty ⇒ "<name>-results.jsonl".
	Out string
	// Version is shown as the root command's version string.
	Version string
}

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

	// loadOrNew returns the session at sessionPath, creating a fresh one if absent.
	loadOrNew := func() (*quest.Session, error) {
		s, err := quest.Load(sessionPath)
		if os.IsNotExist(err) {
			return quest.New(), nil
		}
		return s, err
	}

	root.AddCommand(
		scanCmd(def, &sessionPath, loadOrNew),
		nextCmd(def, loadOrNew),
		submitCmd(def, &sessionPath, &outPath, loadOrNew),
		statusCmd(loadOrNew),
		exportCmd(&sessionPath, &outPath, loadOrNew),
		rulesCmd(def),
	)
	return root
}

// scanCmd seeds new TODO items from args and saves the session.
func scanCmd(def gate.Definition, sessionPath *string, load func() (*quest.Session, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "scan [args...]",
		Short: "seed quest items from the input",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := load()
			if err != nil {
				return err
			}
			items, err := def.Seed(args)
			if err != nil {
				return err
			}
			s.Items = append(s.Items, items...)
			if err := s.Save(*sessionPath); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "seeded %d item(s); %d total\n", len(items), len(s.Items))
			return nil
		},
	}
}

// nextCmd prints the next TODO item's authoring prompt + verification context.
func nextCmd(def gate.Definition, load func() (*quest.Session, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "next",
		Short: "show the next TODO item (read-only)",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := load()
			if err != nil {
				return err
			}
			it := s.NextTODO()
			if it == nil {
				fmt.Fprintln(cmd.OutOrStdout(), "no TODO items remaining")
				return nil
			}
			prompt, err := def.Render(it)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), prompt)
			return nil
		},
	}
}

// submitCmd evaluates a submission for one item through the gate, applies the ratchet
// transition, exports terminal items, and prints the outcome plus any Facts.
func submitCmd(def gate.Definition, sessionPath, outPath *string, load func() (*quest.Session, error)) *cobra.Command {
	var (
		key    string
		inPath string
	)
	cmd := &cobra.Command{
		Use:   "submit --key <k> [--in <file>|-]",
		Short: "submit an item for gate evaluation",
		RunE: func(cmd *cobra.Command, args []string) error {
			if key == "" {
				return fmt.Errorf("--key is required")
			}
			s, err := load()
			if err != nil {
				return err
			}
			it, err := s.Find(key)
			if err != nil {
				return err
			}
			if it.State != quest.TODO {
				return fmt.Errorf("item %s is %s, not TODO", key, it.State)
			}
			raw, err := readSubmission(cmd, inPath)
			if err != nil {
				return err
			}
			ctx, short, err := def.Prepare(it, raw)
			if err != nil {
				return err
			}
			var verdict quest.Verdict
			if short != nil {
				verdict = *short
			} else {
				verdict = gate.Evaluate(def.Rules(), ctx)
			}
			now := time.Now().UTC().Format(time.RFC3339)
			quest.Apply(it, verdict, now)
			if err := s.Save(*sessionPath); err != nil {
				return err
			}
			sink, err := newJSONLSink(*outPath)
			if err != nil {
				return err
			}
			if _, err := quest.Export(s, sink); err != nil {
				return err
			}
			if err := s.Save(*sessionPath); err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "%s -> %s (state %s)\n", key, verdict.Outcome, it.State)
			for _, f := range verdict.Facts {
				fmt.Fprintf(out, "  - %s: %s expected=%q actual=%q\n", f.Rule, f.Where, f.Expected, f.Actual)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&key, "key", "", "item key to submit (required)")
	cmd.Flags().StringVar(&inPath, "in", "-", "submission file, or - for stdin")
	return cmd
}

// statusCmd prints the per-state tally from the session.
func statusCmd(load func() (*quest.Session, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "show the progress tally",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := load()
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			prog := s.Progress()
			for _, st := range []quest.State{quest.TODO, quest.PASS, quest.REVIEW, quest.DONE, quest.SKIPPED, quest.BLOCKED} {
				fmt.Fprintf(out, "%-8s %d\n", st, prog[st])
			}
			fmt.Fprintf(out, "%-8s %d\n", "TOTAL", len(s.Items))
			return nil
		},
	}
}

// exportCmd emits terminal items to the JSONL sink and saves the export ratchet.
func exportCmd(sessionPath, outPath *string, load func() (*quest.Session, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: "export terminal results to JSONL (originals preserved)",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := load()
			if err != nil {
				return err
			}
			sink, err := newJSONLSink(*outPath)
			if err != nil {
				return err
			}
			n, err := quest.Export(s, sink)
			if err != nil {
				return err
			}
			if err := s.Save(*sessionPath); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "exported %d new record(s) to %s\n", n, *outPath)
			return nil
		},
	}
}

// rulesCmd prints the gate's rule catalog — the auto rulebook.
func rulesCmd(def gate.Definition) *cobra.Command {
	return &cobra.Command{
		Use:   "rules",
		Short: "print the gate's rule catalog (auto rulebook)",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			for _, m := range gate.Catalog(def.Rules()) {
				fmt.Fprintf(out, "%-6s %-24s %s\n", m.Level, m.ID, m.Desc)
			}
			return nil
		},
	}
}

// readSubmission reads the raw submission from a file path, or from stdin when path
// is "-" (or empty).
func readSubmission(cmd *cobra.Command, path string) ([]byte, error) {
	if path == "" || path == "-" {
		return io.ReadAll(cmd.InOrStdin())
	}
	return os.ReadFile(path)
}
