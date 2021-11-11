// ================================================================
// Mlrval implements the Unmarshaler and Marshaler interfaces needed for
// marshaling/unmarshaling to/from JSON, via the UnmarshalJSON and MarshalJSON
// methods.
//
// Please see also https://golang.org/pkg/encoding/json/
// ================================================================

package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"mlr/internal/pkg/colorizer"
	"mlr/internal/pkg/lib"
)

const MLRVAL_JSON_INDENT_STRING string = "  "

type TJSONFormatting int

const (
	JSON_SINGLE_LINE = 1
	JSON_MULTILINE   = 2
)

// ================================================================
// The JSON decoder (https://golang.org/pkg/encoding/json/#Decoder) is quite
// nice. What we can have is:
//
// * Values: string, number, true, false, null
// * Arrays: '[' ... ']'
// * Objects: '{' ... '}'
//
// When we get the start token (if any): if it's a value that's it.
//
// Otherwise we got either array-start or object start.
//
// In the array case the decoder will deal with whitespace, commas, etc and
// return the array-values tokens while decoder.More(). Then we need to get the
// next token and verify it's close-bracket.
//
// For example on input '[1,2,3]', the start token is '['. While decoder.More()
// we get the 1,2,3. Then we read the next token and make sure it's ']'.
//
// In the object case the decoder will deal with whitespace, commas, etc and
// return the element tokens alternating key, value, key, value, while
// decoder.More(). Then we need to get the next token and verify it's
// close-bracket.
//
// For example on input '{"a":1,"b":2,"c"3]', the start token is '{'. While
// decoder.More() we get the "a",1,"b",2,"c",3. Then we read the next token and
// make sure it's '}'.
//
// In the case the input is of the form '[1,2,3][4,5,6]' -- i.e. a stream of
// valid JSON items -- we can use decoder.Buffered() to continue stream
// processing.
//
// ----------------------------------------------------------------
// json.Token holds one of the following types:
//
// * Delim, for the four JSON delimiters [ ] { }
// * bool, for JSON booleans
// * float64, for JSON numbers
// * Number, for JSON numbers
// * string, for JSON string literals
// * nil, for JSON null
//
// ----------------------------------------------------------------
// Note: we accept a sequence of valid JSON items, not just a JSON item.
// E.g. either
//
//   {
//     "a": 1,
//     "b": 2
//   }
//   {
//     "a": 3,
//     "b": 4
//   }
//
// or
//
//   [
//     {
//       "a": 1,
//       "b": 2
//     },
//     {
//       "a": 3,
//       "b": 4
//     }
//   ]
//
// This is so the Miller JSON record-reader can be streaming, not needing to
// ingest all records at once, and operable within a tail -f context.
// ================================================================

// ----------------------------------------------------------------
func (mv *Mlrval) UnmarshalJSON(inputBytes []byte) error {
	*mv = MlrvalFromPending()
	decoder := json.NewDecoder(bytes.NewReader(inputBytes))
	mlrval, eof, err := MlrvalDecodeFromJSON(decoder)
	if eof {
		return errors.New("Miller JSON parser: unexpected premature EOF.")
	}
	if err != nil {
		return err
	}
	*mv = *mlrval
	return nil
}

