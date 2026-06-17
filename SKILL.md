---
name: reins-quest
description: Build a quest CLI in Go with the reins framework — it moves the authority to declare "done" from the AI to a deterministic machine gate ("generation is probabilistic; verification is deterministic"). Use this skill when asked to build a quest CLI, a reins consumer, a generate→verify→retry loop, a deterministic completion gate, or an unattended LLM generation loop whose output a machine judges. Plug in one domain contract (gate.Definition) and reins supplies the ratchet, command skeleton, aggregation, feedback, and export. Triggers on keywords like reins, quest CLI, gate.Definition, gate.Rule, deterministic gate, quest-CLI, or unattended loop.
license: MIT
metadata:
  author: park-jun-woo
  version: "0.2.2"
---

# reins-quest — Build a Quest CLI Whose "Done" Is Judged by a Machine

reins is a quest-CLI framework (Go). **It moves the authority to declare "done" from the AI to a deterministic machine gate.** *"Generation is probabilistic; verification is deterministic."* A consumer plugs in only the domain logic (`gate.Definition`); reins supplies the ratchet, the command skeleton, aggregation, feedback, and export. Agents are disposable; progress is cumulative and irreversible.

## When to Use This Skill

- Building a new quest CLI / reins consumer in Go
- Designing a deterministic completion gate (a set of violation-detecting rules)
- Wiring an unattended generate→gate→retry loop (LLM generates, the gate judges)
- Adding rules, a defeat-graph backend, or network-ground verification to an existing reins quest

## Do NOT Use This Skill When

- The "done" check **cannot** be made machine-deterministic and genuinely needs an agent's tool-using exploration (open-ended coding, "fix the repo"). reins generators are pure L0 — no tools, no Act→Observe. Use a normal agentic loop instead.
- You want the LLM to judge its own completion. reins **forbids** this by design (only the machine locks PASS). If you need LLM-as-judge, reins is the wrong tool.

## Install

reins is a library you import from your own Go module:

```bash
go get github.com/park-jun-woo/reins@latest
```

**Prerequisites:** Go 1.25+. reins pulls `github.com/park-jun-woo/toulmin` (only when you use the graph backend).

## The Mental Model (read first)

- **Ratchet** — a one-way state machine. Once PASS, it is irreversible; remaining work only decreases.
- **Gate** — a set of violation-detecting rules. A rule *fires* (true) when it finds a problem and carries a `Fact`. Severity is a **level** (Fail/Review), never a weight. One decisive Fail ⇒ FAIL.
- **Authority asymmetry** — only the machine locks PASS. **L1 machine** (deterministic, sole PASS authority) / **L2 AI** (skeptic, REVIEW only) / **L3 human** (the remainder).
- **Fact feedback** — a FAIL is not an opinion but a located, quantified value (`Fact{Where,Expected,Actual}`). It turns a sycophantic model toward *convergence*.

## Quickstart: the simplest quest (4 methods + one line of `main`)

Implement the four `gate.Definition` methods; reins supplies the rest.

```go
// definition.go
type myDef struct{}

func (myDef) Seed(args []string) ([]*quest.Item, error) {
	// input → initial TODO items (one per line/file/record)
	items := make([]*quest.Item, len(args))
	for i, a := range args {
		it := &quest.Item{Key: fmt.Sprintf("item-%d", i), State: quest.TODO}
		it.SetPayload(map[string]string{"source": a}) // never write Payload directly
		items[i] = it
	}
	return items, nil
}

func (myDef) Render(s *quest.Session, it *quest.Item) (string, error) {
	var p map[string]string
	it.DecodePayload(&p)
	// the authoring prompt `next` shows; READ-ONLY on s.Meta (next never Saves)
	return "From this source, output ONLY a JSON {\"summary\": \"...\"}:\n" + p["source"], nil
}

func (myDef) Prepare(s *quest.Session, it *quest.Item, raw []byte) (gate.Context, *quest.Verdict, error) {
	var sub struct{ Summary string `json:"summary"` }
	if err := json.Unmarshal(raw, &sub); err != nil {
		return gate.Context{}, nil, err // decode failure → caller sees the error
	}
	var p map[string]string
	it.DecodePayload(&p)
	// short verdict (3rd return) != nil short-circuits the gate (e.g. OutSkip for untrusted input)
	return gate.Context{Item: it, Submission: &sub, Source: p["source"]}, nil, nil
}

func (myDef) Rules() []gate.Rule { return []gate.Rule{summaryPresent, summaryGrounded} }

// main.go — one line wires the whole CLI
func main() {
	cli.NewQuestCmd("myquest", myDef{}, cli.Options{Version: "0.1"}).Execute()
}
```

