# reins

[![Version](https://img.shields.io/badge/version-v0.2.0-blue.svg)](https://github.com/park-jun-woo/reins/releases)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**A quest-CLI development framework** (Go). reins — it moves the authority to judge completion from the AI to a machine gate.
*"Generation is probabilistic, verification is deterministic."* (A reusable implementation of the [how-make-quest](https://www.parkjunwoo.com/tech/how-make-quest.md) methodology.)

On multi-step tasks, AI agents *judge their own completion* and quit early. reins hands the completion verdict
to a deterministic **gate**, so even an unreliable generator yields trustworthy completion. Agents are
disposable; progress accumulates.

> **For the detailed API and patterns, see [MANUAL.md](MANUAL.md)** (Manual for AI Agents). This document is an overview.

## Core model

- **Ratchet** — a one-way state machine. Once PASS, it is irreversible; remaining work decreases monotonically
  (`remaining(t+1) ≤ remaining(t)`).
- **Gate = a catalog of rules** — a set of violation-detection rules. Each rule fires (true) when it finds a problem
  and carries a fact (`Fact`). Severity is a **level** (Fail/Review) — a level, not a weight, so a single
  decisive violation means FAIL. `Evaluate` aggregates the fired rules by level into PASS/REVIEW/FAIL.
- **Authority asymmetry** — the PASS lock belongs to the **machine only**. L1 machine (deterministic, sole PASS
  authority) / L2 AI (skeptic, REVIEW only) / L3 human (the residual).
- **Fact feedback** — a FAIL is not an opinion but a location, expectation, and actual value (`Fact`). It turns the
  model's sycophancy into *convergence*.

## Two gate backends

| Backend | When | What |
|---|---|---|
| **Level aggregation** (`pkg/gate`) | Independent rules, simple gates | A `Rule` catalog + `Evaluate`. any Fail→FAIL, else Review→REVIEW, else PASS |
| **Defeat graph** (`pkg/graph`) | Inter-rule priority, root-cause feedback | A toulmin h-Categoriser backend. A tautology PASS warrant + violation Counters + **`Supersedes`** priority. Hand-written guards evaporate into declarative edges. With zero edges it is equivalent to level aggregation |

The defeat graph turns on when a `Definition` implements `gate.Evaluator` (opt-in). Side effects (HTTP/DNS) are
isolated via `pkg/ground` primitives + **staged evaluation** — when a cheap check fails, the network fetch never
happens. Graph evaluation produces a **walkthrough** straight to the agent (`Verdict.Feedback`: "why you lost +
what to change to win").

> toulmin is coupled only to `pkg/graph`/`pkg/ground`, and one-way (toulmin does not know about reins). `pkg/gate`
> and `pkg/cli` are toulmin-free, so a consumer that does not use the graph never links toulmin.

## Architecture (`pkg/`)

| Package | Role | Depends on |
|---|---|---|
| `pkg/textmatch` | Body-containment verifier — `Normalize`(NFC)·`Contains`·`MissingTokens`. Hallucination-block primitive | x/text |
| `pkg/temporal` | Time-spec normalization — a structured `Spec` (calendar/components/offset) → Gregorian ISO | (pure) |
| `pkg/quest` | Ratchet core — `State`·`Item`·`Verdict`/`Fact`·`Apply`·`Session`·`Export` | (pure) |
| `pkg/gate` | Gate contract — `Definition`·`Rule`·`Context`·`Evaluate`(level aggregation)·`Evaluator`(graph hook) | quest |
| `pkg/graph` | Defeat-graph backend — `Graph`·`Warrant`·`Counter`·`Attacks`·`Supersedes`·`EvaluateStaged` | gate, quest, toulmin |
| `pkg/ground` | Network ground primitives — `HTTPBody`·`MXResolves` (injectable `Resolver`, per-request snapshot) | (pure net) |
| `pkg/llm` | LLM call adapters — ollama/xai/gemini (HTTP) + claude/grok/codex (CLI subprocess). **Generation (L0) only; nothing to do with judging/ratchet** (authority asymmetry). Auto-sized num_ctx, env-only keys (HTTP) or delegated CLI login (subprocess) | net/http, os/exec |
| `pkg/cli` | Cobra scaffold — `NewQuestCmd` → scan/next/submit/status/export/rules (+ opt-in `loop`) | cobra, quest, gate, llm |

## Command skeleton (the how-make-quest canon)

```
scan    Seed N quests from input + initialize Progress (for streaming sources the consumer adds a run variant)
next    Emit one TODO + its authoring prompt and verification context
submit  Submit → gate evaluation → verdict → PASS lock / on FAIL, Fact + walkthrough feedback
status  Progress aggregate (PASS/REVIEW/DONE/TODO/SKIPPED/BLOCKED …)
export  Emit terminal results as JSONL (source preserved, emit-once ratchet)
rules   Print the gate's rule catalog (auto-rulebook — audit the list of cheeses it blocks)
loop    (opt-in) the automatic repeat of submit — the LLM generates, the gate judges and locks (below)
```

## Unattended drive — the `loop` command (opt-in)

The flow where an external agent runs `next`→`submit` by hand is closed into **a loop where, inside the CLI, an LLM
generates and the gate judges**. Opt in with `cli.Options{Loop: &cli.LoopOptions{…}}` and the `loop` command is
attached (when unset it is not attached — fully backward compatible).

```
for each remaining TODO:
  system = global + RuleSystem[the root-cause rule of the previous FAIL]   # per-rule coaching
  raw    = backend.Complete(system, Render(it)+feedback)                    # LLM generates (L0)
  verdict = gate verdict → ratchet Apply → export                          # same path as submit
  on FAIL, feed the feedback (identical to submit's output) back and retry (<MaxTries); otherwise lock → next
```

- **Authority asymmetry preserved** — the LLM is merely the generator (L0). **The PASS lock still belongs to the
  gate (machine) only.** On exceeding MaxTries it locks DONE → the loop terminates monotonically.
- **Per-root-cause system coaching** — via `Verdict.RootCause` (exposed deterministically by both the flat and
  graph backends), rule-specific instructions for the rule just missed are fed back.
- **Backends** — HTTP: `--model ollama:gemma4:e4b` (default) / `xai:…` / `gemini:…` (local ollama needs no key;
  num_ctx auto-sized from prompt length). CLI subprocess: `claude:…` / `grok:…` / `codex:…` — single-shot L0
  generators over the CLI's own login (no API key); opt into session continuity with `REINS_<NAME>_SESSION=continue`.

```bash
ccnews run --max-warcs 1                 # seed (streaming ingestion)
ccnews loop --model ollama:gemma4:e4b    # gemma4 generates the remaining TODOs → gate judges
```

## Making a quest

Implement just the 4 methods of `gate.Definition` and reins supplies the ratchet, command skeleton, aggregation,
and export:

```go
type Definition interface {
    Seed(args []string) ([]*quest.Item, error)            // input → initial TODO seeds
    Render(s *quest.Session, it *quest.Item) (string, error)                // the authoring prompt + verification context that next shows
    Prepare(s *quest.Session, it *quest.Item, raw []byte) (gate.Context, *quest.Verdict, error) // decode the submission (short-circuit if short)
    Rules() []gate.Rule                                   // the gate's violation-rule catalog
}

func main() { cli.NewQuestCmd("myquest", myDef{}, cli.Options{}).Execute() }
```

```go
// One cheese-defense rule = one violation detector. Found a new cheese → add one rule and the gate grows.
var whoAnchorPresent = gate.Rule{
    Meta: gate.RuleMeta{ID: "who-anchor-present", Level: gate.LevelFail, Desc: "the required who anchor exists in the source"},
    Check: func(ctx gate.Context) (bool, quest.Fact) {
        sub := ctx.Submission.(*Event6)
        if miss := textmatch.MissingTokens(ctx.Source, sub.Who.Anchors); len(miss) > 0 {
            return true, quest.Fact{Where: "who.anchors", Expected: "a source substring", Actual: miss[0]}
        }
        return false, quest.Fact{}
    },
}
```

When you need inter-rule priority, root-cause feedback, or network verification, use the graph backend
(`pkg/graph` + `gate.Evaluator` + `pkg/ground`) — see [MANUAL.md](MANUAL.md).

## Status

v1 build complete — 8 packages pass `go build`/`go test`, `filefunc validate` 0/0, `tsma` 0 TODO (all functions
covered), gofmt clean. Implemented through level aggregation + the defeat-graph backend (toulmin) + ground/staged
evaluation + **the `loop` unattended generate-verify loop (`pkg/llm`, ollama/xai/gemini)**.

**The first real consumer `comail` (email collection) runs end-to-end on top of reins' graph gate** — scan/next/
submit/status/export, graph verdicts, Supersedes, walkthrough feedback, real-network staged evaluation, ratchet
locking. (Design: `plans/reins/Phase007-toulmin-gate.md`)

## Repository layout

- `pkg/` — the framework Go module (`github.com/park-jun-woo/reins`)
- `MANUAL.md` — the detailed manual for AI agents
- `plans/reins/` — design documents (Phase001–006 v1 scaffold, Phase007 toulmin graph backend, Phase008 payload rehydration, Phase009 the loop (formerly agent) command·`pkg/llm`)
- `plans/ccnews/`·`plans/comail/` — instance designs + reins-porting Phases
- `comail/`, `ccnews/` — quest instances as separate modules (their own go.mod). They import reins and implement
  only the domain (`comail` is ported to the graph gate; `ccnews` is mid-port in design)

## Conventions

- State the deterministic gate explicitly — verdicts from input alone, the PASS lock by machine only.
- Rules are violation detectors, severity is a level — continuous weighting is reserved for the graph's *real
  contests* (L2 consensus).
- Cheese defense first — for every answer to "how would I cheat this gate?", one rule (auditable via the
  auto-catalog).
- Side effects via ground, rules stay pure — the network is owned by reins ground primitives, isolated by staging.
- No N=1 abstractions — a new abstraction is frozen only after a second consumer (`ccnews`) validates it.
