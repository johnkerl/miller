package types

// See mlrval_json.go for details. This is the unmarshal/marshal solely for Mlrmap.

import (
	"bytes"
	//"encoding/json"
)

// ----------------------------------------------------------------
func (this *Mlrmap) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	mapBytes, err := this.marshalJSONAux(1)
	if err != nil {
		return nil, err
	}
	buffer.Write(mapBytes)
	buffer.WriteString("\n")
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

func (this *Mlrmap) marshalJSONAux(elementNestingDepth int) ([]byte, error) {
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
		valueBytes, err := pe.Value.marshalJSONAux(elementNestingDepth + 1)
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
