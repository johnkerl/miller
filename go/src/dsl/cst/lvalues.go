// ================================================================
// This is for Lvalues, i.e. things on the left-hand-side of an assignment
// statement.
// ================================================================

package cst

import (
	"errors"
	"fmt"
	"os"

	"miller/src/dsl"
	"miller/src/lib"
	"miller/src/runtime"
	"miller/src/types"
)

// ----------------------------------------------------------------
func (root *RootNode) BuildAssignableNode(
	astNode *dsl.ASTNode,
) (IAssignable, error) {

	switch astNode.Type {

	case dsl.NodeTypeDirectFieldValue:
		return root.BuildDirectFieldValueLvalueNode(astNode)
		break
	case dsl.NodeTypeIndirectFieldValue:
		return root.BuildIndirectFieldValueLvalueNode(astNode)
		break
	case dsl.NodeTypePositionalFieldName:
		return root.BuildPositionalFieldNameLvalueNode(astNode)
		break
	case dsl.NodeTypePositionalFieldValue:
		return root.BuildPositionalFieldValueLvalueNode(astNode)
		break

	case dsl.NodeTypeFullSrec:
		return root.BuildFullSrecLvalueNode(astNode)
		break

	case dsl.NodeTypeDirectOosvarValue:
		return root.BuildDirectOosvarValueLvalueNode(astNode)
		break
	case dsl.NodeTypeIndirectOosvarValue:
		return root.BuildIndirectOosvarValueLvalueNode(astNode)
		break
	case dsl.NodeTypeFullOosvar:
		return root.BuildFullOosvarLvalueNode(astNode)
		break
	case dsl.NodeTypeLocalVariable:
		return root.BuildLocalVariableLvalueNode(astNode)
		break

	case dsl.NodeTypeArrayOrMapPositionalNameAccess:
		return nil, errors.New(
			"Miller: '[[...]]' is allowed on assignment left-hand sides only when immediately preceded by '$'.",
		)
		break
	case dsl.NodeTypeArrayOrMapPositionalValueAccess:
		return nil, errors.New(
			"Miller: '[[[...]]]' is allowed on assignment left-hand sides only when immediately preceded by '$'.",
		)
		break

	case dsl.NodeTypeArrayOrMapIndexAccess:
		return root.BuildIndexedLvalueNode(astNode)
		break

	case dsl.NodeTypeEnvironmentVariable:
		return root.BuildEnvironmentVariableLvalueNode(astNode)
		break
	}

	return nil, errors.New(
		"CST BuildAssignableNode: unhandled AST node " + string(astNode.Type),
	)
}

// ----------------------------------------------------------------
type DirectFieldValueLvalueNode struct {
	lhsFieldName *types.Mlrval
}

func (root *RootNode) BuildDirectFieldValueLvalueNode(
	astNode *dsl.ASTNode,
) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDirectFieldValue)

	lhsFieldName := types.MlrvalFromString(string(astNode.Token.Lit))
	return NewDirectFieldValueLvalueNode(&lhsFieldName), nil
}

func NewDirectFieldValueLvalueNode(lhsFieldName *types.Mlrval) *DirectFieldValueLvalueNode {
	return &DirectFieldValueLvalueNode{
		lhsFieldName: lhsFieldName,
	}
}

func (this *DirectFieldValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *DirectFieldValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {

	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		return errors.New("There is no current record to assign to.")
	}

	// AssignmentNode checks for absent, so we just assign whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())
	if indices == nil {
		err := state.Inrec.PutCopyWithMlrvalIndex(this.lhsFieldName, rvalue)
		if err != nil {
			return err
		}
		return nil
	} else {
		return state.Inrec.PutIndexed(
			append([]*types.Mlrval{this.lhsFieldName}, indices...),
			rvalue,
		)
	}
}

func (this *DirectFieldValueLvalueNode) Unassign(
	state *runtime.State,
) {
	this.UnassignIndexed(nil, state)
}

