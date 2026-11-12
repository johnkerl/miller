package mlrval

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//----------------------------------------------------------------
// TODO
//* need int/float
//  llx -> x etc
//  https://golang.org/pkg/fmt/
//
//  pre-stuff
//
//  %
//
//  +-0' space
//
//  ll|l
//  %%
//  bdiouxDOUX fegFEG s
//
//  post-stuff

//* callsites:
//  o fmtnum($mv, "%d")
//    - numeric only
//  o format($mv, "%s")
//    - make this new DSL function
//  o --ofmt
//    - numeric only
//  k format-values verb
//    - -i, -f, -s

// Nil means use default format.
// Set from the CLI parser using mlr --ofmt.
var floatOutputFormatter IFormatter = nil

func SetFloatOutputFormat(formatString string) error {
	formatter, err := GetFormatter(formatString)
	if err != nil {
		return err
	}
	floatOutputFormatter = formatter
	return nil
}

var formatterCache map[string]IFormatter = make(map[string]IFormatter)

type IFormatter interface {
	Format(mlrval *Mlrval) *Mlrval
	FormatFloat(floatValue float64) string // for --ofmt
}

func GetFormatter(
	userLevelFormatString string,
) (IFormatter, error) {
	// Cache hit
	formatter, ok := formatterCache[userLevelFormatString]
	if ok {
		return formatter, nil
	}

	// Cache miss
	formatter, err := newFormatter(userLevelFormatString)
	if err != nil {
		// TODO: temp exit
		fmt.Printf("mlr: %v\n", err)
		return nil, err
	}

	formatterCache[userLevelFormatString] = formatter
	return formatter, nil
}

// People can pass in things like "X%sX" unfortunately :(
func newFormatter(
	userLevelFormatString string,
) (IFormatter, error) {
	numPercents := strings.Count(userLevelFormatString, "%")
	if numPercents < 1 {
		return nil, fmt.Errorf("unhandled format string \"%s\": no leading \"%%\"", userLevelFormatString)
	}
	if numPercents > 1 {
		return nil, fmt.Errorf(
			"unhandled format string \"%s\": needs no \"%%\" after the first", userLevelFormatString,
		)
	}

	// Miller 5 and below required C format strings compatible with 64-bit ints
	// and double-precision floats: e.g. "%08lld" and "%9.6lf". For Miller 6,
	// we must still accept these for backward compatibility.
	//
	// The kind of formatter depends on the conversion verb after the '%', not
	// on the end of the string: there may be text after the verb, as in
	// "%d items" or "[%.3f]".
	percentIndex := strings.Index(userLevelFormatString, "%")
	prefix := userLevelFormatString[:percentIndex+1]
	rest := userLevelFormatString[percentIndex+1:]

	// Flags, width, and precision
	n := 0
	for n < len(rest) && strings.IndexByte("+-# 0123456789.", rest[n]) >= 0 {
		n++
	}
	spec := rest[:n]
	rest = rest[n:]

	// Miller-specific "%_d" and "%_f" for thousands separators
	separated := strings.HasPrefix(rest, "_")
	if separated {
		rest = rest[1:]
	}

	// C length modifiers such as the "ll" in "%08lld" and the "l" in "%9.6lf"
	modifiers := strings.TrimLeft(rest, "l")
	if modifiers != rest && modifiers != "" && strings.IndexByte("diuUDxXoObeEfFgG", modifiers[0]) >= 0 {
		rest = modifiers
	}

	// C's "%i" is a synonym for "%d", and "%D" and "%U" are obsolete synonyms for "%ld" and "%lu". Go has no
	// "%i", "%u", "%D", or "%U". Miller ints are signed 64-bit, so "%u" and "%U" are treated as "%d", the same
	// way "%x" already formats a negative int with a minus sign rather than as two's complement.
	if rest != "" && strings.IndexByte("iuUD", rest[0]) >= 0 {
		rest = "d" + rest[1:]
	}

	if rest == "" || (separated && rest[0] != 'd' && rest[0] != 'f') {
		return newFormatterToString(userLevelFormatString), nil
	}
	goFormatString := prefix + spec + rest

	switch rest[0] {
	case 'd':
		if separated {
			return newFormatterToSeparatedInt(goFormatString), nil
		}
		return newFormatterToInt(goFormatString), nil
	case 'x', 'X', 'o', 'O', 'b':
		return newFormatterToInt(goFormatString), nil
	case 'f':
		if separated {
			return newFormatterToSeparatedFloat(goFormatString), nil
		}
		return newFormatterToFloat(goFormatString), nil
	case 'F', 'e', 'E', 'g', 'G':
		return newFormatterToFloat(goFormatString), nil
	}

	// TODO:
	// return nil, errors.New(fmt.Sprintf("unhandled format string \"%s\"", userLevelFormatString))
	return newFormatterToString(userLevelFormatString), nil
}