// ----------------------------------------------------------------
func MlrvalDecodeFromJSON(decoder *json.Decoder) (
	mlrval *Mlrval,
	eof bool,
	err error,
) {
	// Causes the decoder to unmarshal a number into an interface{} as a Number
	// instead of as a float64.
	decoder.UseNumber()

	startToken, err := decoder.Token()
	if err == io.EOF {
		return nil, true, nil
	}
	if err != nil {
		return nil, false, err
	}

	delimiter, isDelim := startToken.(json.Delim)
	if !isDelim {
		if startToken == nil {
			return MLRVAL_NULL, false, nil
		}

		sval, ok := startToken.(string)
		if ok {
			mlrval := MlrvalFromString(sval)
			return mlrval, false, nil
		}

		bval, ok := startToken.(bool)
		if ok {
			return MlrvalFromBool(bval), false, nil
		}

		nval, ok := startToken.(json.Number)
		if ok {
			mlrval := MlrvalFromInferredType(nval.String())
			return mlrval, false, nil
		}

		return nil, false, errors.New(
			"Miller JSON reader internal coding error: non-delimiter token unhandled",
		)

	} else {
		isArray := false
		var expectedClosingDelimiter = ' '
		var collectionType = ""

		if delimiter == '[' {
			isArray = true
			expectedClosingDelimiter = ']'
			collectionType = "JSON array"
		} else if delimiter == '{' {
			isArray = false
			expectedClosingDelimiter = '}'
			collectionType = "JSON object`"
		} else {
			return nil, false, errors.New(
				"Miller JSON reader: Unhandled opening delimiter \"" + string(delimiter) + "\"",
			)
		}

		mlrval := MlrvalFromPending()
		if isArray {
			mlrval = MlrvalEmptyArray()

			for decoder.More() {
				element, eof, err := MlrvalDecodeFromJSON(decoder)
				if eof {
					// xxx constify
					return nil, false, errors.New("Miller JSON parser: unexpected premature EOF.")
				}
				if err != nil {
					return nil, false, err
				}
				mlrval.ArrayAppend(element)
			}

		} else {
			mlrval = *MlrvalFromEmptyMap()

			for decoder.More() {
				key, eof, err := MlrvalDecodeFromJSON(decoder)
				if eof {
					// xxx constify
					return nil, false, errors.New("Miller JSON parser: unexpected premature EOF.")
				}
				if err != nil {
					return nil, false, err
				}
				if !key.IsString() {
					return nil, false, errors.New(
						// TODO: print out what was gotten
						"Miller JSON reader: obejct keys must be string-valued.",
					)
				}

				value, eof, err := MlrvalDecodeFromJSON(decoder)
				if eof {
					// xxx constify
					return nil, false, errors.New("Miller JSON parser: unexpected premature EOF.")
				}
				if err != nil {
					return nil, false, err
				}

				// xxx check here string-valued key
				mlrval.MapPut(key, value)
			}
		}

		imbalanceError := errors.New(
			"Miller JSON reader: did not find closing token '" +
				string(expectedClosingDelimiter) +
				"' for " +
				collectionType,
		)

		endToken, err := decoder.Token()
		if err == io.EOF {
			return nil, false, errors.New("Miller JSON parser: unexpected premature EOF.")
		}
		if err != nil {
			return nil, false, err
		}
		if endToken == nil {
			return nil, false, imbalanceError
		}
		dval, ok := endToken.(json.Delim)
		if !ok {
			return nil, false, imbalanceError
		}
		if rune(dval) != expectedClosingDelimiter {
			return nil, false, imbalanceError
		}

		return &mlrval, false, nil
	}

	return nil, false, errors.New("unimplemented")
}

// ================================================================
func (mv *Mlrval) MarshalJSON(
	jsonFormatting TJSONFormatting,
	outputIsStdout bool,
) (string, error) {
	return mv.marshalJSONAux(jsonFormatting, 1, outputIsStdout)
}

func (mv *Mlrval) marshalJSONAux(
	jsonFormatting TJSONFormatting,
	elementNestingDepth int,
	outputIsStdout bool,
) (string, error) {
	switch mv.mvtype {
	case MT_PENDING:
		return mv.marshalJSONPending(outputIsStdout)
		break
	case MT_ERROR:
		return mv.marshalJSONError(outputIsStdout)
		break
	case MT_ABSENT:
		return mv.marshalJSONAbsent(outputIsStdout)
		break
	case MT_NULL:
		return mv.marshalJSONNull(outputIsStdout)
		break
	case MT_VOID:
		return mv.marshalJSONVoid(outputIsStdout)
		break
	case MT_STRING:
		return mv.marshalJSONString(outputIsStdout)
		break
	case MT_INT:
		return mv.marshalJSONInt(outputIsStdout)
		break
	case MT_FLOAT:
		return mv.marshalJSONFloat(outputIsStdout)
		break
	case MT_BOOL:
		return mv.marshalJSONBool(outputIsStdout)
		break
	case MT_ARRAY:
		return mv.marshalJSONArray(jsonFormatting, elementNestingDepth, outputIsStdout)
		break
	case MT_MAP:
		return mv.marshalJSONMap(jsonFormatting, elementNestingDepth, outputIsStdout)
		break
	case MT_DIM: // MT_DIM is one past the last valid type
		return "", errors.New("mlr: internal coding error detected")
	}
	return "", errors.New("mlr: Internal coding error detected")
}