func (this *DirectFieldValueLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {

	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		return
	}

	if indices == nil {
		lib.InternalCodingErrorIf(!this.lhsFieldName.IsString())
		name := this.lhsFieldName.String()
		state.Inrec.Remove(name)
	} else {
		state.Inrec.RemoveIndexed(
			append([]*types.Mlrval{this.lhsFieldName}, indices...),
		)
	}
}

// ----------------------------------------------------------------
type IndirectFieldValueLvalueNode struct {
	lhsFieldNameExpression IEvaluable
}

func (root *RootNode) BuildIndirectFieldValueLvalueNode(
	astNode *dsl.ASTNode,
) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectFieldValue)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	lhsFieldNameExpression, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return NewIndirectFieldValueLvalueNode(lhsFieldNameExpression), nil
}

func NewIndirectFieldValueLvalueNode(
	lhsFieldNameExpression IEvaluable,
) *IndirectFieldValueLvalueNode {
	return &IndirectFieldValueLvalueNode{
		lhsFieldNameExpression: lhsFieldNameExpression,
	}
}

func (this *IndirectFieldValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *IndirectFieldValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		return errors.New("There is no current record to assign to.")
	}

	lhsFieldName := this.lhsFieldNameExpression.Evaluate(state)

	if indices == nil {
		err := state.Inrec.PutCopyWithMlrvalIndex(lhsFieldName, rvalue)
		if err != nil {
			return err
		}
		return nil
	} else {
		return state.Inrec.PutIndexed(
			append([]*types.Mlrval{lhsFieldName.Copy()}, indices...),
			rvalue,
		)
	}
}

func (this *IndirectFieldValueLvalueNode) Unassign(
	state *runtime.State,
) {
	this.UnassignIndexed(nil, state)
}

func (this *IndirectFieldValueLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		return
	}

	lhsFieldName := this.lhsFieldNameExpression.Evaluate(state)
	if indices == nil {
		name := lhsFieldName.String()
		state.Inrec.Remove(name)
	} else {
		state.Inrec.RemoveIndexed(
			append([]*types.Mlrval{lhsFieldName.Copy()}, indices...),
		)
	}
}

// ----------------------------------------------------------------
// Set the name at 2nd positional index in the current stream record: e.g.
// '$[[2]] = "abc"

type PositionalFieldNameLvalueNode struct {
	lhsFieldIndexExpression IEvaluable
}

func (root *RootNode) BuildPositionalFieldNameLvalueNode(
	astNode *dsl.ASTNode,
) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypePositionalFieldName)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	lhsFieldIndexExpression, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return NewPositionalFieldNameLvalueNode(lhsFieldIndexExpression), nil
}

func NewPositionalFieldNameLvalueNode(
	lhsFieldIndexExpression IEvaluable,
) *PositionalFieldNameLvalueNode {
	return &PositionalFieldNameLvalueNode{
		lhsFieldIndexExpression: lhsFieldIndexExpression,
	}
}

func (this *PositionalFieldNameLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		return errors.New("There is no current record to assign to.")
	}

	lhsFieldIndex := this.lhsFieldIndexExpression.Evaluate(state)

	index, ok := lhsFieldIndex.GetIntValue()
	if ok {
		// TODO: incorporate error-return into this API
		state.Inrec.PutNameWithPositionalIndex(index, rvalue)
		return nil
	} else {
		return errors.New(
			fmt.Sprintf(
				"Miller: positional index for $[[...]] assignment must be integer; got %s.",
				lhsFieldIndex.GetTypeName(),
			),
		)
	}
}

func (this *PositionalFieldNameLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	// TODO: reconsider this if /when we decide to allow string-slice
	// assignments.
	return errors.New(
		"Miller: $[[...]] = ... expressions are not indexable.",
	)
}

func (this *PositionalFieldNameLvalueNode) Unassign(
	state *runtime.State,
) {
	this.UnassignIndexed(nil, state)
}

