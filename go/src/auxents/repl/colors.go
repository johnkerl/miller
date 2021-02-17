// ================================================================
// This is for doing color-highlighting for on-line help
// ================================================================

package repl

import (
	"github.com/fatih/color"
)

func HighlightString(input string) string {
	return color.HiRedString(input)
}
