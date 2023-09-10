package mlrval

import (
	"fmt"
	"strconv"
	"strings"
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
// ----------------------------------------------------------------

// ----------------------------------------------------------------
//* callsites:
//  o fmtnum($mv, "%d")
//    - numeric only
//  o format($mv, "%s")
//    - make this new DSL function
//  o --ofmt
//    - numeric only
//  k format-values verb
//    - -i, -f, -s
// ----------------------------------------------------------------

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

	// TODO: perhaps a full format-string parser. At present, there's nothing to stop people
	// from doing silly things like "%lllld".
	goFormatString := userLevelFormatString
	goFormatString = strings.ReplaceAll(goFormatString, "lld", "d")
	goFormatString = strings.ReplaceAll(goFormatString, "llx", "x")
	goFormatString = strings.ReplaceAll(goFormatString, "ld", "d")
	goFormatString = strings.ReplaceAll(goFormatString, "lx", "x")
	goFormatString = strings.ReplaceAll(goFormatString, "lf", "f")
	goFormatString = strings.ReplaceAll(goFormatString, "le", "e")
	goFormatString = strings.ReplaceAll(goFormatString, "lg", "g")

	// MIller 5 and below required C format strings compatible with 64-bit ints
	// and double-precision floats: e.g. "%08lld" and "%9.6lf". For Miller 6,
	// We must still accept these for backward compatibility.
	if strings.HasSuffix(goFormatString, "d") {
		return newFormatterToInt(goFormatString), nil
	}
	if strings.HasSuffix(goFormatString, "x") {
		return newFormatterToInt(goFormatString), nil
	}

	if strings.HasSuffix(goFormatString, "f") {
		return newFormatterToFloat(goFormatString), nil
	}
	if strings.HasSuffix(goFormatString, "e") {
		return newFormatterToFloat(goFormatString), nil
	}
	if strings.HasSuffix(goFormatString, "g") {
		return newFormatterToFloat(goFormatString), nil
	}

	if strings.HasSuffix(goFormatString, "s") {
		return newFormatterToString(goFormatString), nil
	}

	// TODO:
	// return nil, errors.New(fmt.Sprintf("unhandled format string \"%s\"", userLevelFormatString))
	return newFormatterToString(goFormatString), nil
}

// ----------------------------------------------------------------

type formatterToFloat struct {
	goFormatString string
}

func newFormatterToFloat(goFormatString string) IFormatter {
	return &formatterToFloat{
		goFormatString: goFormatString,
	}
}

func (formatter *formatterToFloat) Format(mv *Mlrval) *Mlrval {
	floatValue, isFloat := mv.GetFloatValue()
	if isFloat {
		formatted := fmt.Sprintf(formatter.goFormatString, floatValue)
		return TryFromFloatString(formatted)
	}
	intValue, isInt := mv.GetIntValue()
	if isInt {
		formatted := fmt.Sprintf(formatter.goFormatString, float64(intValue))
		return TryFromFloatString(formatted)
	}
	return mv
}

func (formatter *formatterToFloat) FormatFloat(floatValue float64) string {
	return fmt.Sprintf(formatter.goFormatString, floatValue)
}

// ----------------------------------------------------------------

type formatterToInt struct {
	goFormatString string
}

func newFormatterToInt(goFormatString string) IFormatter {
	return &formatterToInt{
		goFormatString: goFormatString,
	}
}

func (formatter *formatterToInt) Format(mv *Mlrval) *Mlrval {
	intValue, isInt := mv.GetIntValue()
	if isInt {
		formatted := fmt.Sprintf(formatter.goFormatString, intValue)
		return TryFromIntString(formatted)
	}
	floatValue, isFloat := mv.GetFloatValue()
	if isFloat {
		formatted := fmt.Sprintf(formatter.goFormatString, int(floatValue))
		return TryFromIntString(formatted)
	}
	return mv
}

func (formatter *formatterToInt) FormatFloat(floatValue float64) string {
	return fmt.Sprintf(formatter.goFormatString, int(floatValue))
}

// ----------------------------------------------------------------

type formatterToString struct {
	goFormatString string
}

func newFormatterToString(goFormatString string) IFormatter {
	return &formatterToString{
		goFormatString: goFormatString,
	}
}

func (formatter *formatterToString) Format(mv *Mlrval) *Mlrval {
	return FromString(
		fmt.Sprintf(
			formatter.goFormatString,
			mv.String(),
		),
	)
}

func (formatter *formatterToString) FormatFloat(floatValue float64) string {
	return strconv.FormatFloat(floatValue, 'g', -1, 64)
}