func (this *PositionalFieldNameLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	lhsFieldIndex := this.lhsFieldIndexExpression.Evaluate(state)

	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		return
	}

	if indices == nil {
		index, ok := lhsFieldIndex.GetIntValue()
		if ok {
			state.Inrec.RemoveWithPositionalIndex(index)
		} else {
			// TODO: incorporate error-return into this API
		}
	} else {
		// xxx positional
		state.Inrec.RemoveIndexed(
			append([]*types.Mlrval{lhsFieldIndex}, indices...),
		)
	}
}

// ----------------------------------------------------------------
// Set the value at 2nd positional index in the current stream record: e.g.
// '$[[[2]]] = "abc"

type PositionalFieldValueLvalueNode struct {
	lhsFieldIndexExpression IEvaluable
}

func (root *RootNode) BuildPositionalFieldValueLvalueNode(
	astNode *dsl.ASTNode,
) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypePositionalFieldValue)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	lhsFieldIndexExpression, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return NewPositionalFieldValueLvalueNode(lhsFieldIndexExpression), nil
}

func NewPositionalFieldValueLvalueNode(
	lhsFieldIndexExpression IEvaluable,
) *PositionalFieldValueLvalueNode {
	return &PositionalFieldValueLvalueNode{
		lhsFieldIndexExpression: lhsFieldIndexExpression,
	}
}

func (this *PositionalFieldValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *PositionalFieldValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		return errors.New("There is no current record to assign to.")
	}

	lhsFieldIndex := this.lhsFieldIndexExpression.Evaluate(state)

	if indices == nil {
		index, ok := lhsFieldIndex.GetIntValue()
		if ok {
			// TODO: incorporate error-return into this API
			//err := state.Inrec.PutCopyWithPositionalIndex(&lhsFieldIndex, rvalue)
			//if err != nil {
			//return err
			//}
			//return nil
			state.Inrec.PutCopyWithPositionalIndex(index, rvalue)
			return nil
		} else {
			return errors.New(
				fmt.Sprintf(
					"Miller: positional index for $[[[...]]] assignment must be integer; got %s.",
					lhsFieldIndex.GetTypeName(),
				),
			)
		}
	} else {
		// xxx positional
		return state.Inrec.PutIndexed(
			append([]*types.Mlrval{lhsFieldIndex}, indices...),
			rvalue,
		)
	}
}

// Same code as PositionalFieldNameLvalueNode.
// May as well let them do 'unset $[[[7]]]' as well as $[[7]]'.
func (this *PositionalFieldValueLvalueNode) Unassign(
	state *runtime.State,
) {
	this.UnassignIndexed(nil, state)
}

func (this *PositionalFieldValueLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		return
	}

	lhsFieldIndex := this.lhsFieldIndexExpression.Evaluate(state)
	if indices == nil {
		index, ok := lhsFieldIndex.GetIntValue()
		if ok {
			state.Inrec.RemoveWithPositionalIndex(index)
		} else {
			// TODO: incorporate error-return into this API
		}
	} else {
		// xxx positional
		state.Inrec.RemoveIndexed(
			append([]*types.Mlrval{lhsFieldIndex}, indices...),
		)
	}
}

// ----------------------------------------------------------------
type FullSrecLvalueNode struct {
}

func (root *RootNode) BuildFullSrecLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeFullSrec)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(astNode.Children != nil)
	return NewFullSrecLvalueNode(), nil
}

func NewFullSrecLvalueNode() *FullSrecLvalueNode {
	return &FullSrecLvalueNode{}
}

func (this *FullSrecLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *FullSrecLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		return errors.New("There is no current record to assign to.")
	}

	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	// The input record is a *Mlrmap so just invoke its PutIndexed.
	err := state.Inrec.PutIndexed(indices, rvalue)
	if err != nil {
		return err
	}
	return nil
}

func (this *FullSrecLvalueNode) Unassign(
	state *runtime.State,
) {
	this.UnassignIndexed(nil, state)
}

func (this *FullSrecLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		return
	}

	if indices == nil {
		state.Inrec.Clear()
	} else {
		state.Inrec.RemoveIndexed(indices)
	}
}

