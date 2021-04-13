# Miller regression tests

There are a few files unit-tested with Go's `testing` package -- a few dozen cases total.

The vast majority of Miller tests, though -- thousands of cases -- are tested by running scripted invocations of `mlr` with various flags and inputs, comparing against expected output, and checking the exit code back to the shell.

## How to run the regression tests, in brief

*Note: while this `README.md` file is within the `go/reg-test/` subdirectory, all paths in this file are written from the perspective of the user being cd'ed into the `go/` directory, i.e. this directory's parent directory.*

* `mlr regtest --help`

* `go test` -- TODO -- also comment

## Items for the duration of the Go port

* `mlr regtest -c ...` runs the C version of Miller from the local checkout

## More details

TODO: needs to be written up

## Creating new cases

TODO: needs to be written up