## One rule = one violation detector

A rule fires `(true, Fact{...})` when it finds a problem; otherwise `(false, _)`.

```go
var summaryGrounded = gate.Rule{
	Meta: gate.RuleMeta{ID: "summary-grounded", Level: gate.LevelFail, Desc: "summary tokens exist in source (no hallucination)"},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		sub := ctx.Submission.(*struct{ Summary string `json:"summary"` })
		if miss := textmatch.MissingTokens(ctx.Source, strings.Fields(sub.Summary)); len(miss) > 0 {
			return true, quest.Fact{Where: "summary", Expected: "source substring", Actual: miss[0]}
		}
		return false, quest.Fact{}
	},
}
```

`gate.Evaluate(rules, ctx)` aggregates fired rules **by level: any Fail→FAIL, else any Review→REVIEW, else PASS.** Deterministic: same `(rules, ctx)` → same `Verdict`.

## Commands (auto-generated by NewQuestCmd)

| Command | Purpose |
|---|---|
| `scan <input>` | Seed N quests from input |
| `next` | Show one TODO + its authoring prompt / verification context |
| `submit --key <k> --in <file>\|-` | Decode (`Prepare`) → gate eval → verdict; lock PASS, or on FAIL emit Fact feedback |
| `status` | Progress tally (PASS/REVIEW/DONE/TODO/SKIPPED/BLOCKED) |
| `export` | Emit terminal results as JSONL (originals preserved, emit-once) |
| `rules` | Print the gate's rule catalog (the auto rulebook — audit what it blocks) |
| `loop` | Opt-in unattended drive (LLM generates → gate judges → retry). Attached only if `Options.Loop != nil` |

Every `submit` auto-emits terminal items to `--out` (default `<name>-results.jsonl`). Tune with `Options{Out, Version, ExtraCommands, Loop}`.

## Core types

| Type | Shape |
|---|---|
| `quest.Item` | `{ Key; State; Tries; Payload json.RawMessage; Log; Emitted }` — use `it.SetPayload(v)`/`it.DecodePayload(&v)`, never the field |
| `quest.State` | `TODO PASS REVIEW DONE SKIPPED BLOCKED` (terminal = all but TODO) |
| `quest.Verdict` | `{ Outcome; Facts []Fact; Feedback; RootCause }` — `RootCause` names the rule that caused FAIL/REVIEW (agent coaching) |
| `quest.Fact` | `{ Rule, Where, Expected, Actual string }` |
| `gate.Context` | `{ Item; Submission any; Source string; Grounds map[string]string }` |
| `gate.Level` | `LevelFail \| LevelReview` |
| const | `quest.MaxTries = 3` — FAIL accrued to MaxTries ⇒ lock to DONE (monotone convergence) |

`quest.Apply(it, v, now)` applies a verdict to the ratchet (PASS/REVIEW/SKIPPED/BLOCKED lock, FAIL is `Tries++`).

## Verification primitives

```go
textmatch.Normalize(s)                 // NFC + whitespace fold + Trim (no case-fold — ToLower first if needed)
textmatch.Contains(source, token)      // substring after normalization
textmatch.MissingTokens(source, toks)  // tokens absent from source — the hallucination block
temporal.Resolve(spec, ref)            // structured Spec + ref time.Time → Gregorian ISO (undetermined ⇒ Determined=false)
temporal.ComponentsInAnchor(...)       // extract time components from an anchor string
```