// ----------------------------------------------------------------
type DirectOosvarValueLvalueNode struct {
	lhsOosvarName *types.Mlrval
}

func (root *RootNode) BuildDirectOosvarValueLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDirectOosvarValue)

	lhsOosvarName := types.MlrvalFromString(string(astNode.Token.Lit))
	return NewDirectOosvarValueLvalueNode(&lhsOosvarName), nil
}

func NewDirectOosvarValueLvalueNode(lhsOosvarName *types.Mlrval) *DirectOosvarValueLvalueNode {
	return &DirectOosvarValueLvalueNode{
		lhsOosvarName: lhsOosvarName,
	}
}

func (this *DirectOosvarValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *DirectOosvarValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	// AssignmentNode checks for absent, so we just assign whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())
	if indices == nil {
		err := state.Oosvars.PutCopyWithMlrvalIndex(this.lhsOosvarName, rvalue)
		if err != nil {
			return err
		}
		return nil
	} else {
		return state.Oosvars.PutIndexed(
			append([]*types.Mlrval{this.lhsOosvarName}, indices...),
			rvalue,
		)
	}
}

func (this *DirectOosvarValueLvalueNode) Unassign(
	state *runtime.State,
) {
	this.UnassignIndexed(nil, state)
}

func (this *DirectOosvarValueLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	if indices == nil {
		name := this.lhsOosvarName.String()
		state.Oosvars.Remove(name)
	} else {
		state.Oosvars.RemoveIndexed(
			append([]*types.Mlrval{this.lhsOosvarName}, indices...),
		)
	}
}

// ----------------------------------------------------------------
type IndirectOosvarValueLvalueNode struct {
	lhsOosvarNameExpression IEvaluable
}

func (root *RootNode) BuildIndirectOosvarValueLvalueNode(
	astNode *dsl.ASTNode,
) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectOosvarValue)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	lhsOosvarNameExpression, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return NewIndirectOosvarValueLvalueNode(lhsOosvarNameExpression), nil
}

func NewIndirectOosvarValueLvalueNode(
	lhsOosvarNameExpression IEvaluable,
) *IndirectOosvarValueLvalueNode {
	return &IndirectOosvarValueLvalueNode{
		lhsOosvarNameExpression: lhsOosvarNameExpression,
	}
}

func (this *IndirectOosvarValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *IndirectOosvarValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	lhsOosvarName := this.lhsOosvarNameExpression.Evaluate(state)

	if indices == nil {
		err := state.Oosvars.PutCopyWithMlrvalIndex(lhsOosvarName, rvalue)
		if err != nil {
			return err
		}
		return nil
	} else {
		return state.Oosvars.PutIndexed(
			append([]*types.Mlrval{lhsOosvarName.Copy()}, indices...),
			rvalue,
		)
	}
}

func (this *IndirectOosvarValueLvalueNode) Unassign(
	state *runtime.State,
) {
	this.UnassignIndexed(nil, state)
}

func (this *IndirectOosvarValueLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	lhsOosvarName := this.lhsOosvarNameExpression.Evaluate(state)

	if indices == nil {
		sname := lhsOosvarName.String()
		state.Oosvars.Remove(sname)
	} else {
		state.Oosvars.RemoveIndexed(
			append([]*types.Mlrval{lhsOosvarName}, indices...),
		)
	}
}

// ----------------------------------------------------------------
type FullOosvarLvalueNode struct {
}

func (root *RootNode) BuildFullOosvarLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeFullOosvar)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(astNode.Children != nil)
	return NewFullOosvarLvalueNode(), nil
}

func NewFullOosvarLvalueNode() *FullOosvarLvalueNode {
	return &FullOosvarLvalueNode{}
}

func (this *FullOosvarLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *FullOosvarLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	// The input record is a *Mlrmap so just invoke its PutIndexed.
	err := state.Oosvars.PutIndexed(indices, rvalue)
	if err != nil {
		return err
	}
	return nil
}

func (this *FullOosvarLvalueNode) Unassign(
	state *runtime.State,
) {
	this.UnassignIndexed(nil, state)
}

