# gosec scoping notes

`gosec` (github.com/securego/gosec/v2), run against `./pkg/... ./cmd/mlr/...`
(same scope as `make lint`/`make staticcheck`, which excludes the stale,
non-shipped example code under `cmd/experiments/`, `docs/site/`, and
`docs/src/miller-as-library/`).

Run via `make gosec`. Background: this replaced an earlier attempt at using
`gokart` (github.com/praetorian-inc/gokart), which is unmaintained/archived
and can't parse export data from Go 1.24+ toolchains; `gosec` is actively
maintained and its taint-tracking rules (G702/G703) cover the same ground
gokart did, plus more.

## Baseline (first run)

- **237 raw issues** before any exclusions or `#nosec` annotations
- **53 issues** after scoping (see below)

### Rules excluded globally (via `-exclude=` in the `gosec` Makefile target)

These don't represent real risk for a CLI data-processing tool and were the
overwhelming majority of raw noise; every sampled instance below was
reviewed before excluding the rule wholesale.

| Rule | Count | Why excluded |
|------|------:|--------------|
| G104 | 151 | Unhandled errors -- already tuned precisely via `errcheck.exclude-functions` in `.golangci.yml` (bufio/strings.Builder writes, Fprint family) |
| G115 | 17 | Integer overflow on numeric conversions -- inherent to `Mlrval`'s int64/uint64/rune/byte representation handling; every sampled site (`mlrval_infer.go`, `lib/util.go`, `bifs/bits.go`, `latin1.go`) is a benign value-representation conversion, not attacker-controlled-length arithmetic |
| G602 | 5 | "slice index out of range" -- gosec can't prove loop-bound invariants; every sampled site (`emit_emitp.go`, `help/entry.go`) is a loop-bounded index |
| G301/G302/G306 | 8 | File/dir permissions (wants 0600/0750) -- output data files and directories are meant to be group/world-readable by design (0644/0755); these aren't secrets |

### Rules `#nosec`'d at specific sites (not excluded globally)

Kept enabled everywhere else, so any *new* misuse still gets flagged.

| Rule | Site | Why |
|------|------|-----|
| G401/G501/G505 | `pkg/bifs/hashing.go` (`md5()`/`sha1()` DSL builtins) | Exposed to users as data-fingerprinting/checksum functions, never for authentication or integrity-of-secrets purposes |
| G404 | `pkg/lib/rand.go` (`SeedRandom`, package-level generator) | Must be deterministically seedable for `--seed`/reproducible `urand()` output; `crypto/rand` can't do that by definition |

### Remaining after scoping (53) -- live findings, triage as normal work

| Rule | Count | Category |
|------|------:|----------|
| G703 | 22 | Path traversal via taint analysis |
| G304 | 20 | File path from variable (same sites as G703, mostly) |
| G204 | 6 | Subprocess launched with tainted args |
| G702 | 2 | Command injection via taint analysis (same sites as G204) |
| G107 | 1 | HTTP request with variable URL (`pkg/lib/file_readers.go:74`, the `mlr --from http://...` feature) |

Most of the G304/G703 file-path findings are the expected shape for a CLI
tool: the "taint source" is `os.Args`/CLI flags, which is Miller's own
trust boundary -- a user passing `mlr --from ./whatever` is not an attacker.
The two sites worth a closer look precisely because they're *not* just the
regtest/REPL harnesses:

- `pkg/bifs/system.go` -- the DSL's `system()` builtin (subprocess from DSL
  expression, i.e. from data the user's *script* constructs, not just argv)
- `pkg/terminals/mcp/tools_exec.go:94` -- MCP tool execution; worth
  double-checking exactly what can reach this before treating it as
  low-risk by the "it's just argv" argument above

`G107` (`file_readers.go:74`) is the documented `--from http://...` /
`--from https://...` feature -- also CLI-argv-driven today, but worth
revisiting if Miller's HTTP-fetch path is ever reachable from anything
other than the operator's own command line (e.g. if driven by data from a
record rather than a flag).
