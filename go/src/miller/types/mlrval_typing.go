package types

import (
	"errors"
	"fmt"
)

// ================================================================
// Support for things like 'num x = $a + $b' in the DSL, wherein we check types
// at assigment time.
// ================================================================

type TypeGatedMlrvalName struct {
	Name     string
	TypeName string
	TypeMask int
}

func NewTypeGatedMlrvalName(
	name string, // e.g. "x"
	typeName string, // e.g. "num"
) (*TypeGatedMlrvalName, error) {
	typeMask, ok := TypeNameToMask(typeName)
	if !ok {
		return nil, errors.New(
			fmt.Sprintf(
				"Miller: couldn't resolve type name \"%s\".", typeName,
			),
		)
	}
	return &TypeGatedMlrvalName{
		Name:     name,
		TypeName: typeName,
		TypeMask: typeMask,
	}, nil
}

func (this *TypeGatedMlrvalName) Check(value *Mlrval) error {
	bit := value.GetTypeBit()
	if bit&this.TypeMask != 0 {
		return nil
	} else {
		return errors.New(
			fmt.Sprintf(
				"Miller: couldn't assign variable %s %s from value %s %s\n",
				this.TypeName, this.Name, value.GetTypeName(), value.String(),
			),
		)
	}
}
