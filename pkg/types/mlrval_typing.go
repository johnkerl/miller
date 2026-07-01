// Support for things like 'num x = $a + $b' in the DSL, wherein we check types
// at assignment time.

package types

import (
	"fmt"

	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

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
		return nil, fmt.Errorf(`couldn't resolve type name "%s"`, typeName)
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
	}
	return fmt.Errorf(
		"couldn't assign variable %s %s from value %s %s",
		tname.TypeName, tname.Name, value.GetTypeName(), value.String(),
	)
}

type TypeGatedMlrvalVariable struct {
	typeGatedMlrvalName *TypeGatedMlrvalName
	value               *mlrval.Mlrval
}

// copyForBind returns the value to store when binding a local variable or
// function parameter. Scalars are stored by reference (no copy, no allocation):
// assignment everywhere replaces pointers rather than mutating Mlrvals in place
// (Mlrmap.PutCopy reassigns pe.Value; Assign reassigns tvar.value), and the
// only in-place mutation a scalar undergoes is idempotent type-inference
// caching -- so an aliased scalar can never be observed to change underneath
// its source. Maps and arrays, however, ARE mutated in place by indexed
// assignment (m[k]=v), so they must be deep-copied to keep the binding
// independent of its source.
func copyForBind(value *mlrval.Mlrval) *mlrval.Mlrval {
	if value.IsArrayOrMap() {
		return value.Copy()
	}
	return value
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
		copyForBind(value),
	}, nil
}

func (tvar *TypeGatedMlrvalVariable) GetValue() *mlrval.Mlrval {
	return tvar.value
}

func (tvar *TypeGatedMlrvalVariable) Assign(value *mlrval.Mlrval) error {
	err := tvar.typeGatedMlrvalName.Check(value)
	if err != nil {
		return err
	}

	tvar.value = copyForBind(value)
	return nil
}

func (tvar *TypeGatedMlrvalVariable) Unassign() {
	tvar.value = mlrval.ABSENT
}
