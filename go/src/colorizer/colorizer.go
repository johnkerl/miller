// ================================================================
// Output-coloring for Miller
//
// Please see mlr --usage-output-colorization for context.
// ================================================================

package colorizer

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mattn/go-isatty"
)

// ================================================================
// External API

// Enum-ish type for when to apply output-coloring
type TOutputColorization int

const (
	ColorizeOutputNever = iota
	ColorizeOutputIfTTY
	ColorizeOutputAlways
)

// For command-line flags like --no-color and --always-color
func SetColorization(arg TOutputColorization) {
	colorization = arg
}

// For command-line flags like --pass-color 208 etc
func SetKeyColor(i int) {
	keyColorString = makeColorString(i)
}
func SetValueColor(i int) {
	valueColorString = makeColorString(i)
}
func SetPassColor(i int) {
	passColorString = makeColorString(i)
}
func SetFailColor(i int) {
	failColorString = makeColorString(i)
}
func SetREPLPS1Color(i int) {
	replPS1ColorString = makeColorString(i)
}
func SetREPLPS2Color(i int) {
	replPS2ColorString = makeColorString(i)
}
func SetHelpColor(i int) {
	helpColorString = makeColorString(i)
}

// For record-writer, regression-test, and online-help callsites
func MaybeColorizeKey(key string, outputIsStdout bool) string {
	return maybeColorize(key, keyColorString, outputIsStdout)
}

func MaybeColorizeValue(value string, outputIsStdout bool) string {
	return maybeColorize(value, valueColorString, outputIsStdout)
}

func MaybeColorizePass(text string, outputIsStdout bool) string {
	return maybeColorize(text, passColorString, outputIsStdout)
}

func MaybeColorizeFail(text string, outputIsStdout bool) string {
	return maybeColorize(text, failColorString, outputIsStdout)
}

func MaybeColorizeREPLPS1(text string, outputIsStdout bool) string {
	return maybeColorize(text, replPS1ColorString, outputIsStdout)
}

func MaybeColorizeREPLPS2(text string, outputIsStdout bool) string {
	return maybeColorize(text, replPS2ColorString, outputIsStdout)
}

func MaybeColorizeHelp(text string, outputIsStdout bool) string {
	return maybeColorize(text, helpColorString, outputIsStdout)
}

// For --list-colors command-line flag
func ListColors(o *os.File) {
	fmt.Println("Available colors:")
	for i := 0; i <= 255; i++ {
		fmt.Printf("%s%3d%s", makeColorString(i), i, defaultColorString)
		if i%16 < 15 {
			fmt.Print(" ")
		} else {
			fmt.Print("\n")
		}
	}
}

// ================================================================
// Internal implementation

// Default ANSI color codes
var keyColorString = make256ColorString(208)  // orange
var valueColorString = make256ColorString(33) // blue
var passColorString = make16ColorString(10)   // bold green
var failColorString = make16ColorString(9)    // bold red
var replPS1ColorString = make16ColorString(9) // bold red
var replPS2ColorString = make16ColorString(1) // red
var helpColorString = make16ColorString(9)    // bold red

// Used to switch back to default color
var defaultColorString = "\033[0m"

// Default: colorize if writing to stdout and if stdout is a TTY
var colorization TOutputColorization = ColorizeOutputIfTTY
var stdoutIsATTY = getStdoutIsATTY()

// Read environment variables at startup time. These can be overridden
// afterward using command-line flags.
func init() {
	if os.Getenv("MLR_NO_COLOR") != "" {
		colorization = ColorizeOutputNever
	} else if os.Getenv("MLR_ALWAYS_COLOR") != "" {
		colorization = ColorizeOutputAlways
	}

	var temp string
	var ok bool

	temp, ok = makeColorStringFromEnv("MLR_KEY_COLOR")
	if ok {
		keyColorString = temp
	}
	temp, ok = makeColorStringFromEnv("MLR_VALUE_COLOR")
	if ok {
		valueColorString = temp
	}
	temp, ok = makeColorStringFromEnv("MLR_PASS_COLOR")
	if ok {
		passColorString = temp
	}
	temp, ok = makeColorStringFromEnv("MLR_FAIL_COLOR")
	if ok {
		failColorString = temp
	}
	temp, ok = makeColorStringFromEnv("MLR_REPL_PS1_COLOR")
	if ok {
		replPS1ColorString = temp
	}
	temp, ok = makeColorStringFromEnv("MLR_REPL_PS2_COLOR")
	if ok {
		replPS2ColorString = temp
	}
	temp, ok = makeColorStringFromEnv("MLR_HELP_COLOR")
	if ok {
		helpColorString = temp
	}
}

func getStdoutIsATTY() bool {
	// Don't try ANSI color on Windows (except Cygwin)
	if os.Getenv("TERM") == "" {
		return false
	}
	if isatty.IsTerminal(os.Stdout.Fd()) {
		return true
	}
	if isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return true
	}
	return false
}

// make16ColorString constructs an ANSI-16-color-mode escape sequence
func make16ColorString(i int) string {
	i &= 15
	boldBit := (i >> 3) & 1
	colorBits := i & 7
	return fmt.Sprintf("\033[%d;%dm", boldBit, 30+colorBits)
}

// make256ColorString constructs an ANSI-256-color-mode escape sequence
func make256ColorString(i int) string {
	i &= 255
	return fmt.Sprintf("\033[1;38;5;%dm", i&255)
}

// makeColorString constructs an ANSI-16-color-mode escape sequence if arg is
// in 0..15, else ANSI-256-color-mode escape sequence
func makeColorString(i int) string {
	if 0 <= i && i <= 15 {
		return make16ColorString(i)
	} else {
		return make256ColorString(i)
	}
}

func makeColorStringFromEnv(envName string) (string, bool) {
	envValue := os.Getenv(envName)
	if envValue == "" {
		return "", false
	}
	i, err := strconv.Atoi(envValue)
	if err != nil {
		return "", false // TODO: return error?
	}
	if i < 0 {
		return "", false // TODO: return error?
	}
	i &= 255
	return makeColorString(i), true
}

func maybeColorize(text string, colorString string, outputIsStdout bool) string {
	if outputIsStdout && stdoutIsATTY {
		if colorization == ColorizeOutputNever {
			return text
		} else {
			return colorize(text, colorString)
		}
	} else {
		if colorization == ColorizeOutputAlways {
			return colorize(text, colorString)
		} else {
			return text
		}
	}
}

func colorize(text string, colorString string) string {
	return colorString + text + defaultColorString
}
