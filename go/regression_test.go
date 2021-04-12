// ================================================================
// TODO: comment
// ================================================================

package main

import (
	"testing"

	"miller/src/lib"
	"miller/src/auxents/regtest"
)

func TestFoo(t *testing.T) {

	regtester := regtest.NewRegTester(
		lib.MlrExeName(),
		"regression",
		false, // doPopulate
		0, // verbosityLevel
		0, // firstNFailsToShow
	)

	paths := []string{}

	ok := regtester.Execute(paths)
	if !ok {
		t.Fatal()
	}
}