## Unattended drive: the `loop` command (opt-in)

Closes the `next`→`submit` cycle in-process: an LLM generates each TODO's payload, the gate judges, FAIL feedback is fed back until PASS or `MaxTries`. **The LLM is only the generator (L0); only the gate locks PASS** — there is no API to grant the LLM PASS authority.

```go
cli.NewQuestCmd("myquest", myDef{}, cli.Options{
	Version: "0.1",
	Loop: &cli.LoopOptions{
		DefaultModel: "ollama:gemma4:e4b",          // "" ⇒ this default
		System:       "You are a strict generator.", // global generation system prompt
		RuleSystem:   map[string]string{             // per-rule coaching keyed by Verdict.RootCause
			"summary-grounded": "Use ONLY words present in the source.",
		},
		// LLM: injected llm.Backend (tests); when set, --model is ignored
	},
}).Execute()
```

Run: `myquest loop [--model backend:model] [--max-items N]`. On the `MaxTries`-th FAIL the item locks DONE → `NextTODO` drops it → the loop terminates (monotone convergence).

### LLM backends (`pkg/llm`)

| Token | Transport | Auth |
|---|---|---|
| `ollama:<model>` | HTTP (local, `num_ctx` auto-sized) | none |
| `xai:<model>` / `gemini:<model>` | HTTP (OpenAI-compat / Gemini) | **env-only** API key |
| `claude:<model>` | subprocess `claude -p` (`--max-turns 1 --tools ""`) | CLI login — **no API key read** |
| `grok:<model>` | subprocess `grok -p` (single-turn) | CLI login — **no API key read** |
| `codex:<model>` | subprocess `codex exec` (`-s read-only`) | CLI login — **no API key read** |
| `geminicli:<model>` | subprocess `gemini -p` (`--approval-mode plan`) | Google login — **no API key read** (separate token; `gemini:` is the HTTP backend) |

