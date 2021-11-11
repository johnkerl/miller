// ================================================================
// This is for Lvalues, i.e. things on the left-hand-side of an assignment
// statement.
// ================================================================

package cst

import (
	"errors"
	"fmt"
	"os"

	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/runtime"
	"mlr/internal/pkg/types"
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
			"mlr: '[[...]]' is allowed on assignment left-hand sides only when immediately preceded by '$'.",
		)
		break
	case dsl.NodeTypeArrayOrMapPositionalValueAccess:
		return nil, errors.New(
			"mlr: '[[[...]]]' is allowed on assignment left-hand sides only when immediately preceded by '$'.",
		)
		break

	case dsl.NodeTypeArrayOrMapIndexAccess:
		return root.BuildIndexedLvalueNode(astNode)
		break

	case dsl.NodeTypeDotOperator:
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
	return NewDirectFieldValueLvalueNode(lhsFieldName), nil
}

func NewDirectFieldValueLvalueNode(lhsFieldName *types.Mlrval) *DirectFieldValueLvalueNode {
	return &DirectFieldValueLvalueNode{
		lhsFieldName: lhsFieldName,
	}
}

func (node *DirectFieldValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return node.AssignIndexed(rvalue, nil, state)
}

func (node *DirectFieldValueLvalueNode) AssignIndexed(
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
		err := state.Inrec.PutCopyWithMlrvalIndex(node.lhsFieldName, rvalue)
		if err != nil {
			return err
		}
		return nil
	} else {
		return state.Inrec.PutIndexed(
			append([]*types.Mlrval{node.lhsFieldName}, indices...),
			rvalue,
		)
	}
}

func (node *DirectFieldValueLvalueNode) Unassign(
	state *runtime.State,
) {
	node.UnassignIndexed(nil, state)
}

