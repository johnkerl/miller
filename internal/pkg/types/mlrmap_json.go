// ================================================================
// See mlrval_json.go for details. This is the unmarshal/marshal solely for Mlrmap.
// ================================================================

package types

import (
	"bytes"

	"mlr/internal/pkg/colorizer"
	"mlr/internal/pkg/lib"
)

// ----------------------------------------------------------------
func (mlrmap *Mlrmap) MarshalJSON(
	jsonFormatting TJSONFormatting,
	outputIsStdout bool,
) (string, error) {
	var buffer bytes.Buffer
	s, err := mlrmap.marshalJSONAux(jsonFormatting, 1, outputIsStdout)
	if err != nil {
		return "", err
	}
	buffer.WriteString(s)
	// Do not write the final newline here, so the caller can write commas
	// in the right place if desired.
	// buffer.WriteString("\n")
	return buffer.String(), nil
}

// For a map we only write from opening curly brace to closing curly brace.  In
// nested-map contexts, mlrmap particular map might be written with a comma
// immediately after its closing curly brace, or a newline, and only the caller
// can know that.
//
// The element nesting depth is how deeply our element should be indented. Our
// closing curly brace is indented one less than that. For example, a
// root-level record '{"a":1,"b":2}' should be formatted as
//
// {
//   "a": 1, <-- element nesting depth is 1 for root-level map
//   "b": 2  <-- element nesting depth is 1 for root-level map
// }         <-- closing curly brace nesting depth is 0 for root-level map

func (mlrmap *Mlrmap) marshalJSONAux(
	jsonFormatting TJSONFormatting,
	elementNestingDepth int,
	outputIsStdout bool,
) (string, error) {
	if jsonFormatting == JSON_MULTILINE {
		return mlrmap.marshalJSONAuxMultiline(jsonFormatting, elementNestingDepth, outputIsStdout)
	} else if jsonFormatting == JSON_SINGLE_LINE {
		return mlrmap.marshalJSONAuxSingleLine(jsonFormatting, elementNestingDepth, outputIsStdout)
	} else {
		lib.InternalCodingErrorIf(true)
		return "", nil // not reached
	}
}

func (mlrmap *Mlrmap) marshalJSONAuxMultiline(
	jsonFormatting TJSONFormatting,
	elementNestingDepth int,
	outputIsStdout bool,
) (string, error) {
	var buffer bytes.Buffer

	buffer.WriteString("{")
	// Write empty map as '{}'. For anything else, opening curly brace in a
	// line of its own, one key-value pair per line, closing curly brace on a
	// line of its own.
	if mlrmap.Head != nil {
		buffer.WriteString("\n")
	}

	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		// Write the key which is necessarily string-valued in Miller, and in
		// JSON for that matter :)
		for i := 0; i < elementNestingDepth; i++ {
			buffer.WriteString(MLRVAL_JSON_INDENT_STRING)
		}
		encoded := string(millerJSONEncodeString(pe.Key))
		colorized := colorizer.MaybeColorizeKey(encoded, outputIsStdout)
		buffer.WriteString(colorized)
		buffer.WriteString(": ")

		// Write the value which is a mlrval
		valueString, err := pe.Value.marshalJSONAux(jsonFormatting, elementNestingDepth+1, outputIsStdout)
		if err != nil {
			return "", err
		}
		_, err = buffer.WriteString(valueString)
		if err != nil {
			return "", err
		}

		if pe.Next != nil {
			buffer.WriteString(",")
		}
		buffer.WriteString("\n")
	}

	// Write empty map as '{}'.
	if mlrmap.Head != nil {
		for i := 0; i < elementNestingDepth-1; i++ {
			buffer.WriteString(MLRVAL_JSON_INDENT_STRING)
		}
	}
	buffer.WriteString("}")

	return buffer.String(), nil
}

func (mlrmap *Mlrmap) marshalJSONAuxSingleLine(
	jsonFormatting TJSONFormatting,
	elementNestingDepth int,
	outputIsStdout bool,
) (string, error) {
	var buffer bytes.Buffer

	buffer.WriteString("{")

	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		// Write the key which is necessarily string-valued in Miller, and in
		// JSON for that matter :)
		buffer.WriteString("\"")
		buffer.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
		buffer.WriteString("\": ")

		// Write the value which is a mlrval
		valueString, err := pe.Value.marshalJSONAux(jsonFormatting, elementNestingDepth+1, outputIsStdout)
		if err != nil {
			return "", err
		}
		_, err = buffer.WriteString(valueString)
		if err != nil {
			return "", err
		}

		if pe.Next != nil {
			buffer.WriteString(", ")
		}
	}

	buffer.WriteString("}")

	return buffer.String(), nil
}

// ----------------------------------------------------------------
// JSON-stringifies a single field of a record
func (entry *MlrmapEntry) JSONStringifyInPlace(
	jsonFormatting TJSONFormatting,
) {
	outputBytes, err := entry.Value.MarshalJSON(jsonFormatting, false)
	if err != nil {
		entry.Value = MLRVAL_ERROR
	} else {
		entry.Value = MlrvalFromString(string(outputBytes))
	}
}

// ----------------------------------------------------------------
// JSON-parses a single field of a record
func (entry *MlrmapEntry) JSONParseInPlace() {
	input := entry.Value.String()
	err := entry.Value.UnmarshalJSON([]byte(input))
	if err != nil {
		entry.Value = MLRVAL_ERROR
	}
}
