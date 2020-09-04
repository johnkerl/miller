package lib

// Mlrval implements the interfaces needed for marshaling/unmarshaling to/from JSON:
//
// type Unmarshaler interface {
// 	UnmarshalJSON([]byte) error
// }
// type Marshaler interface {
// 	MarshalJSON() ([]byte, error)
// }
//
// Please see also https://golang.org/pkg/encoding/json/

import (
	"bytes"
	//"encoding/json"
	"errors"
	"strings"
)

// ================================================================
func (this *Mlrval) UnmarshalJSON(bytes []byte) error {
	return errors.New("unimplemented")
}

// ================================================================
func (this *Mlrval) MarshalJSON() ([]byte, error) {
	switch this.mvtype {
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
