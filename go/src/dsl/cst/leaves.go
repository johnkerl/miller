// ================================================================
// CST build/execute for AST leaf nodes
// ================================================================

package cst

import (
	"errors"
	"math"

	"miller/src/dsl"
	"miller/src/lib"
	"miller/src/runtime"
	"miller/src/types"
)

// ----------------------------------------------------------------
func (this *RootNode) BuildLeafNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Children != nil)
	sval := string(astNode.Token.Lit)

	switch astNode.Type {

	case dsl.NodeTypeDirectFieldValue:
		return this.BuildDirectFieldRvalueNode(sval), nil
		break
	case dsl.NodeTypeFullSrec:
		return this.BuildFullSrecRvalueNode(), nil
		break

	case dsl.NodeTypeDirectOosvarValue:
		return this.BuildDirectOosvarRvalueNode(sval), nil
		break
	case dsl.NodeTypeFullOosvar:
		return this.BuildFullOosvarRvalueNode(), nil
		break

	case dsl.NodeTypeLocalVariable:
		return this.BuildLocalVariableNode(sval), nil
		break

	case dsl.NodeTypeStringLiteral:
		return this.BuildStringLiteralNode(sval), nil
		break
	case dsl.NodeTypeRegexCaseInsensitive:
		// StringLiteral nodes like '"abc"' entered by the user come in from
		// the AST as 'abc', with double quotes removed. Case-insensitive
		// regexes like '"a.*b"i' come in with initial '"' and final '"i'
		// intact. We let the sub/regextract/etc functions deal with this.
		// (The alternative would be to make a separate Mlrval type separate
		// from string.)
		return this.BuildStringLiteralNode(sval), nil
		break
	case dsl.NodeTypeIntLiteral:
		return this.BuildIntLiteralNode(sval), nil
		break
	case dsl.NodeTypeFloatLiteral:
		return this.BuildFloatLiteralNode(sval), nil
		break
	case dsl.NodeTypeBoolLiteral:
		return this.BuildBoolLiteralNode(sval), nil
		break
	case dsl.NodeTypeContextVariable:
		return this.BuildContextVariableNode(astNode)
		break
	case dsl.NodeTypeConstant:
		return this.BuildConstantNode(astNode)
		break

	case dsl.NodeTypeArraySliceEmptyLowerIndex:
		return this.BuildArraySliceEmptyLowerIndexNode(astNode)
		break
	case dsl.NodeTypeArraySliceEmptyUpperIndex:
		return this.BuildArraySliceEmptyUpperIndexNode(astNode)
		break

	case dsl.NodeTypePanic:
		return this.BuildPanicNode(astNode)
		break
	}

	return nil, errors.New(
		"CST BuildLeafNode: unhandled AST node " + string(astNode.Type),
	)
}

// ----------------------------------------------------------------
type DirectFieldRvalueNode struct {
	fieldName string
}

func (this *RootNode) BuildDirectFieldRvalueNode(fieldName string) *DirectFieldRvalueNode {
	return &DirectFieldRvalueNode{
		fieldName: fieldName,
	}
}
func (this *DirectFieldRvalueNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		output.SetFromAbsent()
		return
	}
	value := state.Inrec.Get(this.fieldName)
	if value == nil {
		output.SetFromAbsent()
	} else {
		output.CopyFrom(value)
	}
}

// ----------------------------------------------------------------
type FullSrecRvalueNode struct {
}

func (this *RootNode) BuildFullSrecRvalueNode() *FullSrecRvalueNode {
	return &FullSrecRvalueNode{}
}
func (this *FullSrecRvalueNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		output.SetFromAbsent()
	} else {
		output.SetFromMap(state.Inrec)
	}
}

// ----------------------------------------------------------------
type DirectOosvarRvalueNode struct {
	variableName string
}

func (this *RootNode) BuildDirectOosvarRvalueNode(variableName string) *DirectOosvarRvalueNode {
	return &DirectOosvarRvalueNode{
		variableName: variableName,
	}
}
func (this *DirectOosvarRvalueNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	value := state.Oosvars.Get(this.variableName)
	if value == nil {
		output.SetFromAbsent()
	} else {
		output.CopyFrom(value)
	}
}

// ----------------------------------------------------------------
type FullOosvarRvalueNode struct {
}