func (this *FullOosvarLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	if indices == nil {
		state.Oosvars.Clear()
	} else {
		state.Oosvars.RemoveIndexed(indices)
	}
}

// ----------------------------------------------------------------
type LocalVariableLvalueNode struct {
	stackVariable *runtime.StackVariable
	typeName      string

	// a = 1;
	// b = 1;
	// if (true) {
	//   a = 3;     <-- defineTypedAtScope is false; updates outer a
	//   var b = 4; <-- defineTypedAtScope is true;  creates new inner b
	// }
	defineTypedAtScope bool
}

func (root *RootNode) BuildLocalVariableLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeLocalVariable)

	variableName := string(astNode.Token.Lit)
	typeName := "any"
	defineTypedAtScope := false
	if astNode.Children != nil { // typed, like 'num x = 3'
		typeNode := astNode.Children[0]
		lib.InternalCodingErrorIf(typeNode.Type != dsl.NodeTypeTypedecl)
		typeName = string(typeNode.Token.Lit)
		defineTypedAtScope = true
	}
	return NewLocalVariableLvalueNode(
		runtime.NewStackVariable(variableName),
		typeName,
		defineTypedAtScope,
	), nil
}

func NewLocalVariableLvalueNode(
	stackVariable *runtime.StackVariable,
	typeName string,
	defineTypedAtScope bool,
) *LocalVariableLvalueNode {
	return &LocalVariableLvalueNode{
		stackVariable:      stackVariable,
		typeName:           typeName,
		defineTypedAtScope: defineTypedAtScope,
	}
}

func (this *LocalVariableLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *LocalVariableLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	// AssignmentNode checks for absent, so we just assign whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())
	var err error = nil
	if indices == nil {
		if this.defineTypedAtScope {
			err = state.Stack.DefineTypedAtScope(this.stackVariable, this.typeName, rvalue)
		} else {
			err = state.Stack.Set(this.stackVariable, rvalue)
		}
	} else {
		// There is no 'map x[1] = {}' in the DSL grammar.
		lib.InternalCodingErrorIf(this.defineTypedAtScope)

		err = state.Stack.SetIndexed(this.stackVariable, indices, rvalue)
	}
	return err
}

func (this *LocalVariableLvalueNode) Unassign(
	state *runtime.State,
) {
	this.UnassignIndexed(nil, state)
}

func (this *LocalVariableLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	if indices == nil {
		state.Stack.Unset(this.stackVariable)
	} else {
		state.Stack.UnsetIndexed(this.stackVariable, indices)
	}
}

// ----------------------------------------------------------------
// IndexedValueNode is a delegator to base-lvalue types.
// * The baseLvalue is some IAssignable
// * The indexEvaluables are an array of IEvaluables
// * Each needs to evaluate to int or string
// * Assignment needs to walk each level:
//   o error if ith mlrval is int and that level isn't an array
//   o error if ith mlrval is string and that level isn't a map
//   o error for any other types -- maybe absent-handling for heterogeneity ...

// ----------------------------------------------------------------
type IndexedLvalueNode struct {
	baseLvalue      IAssignable
	indexEvaluables []IEvaluable
}