- **HTTP options** (struct fields, zero ⇒ prior default = backward-compatible): all three HTTP backends take `MaxOutputTokens int` (0 ⇒ 2048) and `Temperature *float64` (nil ⇒ 0); ollama also takes `Think *bool`. Or pass them in the `--model` query: `ollama:qwen3:8b?max_output_tokens=8192&think=false`. Raise `max_output_tokens` so reasoning models aren't truncated (ollama grows `num_ctx` to match). An unknown key for a backend is a loud error (allowed: ollama `max_output_tokens`/`num_ctx`/`temperature`/`think`, xai & gemini `max_output_tokens`/`temperature`, subprocess none).
- Subprocess backends: token `:<model>` or `:default` (CLI's configured model). `REINS_<NAME>_BIN` overrides the binary.
- **Session is fully stateless by default** (matches HTTP backends + reins' deterministic FAIL-feedback convergence). Opt into carrying the CLI's own conversation with `REINS_<NAME>_SESSION=continue` (stateless recommended — session mode double-exposes the prior attempt).
- Inject `llm.CallFunc` (HTTP) or the `exec<Name>` package-var seam (subprocess) for network-free tests.

## Advanced: defeat-graph backend (`pkg/graph`)

Use **only** when rules are not independent (one violation makes another moot) and you need inter-rule precedence or root-cause feedback. If level aggregation suffices, do not use the graph (it is overkill). Implement `gate.Evaluator` (`Evaluate(ctx) quest.Verdict`) and reins takes that path.

```go
g := graph.NewGraph("myquest")
// tautology PASS warrant — supply your own always-true fn (graph has no exported helper)
pass := g.Warrant(func(toulmin.Context, toulmin.Specs) (bool, any) { return true, nil })
fmtR := g.Counter(ruleEmailFormat, gate.LevelFail).Attacks(pass)
holder := g.Counter(ruleSourceLacksEmail, gate.LevelFail).Attacks(pass).Needs("source-body")
fmtR.Supersedes(holder)                                      // deterministic precedence (replaces hand-rolled guards)
v := g.EvaluateStaged(ctx, snap, provider)                  // tier-0 (no ground) first; resolve ground only if clean
```

- **`.Attacks(target)`** = toulmin graph edge (verdict/contest). **`.Supersedes(...)`** = reins-side deterministic precedence (excludes a downstream counter from the tally).
- **Side effects through ground, rules stay pure**: declare `.Needs("name")`, map it in a provider via `ground.Snapshot` (`HTTPBody`/`MXResolves`), and `EvaluateStaged` resolves it once (cached) into `ctx.Grounds["name"]`. Inject the `ground.Resolver` interface for network-free tests.

## Conventions (the philosophy — follow these)

1. **State the deterministic gate** — judge from input alone; only the machine locks PASS.
2. **Rules are violation detectors; severity is a level** — never fake a hard check with a weight (continuous weighting is for the graph's genuine contest only).
3. **Cheese defense first** — for every answer to "how would I fool this gate?", add one rule (audited by `rules`).
4. **Side effects through ground, rules stay pure** — the network is owned by reins ground primitives, isolated by staged eval.
5. **No N=1 abstraction** — freeze a new abstraction only after a second consumer validates it.
6. **Follow filefunc conventions** if the project uses them — 1 file = 1 func/type (tests included), `//ff:func`/`//ff:type` + `//ff:what` annotations at the top of each file.

## Quick decision guide

| Situation | Use |
|---|---|
| Independent rules, simple gate | `gate.Rule` + `Rules()` (level aggregation) |
| Inter-rule precedence / root-cause feedback | `pkg/graph` + `Evaluator` + `Supersedes` |
| Side-effect verification (HTTP/DNS) | `pkg/ground` + `.Needs()` + `EvaluateStaged` |
| Block body hallucination | `textmatch.MissingTokens` |
| Date/time normalization | `pkg/temporal` |
| Short-circuit an untrusted submission | `Prepare`'s `short` verdict |
| Unattended drive (LLM generates, gate judges) | `Options{Loop}` + `pkg/llm` |

## Common Errors and Fixes

| Symptom | Cause | Fix |
|---|---|---|
| `loop` subcommand missing | `Options.Loop == nil` | Set `Options{Loop: &cli.LoopOptions{...}}` |
| `invalid --model ...: model name is empty` | empty model token | Use `backend:model` or the `:default` sentinel (e.g. `claude:default`) |
| Gate never reaches PASS in `loop` | rule fires forever; sycophantic generator | Add `RuleSystem[ruleID]` coaching keyed on `Verdict.RootCause`; verify the rule is satisfiable |
| Item locks DONE without PASS | hit `MaxTries` (3) FAILs | Expected (monotone convergence); inspect Facts — the gate or prompt is mis-specified |
| Payload reads wrong/empty | wrote `it.Payload` directly | Use `it.SetPayload(v)` / `it.DecodePayload(&v)` only |
| `Render` mutated state but it vanished | `next` does not Save; Render is read-only | Mutate `s.Meta` in `Prepare` (submit Saves after Prepare), not Render |
| Subprocess backend "command not found" | CLI not on PATH | Install the CLI (`claude`/`grok`/`codex`) and log in, or set `REINS_<NAME>_BIN` |
| Linking toulmin you don't use | imported `pkg/graph` | `pkg/gate`+`pkg/cli` don't import toulmin — skip the graph if level aggregation suffices |

## Full Documentation

- **`MANUAL.md`** (repo root) — the complete manual for AI agents: package map, defeat-graph topology, staged-eval ground, walkthrough feedback, every backend.
- **Reference consumer** — `comail/main.go` shows the one-line `cli.NewQuestCmd` wiring.
- **Quest philosophy** — https://www.parkjunwoo.com/tech/how-make-quest.md ("generation is probabilistic, verification is deterministic").