func (this *RootNode) BuildFullOosvarRvalueNode() *FullOosvarRvalueNode {
	return &FullOosvarRvalueNode{}
}
func (this *FullOosvarRvalueNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromMap(state.Oosvars)
}

// ----------------------------------------------------------------
type LocalVariableNode struct {
	variableName string
}

func (this *RootNode) BuildLocalVariableNode(variableName string) *LocalVariableNode {
	return &LocalVariableNode{
		variableName: variableName,
	}
}
func (this *LocalVariableNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	value := state.Stack.Get(this.variableName)
	if value == nil {
		output.SetFromAbsent()
	} else {
		output.CopyFrom(value)
	}
}

// ----------------------------------------------------------------
type StringLiteralNode struct {
	literal types.Mlrval
}

func (this *RootNode) BuildStringLiteralNode(literal string) *StringLiteralNode {
	return &StringLiteralNode{
		literal: types.MlrvalFromString(literal),
	}
}
func (this *StringLiteralNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.CopyFrom(&this.literal)
}

// ----------------------------------------------------------------
type IntLiteralNode struct {
	literal types.Mlrval
}

func (this *RootNode) BuildIntLiteralNode(literal string) *IntLiteralNode {
	return &IntLiteralNode{
		literal: types.MlrvalFromIntString(literal),
	}
}
func (this *IntLiteralNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.CopyFrom(&this.literal)
}

// ----------------------------------------------------------------
type FloatLiteralNode struct {
	literal types.Mlrval
}

func (this *RootNode) BuildFloatLiteralNode(literal string) *FloatLiteralNode {
	return &FloatLiteralNode{
		literal: types.MlrvalFromFloat64String(literal),
	}
}
func (this *FloatLiteralNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.CopyFrom(&this.literal)
}

// ----------------------------------------------------------------
type BoolLiteralNode struct {
	literal types.Mlrval
}

func (this *RootNode) BuildBoolLiteralNode(literal string) *BoolLiteralNode {
	return &BoolLiteralNode{
		literal: types.MlrvalFromBoolString(literal),
	}
}
func (this *BoolLiteralNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.CopyFrom(&this.literal)
}

// ================================================================
func (this *RootNode) BuildContextVariableNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Token == nil)
	sval := string(astNode.Token.Lit)

	switch sval {

	case "FILENAME":
		return this.BuildFILENAMENode(), nil
		break
	case "FILENUM":
		return this.BuildFILENUMNode(), nil
		break

	case "NF":
		return this.BuildNFNode(), nil
		break
	case "NR":
		return this.BuildNRNode(), nil
		break
	case "FNR":
		return this.BuildFNRNode(), nil
		break

	case "IRS":
		return this.BuildIRSNode(), nil
		break
	case "IFS":
		return this.BuildIFSNode(), nil
		break
	case "IPS":
		return this.BuildIPSNode(), nil
		break

	case "ORS":
		return this.BuildORSNode(), nil
		break
	case "OFS":
		return this.BuildOFSNode(), nil
		break
	case "OPS":
		return this.BuildOPSNode(), nil
		break
	case "OFLATSEP":
		return this.BuildOFLATSEPNode(), nil
		break

	}

	return nil, errors.New(
		"CST BuildContextVariableNode: unhandled context variable " + sval,
	)
}

// ----------------------------------------------------------------
type FILENAMENode struct {
}

func (this *RootNode) BuildFILENAMENode() *FILENAMENode {
	return &FILENAMENode{}
}
func (this *FILENAMENode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromString(state.Context.FILENAME)
}

// ----------------------------------------------------------------
type FILENUMNode struct {
}

func (this *RootNode) BuildFILENUMNode() *FILENUMNode {
	return &FILENUMNode{}
}
func (this *FILENUMNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromInt(state.Context.FILENUM)
}

// ----------------------------------------------------------------
type NFNode struct {
}

func (this *RootNode) BuildNFNode() *NFNode {
	return &NFNode{}
}
func (this *NFNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromInt(state.Inrec.FieldCount)
}

// ----------------------------------------------------------------
type NRNode struct {
}

func (this *RootNode) BuildNRNode() *NRNode {
	return &NRNode{}
}
func (this *NRNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromInt(state.Context.NR)
}

