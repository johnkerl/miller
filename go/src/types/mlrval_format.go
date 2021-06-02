package types

import (
	"errors"
	"fmt"
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
		return nil, err
	}

	mlrvalFormatterCache[userLevelFormatString] = formatter
	return formatter, nil
}

// ----------------------------------------------------------------
type IMlrvalFormatter interface {
	Format(mlrval *Mlrval) *Mlrval
}

func newMlrvalFormatter(
	userLevelFormatString string,
) (IMlrvalFormatter, error) {
	// TODO: very temporary. Pending full parse.
	// Including but not limited to "%08lld" -> "%08d" C-impl back-compat etc.
	if strings.HasSuffix(userLevelFormatString, "d") {
		return newMlrvalFormatterToInt(userLevelFormatString), nil
	}
	if strings.HasSuffix(userLevelFormatString, "x") {
		return newMlrvalFormatterToInt(userLevelFormatString), nil
	}

	if strings.HasSuffix(userLevelFormatString, "f") {
		return newMlrvalFormatterToFloat(userLevelFormatString), nil
	}
	if strings.HasSuffix(userLevelFormatString, "e") {
		return newMlrvalFormatterToFloat(userLevelFormatString), nil
	}
	if strings.HasSuffix(userLevelFormatString, "g") {
		return newMlrvalFormatterToFloat(userLevelFormatString), nil
	}

	if strings.HasSuffix(userLevelFormatString, "s") {
		return newMlrvalFormatterToString(userLevelFormatString), nil
	}

	return nil, errors.New("TBD") // TODO
}

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
	return MlrvalPointerFromString(
		fmt.Sprintf(
			formatter.goFormatString,
			mlrval.String(),
		),
	)
}
