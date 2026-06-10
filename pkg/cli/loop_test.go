//ff:func feature=cli type=command control=sequence level=error
//ff:what TestNewLoopCmd — newLoopCmd가 loop 명령을 만들고 flag 기본값(--model: DefaultModel 미설정 ⇒ ollama:gemma4:e4b, 설정 ⇒ 그 값 / --max-items 0=전부)을 노출하며, backend 해석(주입 LLM 우선, 없으면 --model을 FromFlag로 — 잘못된 flag는 루프 전에 에러), load 실패 전파, --max-items 상한, backend 에러 중단, defer Save 실패의 stderr 경고 표면화를 검증한다.

package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// TestNewLoopCmd: newLoopCmd builds the `loop` command with the documented flag
// defaults (--model falls back to ollama:gemma4:e4b unless opts.DefaultModel is
// set; --max-items defaults to 0 = all) and its RunE resolves the backend
// (injected LLM wins, otherwise --model via llm.FromFlag — a bad flag errors
// before the loop), propagates load failures, caps work at --max-items, aborts
// on a backend error, and surfaces a deferred Save failure as a stderr warning.
func TestNewLoopCmd(t *testing.T) {
	dir := t.TempDir()
	session := dir + "/session.json"
	out := dir + "/out.jsonl"

	// runLoop executes a standalone loop command built by newLoopCmd against sess.
	runLoop := func(opts *LoopOptions, sess string, args ...string) (string, string, error) {
		c := newLoopCmd(stubDef{}, opts, &sess, &out, func() (*quest.Session, error) { return loadSession(sess) })
		var stdout, stderr bytes.Buffer
		c.SetOut(&stdout)
		c.SetErr(&stderr)
		c.SetArgs(append([]string{}, args...))
		err := c.Execute()
		return stdout.String(), stderr.String(), err
	}

	// Flag defaults: --model falls back to the package default and --max-items to 0…
	cmd := newLoopCmd(stubDef{}, &LoopOptions{}, &session, &out, func() (*quest.Session, error) { return loadSession(session) })
	if got := cmd.Flags().Lookup("model").DefValue; got != defaultLoopModel {
		t.Fatalf("model default = %q, want %q", got, defaultLoopModel)
	}
	if got := cmd.Flags().Lookup("max-items").DefValue; got != "0" {
		t.Fatalf("max-items default = %q, want %q", got, "0")
	}
	// …and opts.DefaultModel overrides the --model default.
	custom := newLoopCmd(stubDef{}, &LoopOptions{DefaultModel: "stub:model"}, &session, &out, func() (*quest.Session, error) { return loadSession(session) })
	if got := custom.Flags().Lookup("model").DefValue; got != "stub:model" {
		t.Fatalf("model default with DefaultModel = %q, want %q", got, "stub:model")
	}

	// With no injected backend a malformed --model errors before the loop runs.
	if _, _, err := runLoop(&LoopOptions{}, session, "--model", "nocolon"); err == nil {
		t.Fatal("loop = nil error, want FromFlag error for bad --model")
	}

	// A valid --model resolves via FromFlag, then a corrupt session fails load.
	badSession := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(badSession, []byte("{not json"), 0o644); err != nil {
		t.Fatalf("write bad session: %v", err)
	}
	if _, _, err := runLoop(&LoopOptions{}, badSession); err == nil {
		t.Fatal("loop = nil error, want load error for corrupt session")
	}

	// Injected backend (--model ignored): --max-items 1 caps a 2-item session…
	backend := llm.CallFunc(func(system, user string) (string, error) { return "good", nil })
	opts := Options{Loop: &LoopOptions{LLM: backend}}
	if _, err := newLoopRoot(t, opts, session, out, "scan", "a", "b"); err != nil {
		t.Fatalf("scan: %v", err)
	}
	stdout, _, err := runLoop(opts.Loop, session, "--max-items", "1")
	if err != nil {
		t.Fatalf("loop --max-items 1: %v", err)
	}
	if !strings.Contains(stdout, "processed 1 item(s)") {
		t.Fatalf("loop output = %q, want 'processed 1 item(s)'", stdout)
	}
	// …and a second uncapped run drains the remaining TODO through the gate path.
	stdout, _, err = runLoop(opts.Loop, session)
	if err != nil {
		t.Fatalf("loop: %v", err)
	}
	if !strings.Contains(stdout, "processed 1 item(s)") {
		t.Fatalf("second loop output = %q, want 'processed 1 item(s)'", stdout)
	}

	// A backend error aborts the loop with that error.
	failSession := filepath.Join(dir, "fail.json")
	if _, err := newLoopRoot(t, opts, failSession, out, "scan", "x"); err != nil {
		t.Fatalf("scan fail session: %v", err)
	}
	failing := &LoopOptions{LLM: llm.CallFunc(func(system, user string) (string, error) { return "", errBackend })}
	if _, _, err := runLoop(failing, failSession); err == nil {
		t.Fatal("loop = nil error, want backend error")
	}

	// A deferred Save failure is surfaced as a stderr warning, not a hard error.
	roDir := filepath.Join(dir, "ro")
	if err := os.Mkdir(roDir, 0o555); err != nil {
		t.Fatalf("mkdir ro: %v", err)
	}
	t.Cleanup(func() { os.Chmod(roDir, 0o755) })
	_, stderr, err := runLoop(opts.Loop, filepath.Join(roDir, "session.json"))
	if err != nil {
		t.Fatalf("loop on read-only dir: %v", err)
	}
	if !strings.Contains(stderr, "warning: save session after loop") {
		t.Fatalf("stderr = %q, want save warning", stderr)
	}
}
