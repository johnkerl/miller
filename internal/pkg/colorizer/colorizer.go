// ================================================================
// Output-coloring for Miller
//
// Please see mlr --usage-output-colorization for context.
//
// Note: code-share with github.com/johnkerl/lumin.
// ================================================================

package colorizer

import (
	"os"

	lumin "github.com/johnkerl/lumin/pkg/colors"
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
func SetKeyColor(name string) bool {
	escape, ok := lumin.MakeANSIEscapesFromName(name)
	if ok {
		keyColorString = escape
	}
	return ok
}
func SetValueColor(name string) bool {
	escape, ok := lumin.MakeANSIEscapesFromName(name)
	if ok {
		valueColorString = escape
	}
	return ok
}
func SetPassColor(name string) bool {
	escape, ok := lumin.MakeANSIEscapesFromName(name)
	if ok {
		passColorString = escape
	}
	return ok
}
func SetFailColor(name string) bool {
	escape, ok := lumin.MakeANSIEscapesFromName(name)
	if ok {
		failColorString = escape
	}
	return ok
}
func SetREPLPS1Color(name string) bool {
	escape, ok := lumin.MakeANSIEscapesFromName(name)
	if ok {
		replPS1ColorString = escape
	}
	return ok
}
func SetREPLPS2Color(name string) bool {
	escape, ok := lumin.MakeANSIEscapesFromName(name)
	if ok {
		replPS2ColorString = escape
	}
	return ok
}
func SetHelpColor(name string) bool {
	escape, ok := lumin.MakeANSIEscapesFromName(name)
	if ok {
		helpColorString = escape
	}
	return ok
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

// ListColorCodes shows codes in the range 0..255.
// For --list-color-codes command-line flag.
func ListColorCodes() {
	lumin.ListColorCodes()
}

// ListColorNames shows names for codes in the range 0..255.
// For --list-color-names command-line flag.
func ListColorNames() {
	lumin.ListColorNames()
}

// ================================================================
// Internal implementation

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

// GetColorization is for the CSV writer
func GetColorization(outputIsStdout bool, isKey bool) (string, string) {
	if outputIsStdout && stdoutIsATTY {
		if colorization == ColorizeOutputNever {
			return "", ""
		} else {
			if isKey {
				return keyColorString, defaultColorString
			} else {
				return valueColorString, defaultColorString
			}
		}
	} else {
		if colorization == ColorizeOutputAlways {
			if isKey {
				return keyColorString, defaultColorString
			} else {
				return valueColorString, defaultColorString
			}
		} else {
			return "", ""
		}
	}
}

// ================================================================
// Internal implementation

// Default ANSI color codes
// 6.0.0:
// var keyColorString = lumin.MakeANSIEscapesFromNameUnconditionally("orange")
// var valueColorString = lumin.MakeANSIEscapesFromNameUnconditionally("blue")
// 6.1.0:
var keyColorString = lumin.MakeANSIEscapesFromNameUnconditionally("bold-underline")
var valueColorString = lumin.MakeANSIEscapesFromNameUnconditionally("plain")
var passColorString = lumin.MakeANSIEscapesFromNameUnconditionally("bold-lime")
var failColorString = lumin.MakeANSIEscapesFromNameUnconditionally("bold-red")
var replPS1ColorString = lumin.MakeANSIEscapesFromNameUnconditionally("bold-red")
var replPS2ColorString = lumin.MakeANSIEscapesFromNameUnconditionally("red")
var helpColorString = lumin.MakeANSIEscapesFromNameUnconditionally("bold-red")

// Used to switch back to default color
var defaultColorString = "\u001b[0m"

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

func makeColorStringFromEnv(envName string) (string, bool) {
	envValue := os.Getenv(envName)
	if envValue == "" {
		return "", false
	}

	return lumin.MakeANSIEscapesFromName(envValue)
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