type formatterToFloat struct {
	goFormatString string
}

func newFormatterToFloat(goFormatString string) IFormatter {
	return &formatterToFloat{
		goFormatString: goFormatString,
	}
}

func (fmtr *formatterToFloat) Format(mv *Mlrval) *Mlrval {
	floatValue, isFloat := mv.GetFloatValue()
	if isFloat {
		formatted := fmt.Sprintf(fmtr.goFormatString, floatValue)
		return TryFromFloatString(formatted)
	}
	intValue, isInt := mv.GetIntValue()
	if isInt {
		formatted := fmt.Sprintf(fmtr.goFormatString, float64(intValue))
		return TryFromFloatString(formatted)
	}
	return mv
}

func (fmtr *formatterToFloat) FormatFloat(floatValue float64) string {
	return fmt.Sprintf(fmtr.goFormatString, floatValue)
}

func getLanguageTag() language.Tag {
	v, ok := os.LookupEnv("LANG")
	if ok {
		return language.Make(v)
	}
	return language.Make("en")
}

type formatterToSeparatedInt struct {
	goFormatString string
	printer        *message.Printer
}

func newFormatterToSeparatedInt(goFormatString string) IFormatter {
	return &formatterToSeparatedInt{
		goFormatString: goFormatString,
		printer:        message.NewPrinter(getLanguageTag()),
	}
}

func (fmtr *formatterToSeparatedInt) Format(mv *Mlrval) *Mlrval {
	intValue, isInt := mv.GetIntValue()
	if isInt {
		formatted := fmtr.printer.Sprintf(fmtr.goFormatString, intValue)
		return TryFromIntString(formatted)
	}
	floatValue, isFloat := mv.GetFloatValue()
	if isFloat {
		formatted := fmtr.printer.Sprintf(fmtr.goFormatString, int(floatValue))
		return TryFromIntString(formatted)
	}
	return mv
}

func (fmtr *formatterToSeparatedInt) FormatFloat(floatValue float64) string {
	return fmtr.printer.Sprintf(fmtr.goFormatString, int(floatValue))
}

type formatterToSeparatedFloat struct {
	goFormatString string
	printer        *message.Printer
}

func newFormatterToSeparatedFloat(goFormatString string) IFormatter {
	return &formatterToSeparatedFloat{
		goFormatString: goFormatString,
		printer:        message.NewPrinter(getLanguageTag()),
	}
}

func (fmtr *formatterToSeparatedFloat) Format(mv *Mlrval) *Mlrval {
	floatValue, isFloat := mv.GetFloatValue()
	if isFloat {
		formatted := fmtr.printer.Sprintf(fmtr.goFormatString, floatValue)
		return TryFromFloatString(formatted)
	}
	intValue, isInt := mv.GetIntValue()
	if isInt {
		formatted := fmtr.printer.Sprintf(fmtr.goFormatString, float64(intValue))
		return TryFromFloatString(formatted)
	}
	return mv
}

func (fmtr *formatterToSeparatedFloat) FormatFloat(floatValue float64) string {
	return fmtr.printer.Sprintf(fmtr.goFormatString, floatValue)
}

type formatterToInt struct {
	goFormatString string
}

func newFormatterToInt(goFormatString string) IFormatter {
	return &formatterToInt{
		goFormatString: goFormatString,
	}
}

func (fmtr *formatterToInt) Format(mv *Mlrval) *Mlrval {
	intValue, isInt := mv.GetIntValue()
	if isInt {
		formatted := fmt.Sprintf(fmtr.goFormatString, intValue)
		return TryFromIntString(formatted)
	}
	floatValue, isFloat := mv.GetFloatValue()
	if isFloat {
		formatted := fmt.Sprintf(fmtr.goFormatString, int(floatValue))
		return TryFromIntString(formatted)
	}
	return mv
}

func (fmtr *formatterToInt) FormatFloat(floatValue float64) string {
	return fmt.Sprintf(fmtr.goFormatString, int(floatValue))
}

type formatterToString struct {
	goFormatString string
}

func newFormatterToString(goFormatString string) IFormatter {
	return &formatterToString{
		goFormatString: goFormatString,
	}
}

func (fmtr *formatterToString) Format(mv *Mlrval) *Mlrval {
	return FromString(
		fmt.Sprintf(
			fmtr.goFormatString,
			mv.String(),
		),
	)
}

func (fmtr *formatterToString) FormatFloat(floatValue float64) string {
	return strconv.FormatFloat(floatValue, 'g', -1, 64)
}
