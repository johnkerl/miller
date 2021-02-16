// ================================================================
// Support for things like 'num x = $a + $b' in the DSL, wherein we check types
// at assigment time.
// ================================================================

package types

import (
	"errors"
	"fmt"
)

// ----------------------------------------------------------------
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

// ----------------------------------------------------------------
type TypeGatedMlrvalVariable struct {
	typeGatedMlrvalName *TypeGatedMlrvalName
	value               *Mlrval
}

func NewTypeGatedMlrvalVariable(
	name string, // e.g. "x"
	typeName string, // e.g. "num"
	value *Mlrval,
) (*TypeGatedMlrvalVariable, error) {
	typeGatedMlrvalName, err := NewTypeGatedMlrvalName(name, typeName)
	if err != nil {
		return nil, err
	}

	err = typeGatedMlrvalName.Check(value)
	if err != nil {
		return nil, err
	}

	return &TypeGatedMlrvalVariable{
		typeGatedMlrvalName,
		value.Copy(),
	}, nil
}

func (this *TypeGatedMlrvalVariable) GetValue() *Mlrval {
	return this.value
}

func (this *TypeGatedMlrvalVariable) ValueString() string {
	return this.value.String()
}

func (this *TypeGatedMlrvalVariable) Assign(value *Mlrval) error {
	err := this.typeGatedMlrvalName.Check(value)
	if err != nil {
		return err
	}

	this.value = value.Copy()
	return nil
}
