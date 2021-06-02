// ================================================================
// Most Miller tests (thousands of them) are command-line-driven via
// mlr regtest. Here are some cases needing special focus.
// ================================================================

package lib

import (
	"testing"
)

func TestRegexReplaceOnce(t *testing.T) {
	regexString := "[a-z]"
	regex := CompileMillerRegexOrDie(regexString)
	replacement := "X"

	input := "abcde"

	gsubOutput := regex.ReplaceAllString(input, replacement)
	if gsubOutput != "XXXXX" {
		t.Fatal()
	}

	subOutput := RegexReplaceOnce(regex, input, replacement)
	if subOutput != "Xbcde" {
		t.Fatal()
	}
}
