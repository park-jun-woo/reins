// Package cli is reins' Cobra command scaffold (reins Phase004). It implements the
// how-make-quest canonical command skeleton — scan / next / submit / status / export
// — plus reins' own `rules` (the auto rulebook). A new quest gets the whole standard
// CLI by plugging in a gate.Definition (Seed/Render/Prepare/Rules); everything else
// (the ratchet, level aggregation, export) is driven by reins.
//
// cli does IO, parsing, and output formatting only — it never judges state. The
// verdict comes from gate (deterministic rule aggregation) and the transition from
// quest.Apply (pure). Commands use RunE and cmd.OutOrStdout/InOrStdin for testability.
package cli
