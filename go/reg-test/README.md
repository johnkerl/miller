# Miller regression tests

There are a few files unit-tested with Go's `testing` package -- a few dozen cases total.

The vast majority of Miller tests, though -- thousands of cases -- are tested by running scripted invocations of `mlr` with various flags and inputs, comparing against expected output, and checking the exit code back to the shell.

## How to run them, in brief

Note: while this `README.md` file is within the `go/reg-test/` subdirectory, all paths in this file are written from the perspective of the user being cd'ed into the `go/` directory, i.e. this directory's parent directory.

* `rr` in the Miller `go/` subdirectory is an alias-like script for `reg-test/run`
* `rr --help` for help
* Without `-v`, this runs all cases with a pass/fail indication per case, and an overall pass/fail indication -- overall pass only if all cases pass.

## More details

* `reg-test/cases/case*.sh` files consist of "invocations" of Miller
* Each case `./reg-test/cases/case-foo.sh` has expected output `./reg-test/expected/case-foo.sh.out` and `./reg-test/expected/case-foo.sh.err`, along with actual outputs `output-reg-test/case-c-cat.sh.out` and `output-reg-test/case-c-cat.sh.err`.
* The `reg-test/run` script loops over all of the `case-*.sh` files and executes them via
  sourcing them with the Bash `.` operator.
* Each has lines of the form `run_mlr ...` or `mlr_expect_fail ...`. Those functions are defined in `reg_test/run`. They take `mlr` command-lines as arguments.
* Each case can fail in the following ways:
  * Zero invocations were attempted.
  * A given `run_mlr ...` invocation exits with non-zero when it should exit with zero.
  * A given `mlr_expect_fail ...` invocation exits with zero when it should exit with non-zero.
  * The output of the invocations in the case's actual-output file differs from the case's expected-output file.

## Debugging failures of existing cases

* If a case fails, you can run it by itself with `-v` if you like: e.g. `./reg-test/run -v reg-test/cases/case-cat.sh`.
* Also `-C 1` or `-C 5`, etc. (note the space) to control number of context lines in the diff output.

## Creating new cases

* Edit `reg-test/cases/case-new-name-goes-here.sh`. Note the `reg-test/cases` directory path, the filename starting with `case-`, and the filename ending with `.sh` are all required.
* Run `reg-test/run reg-test/cases/case-new-name-goes-here.sh`
* That will create `output-regtest/case-new-name-goes-here.sh.out`
* If this all looks OK, `accept-case case-new-name-goes-here.sh` which will copy actual output to `reg-test/expected/case-new-name-goes-here.sh.out`
* Add the `case*sh` and the expected-output file to source control.
