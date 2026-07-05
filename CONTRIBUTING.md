# Contributing to Miller

Thanks for your interest in contributing! This page is a quick-start guide for code, test, and
documentation contributions. It links out to more detailed docs rather than repeating them, so if
something below feels thin, follow the link.

## Asking questions / reporting issues

* Questions and general discussion: [GitHub Discussions](https://github.com/johnkerl/miller/discussions)
* Bug reports and feature requests: [GitHub Issues](https://github.com/johnkerl/miller/issues)
* See also [notes on issue-labeling](https://github.com/johnkerl/miller/wiki/Issue-labeling).

When filing a bug, please include your Miller version (`mlr --version`), OS, and a minimal repro
(command line plus a few lines of sample input).

## Getting the source and building

```
git clone https://github.com/johnkerl/miller
cd miller
make        # builds ./mlr
make check  # runs unit + regression tests
```

Miller is written in Go (1.25+ as of July 2026) with no non-standard-library runtime dependencies.
For build details beyond `make`/`go build`, see the [Building from
source](https://miller.readthedocs.io/en/latest/build) doc and [README-dev.md](./README-dev.md).

## Before opening a PR

Run the full developer workflow:

```
make dev
```

This formats code, builds, runs unit and regression tests, and regenerates
docs/man-pages, in the right order. A PR that passes this is in good shape to
submit.

A few conventions to follow:

* Code should read like idiomatic, `go fmt`-clean Go. Miller favors clarity
  over cleverness, including some acceptable code duplication: see
  [README-dev.md: Source-code goals](./README-dev.md#source-code-goals).
* Avoid line wraps at a terminal width of 120 columns, tabwidth 4.
* Prose in docs/comments uses the Oxford comma.
* Reference relevant issue numbers in commit messages where applicable.

## Where to look depending on what you're changing

* **Adding or fixing a built-in DSL function**: implement in `pkg/bifs/`,
  add unit tests alongside, update `docs/src/`.
* **Adding/changing a verb (transformer)**: see `pkg/transformers/` and the
  directory-structure notes in [README-dev.md](./README-dev.md).
* **Anything touching the DSL grammar** (`pkg/parsing/mlr.bnf`): see
  [README-dev.md](./README-dev.md) for the `tools/build-dsl` regeneration step.
* **Tests**: Miller's test suite is mostly scripted CLI invocations compared
  against expected output, plus a smaller set of Go unit tests. See
  [test/README.md](./test/README.md) for how to run and add cases.
* **Documentation**: the published docs are generated from `docs/src/*.md.in`
  files, not hand-edited `.md`. See [README-docs.md](./README-docs.md) for the
  edit/preview/build loop.
* **Performance work**: see [README-profiling.md](./README-profiling.md) for
  profiling Miller itself, and [scripts/perf/README.md](./scripts/perf/README.md)
  for the benchmark scripts used to produce the performance graphs in the docs.
* **Background on the Go port** (history, design rationale, why Go over C):
  see [README-go-port.md](./README-go-port.md). Not required reading to
  contribute, but useful context for anything touching core internals like
  `Mlrval` or the record-stream architecture.

## Using an AI coding assistant

If you're using Claude Code against this repo, [CLAUDE.md](./CLAUDE.md) has
build/test commands and conventions written for that workflow. It's a
supplement for AI-assisted contributions, not a replacement for this page.

## License

Miller is licensed under the [two-clause BSD license](./LICENSE.txt).
Contributions are accepted under the same license.
