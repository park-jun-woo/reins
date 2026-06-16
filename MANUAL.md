# reins — Manual for AI Agents

A quest-CLI development framework (Go). **It moves the authority to declare "done" from the AI to a
deterministic machine gate.** *"Generation is probabilistic; verification is deterministic."*

A consumer plugs in only the domain logic (`gate.Definition`); reins supplies the ratchet, the command
skeleton, aggregation, feedback, and export. Agents are disposable; progress is cumulative and
irreversible.

---

## Core model

- **Ratchet** — a one-way state machine. Once PASS, it is irreversible; the remaining work
  monotonically decreases.
- **Gate** — a set of violation-detecting rules. Each rule fires (true) when it finds a problem and
  carries a `Fact`. Severity is a **level** (Fail/Review): a single decisive violation is FAIL.
- **Authority asymmetry** — only the machine locks PASS. L1 machine (deterministic, sole PASS
  authority) / L2 AI (skeptic, REVIEW only) / L3 human (the remainder).
- **Fact feedback** — a FAIL is not an opinion but a located, quantified value (`Fact`). It turns a
  sycophantic model toward *convergence*.

## Package map

| Package | Role | Deps |
|---|---|---|
| `pkg/quest` | Ratchet core — `Item`·`State`·`Verdict`/`Fact`·`Session`·`Apply`·`Export` | (pure) |
| `pkg/gate` | Gate contract — `Definition`·`Rule`·`Level`·`Context`·`Evaluate`(level aggregation)·`Evaluator` | quest |
| `pkg/graph` | **Defeat-graph backend** — toulmin h-Categoriser. `Graph`·`Counter`·`Supersedes`·staged eval | gate, quest, **toulmin** |
| `pkg/ground` | Network ground primitives — `HTTPBody`·`MXResolves` (injectable, snapshot) | (pure net) |
| `pkg/textmatch` | Body-containment verification — `Normalize`(NFC)·`Contains`·`MissingTokens`. Hallucination block | x/text |
| `pkg/temporal` | Time normalization — structured `Spec` → Gregorian ISO | (pure) |
| `pkg/llm` | LLM call adapters — `Backend`(ollama/xai/gemini HTTP + `claude`/`grok`/`codex` subprocess)·`CallFunc`·`FromFlag`·shared `runSubprocess`/no-tools preamble·auto `num_ctx`. Generation (L0) only; never judges/locks | net/http, os/exec |
| `pkg/cli` | Cobra scaffold — `NewQuestCmd` → scan/next/submit/status/export/rules (+ opt-in `loop`) | cobra, llm |

> **toulmin isolation**: only `pkg/graph`/`pkg/ground` are heavy. `pkg/gate`·`pkg/cli` do not import
> toulmin, so a consumer that doesn't use the graph never links toulmin.

---

## The simplest quest (4 methods + one line of main)

Implement the four `gate.Definition` methods and reins supplies the rest:

```go
type Definition interface {
    Seed(args []string) ([]*quest.Item, error)            // input → initial TODOs
    Render(s *quest.Session, it *quest.Item) (string, error)                // the authoring prompt `next` shows (read-only s.Meta)
    Prepare(s *quest.Session, it *quest.Item, raw []byte) (gate.Context, *quest.Verdict, error) // decode a submission (short-circuit if non-nil; may update s.Meta)
    Rules() []gate.Rule                                   // the gate's violation-rule catalog
}

func main() { cli.NewQuestCmd("myquest", myDef{}, cli.Options{Version: "0.1"}).Execute() }
```

**One rule = one violation detector.** On fire it returns `(true, Fact{Where,Expected,Actual})`:

```go
var whoAnchorPresent = gate.Rule{
    Meta: gate.RuleMeta{ID: "who-anchor-present", Level: gate.LevelFail, Desc: "required who anchor is real"},
    Check: func(ctx gate.Context) (bool, quest.Fact) {
        sub := ctx.Submission.(*Event6)
        if miss := textmatch.MissingTokens(ctx.Source, sub.Who.Anchors); len(miss) > 0 {
            return true, quest.Fact{Where: "who.anchors", Expected: "source substring", Actual: miss[0]}
        }
        return false, quest.Fact{}
    },
}
```

`gate.Evaluate(rules, ctx)` aggregates fired rules by level: **any Fail→FAIL, else any Review→REVIEW,
else PASS**. Deterministic (same `(rules, ctx)` → same `Verdict`).

## Core types

