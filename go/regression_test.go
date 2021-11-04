package main

import (
	"fmt"
	"os"
	"testing"

	"mlr/src/auxents/regtest"
)

// TestRegression is a familiar entry point for regression testing.  Miller
// regression tests are more flexibly invoked via 'mlr regtest'.  However here
// is a standard location so people can get at them via 'go test'.  Please see
// (as of this writing) src/auxents/regtest for the Miller regtest package.
func TestRegression(t *testing.T) {
	// How much detail to show?  There are thousands of cases, organized into a
	// few hundred top-level directories under ./regtest/cases.
	//
	// Default behavior is to show PASS/FAIL for those top-level directories.
	// If (for whatever reason) lots of tests are systematically failing then
	// verbosityLevel = 3 for all cases is probably too much output to be
	// useful.
	//
	// What we do is:
	// * go test:    verbosityLevel = 0
	// * go test -v: firstNFailsToShow = 20 which *re-runs* up to 20 failing
	//               tests (after having failed once with verbosityLevel = 0)
	//               with verbosityLevel = 3 on the re-run.
	//
	// This is (I hope) sufficient flexibility for use in GitHub Actions
	// continuous-integration jobs. If more detail is needed then one may
	// run 'mlr regtest -v', 'mlr regtest -vv', 'mlr regtest -vvv' instead
	// of going through 'go test'.
	firstNFailsToShow := 0
	if testing.Verbose() {
		firstNFailsToShow = 20
	}

	// Let the tests find ./mlr
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "mlr: could not find current working directory.")
		os.Exit(1)
	}
	path := os.Getenv("PATH")
	os.Setenv("PATH", cwd + ":" + path)

	regtester := regtest.NewRegTester(
		"mlr", // exeName
		false, // doPopulate
		0,     // verbosityLevel
		false, // plainMode
		firstNFailsToShow,
	)

	paths := []string{} // use default

	ok := regtester.Execute(paths)
	if !ok {
		t.Fatal()
	}
}
