package types

import (
	"errors"
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

// ----------------------------------------------------------------
var mlrvalFormatterCache map[string]IMlrvalFormatter = make(map[string]IMlrvalFormatter)

func GetMlrvalFormatter(
	userLevelFormatString string,
) (IMlrvalFormatter, error) {
	// Cache hit
	formatter, ok := mlrvalFormatterCache[userLevelFormatString]
	if ok {
		return formatter, nil
	}

	// Cache miss
	formatter, err := newMlrvalFormatter(userLevelFormatString)
	if err != nil {
		// TODO: temp exit
		fmt.Printf("mlr: %v\n", err)
		return nil, err
	}

	mlrvalFormatterCache[userLevelFormatString] = formatter
	return formatter, nil
}

// ----------------------------------------------------------------
type IMlrvalFormatter interface {
	Format(mlrval *Mlrval) *Mlrval
	FormatFloat(floatValue float64) string // for --ofmt
}

// People can pass in things like "X%sX" unfortunately :(
func newMlrvalFormatter(
	userLevelFormatString string,
) (IMlrvalFormatter, error) {
	numPercents := strings.Count(userLevelFormatString, "%")
	if numPercents < 1 {
		return nil, errors.New(
			fmt.Sprintf("unhandled format string \"%s\": no leading \"%%\"", userLevelFormatString),
		)
	}
	if numPercents > 1 {
		return nil, errors.New(
			fmt.Sprintf("unhandled format string \"%s\": needs no \"%%\" after the first", userLevelFormatString),
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
		return newMlrvalFormatterToInt(goFormatString), nil
	}
	if strings.HasSuffix(goFormatString, "x") {
		return newMlrvalFormatterToInt(goFormatString), nil
	}

	if strings.HasSuffix(goFormatString, "f") {
		return newMlrvalFormatterToFloat(goFormatString), nil
	}
	if strings.HasSuffix(goFormatString, "e") {
		return newMlrvalFormatterToFloat(goFormatString), nil
	}
	if strings.HasSuffix(goFormatString, "g") {
		return newMlrvalFormatterToFloat(goFormatString), nil
	}

	if strings.HasSuffix(goFormatString, "s") {
		return newMlrvalFormatterToString(goFormatString), nil
	}

	// TODO:
	// return nil, errors.New(fmt.Sprintf("unhandled format string \"%s\"", userLevelFormatString))
	return newMlrvalFormatterToString(goFormatString), nil
}

//func regularizeFormat

// ----------------------------------------------------------------
type mlrvalFormatterToFloat struct {
	goFormatString string
}

func newMlrvalFormatterToFloat(goFormatString string) IMlrvalFormatter {
	return &mlrvalFormatterToFloat{
		goFormatString: goFormatString,
	}
}

func (formatter *mlrvalFormatterToFloat) Format(mlrval *Mlrval) *Mlrval {
	floatValue, isFloat := mlrval.GetFloatValue()
	if isFloat {
		formatted := fmt.Sprintf(formatter.goFormatString, floatValue)
		return MlrvalTryPointerFromFloatString(formatted)
	}
	intValue, isInt := mlrval.GetIntValue()
	if isInt {
		formatted := fmt.Sprintf(formatter.goFormatString, float64(intValue))
		return MlrvalTryPointerFromFloatString(formatted)
	}
	return mlrval
}

func (formatter *mlrvalFormatterToFloat) FormatFloat(floatValue float64) string {
	return fmt.Sprintf(formatter.goFormatString, floatValue)
}

// ----------------------------------------------------------------
type mlrvalFormatterToInt struct {
	goFormatString string
}

func newMlrvalFormatterToInt(goFormatString string) IMlrvalFormatter {
	return &mlrvalFormatterToInt{
		goFormatString: goFormatString,
	}
}

func (formatter *mlrvalFormatterToInt) Format(mlrval *Mlrval) *Mlrval {
	intValue, isInt := mlrval.GetIntValue()
	if isInt {
		formatted := fmt.Sprintf(formatter.goFormatString, intValue)
		return MlrvalTryPointerFromIntString(formatted)
	}
	floatValue, isFloat := mlrval.GetFloatValue()
	if isFloat {
		formatted := fmt.Sprintf(formatter.goFormatString, int(floatValue))
		return MlrvalTryPointerFromIntString(formatted)
	}
	return mlrval
}

func (formatter *mlrvalFormatterToInt) FormatFloat(floatValue float64) string {
	return fmt.Sprintf(formatter.goFormatString, int(floatValue))
}

// ----------------------------------------------------------------
type mlrvalFormatterToString struct {
	goFormatString string
}

func newMlrvalFormatterToString(goFormatString string) IMlrvalFormatter {
	return &mlrvalFormatterToString{
		goFormatString: goFormatString,
	}
}

func (formatter *mlrvalFormatterToString) Format(mlrval *Mlrval) *Mlrval {
	return MlrvalFromString(
		fmt.Sprintf(
			formatter.goFormatString,
			mlrval.String(),
		),
	)
}

func (formatter *mlrvalFormatterToString) FormatFloat(floatValue float64) string {
	return strconv.FormatFloat(floatValue, 'g', -1, 64)
}
