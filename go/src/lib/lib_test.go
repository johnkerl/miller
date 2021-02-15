// ================================================================
// Most Miller tests (thousands of them) are command-line-driven via
// reg0test/run. Here are some cases needing special focus.
// ================================================================

package lib

import (
	"miller/src/lib"
	"testing"
)

func TestRegexReplaceOnce(t *testing.T) {
	regexString := "[a-z]"
	regex := lib.CompileMillerRegexOrDie(regexString)
	replacement := "X"

	input := "abcde"

	gsubOutput := regex.ReplaceAllString(input, replacement)
	if gsubOutput != "XXXXX" {
		t.Fatal()
	}

	subOutput := lib.RegexReplaceOnce(regex, input, replacement)
	if subOutput != "Xbcde" {
		t.Fatal()
	}
}
