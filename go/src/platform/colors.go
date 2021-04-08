package platform

import (
	"fmt"

	"github.com/fatih/color"
)

// Outside of Windows, we don't need the Fprintf to color.Output (and in fact
// color.Output is simply os.Stdout inside fatih/color). But inside Windows we
// do. That's why these functions take a string and print it, returning void --
// rather than taking a string and color-decorating it and returning the
// decoration as a string.

func PrintHiGreen(s string) {
	fmt.Fprintf(color.Output, color.HiGreenString(s))
}

func PrintHiRed(s string) {
	fmt.Fprintf(color.Output, color.HiRedString(s))
}