func (node *DirectFieldValueLvalueNode) UnassignIndexed(
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
		lib.InternalCodingErrorIf(!node.lhsFieldName.IsString())
		name := node.lhsFieldName.String()
		state.Inrec.Remove(name)
	} else {
		state.Inrec.RemoveIndexed(
			append([]*types.Mlrval{node.lhsFieldName}, indices...),
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

func (node *IndirectFieldValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return node.AssignIndexed(rvalue, nil, state)
}

func (node *IndirectFieldValueLvalueNode) AssignIndexed(
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

	lhsFieldName := node.lhsFieldNameExpression.Evaluate(state)

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

func (node *IndirectFieldValueLvalueNode) Unassign(
	state *runtime.State,
) {
	node.UnassignIndexed(nil, state)
}

func (node *IndirectFieldValueLvalueNode) UnassignIndexed(
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

	lhsFieldName := node.lhsFieldNameExpression.Evaluate(state)
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

func (node *PositionalFieldNameLvalueNode) Assign(
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

	lhsFieldIndex := node.lhsFieldIndexExpression.Evaluate(state)

	index, ok := lhsFieldIndex.GetIntValue()
	if ok {
		// TODO: incorporate error-return into this API
		state.Inrec.PutNameWithPositionalIndex(index, rvalue)
		return nil
	} else {
		return errors.New(
			fmt.Sprintf(
				"mlr: positional index for $[[...]] assignment must be integer; got %s.",
				lhsFieldIndex.GetTypeName(),
			),
		)
	}
}

func (node *PositionalFieldNameLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	// TODO: reconsider this if /when we decide to allow string-slice
	// assignments.
	return errors.New(
		"mlr: $[[...]] = ... expressions are not indexable.",
	)
}

func (node *PositionalFieldNameLvalueNode) Unassign(
	state *runtime.State,
) {
	node.UnassignIndexed(nil, state)
}

func (node *PositionalFieldNameLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	lhsFieldIndex := node.lhsFieldIndexExpression.Evaluate(state)

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

func (node *PositionalFieldValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return node.AssignIndexed(rvalue, nil, state)
}

func (node *PositionalFieldValueLvalueNode) AssignIndexed(
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

	lhsFieldIndex := node.lhsFieldIndexExpression.Evaluate(state)

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
					"mlr: positional index for $[[[...]]] assignment must be integer; got %s.",
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
func (node *PositionalFieldValueLvalueNode) Unassign(
	state *runtime.State,
) {
	node.UnassignIndexed(nil, state)
}

func (node *PositionalFieldValueLvalueNode) UnassignIndexed(
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

	lhsFieldIndex := node.lhsFieldIndexExpression.Evaluate(state)
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

func (node *FullSrecLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return node.AssignIndexed(rvalue, nil, state)
}

func (node *FullSrecLvalueNode) AssignIndexed(
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

func (node *FullSrecLvalueNode) Unassign(
	state *runtime.State,
) {
	node.UnassignIndexed(nil, state)
}

func (node *FullSrecLvalueNode) UnassignIndexed(
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
	return NewDirectOosvarValueLvalueNode(lhsOosvarName), nil
}

func NewDirectOosvarValueLvalueNode(lhsOosvarName *types.Mlrval) *DirectOosvarValueLvalueNode {
	return &DirectOosvarValueLvalueNode{
		lhsOosvarName: lhsOosvarName,
	}
}

func (node *DirectOosvarValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return node.AssignIndexed(rvalue, nil, state)
}

func (node *DirectOosvarValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	// AssignmentNode checks for absent, so we just assign whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())
	if indices == nil {
		err := state.Oosvars.PutCopyWithMlrvalIndex(node.lhsOosvarName, rvalue)
		if err != nil {
			return err
		}
		return nil
	} else {
		return state.Oosvars.PutIndexed(
			append([]*types.Mlrval{node.lhsOosvarName}, indices...),
			rvalue,
		)
	}
}

func (node *DirectOosvarValueLvalueNode) Unassign(
	state *runtime.State,
) {
	node.UnassignIndexed(nil, state)
}

func (node *DirectOosvarValueLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	if indices == nil {
		name := node.lhsOosvarName.String()
		state.Oosvars.Remove(name)
	} else {
		state.Oosvars.RemoveIndexed(
			append([]*types.Mlrval{node.lhsOosvarName}, indices...),
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

func (node *IndirectOosvarValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return node.AssignIndexed(rvalue, nil, state)
}

func (node *IndirectOosvarValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	lhsOosvarName := node.lhsOosvarNameExpression.Evaluate(state)

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

func (node *IndirectOosvarValueLvalueNode) Unassign(
	state *runtime.State,
) {
	node.UnassignIndexed(nil, state)
}

func (node *IndirectOosvarValueLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	lhsOosvarName := node.lhsOosvarNameExpression.Evaluate(state)

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

func (node *FullOosvarLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return node.AssignIndexed(rvalue, nil, state)
}

func (node *FullOosvarLvalueNode) AssignIndexed(
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

func (node *FullOosvarLvalueNode) Unassign(
	state *runtime.State,
) {
	node.UnassignIndexed(nil, state)
}

func (node *FullOosvarLvalueNode) UnassignIndexed(
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

func (node *LocalVariableLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	return node.AssignIndexed(rvalue, nil, state)
}

func (node *LocalVariableLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	// AssignmentNode checks for absent, so we just assign whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())
	var err error = nil
	if indices == nil {
		if node.defineTypedAtScope {
			err = state.Stack.DefineTypedAtScope(node.stackVariable, node.typeName, rvalue)
		} else {
			err = state.Stack.Set(node.stackVariable, rvalue)
		}
	} else {
		// There is no 'map x[1] = {}' in the DSL grammar.
		lib.InternalCodingErrorIf(node.defineTypedAtScope)

		err = state.Stack.SetIndexed(node.stackVariable, indices, rvalue)
	}
	return err
}

func (node *LocalVariableLvalueNode) Unassign(
	state *runtime.State,
) {
	node.UnassignIndexed(nil, state)
}

func (node *LocalVariableLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	if indices == nil {
		state.Stack.Unset(node.stackVariable)
	} else {
		state.Stack.UnsetIndexed(node.stackVariable, indices)
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

// Either 'mymap["attr"]' or 'mymap.attr'. Furthermore they can be mixed as in
// 'mymap["foo"].bar' or 'mymap.foo["bar"]'.
func (root *RootNode) BuildIndexedLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayOrMapIndexAccess && astNode.Type != dsl.NodeTypeDotOperator)
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
		} else if walkerNode.Type == dsl.NodeTypeDotOperator {
			lib.InternalCodingErrorIf(walkerNode == nil)
			lib.InternalCodingErrorIf(len(walkerNode.Children) != 2)
			indexEvaluable := root.BuildStringLiteralNode(string(walkerNode.Children[1].Token.Lit))
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

// BuildDottedLvalueNode is basically the same as BuildIndexedLvalueNode except
// at the syntax level:
// 'mymap["x"]["y"]' is the same as 'mymap.x.y'.
func (root *RootNode) BuildDottedLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDotOperator)
	lib.InternalCodingErrorIf(astNode == nil)

	var baseLvalue IAssignable = nil
	indexEvaluables := make([]IEvaluable, 0)
	var err error = nil

	// $ mlr -n put -v '$x.a.b=3'
	// DSL EXPRESSION:
	// $x.a.b=3
	// AST:
	// * statement block
	//     * assignment "="
	//         * dot operator "."
	//             * dot operator "."
	//                 * direct field value "x"
	//                 * local variable "a"
	//             * local variable "b"
	//         * int literal "3"

	// In the AST, the indices come in last-shallowest, down to first-deepest,
	// then the base Lvalue.
	walkerNode := astNode
	for {
		if walkerNode.Type == dsl.NodeTypeDotOperator {
			lib.InternalCodingErrorIf(walkerNode == nil)
			lib.InternalCodingErrorIf(len(walkerNode.Children) != 2)
			indexEvaluable, err := root.BuildEvaluableNode(walkerNode.Children[1])
			walkerNode.Children[1].Print()
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
func (node *IndexedLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	indices := make([]*types.Mlrval, len(node.indexEvaluables))

	for i := range node.indexEvaluables {
		indices[i] = node.indexEvaluables[i].Evaluate(state)
		if indices[i].IsAbsent() {
			return nil
		}
	}

	// This lets the user do '$y[ ["a", "b", "c"] ] = $x' in lieu of
	// '$y["a"]["b"]["c"] = $x'.
	if len(indices) == 1 && indices[0].IsArray() {
		indices = types.MakePointerArray(indices[0].GetArray())
	}

	return node.baseLvalue.AssignIndexed(rvalue, indices, state)
}

func (node *IndexedLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	// We are the delegator, not the delegatee
	lib.InternalCodingErrorIf(true)
	return nil // not reached
}

func (node *IndexedLvalueNode) Unassign(
	state *runtime.State,
) {
	indices := make([]*types.Mlrval, len(node.indexEvaluables))
	for i := range node.indexEvaluables {
		indices[i] = node.indexEvaluables[i].Evaluate(state)
	}

	node.baseLvalue.UnassignIndexed(indices, state)
}

func (node *IndexedLvalueNode) UnassignIndexed(
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

func (node *EnvironmentVariableLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *runtime.State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	name := node.nameExpression.Evaluate(state)
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

	sname := name.String()
	svalue := rvalue.String()
	os.Setenv(sname, svalue)
	if sname == "TZ" {
		lib.SetTZFromEnv() // affects the time library; notify it
	}
	return nil
}

func (node *EnvironmentVariableLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *runtime.State,
) error {
	return errors.New("mlr: ENV[...] cannot be indexed.")
}

func (node *EnvironmentVariableLvalueNode) Unassign(
	state *runtime.State,
) {
	name := node.nameExpression.Evaluate(state)
	if name.IsAbsent() {
		return
	}

	if !name.IsString() {
		// TODO: needs error-return
		return
	}

	os.Unsetenv(name.String())
}

func (node *EnvironmentVariableLvalueNode) UnassignIndexed(
	indices []*types.Mlrval,
	state *runtime.State,
) {
	// TODO: needs error return
	//return errors.New("mlr: ENV[...] cannot be indexed.")
}