```go
// quest
type Item struct { Key string; State State; Tries int; Payload json.RawMessage; Log []Attempt; Emitted bool; … }
// Payload is raw JSON — write/read it via it.SetPayload(v) / it.DecodePayload(&v), never the field directly.
type State string  // TODO PASS REVIEW DONE SKIPPED BLOCKED  (terminal = PASS/REVIEW/DONE/SKIPPED/BLOCKED)
type Verdict struct { Outcome Outcome; Facts []Fact; Feedback string; RootCause string } // Outcome: PASS REVIEW FAIL SKIPPED BLOCKED; RootCause = the rule that caused FAIL/REVIEW (agent coaching)
type Fact struct { Rule, Where, Expected, Actual string }
const MaxTries = 3  // FAIL accrued to MaxTries → lock to DONE

// gate
type Context struct { Item *quest.Item; Submission any; Source string; Grounds map[string]string }
type Level int  // LevelFail | LevelReview
```

- `quest.Apply(it, v, now)` — applies a verdict to the ratchet (PASS/REVIEW/SKIPPED/BLOCKED lock, FAIL
  is Tries++).
- `quest.Export(s, sink)` — emits terminal, not-yet-emitted Items to the sink once (the export ratchet).
- If `Prepare`'s `short *quest.Verdict != nil`, the gate is skipped and that verdict short-circuits
  (e.g. an untrusted submission → `OutSkip`).

## Command skeleton (NewQuestCmd)

```
scan    input → seed N quests (for a streaming source, the consumer adds a run variant)
next    one TODO + the authoring prompt / verification context
submit  submit → gate eval → verdict → lock PASS / on FAIL emit Fact feedback
status  progress tally (PASS/REVIEW/DONE/TODO/SKIPPED/BLOCKED)
export  emit terminal results as JSONL (originals preserved, emit-once)
rules   the gate's rule catalog (auto rulebook — audit the cheese it blocks)
```

`submit` takes `--key <k>` + `--in <file>|-` (raw bytes → `Prepare` decodes). Every submit auto-emits
terminal items to `--out` (default `<name>-results.jsonl`). Tune via `Options{Out, Version}`.

---

## Unattended drive: the `loop` command (opt-in)

The same `next`→`submit` an external agent runs by hand, closed **in-process** as a generate→gate→retry
loop: an LLM generates each TODO's payload, the gate judges, FAIL feedback is fed back until PASS or
`MaxTries`. Opt in with `Options{Loop: &LoopOptions{…}}` (nil ⇒ the command is not attached, fully
backward-compatible).

```go
type LoopOptions struct {
    DefaultModel string            // "" ⇒ "ollama:gemma4:e4b"
    System       string            // global generation system prompt
    RuleSystem   map[string]string // toulmin rule ID → extra system coaching when that rule was the FAIL root cause
    LLM          llm.Backend       // injected backend (tests); when set, --model is ignored
}
```

The loop (`loop [--model backend:model] [--max-items N]`):

```
for it := s.NextTODO(); it != nil; it = s.NextTODO() {
    system := opts.System + RuleSystem[verdict.RootCause]   // global + last-FAIL rule coaching
    raw    := backend.Complete(system, def.Render(s,it)+feedback)  // LLM generates (L0)
    verdict := evaluateAndApply(...)                        // SAME path as submit: gate→Apply→export
    if verdict.Outcome != quest.OutFail { break }           // PASS/REVIEW/SKIP/BLOCK → lock, next item
    feedback = renderVerdictText(...)                       // identical to what submit prints
}
```

- **Authority asymmetry holds** — the LLM is only the generator (L0). **Only the gate locks PASS.** The
  loop calls `quest.Apply`; on the `MaxTries`-th FAIL it locks DONE → `NextTODO` drops it → the loop
  terminates (monotone convergence). The framework exposes no API to grant the LLM PASS authority.
