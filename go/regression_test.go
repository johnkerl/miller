// ================================================================
// Miller regression tests are flexibly invoked via 'mlr regtest'.
// However here is a standard location so people can get at them
// via 'go test'.
//
// Please see (as of this writing) src/auxents/regtest for the Miller
// regtest package.
// ================================================================

package main

import (
	"testing"

	"miller/src/auxents/regtest"
	"miller/src/lib"
)

func TestRegression(t *testing.T) {
	regtester := regtest.NewRegTester(
		lib.MlrExeName(),
		false, // doPopulate
		0,     // verbosityLevel
		0,     // firstNFailsToShow
	)

	paths := []string{} // use default

	ok := regtester.Execute(paths)
	if !ok {
		t.Fatal()
	}
}
