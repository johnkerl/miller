// ================================================================
// Support for things like 'num x = $a + $b' in the DSL, wherein we check types
// at assignment time.
// ================================================================

package types

import (
	"fmt"

	"github.com/johnkerl/miller/internal/pkg/mlrval"
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
	typeMask, ok := mlrval.TypeNameToMask(typeName)
	if !ok {
		return nil, fmt.Errorf("mlr: couldn't resolve type name \"%s\".", typeName)
	}
	return &TypeGatedMlrvalName{
		Name:     name,
		TypeName: typeName,
		TypeMask: typeMask,
	}, nil
}

func (tname *TypeGatedMlrvalName) Check(value *mlrval.Mlrval) error {
	bit := value.GetTypeBit()
	if bit&tname.TypeMask != 0 {
		return nil
	} else {
		return fmt.Errorf(
			"mlr: couldn't assign variable %s %s from value %s %s\n",
			tname.TypeName, tname.Name, value.GetTypeName(), value.String(),
		)
	}
}

// ----------------------------------------------------------------
type TypeGatedMlrvalVariable struct {
	typeGatedMlrvalName *TypeGatedMlrvalName
	value               *mlrval.Mlrval
}

func NewTypeGatedMlrvalVariable(
	name string, // e.g. "x"
	typeName string, // e.g. "num"
	value *mlrval.Mlrval,
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

func (tvar *TypeGatedMlrvalVariable) GetValue() *mlrval.Mlrval {
	return tvar.value
}

func (tvar *TypeGatedMlrvalVariable) ValueString() string {
	return tvar.value.String()
}

func (tvar *TypeGatedMlrvalVariable) Assign(value *mlrval.Mlrval) error {
	err := tvar.typeGatedMlrvalName.Check(value)
	if err != nil {
		return err
	}

	// TODO: revisit copy-reduction
	tvar.value = value.Copy()
	return nil
}

func (tvar *TypeGatedMlrvalVariable) Unassign() {
	tvar.value = mlrval.ABSENT
}