- **`Verdict.RootCause`** (additive, backward-compatible field) names the rule that caused FAIL/REVIEW —
  set deterministically on **both** paths: flat `gate.Evaluate` (first fired Fail rule's ID) and the
  graph backend (`selectRootCause`'s top counter). `RuleSystem[RootCause]` turns it into rule-specific
  coaching on retry.
- **Feedback parity** — the FAIL text fed to the model is the very string `submit` prints
  (`renderVerdict`/`renderVerdictText` shared), so human-visible and model-visible feedback never drift.
- **Backends** (`pkg/llm`): **HTTP** — `ollama:<model>` (local, no key, `num_ctx` auto-sized from prompt
  length), `xai:<model>`/`gemini:<model>` (OpenAI-compat / Gemini, **env-only** API keys); all three are
  `net/http` at `temperature: 0`. **Subprocess** — `claude:<model>`/`grok:<model>`/`codex:<model>` shell
  out to a CLI via `os/exec`; auth is delegated entirely to that CLI's own login (subscription/OAuth/
  keychain/env key) so reins reads **no API key** for them. All subprocess backends share `runSubprocess`
  and the fixed no-tools preamble (`withNoToolsPreamble`); each exposes a `var exec<Name>` package seam for
  subprocess-free tests, a pointer-receiver adapter that carries an opt-in session id, and a `:<model>` or
  `:default` token (`default` ⇒ the CLI's configured model, since `FromFlag` rejects an empty model).
  Inject `llm.CallFunc` for network-free HTTP tests.
- **`claude:<model>` — Claude CLI** (`claude -p`, headless print): `--max-turns 1` pins a single-shot L0
  generator (no agentic tool loop), `--tools ""` + no-tools preamble block tool narration,
  `--permission-mode dontAsk` keeps it from blocking. Token `claude:opus`/`claude:sonnet`/`claude:default`;
  `REINS_CLAUDE_BIN` overrides the binary; envelope `--output-format json` → `result`/`session_id`/`is_error`.
- **`grok:<model>` — Grok CLI** (`grok -p`, single-turn): the claude twin (`--max-turns 1` +
  `--tools ""`/`--disable-web-search`/`--no-subagents`/`--no-memory` + preamble + `--permission-mode dontAsk`),
  reaching the same xAI models as the `xai:` HTTP backend but over the **CLI login (no API key)**. The user
  prompt travels via `--prompt-file` (grok's `-p` takes the prompt as a value; stdin is not accepted).
  Envelope `--output-format json` → `text`/`sessionId`/`stopReason` (success = `EndTurn`). `REINS_GROK_BIN`
  overrides the binary.
- **`codex:<model>` — Codex CLI** (`codex exec`, headless agent): has **no `--max-turns`**, so the single-shot
  L0 guarantee leans on `-s read-only` (block side effects) + the no-tools preamble + `--ignore-user-config`/
  `--ignore-rules` (keep `CODEX_HOME` auth, drop repo `AGENTS.md`/rules). Output is a `--json` **JSONL event
  stream**: the last `item.completed` with `agent_message` is the result text; `thread.started.thread_id`
  is the session id. The system prompt is prepended into the stdin prompt (codex has no system channel).
  Token `codex:gpt-5`/`codex:o3`/`codex:default`; `REINS_CODEX_BIN` overrides the binary.
- **Session (subprocess backends, default fully stateless)**: each `Complete` is independent (claude
  `--no-session-persistence` / codex `--ephemeral` / grok no-resume), matching the HTTP backends and reins'
  own deterministic FAIL-feedback convergence. Opt into carrying the CLI's own conversation across a run with
  `REINS_<NAME>_SESSION=continue` — the first reply's session id is carried into later calls as `--resume`
  (claude/grok) or the `exec resume <id>` subcommand (codex; the read-only flag stays at exec level, before
  the subcommand). Any other value falls back to stateless (no forged session). Note: session mode
  double-exposes the prior attempt (model history *and* re-fed FAIL text), so stateless is recommended.

---

## Advanced: the defeat-graph backend (pkg/graph)

Use it when rules are **not independent** (one violation makes another moot). If a `Definition`
implements `gate.Evaluator` (`Evaluate(ctx) quest.Verdict`), reins takes that path instead of
`Rules()` (`Rules()` is kept for the `rules` catalog). **If level aggregation suffices, do not use the
graph** (it is overkill).

Topology: **one tautology PASS warrant + every violation = a Counter that attacks the warrant**.

```go
g := graph.NewGraph("myquest")
pass   := g.Warrant(alwaysTrue)                          // always-active PASS warrant
fmtR   := g.Counter(ruleEmailFormat,      gate.LevelFail).Attacks(pass)
holder := g.Counter(ruleSourceLacksEmail, gate.LevelFail).Attacks(pass)
free   := g.Counter(ruleFreemail,         gate.LevelReview).Attacks(pass)

fmtR.Supersedes(holder)     // bad format → drop the source check from the tally (precedence)
free.Supersedes(holder)     // freemail → absorb the source check → preserve REVIEW
```

- **`.Supersedes(...)`** = reins-side deterministic precedence (an active upstream counter excludes a
  downstream one from the tally). It replaces hand-rolled guards. (toulmin's `Attacks` defeat only
  lowers the verdict float and cannot clear *Activated*, so crisp precedence goes through Supersedes.)
- **`.Attacks(target)`** = a toulmin graph edge (for the verdict/contest). Violation→warrant is Attacks.
- Decision: `g.Evaluate(ctx)` takes *active counters − superseded = remaining* and aggregates them by
  Level → `Verdict` (+ the walkthrough `Feedback`). With zero edges it equals `gate.Evaluate`
  (`graph.FromRules(rules)`).

### Side effects / network: ground provider + staged eval (G5)

Do not put side effects (HTTP/DNS) inside rules. reins provides **ground primitives**:

```go
snap := ground.NewSnapshot(nil)                 // nil = real net; tests inject a fake Resolver
// a counter declares its ground dependency → automatic tier classification
holder := g.Counter(ruleSourceLacksEmail, gate.LevelFail).Attacks(pass).Needs("source-body")
mx     := g.Counter(ruleMxMissing,        gate.LevelFail).Attacks(pass).Needs("mx")

// the consumer maps a ground name → the actual resolve
provider := func(name string, ctx gate.Context, snap *ground.Snapshot) (string, error) {
    c := ctx.Submission.(*Candidate)
    switch name {
    case "source-body": return snap.HTTPBody(c.Source)
    case "mx":          b, e := snap.MXResolves(domain(c.Email)); return fmt.Sprint(b), e
    }
    return "", fmt.Errorf("unknown ground %q", name)
}
v := g.EvaluateStaged(ctx, snap, provider)
```

- **Staged**: the no-ground tier 0 is evaluated first → **if a FAIL remains, stop immediately (no
  ground is resolved = zero network)**. If clean, each ground is snapshot-resolved once → injected into
  `ctx.Grounds` → tier 1 is evaluated.
- A ground is **snapshotted/cached once per request** (re-reading the same URL is still one call).
  Rules read `ctx.Grounds["source-body"]` and stay pure. A resolve failure is a deterministic FAIL Fact.
- Inject the `ground.Resolver` interface (`Fetch`/`LookupMX`) so tests are deterministic with no
  network.

### Walkthrough feedback (Verdict.Feedback)

On FAIL/REVIEW, graph evaluation fills `Verdict.Feedback` with an **agent-facing walkthrough** — not a
flat Fact list but *"why you lost + what to change to win"*:

```
FAIL. root cause = email-format (remaining active FAIL, upstream).
  Fact: where=email expected="valid email format" actual="not-an-email"
  source-lacks-email: superseded by email-format → side-branch.
  → to flip the verdict, clear email-format.
```

The cli `submit` prints `Feedback` when present, otherwise the Fact list (backward-compatible for
level-aggregation consumers).

---

## Verification primitives

```go
textmatch.Normalize(s)                // NFC + whitespace fold + Trim (no case-fold — ToLower first if needed)
textmatch.Contains(source, token)     // substring after normalization
textmatch.MissingTokens(source, toks) // tokens absent from source (hallucination block)

temporal.Resolve(spec, ref)           // structured Spec (calendar/components/offset) + ref time.Time (injected now) → Gregorian ISO (undetermined ⇒ Determined=false)
temporal.ComponentsInAnchor(...)      // extract time components from an anchor string
```

Match source-language anchors with `Normalize` (combining marks, NFC); the recommended pattern is to
unify the output `value` in English.

---

## Conventions (philosophy)

- **State the deterministic gate** — judge from input alone; only the machine locks PASS.
- **Rules are violation detectors; severity is a level** — never fake a hard check with a weight
  (continuous weighting is for the graph's genuine contest only).
- **Cheese defense first** — for every answer to "how would I fool this gate?", add one rule (audited
  by the auto catalog).
- **Side effects through ground, rules stay pure** — the network is owned by reins ground primitives,
  isolated by staged eval.
- **No N=1 abstraction** — freeze a new abstraction only after a second consumer validates it.

## Quick decision guide

| Situation | Use |
|---|---|
| Independent rules, simple gate | `gate.Rule` + `Rules()` (level aggregation) |
| Inter-rule precedence / root-cause feedback | `pkg/graph` + `Evaluator` + `Supersedes` |
| Side-effect verification (HTTP/DNS) | `pkg/ground` + `.Needs()` + `EvaluateStaged` |
| Block body hallucination | `textmatch.MissingTokens` |
| Date/time normalization | `pkg/temporal` |
| Short-circuit an untrusted submission | `Prepare`'s `short` verdict (`OutSkip`/`OutBlock`) |
| A streaming source instead of a one-shot seed | the consumer adds a `run` command (not yet shipped by reins) |
| Unattended drive (LLM generates, gate judges) | `Options{Loop}` + `pkg/llm` — the opt-in `loop` command; rule-specific coaching via `RuleSystem`/`Verdict.RootCause` |