func (root *RootNode) BuildIndexedLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayOrMapIndexAccess)
	lib.InternalCodingErrorIf(astNode == nil)

	var baseLvalue IAssignable = nil
	indexEvaluables := make([]IEvaluable, 0)
	var err error = nil

	// $ mlr -n put -v '$x[1][2]=3'
	// DSL EXPRESSION:
	// $x[1][2]=3
	// AST:
	// * StatementBlock
	//     * Assignment "="
	//         * ArrayOrMapIndexAccess "[]"
	//             * ArrayOrMapIndexAccess "[]"
	//                 * DirectFieldValue "x"
	//                 * IntLiteral "1"
	//             * IntLiteral "2"
	//         * IntLiteral "3"

	// In the AST, the indices come in last-shallowest, down to first-deepest,
	// then the base Lvalue.
	walkerNode := astNode
	for {
		if walkerNode.Type == dsl.NodeTypeArrayOrMapIndexAccess {
			lib.InternalCodingErrorIf(walkerNode == nil)
			lib.InternalCodingErrorIf(len(walkerNode.Children) != 2)
			indexEvaluable, err := root.BuildEvaluableNode(walkerNode.Children[1])
			if err != nil {
				return nil, err
			}
			indexEvaluables = append([]IEvaluable{indexEvaluable}, indexEvaluables...)
			walkerNode = walkerNode.Children[0]
		} else {
			baseLvalue, err = root.BuildAssignableNode(walkerNode)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	return NewIndexedLvalueNode(baseLvalue, indexEvaluables), nil
}

func NewIndexedLvalueNode(
	baseLvalue IAssignable,
	indexEvaluables []IEvaluable,
) *IndexedLvalueNode {
	return &IndexedLvalueNode{
		baseLvalue:      baseLvalue,
		indexEvaluables: indexEvaluables,
	}
}

// Computes Lvalue indices and then delegates to the baseLvalue.  E.g. for
// '$x[1][2] = 3' or '@x[1][2] = 3', the indices are [1,2], and the baseLvalue
// is '$x' or '@x' respectively.
func (this *IndexedLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	indices := make([]*types.Mlrval, len(this.indexEvaluables))

	for i, _ := range this.indexEvaluables {
		indices[i] = this.indexEvaluables[i].Evaluate(state)
		if indices[i].IsAbsent() {
			return nil
		}
	}

	// This lets the user do '$y[ ["a", "b", "c"] ] = $x' in lieu of
	// '$y["a"]["b"]["c"] = $x'.
	if len(indices) == 1 && indices[0].IsArray() {
		indices = types.MakePointerArray(indices[0].GetArray())
	}

	return this.baseLvalue.AssignIndexed(rvalue, indices, state)
}

func (this *IndexedLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	// We are the delegator, not the delegatee
	lib.InternalCodingErrorIf(true)
	return nil // not reached
}

func (this *IndexedLvalueNode) Unassign(
	state *runtime.State,
) {
	indices := make([]*types.Mlrval, len(this.indexEvaluables))
	for i, _ := range this.indexEvaluables {
		indices[i] = this.indexEvaluables[i].Evaluate(state)
	}

	this.baseLvalue.UnassignIndexed(indices, state)
}

func (this *IndexedLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	// We are the delegator, not the delegatee
	lib.InternalCodingErrorIf(true)
}

// ----------------------------------------------------------------
type EnvironmentVariableLvalueNode struct {
	nameExpression IEvaluable
}

func (root *RootNode) BuildEnvironmentVariableLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEnvironmentVariable)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	nameExpression, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return NewEnvironmentVariableLvalueNode(nameExpression), nil
}

func NewEnvironmentVariableLvalueNode(
	nameExpression IEvaluable,
) *EnvironmentVariableLvalueNode {
	return &EnvironmentVariableLvalueNode{
		nameExpression: nameExpression,
	}
}

func (this *EnvironmentVariableLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	name := this.nameExpression.Evaluate(state)
	if name.IsAbsent() {
		return nil
	}

	if !name.IsString() {
		return errors.New(
			fmt.Sprintf(
				"ENV[...] assignments must have string names; got %s \"%s\"\n",
				name.GetTypeName(),
				name.String(),
			),
		)
	}

	os.Setenv(name.String(), rvalue.String())
	return nil
}

func (this *EnvironmentVariableLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	return errors.New("Miller: ENV[...] cannot be indexed.")
}

func (this *EnvironmentVariableLvalueNode) Unassign(
	state *runtime.State,
) {
	name := this.nameExpression.Evaluate(state)
	if name.IsAbsent() {
		return
	}

	if !name.IsString() {
		// TODO: needs error-return
		return
	}

	os.Unsetenv(name.String())
}

func (this *EnvironmentVariableLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	// TODO: needs error return
	//return errors.New("Miller: ENV[...] cannot be indexed.")
}