// ----------------------------------------------------------------
type FNRNode struct {
}

func (this *RootNode) BuildFNRNode() *FNRNode {
	return &FNRNode{}
}
func (this *FNRNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromInt(state.Context.FNR)
}

// ----------------------------------------------------------------
type IRSNode struct {
}

func (this *RootNode) BuildIRSNode() *IRSNode {
	return &IRSNode{}
}
func (this *IRSNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromString(state.Context.IRS)
}

// ----------------------------------------------------------------
type IFSNode struct {
}

func (this *RootNode) BuildIFSNode() *IFSNode {
	return &IFSNode{}
}
func (this *IFSNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromString(state.Context.IFS)
}

// ----------------------------------------------------------------
type IPSNode struct {
}

func (this *RootNode) BuildIPSNode() *IPSNode {
	return &IPSNode{}
}
func (this *IPSNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromString(state.Context.IPS)
}

// ----------------------------------------------------------------
type ORSNode struct {
}

func (this *RootNode) BuildORSNode() *ORSNode {
	return &ORSNode{}
}
func (this *ORSNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromString(state.Context.ORS)
}

// ----------------------------------------------------------------
type OFSNode struct {
}

func (this *RootNode) BuildOFSNode() *OFSNode {
	return &OFSNode{}
}
func (this *OFSNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromString(state.Context.OFS)
}

// ----------------------------------------------------------------
type OPSNode struct {
}

func (this *RootNode) BuildOPSNode() *OPSNode {
	return &OPSNode{}
}
func (this *OPSNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromString(state.Context.OPS)
}

// ----------------------------------------------------------------
type OFLATSEPNode struct {
}

func (this *RootNode) BuildOFLATSEPNode() *OFLATSEPNode {
	return &OFLATSEPNode{}
}
func (this *OFLATSEPNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromString(state.Context.OFLATSEP)
}

// ================================================================
func (this *RootNode) BuildConstantNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Token == nil)
	sval := string(astNode.Token.Lit)

	switch sval {

	case "M_PI":
		return this.BuildMathPINode(), nil
		break
	case "M_E":
		return this.BuildMathENode(), nil
		break

	}

	return nil, errors.New(
		"CST BuildContextVariableNode: unhandled context variable " + sval,
	)
}

// ----------------------------------------------------------------
type MathPINode struct {
}

func (this *RootNode) BuildMathPINode() *MathPINode {
	return &MathPINode{}
}
func (this *MathPINode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromFloat64(math.Pi)
}

// ----------------------------------------------------------------
type MathENode struct {
}

func (this *RootNode) BuildMathENode() *MathENode {
	return &MathENode{}
}
func (this *MathENode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromFloat64(math.E)
}

// ================================================================
type LiteralOneNode struct {
}

// In array slices like 'myarray[:4]', the lower index is always 1 since Miller
// user-space indices are 1-up.
func (this *RootNode) BuildArraySliceEmptyLowerIndexNode(
	astNode *dsl.ASTNode,
) (*LiteralOneNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArraySliceEmptyLowerIndex)
	return &LiteralOneNode{}, nil
}
func (this *LiteralOneNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromInt(1)
}

// ================================================================
type LiteralEmptyStringNode struct {
}

// In array slices like 'myarray[4:]', the upper index is always n, where n is
// the length of the array, since Miller user-space indices are 1-up. However,
// we don't have access to the array length in this AST node so we return ""
// so the slice-index CST node can compute it.
func (this *RootNode) BuildArraySliceEmptyUpperIndexNode(
	astNode *dsl.ASTNode,
) (*LiteralEmptyStringNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArraySliceEmptyUpperIndex)
	return &LiteralEmptyStringNode{}, nil
}
func (this *LiteralEmptyStringNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromString("")
}

// ----------------------------------------------------------------
// The panic token is a special token which causes a panic when evaluated.
// This is for testing that AND/OR short-circuiting is implemented correctly:
// output = input1 || panic should NOT panic the process when input1 is true.

type PanicNode struct {
}

func (this *RootNode) BuildPanicNode(astNode *dsl.ASTNode) (*PanicNode, error) {
	return &PanicNode{}, nil
}
func (this *PanicNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	lib.InternalCodingErrorPanic("Panic token was evaluated, not short-circuited.")
}
