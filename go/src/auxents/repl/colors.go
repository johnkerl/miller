// ================================================================
// This is for doing color-highlighting for on-line help.  This is a thin layer
// over the package is uses -- the value-add here is to centralize the choice
// of particular color (as of this writing, red) all in one spot.
// ================================================================

package repl

import (
	"miller/src/platform"
)

func PrintHighlightString(input string) {
	platform.PrintHiRed(input)
}
