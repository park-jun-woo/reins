package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// stubDef is a trivial Definition: one item per arg, gate fails when the submission
// is "bad".
type stubDef struct{}

func (stubDef) Seed(args []string) ([]*quest.Item, error) {
	items := make([]*quest.Item, len(args))
	for i, a := range args {
		items[i] = &quest.Item{Key: a, State: quest.TODO}
	}
	return items, nil
}
func (stubDef) Render(it *quest.Item) (string, error) { return "render:" + it.Key, nil }
func (stubDef) Prepare(it *quest.Item, raw []byte) (gate.Context, *quest.Verdict, error) {
	return gate.Context{Item: it, Submission: string(raw)}, nil, nil
}
func (stubDef) Rules() []gate.Rule {
	return []gate.Rule{{
		Meta: gate.RuleMeta{ID: "not-bad", Level: gate.LevelFail, Desc: "submission must not be bad"},
		Check: func(ctx gate.Context) (bool, quest.Fact) {
			if s, _ := ctx.Submission.(string); strings.TrimSpace(s) == "bad" {
				return true, quest.Fact{Where: "body", Expected: "good", Actual: "bad"}
			}
			return false, quest.Fact{}
		},
	}}
}

func run(t *testing.T, def gate.Definition, session, out string, in string, args ...string) string {
	t.Helper()
	cmd := NewQuestCmd("stub", def, Options{})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetIn(strings.NewReader(in))
	full := append([]string{"--session", session, "--out", out}, args...)
	cmd.SetArgs(full)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute %v: %v\n%s", args, err, buf.String())
	}
	return buf.String()
}

func TestQuestCmdFlow(t *testing.T) {
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	out := filepath.Join(dir, "out.jsonl")
	def := stubDef{}

	run(t, def, session, out, "", "scan", "a", "b")

	if got := run(t, def, session, out, "", "next"); !strings.Contains(got, "render:a") {
		t.Fatalf("next = %q", got)
	}

	// PASS path locks the item.
	if got := run(t, def, session, out, "good", "submit", "--key", "a"); !strings.Contains(got, "PASS") {
		t.Fatalf("submit a = %q", got)
	}

	// FAIL path keeps it TODO and prints the Fact.
	got := run(t, def, session, out, "bad", "submit", "--key", "b")
	if !strings.Contains(got, "FAIL") || !strings.Contains(got, "not-bad") {
		t.Fatalf("submit b = %q", got)
	}

	if got := run(t, def, session, out, "", "rules"); !strings.Contains(got, "not-bad") {
		t.Fatalf("rules = %q", got)
	}
	if got := run(t, def, session, out, "", "status"); !strings.Contains(got, "PASS") {
		t.Fatalf("status = %q", got)
	}
}