// ================================================================
// TYPE-SPECIFIC MARSHALERS

// ----------------------------------------------------------------
func (mv *Mlrval) marshalJSONPending(outputIsStdout bool) (string, error) {
	lib.InternalCodingErrorIf(mv.mvtype != MT_PENDING)
	return "", errors.New(
		"Miller internal coding error: pending-values should not have been produced",
	)
}

// ----------------------------------------------------------------
func (mv *Mlrval) marshalJSONError(outputIsStdout bool) (string, error) {
	lib.InternalCodingErrorIf(mv.mvtype != MT_ERROR)
	return colorizer.MaybeColorizeValue(mv.printrep, outputIsStdout), nil
}

// ----------------------------------------------------------------
func (mv *Mlrval) marshalJSONAbsent(outputIsStdout bool) (string, error) {
	lib.InternalCodingErrorIf(mv.mvtype != MT_ABSENT)
	return "", errors.New(
		"Miller internal coding error: absent-values should not have been assigned",
	)
}

// ----------------------------------------------------------------
func (mv *Mlrval) marshalJSONNull(outputIsStdout bool) (string, error) {
	lib.InternalCodingErrorIf(mv.mvtype != MT_NULL)
	return colorizer.MaybeColorizeValue("null", outputIsStdout), nil
}

// ----------------------------------------------------------------
func (mv *Mlrval) marshalJSONVoid(outputIsStdout bool) (string, error) {
	lib.InternalCodingErrorIf(mv.mvtype != MT_VOID)
	return colorizer.MaybeColorizeValue("\"\"", outputIsStdout), nil
}

// ----------------------------------------------------------------
func (mv *Mlrval) marshalJSONString(outputIsStdout bool) (string, error) {
	lib.InternalCodingErrorIf(mv.mvtype != MT_STRING)

	return colorizer.MaybeColorizeValue(millerJSONEncodeString(mv.printrep), outputIsStdout), nil
}

// Wraps with double-quotes and escape-encoded JSON-special characters.
func millerJSONEncodeString(input string) string {
	var buffer bytes.Buffer

	buffer.WriteByte('"')

	for _, b := range []byte(input) {
		switch b {
		case '\\':
			buffer.WriteByte('\\')
			buffer.WriteByte('\\')
		case '\n':
			buffer.WriteByte('\\')
			buffer.WriteByte('n')
		case '\b':
			buffer.WriteByte('\\')
			buffer.WriteByte('b')
		case '\f':
			buffer.WriteByte('\\')
			buffer.WriteByte('f')
		case '\r':
			buffer.WriteByte('\\')
			buffer.WriteByte('r')
		case '\t':
			buffer.WriteByte('\\')
			buffer.WriteByte('t')
		case '"':
			buffer.WriteByte('\\')
			buffer.WriteByte('"')
		default:
			buffer.WriteByte(b)
		}
	}

	buffer.WriteByte('"')

	return buffer.String()
}

// ----------------------------------------------------------------
func (mv *Mlrval) marshalJSONInt(outputIsStdout bool) (string, error) {
	lib.InternalCodingErrorIf(mv.mvtype != MT_INT)
	return colorizer.MaybeColorizeValue(mv.String(), outputIsStdout), nil
}

// ----------------------------------------------------------------
func (mv *Mlrval) marshalJSONFloat(outputIsStdout bool) (string, error) {
	lib.InternalCodingErrorIf(mv.mvtype != MT_FLOAT)
	return colorizer.MaybeColorizeValue(mv.String(), outputIsStdout), nil
}

