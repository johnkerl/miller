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

	"mlr/src/auxents/regtest"
)

func TestRegression(t *testing.T) {
	regtester := regtest.NewRegTester(
		"mlr", // exeName
		false, // doPopulate
		0,     // verbosityLevel
		false, // plainMode
		0,     // firstNFailsToShow
	)

	paths := []string{} // use default

	ok := regtester.Execute(paths)
	if !ok {
		t.Fatal()
	}
}
