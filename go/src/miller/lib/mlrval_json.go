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

// ----------------------------------------------------------------
func (this *Mlrval) UnmarshalJSON(bytes []byte) error {
	switch this.mvtype {
	case MT_ABSENT:
		return this.unmarshalJSONAbsent(bytes)
		break
	case MT_VOID:
		return this.unmarshalJSONVoid(bytes)
		break
	case MT_STRING:
		return this.unmarshalJSONString(bytes)
		break
	case MT_INT:
		return this.unmarshalJSONInt(bytes)
		break
	case MT_FLOAT:
		return this.unmarshalJSONFloat(bytes)
		break
	case MT_BOOL:
		return this.unmarshalJSONBool(bytes)
		break
	case MT_ARRAY:
		return this.unmarshalJSONArray(bytes)
		break
	case MT_MAP:
		return this.unmarshalJSONMap(bytes)
		break
	case MT_DIM: // MT_DIM is one past the last valid type
		return errors.New("internal coding error detected")
	}
	return errors.New("internal coding error detected")
}

// ----------------------------------------------------------------
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

// ----------------------------------------------------------------
func (this *Mlrval) unmarshalJSONAbsent(bytes []byte) error {
	return errors.New("unimplemented")
}

func (this *Mlrval) unmarshalJSONVoid(bytes []byte) error {
	return errors.New("unimplemented")
}

func (this *Mlrval) unmarshalJSONString(bytes []byte) error {
	return errors.New("unimplemented")
}

func (this *Mlrval) unmarshalJSONInt(bytes []byte) error {
	return errors.New("unimplemented")
}

func (this *Mlrval) unmarshalJSONFloat(bytes []byte) error {
	return errors.New("unimplemented")
}

func (this *Mlrval) unmarshalJSONBool(bytes []byte) error {
	return errors.New("unimplemented")
}

func (this *Mlrval) unmarshalJSONArray(bytes []byte) error {
	return errors.New("unimplemented")
}

func (this *Mlrval) unmarshalJSONMap(bytes []byte) error {
	return errors.New("unimplemented")
}

// ----------------------------------------------------------------
func (this *Mlrval) marshalJSONAbsent() ([]byte, error) {
	return nil, errors.New(
		"Miller: internal coding error: absent-values should not have been assigned",
	)
}

func (this *Mlrval) marshalJSONVoid() ([]byte, error) {
	return []byte(""), nil
}

func (this *Mlrval) marshalJSONString() ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteByte('"')
	buffer.WriteString(strings.ReplaceAll(this.printrep, "\"", "\\\""))
	buffer.WriteByte('"')
	return buffer.Bytes(), nil
}

func (this *Mlrval) marshalJSONInt() ([]byte, error) {
	return []byte(this.String()), nil
}

func (this *Mlrval) marshalJSONFloat() ([]byte, error) {
	return []byte(this.String()), nil
}

func (this *Mlrval) marshalJSONBool() ([]byte, error) {
	return []byte(this.String()), nil
}

func (this *Mlrval) marshalJSONArray() ([]byte, error) {
	//var buffer bytes.Buffer
	//buffer.WriteByte('[')
	//buffer.WriteByte(']')
	return nil, errors.New("unimplemented")
}

func (this *Mlrval) marshalJSONMap() ([]byte, error) {
	return nil, errors.New("unimplemented")
}

// ----------------------------------------------------------------
//	printrep      string
//	printrepValid bool
//	intval        int64
//	floatval      float64
//	boolval       bool
//	arrayval      []Mlrval
//	mapval        *Mlrmap

// MT_ABSENT
// MT_VOID
// MT_STRING
// MT_INT
// MT_FLOAT
// MT_BOOL
// MT_ARRAY
// MT_MAP
