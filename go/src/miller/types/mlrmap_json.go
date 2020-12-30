// ================================================================
// See mlrval_json.go for details. This is the unmarshal/marshal solely for Mlrmap.
// ================================================================

package types

import (
	"bytes"

	"miller/lib"
)

// ----------------------------------------------------------------
func (this *Mlrmap) MarshalJSON(jsonFormatting TJSONFormatting) ([]byte, error) {
	var buffer bytes.Buffer
	mapBytes, err := this.marshalJSONAux(jsonFormatting, 1)
	if err != nil {
		return nil, err
	}
	buffer.Write(mapBytes)
	// Do not write the final newline here, so the caller can write commas
	// in the right place if desired.
	// buffer.WriteString("\n")
	return buffer.Bytes(), nil
}

// For a map we only write from opening curly brace to closing curly brace.  In
// nested-map contexts, this particular map might be written with a comma
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

func (this *Mlrmap) marshalJSONAux(
	jsonFormatting TJSONFormatting,
	elementNestingDepth int,
) ([]byte, error) {
	if jsonFormatting == JSON_MULTILINE {
		return this.marshalJSONAuxMultiline(jsonFormatting, elementNestingDepth)
	} else if jsonFormatting == JSON_SINGLE_LINE {
		return this.marshalJSONAuxSingleLine(jsonFormatting, elementNestingDepth)
	} else {
		lib.InternalCodingErrorIf(true)
		return nil, nil // not reached
	}
}

func (this *Mlrmap) marshalJSONAuxMultiline(
	jsonFormatting TJSONFormatting,
	elementNestingDepth int,
) ([]byte, error) {
	var buffer bytes.Buffer

	buffer.WriteString("{")
	// Write empty map as '{}'. For anything else, opening curly brace in a
	// line of its own, one key-value pair per line, closing curly brace on a
	// line of its own.
	if this.Head != nil {
		buffer.WriteString("\n")
	}

	for pe := this.Head; pe != nil; pe = pe.Next {
		// Write the key which is necessarily string-valued in Miller, and in
		// JSON for that matter :)
		for i := 0; i < elementNestingDepth; i++ {
			buffer.WriteString(MLRVAL_JSON_INDENT_STRING)
		}
		buffer.WriteString("\"")
		buffer.WriteString(*pe.Key)
		buffer.WriteString("\": ")

		// Write the value which is a mlrval
		valueBytes, err := pe.Value.marshalJSONAux(jsonFormatting, elementNestingDepth+1)
		if err != nil {
			return nil, err
		}
		_, err = buffer.Write(valueBytes)
		if err != nil {
			return nil, err
		}

		if pe.Next != nil {
			buffer.WriteString(",")
		}
		buffer.WriteString("\n")
	}

	// Write empty map as '{}'.
	if this.Head != nil {
		for i := 0; i < elementNestingDepth-1; i++ {
			buffer.WriteString(MLRVAL_JSON_INDENT_STRING)
		}
	}
	buffer.WriteString("}")

	return buffer.Bytes(), nil
}

func (this *Mlrmap) marshalJSONAuxSingleLine(
	jsonFormatting TJSONFormatting,
	elementNestingDepth int,
) ([]byte, error) {
	var buffer bytes.Buffer

	buffer.WriteString("{")

	for pe := this.Head; pe != nil; pe = pe.Next {
		// Write the key which is necessarily string-valued in Miller, and in
		// JSON for that matter :)
		buffer.WriteString("\"")
		buffer.WriteString(*pe.Key)
		buffer.WriteString("\": ")

		// Write the value which is a mlrval
		valueBytes, err := pe.Value.marshalJSONAux(jsonFormatting, elementNestingDepth+1)
		if err != nil {
			return nil, err
		}
		_, err = buffer.Write(valueBytes)
		if err != nil {
			return nil, err
		}

		if pe.Next != nil {
			buffer.WriteString(", ")
		}
	}

	buffer.WriteString("}")

	return buffer.Bytes(), nil
}
