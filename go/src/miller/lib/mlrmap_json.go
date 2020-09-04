package lib

// See mlrval_json.go for details. This is the unmarshal/marshal solely for Mlrmap.

import (
	"bytes"
	//"encoding/json"
)

// ----------------------------------------------------------------
func (this *Mlrmap) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer

	// TODO: how to handle indentation for the nested-object case.
	buffer.WriteString("{\n")

	for pe := this.Head; pe != nil; pe = pe.Next {
		// Write the key which is necessarily string-valued in Miller, and in
		// JSON for that matter :)
		buffer.WriteString("  \"")
		buffer.WriteString(*pe.Key)
		buffer.WriteString("\": ")

		// Write the value which is a mlrval
		valueBytes, err := pe.Value.MarshalJSON()
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
	buffer.WriteString("}\n")
	return buffer.Bytes(), nil
}
