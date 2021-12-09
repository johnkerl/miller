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
				"mlr: couldn't resolve type name \"%s\".", typeName,
			),
		)
	}
	return &TypeGatedMlrvalName{
		Name:     name,
		TypeName: typeName,
		TypeMask: typeMask,
	}, nil
}

func (tname *TypeGatedMlrvalName) Check(value *Mlrval) error {
	bit := value.GetTypeBit()
	if bit&tname.TypeMask != 0 {
		return nil
	} else {
		return errors.New(
			fmt.Sprintf(
				"mlr: couldn't assign variable %s %s from value %s %s\n",
				tname.TypeName, tname.Name, value.GetTypeName(), value.String(),
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

func (tvar *TypeGatedMlrvalVariable) GetName() string {
	return tvar.typeGatedMlrvalName.Name
}

func (tvar *TypeGatedMlrvalVariable) GetValue() *Mlrval {
	return tvar.value
}

func (tvar *TypeGatedMlrvalVariable) ValueString() string {
	return tvar.value.String()
}

func (tvar *TypeGatedMlrvalVariable) Assign(value *Mlrval) error {
	err := tvar.typeGatedMlrvalName.Check(value)
	if err != nil {
		return err
	}

	// TODO: revisit copy-reduction
	tvar.value = value.Copy()
	return nil
}

func (tvar *TypeGatedMlrvalVariable) Unassign() {
	tvar.value = MLRVAL_ABSENT
}
