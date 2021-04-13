// ================================================================
// Miller regression tests are flexibly invoked via 'mlr regtest'.
// However here is a standard location so people can get at them
// via 'go test'.
// ================================================================

package main

import (
	"testing"

	"miller/src/lib"
	"miller/src/auxents/regtest"
)

func TestRegression(t *testing.T) {

	regtester := regtest.NewRegTester(
		lib.MlrExeName(),
		"regression",
		false, // doPopulate
		0, // verbosityLevel
		0, // firstNFailsToShow
	)

	paths := []string{} // use default

	ok := regtester.Execute(paths)
	if !ok {
		t.Fatal()
	}
}
