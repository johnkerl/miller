# Miller regression tests

There are a few files unit-tested with Go's `testing` package -- a few dozen cases total.

The vast majority of Miller tests, though -- thousands of cases -- are tested by running scripted invocations of `mlr` with various flags and inputs, comparing against expected output, and checking the exit code back to the shell.

## How to run the regression tests, in brief

*Note: while this `README.md` file is within the `test/` subdirectory, all paths in this file are written from the perspective of the user being cd'ed into the repository base directory, i.e. this directory's parent directory.*

* `mlr regtest --help`

* `go test github.com/johnkerl/miller/v6/pkg/...` — runs the Go unit tests (a few dozen cases).

## Items for the duration of the Go port

* `mlr regtest -c ...` runs the C version of Miller from the local checkout

## More details

You can alias `mr='mlr regtest'` for convenience. With no arguments, `mr` runs all cases under `test/cases/`. Pass one or more paths to run only those directories or specific `.cmd` files.

* **`mr`** — run all regression cases (default path is `test/cases/`).
* **`mr test/cases/foo`** — run only cases under that directory.
* **`mr -v test/cases/foo`** — same, plus per-command pass/fail; use `-vv` or `-vvv` for more detail.
* **`mr -j test/cases/foo/0003`** — show the Miller command, any script, and actual output for that case (handy for debugging).
* **`mr -p test/cases/foo/0003`** — *populate*: write or overwrite `expout` and `experr` from the current run (use when adding or updating expected output).
* **`mr -c ...`** — use the C build of Miller (e.g. `-c` → `../c/mlr`) instead of the current executable.

To review populated files before committing, run `mr -p` on the desired path, then `git diff` to inspect changes and `git reset --hard` to discard them.

## Creating new cases

1. Create a case directory under `test/cases/`, e.g. `test/cases/my-feature/0001`.
2. Add a **`cmd`** file containing the Miller command line (one line), e.g. `mlr cat test/input/simple.dkvp`.
3. Use shared input under `test/input/`, or add a local **`input`** file in the case directory; in `cmd` you can use **`${CASEDIR}`** so the command refers to the case directory (e.g. `mlr cat ${CASEDIR}/input`).
4. Run **`mlr regtest -p test/cases/my-feature/0001`** to generate **`expout`** (and **`experr`** if the command produces stderr). If the command is expected to exit non-zero, add an empty **`should-fail`** file.
5. Run **`mlr regtest test/cases/my-feature/0001`** (without `-p`) to confirm the case passes.

Optional: **`mlr`** — DSL script file when the test uses `-f`/`put`/`filter`; **`env`** — environment variables to set for the case (unset after).