// ----------------------------------------------------------------
func (mv *Mlrval) marshalJSONBool(outputIsStdout bool) (string, error) {
	lib.InternalCodingErrorIf(mv.mvtype != MT_BOOL)
	return colorizer.MaybeColorizeValue(mv.String(), outputIsStdout), nil
}

// ----------------------------------------------------------------
func (mv *Mlrval) marshalJSONArray(
	jsonFormatting TJSONFormatting,
	elementNestingDepth int,
	outputIsStdout bool,
) (string, error) {
	lib.InternalCodingErrorIf(mv.mvtype != MT_ARRAY)
	lib.InternalCodingErrorIf(jsonFormatting != JSON_SINGLE_LINE && jsonFormatting != JSON_MULTILINE)

	// Put an array of all-terminal nodes all on one line, like '[1,2,3,4,5].

	// TODO: libify
	allTerminal := true
	for _, element := range mv.arrayval {
		if element.IsArrayOrMap() {
			allTerminal = false
			break
		}
	}

	if allTerminal || (jsonFormatting == JSON_SINGLE_LINE) {
		return mv.marshalJSONArraySingleLine(elementNestingDepth, outputIsStdout)
	} else {
		return mv.marshalJSONArrayMultipleLines(jsonFormatting, elementNestingDepth, outputIsStdout)
	}
}

func (mv *Mlrval) marshalJSONArraySingleLine(
	elementNestingDepth int,
	outputIsStdout bool,
) (string, error) {
	n := len(mv.arrayval)
	var buffer bytes.Buffer
	buffer.WriteByte('[')

	for i, element := range mv.arrayval {
		elementString, err := element.marshalJSONAux(JSON_SINGLE_LINE, elementNestingDepth+1, outputIsStdout)
		if err != nil {
			return "", err
		}
		buffer.WriteString(elementString)
		if i < n-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteByte(']')
	return buffer.String(), nil
}

// The element nesting depth is how deeply our element should be indented. Our
// closing bracket is indented one less than that. For example, a
// record '{"a":1,"b":[3,[4,5],6]"c":7}' should be formatted as
//
// {
//   "a": 1,
//   "b": [    <-- root-level map element nesting depth is 1
//     3,      <-- this array's element nesting depth is 2
//     [4, 5],
//     6
//   ],        <-- this array's closing-bracket is 1, one less than its element nesting detph
//   "c": 7
// }

func (mv *Mlrval) marshalJSONArrayMultipleLines(
	jsonFormatting TJSONFormatting,
	elementNestingDepth int,
	outputIsStdout bool,
) (string, error) {
	n := len(mv.arrayval)
	var buffer bytes.Buffer

	// Write empty array as '[]'
	buffer.WriteByte('[')
	if n > 0 {
		buffer.WriteByte('\n')
	}

	for i, element := range mv.arrayval {
		elementString, err := element.marshalJSONAux(jsonFormatting, elementNestingDepth+1, outputIsStdout)
		if err != nil {
			return "", err
		}
		for i := 0; i < elementNestingDepth; i++ {
			buffer.WriteString(MLRVAL_JSON_INDENT_STRING)
		}
		buffer.WriteString(elementString)
		if i < n-1 {
			buffer.WriteString(",")
		}
		buffer.WriteString("\n")
	}

	// Write empty array as '[]'
	if n > 0 {
		for i := 0; i < elementNestingDepth-1; i++ {
			buffer.WriteString(MLRVAL_JSON_INDENT_STRING)
		}
	}

	buffer.WriteByte(']')
	return buffer.String(), nil
}

// ----------------------------------------------------------------
func (mv *Mlrval) marshalJSONMap(
	jsonFormatting TJSONFormatting,
	elementNestingDepth int,
	outputIsStdout bool,
) (string, error) {
	lib.InternalCodingErrorIf(mv.mvtype != MT_MAP)
	s, err := mv.mapval.marshalJSONAux(jsonFormatting, elementNestingDepth, outputIsStdout)
	if err != nil {
		return "", err
	}
	return s, nil
}
