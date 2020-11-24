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
	"strings"

	"miller/lib"
)

const MLRVAL_JSON_INDENT_STRING string = "  "

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
func (this *Mlrval) UnmarshalJSON(inputBytes []byte) error {
	*this = MlrvalFromPending()
	decoder := json.NewDecoder(bytes.NewReader(inputBytes))
	mlrval, eof, err := MlrvalDecodeFromJSON(decoder)
	if eof {
		return errors.New("Miller JSON parser: unexpected premature EOF.")
	}
	if err != nil {
		return err
	}
	*this = *mlrval
	return nil
}

// ----------------------------------------------------------------
func MlrvalDecodeFromJSON(decoder *json.Decoder) (mlrval *Mlrval, eof bool, err error) {
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
			mlrval := MlrvalFromVoid()
			return &mlrval, false, nil
		}

		sval, ok := startToken.(string)
		if ok {
			mlrval := MlrvalFromString(sval)
			return &mlrval, false, nil
		}

		bval, ok := startToken.(bool)
		if ok {
			mlrval := MlrvalFromBool(bval)
			return &mlrval, false, nil
		}

		nval, ok := startToken.(json.Number)
		if ok {
			mlrval := MlrvalFromInferredType(nval.String())
			return &mlrval, false, nil
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
			mlrval = MlrvalEmptyMap()

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
func (this *Mlrval) MarshalJSON() ([]byte, error) {
	return this.marshalJSONAux(1)
}

func (this *Mlrval) marshalJSONAux(elementNestingDepth int) ([]byte, error) {
	switch this.mvtype {
	case MT_PENDING:
		return this.marshalJSONPending()
		break
	case MT_ERROR:
		return this.marshalJSONError()
		break
	case MT_ABSENT:
		return this.marshalJSONAbsent()
		break
	case MT_VOID:
		return this.marshalJSONVoid()
		break
	case MT_STRING:
		return this.marshalJSONString()
		break
	case MT_INT:
		return this.marshalJSONInt()
		break
	case MT_FLOAT:
		return this.marshalJSONFloat()
		break
	case MT_BOOL:
		return this.marshalJSONBool()
		break
	case MT_ARRAY:
		return this.marshalJSONArray(elementNestingDepth)
		break
	case MT_MAP:
		return this.marshalJSONMap(elementNestingDepth)
		break
	case MT_DIM: // MT_DIM is one past the last valid type
		return nil, errors.New("Miller: internal coding error detected")
	}
	return nil, errors.New("Miller: iInternal coding error detected")
}

// ================================================================
// TYPE-SPECIFIC MARSHALERS

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONPending() ([]byte, error) {
	lib.InternalCodingErrorIf(this.mvtype != MT_PENDING)
	return nil, errors.New(
		"Miller internal coding error: pending-values should not have been produced",
	)
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONError() ([]byte, error) {
	lib.InternalCodingErrorIf(this.mvtype != MT_ERROR)
	return []byte(this.printrep), nil
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONAbsent() ([]byte, error) {
	lib.InternalCodingErrorIf(this.mvtype != MT_ABSENT)
	return nil, errors.New(
		"Miller internal coding error: absent-values should not have been assigned",
	)
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONVoid() ([]byte, error) {
	lib.InternalCodingErrorIf(this.mvtype != MT_VOID)
	return []byte("\"\""), nil
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONString() ([]byte, error) {
	lib.InternalCodingErrorIf(this.mvtype != MT_STRING)
	var buffer bytes.Buffer
	buffer.WriteByte('"')
	buffer.WriteString(strings.Replace(this.printrep, "\"", "\\\"", -1))
	buffer.WriteByte('"')
	return buffer.Bytes(), nil
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONInt() ([]byte, error) {
	lib.InternalCodingErrorIf(this.mvtype != MT_INT)
	return []byte(this.String()), nil
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONFloat() ([]byte, error) {
	lib.InternalCodingErrorIf(this.mvtype != MT_FLOAT)
	return []byte(this.String()), nil
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONBool() ([]byte, error) {
	lib.InternalCodingErrorIf(this.mvtype != MT_BOOL)
	return []byte(this.String()), nil
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONArray(elementNestingDepth int) ([]byte, error) {
	lib.InternalCodingErrorIf(this.mvtype != MT_ARRAY)

	// Put an array of all-terminal nodes all on one line, like '[1,2,3,4,5].
	allTerminal := true
	for _, element := range this.arrayval {
		if element.IsArrayOrMap() {
			allTerminal = false
			break
		}
	}
	if allTerminal {
		return this.marshalJSONArraySingleLine(elementNestingDepth)
	} else {
		return this.marshalJSONArrayMultipleLines(elementNestingDepth)
	}
}

func (this *Mlrval) marshalJSONArraySingleLine(elementNestingDepth int) ([]byte, error) {
	n := len(this.arrayval)
	var buffer bytes.Buffer
	buffer.WriteByte('[')

	for i, element := range this.arrayval {
		elementBytes, err := element.marshalJSONAux(elementNestingDepth + 1)
		if err != nil {
			return nil, err
		}
		buffer.Write(elementBytes)
		if i < n-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteByte(']')
	return buffer.Bytes(), nil
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

func (this *Mlrval) marshalJSONArrayMultipleLines(elementNestingDepth int) ([]byte, error) {
	n := len(this.arrayval)
	var buffer bytes.Buffer

	// Write empty array as '[]'
	buffer.WriteByte('[')
	if n > 0 {
		buffer.WriteByte('\n')
	}

	for i, element := range this.arrayval {
		elementBytes, err := element.marshalJSONAux(elementNestingDepth + 1)
		if err != nil {
			return nil, err
		}
		for i := 0; i < elementNestingDepth; i++ {
			buffer.WriteString(MLRVAL_JSON_INDENT_STRING)
		}
		buffer.Write(elementBytes)
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
	return buffer.Bytes(), nil
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONMap(elementNestingDepth int) ([]byte, error) {
	lib.InternalCodingErrorIf(this.mvtype != MT_MAP)
	bytes, err := this.mapval.marshalJSONAux(elementNestingDepth)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
