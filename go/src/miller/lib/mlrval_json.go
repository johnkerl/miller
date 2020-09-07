package lib

// ================================================================
// Mlrval implements the Unmarshaler and Marshaler interfaces needed for
// marshaling/unmarshaling to/from JSON, via the UnmarshalJSON and MarshalJSON
// methods.
//
// Please see also https://golang.org/pkg/encoding/json/
// ================================================================

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
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
			"Miller JSON reader: internal coding error: non-delimiter token unhandled",
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
				mlrval.ArrayExtend(element)
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
	switch this.mvtype {
	case MT_PENDING:
		return this.marshalJSONPending()
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
		return this.marshalJSONArray()
		break
	case MT_MAP:
		return this.marshalJSONMap()
		break
	case MT_DIM: // MT_DIM is one past the last valid type
		return nil, errors.New("internal coding error detected")
	}
	return nil, errors.New("internal coding error detected")
}

// ================================================================
// TYPE-SPECIFIC MARSHALERS

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONPending() ([]byte, error) {
	InternalCodingErrorIf(this.mvtype != MT_PENDING)
	return nil, errors.New(
		"Miller: internal coding error: pending-values should not have been produced",
	)
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONAbsent() ([]byte, error) {
	InternalCodingErrorIf(this.mvtype != MT_ABSENT)
	return nil, errors.New(
		"Miller: internal coding error: absent-values should not have been assigned",
	)
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONVoid() ([]byte, error) {
	InternalCodingErrorIf(this.mvtype != MT_VOID)
	return []byte(""), nil
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONString() ([]byte, error) {
	InternalCodingErrorIf(this.mvtype != MT_STRING)
	var buffer bytes.Buffer
	buffer.WriteByte('"')
	buffer.WriteString(strings.ReplaceAll(this.printrep, "\"", "\\\""))
	buffer.WriteByte('"')
	return buffer.Bytes(), nil
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONInt() ([]byte, error) {
	InternalCodingErrorIf(this.mvtype != MT_INT)
	return []byte(this.String()), nil
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONFloat() ([]byte, error) {
	InternalCodingErrorIf(this.mvtype != MT_FLOAT)
	return []byte(this.String()), nil
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONBool() ([]byte, error) {
	InternalCodingErrorIf(this.mvtype != MT_BOOL)
	return []byte(this.String()), nil
}

// ----------------------------------------------------------------
// TODO: find out how to handle indentation in the nested-array/nested-map case ...
func (this *Mlrval) marshalJSONArray() ([]byte, error) {
	InternalCodingErrorIf(this.mvtype != MT_ARRAY)
	n := len(this.arrayval)
	var buffer bytes.Buffer
	buffer.WriteByte('[')
	for i, element := range this.arrayval {
		elementBytes, err := element.MarshalJSON()
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

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONMap() ([]byte, error) {
	InternalCodingErrorIf(this.mvtype != MT_MAP)
	bytes, err := this.mapval.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
