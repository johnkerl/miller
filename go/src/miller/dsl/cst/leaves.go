package cst

import (
	"errors"
	"math"

	"miller/dsl"
	"miller/lib"
	"miller/types"
)

// ================================================================
// CST build/execute for AST leaf nodes
// ================================================================

// ----------------------------------------------------------------
func BuildLeafNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Children != nil)
	sval := string(astNode.Token.Lit)

	switch astNode.Type {

	case dsl.NodeTypeDirectFieldValue:
		return BuildDirectFieldRvalueNode(sval), nil
		break
	case dsl.NodeTypeFullSrec:
		return BuildFullSrecRvalueNode(sval), nil
		break

	case dsl.NodeTypeDirectOosvarValue:
		return BuildDirectOosvarRvalueNode(sval), nil
		break
	case dsl.NodeTypeFullOosvar:
		return BuildFullOosvarRvalueNode(sval), nil
		break

	case dsl.NodeTypeLocalVariable:
		return BuildLocalVariableNode(sval), nil
		break

	case dsl.NodeTypeStringLiteral:
		return BuildStringLiteralNode(sval), nil
		break
	case dsl.NodeTypeIntLiteral:
		return BuildIntLiteralNode(sval), nil
		break
	case dsl.NodeTypeFloatLiteral:
		return BuildFloatLiteralNode(sval), nil
		break
	case dsl.NodeTypeBoolLiteral:
		return BuildBoolLiteralNode(sval), nil
		break
	case dsl.NodeTypeContextVariable:
		return BuildContextVariableNode(astNode)
		break
	case dsl.NodeTypeConstant:
		return BuildConstantNode(astNode)
		break

	case dsl.NodeTypePanic:
		return BuildPanicNode(astNode)
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

func BuildDirectFieldRvalueNode(fieldName string) *DirectFieldRvalueNode {
	return &DirectFieldRvalueNode{
		fieldName: fieldName,
	}
}
func (this *DirectFieldRvalueNode) Evaluate(state *State) types.Mlrval {
	value := state.Inrec.Get(&this.fieldName)
	if value == nil {
		return types.MlrvalFromAbsent()
	} else {
		return *value
	}
}

// ----------------------------------------------------------------
type FullSrecRvalueNode struct {
}

func BuildFullSrecRvalueNode(fieldName string) *FullSrecRvalueNode {
	return &FullSrecRvalueNode{}
}
func (this *FullSrecRvalueNode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromMap(state.Inrec)
}

// ----------------------------------------------------------------
type DirectOosvarRvalueNode struct {
	variableName string
}

func BuildDirectOosvarRvalueNode(variableName string) *DirectOosvarRvalueNode {
	return &DirectOosvarRvalueNode{
		variableName: variableName,
	}
}
func (this *DirectOosvarRvalueNode) Evaluate(state *State) types.Mlrval {
	value := state.Oosvars.Get(&this.variableName)
	if value == nil {
		return types.MlrvalFromAbsent()
	} else {
		return *value
	}
}

// ----------------------------------------------------------------
type FullOosvarRvalueNode struct {
}

func BuildFullOosvarRvalueNode(fieldName string) *FullOosvarRvalueNode {
	return &FullOosvarRvalueNode{}
}
func (this *FullOosvarRvalueNode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromMap(state.Oosvars)
}

// ----------------------------------------------------------------
type LocalVariableNode struct {
	variableName string
}

func BuildLocalVariableNode(variableName string) *LocalVariableNode {
	return &LocalVariableNode{
		variableName: variableName,
	}
}
func (this *LocalVariableNode) Evaluate(state *State) types.Mlrval {
	value := state.stack.ReadVariable(this.variableName)
	//state.stack.Dump()
	if value == nil {
		return types.MlrvalFromAbsent()
	} else {
		return *value
	}
}

// ----------------------------------------------------------------
type StringLiteralNode struct {
	literal types.Mlrval
}

func BuildStringLiteralNode(literal string) *StringLiteralNode {
	return &StringLiteralNode{
		literal: types.MlrvalFromString(literal),
	}
}
func (this *StringLiteralNode) Evaluate(state *State) types.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type IntLiteralNode struct {
	literal types.Mlrval
}

func BuildIntLiteralNode(literal string) *IntLiteralNode {
	return &IntLiteralNode{
		literal: types.MlrvalFromInt64String(literal),
	}
}
func (this *IntLiteralNode) Evaluate(state *State) types.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type FloatLiteralNode struct {
	literal types.Mlrval
}

func BuildFloatLiteralNode(literal string) *FloatLiteralNode {
	return &FloatLiteralNode{
		literal: types.MlrvalFromFloat64String(literal),
	}
}
func (this *FloatLiteralNode) Evaluate(state *State) types.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type BoolLiteralNode struct {
	literal types.Mlrval
}

func BuildBoolLiteralNode(literal string) *BoolLiteralNode {
	return &BoolLiteralNode{
		literal: types.MlrvalFromBoolString(literal),
	}
}
func (this *BoolLiteralNode) Evaluate(state *State) types.Mlrval {
	return this.literal
}

// ================================================================
func BuildContextVariableNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Token == nil)
	sval := string(astNode.Token.Lit)

	switch sval {

	case "FILENAME":
		return BuildFILENAMENode(), nil
		break
	case "FILENUM":
		return BuildFILENUMNode(), nil
		break

	case "NF":
		return BuildNFNode(), nil
		break
	case "NR":
		return BuildNRNode(), nil
		break
	case "FNR":
		return BuildFNRNode(), nil
		break

	case "IRS":
		return BuildIRSNode(), nil
		break
	case "IFS":
		return BuildIFSNode(), nil
		break
	case "IPS":
		return BuildIPSNode(), nil
		break

	case "ORS":
		return BuildORSNode(), nil
		break
	case "OFS":
		return BuildOFSNode(), nil
		break
	case "OPS":
		return BuildOPSNode(), nil
		break

	}

	return nil, errors.New(
		"CST BuildContextVariableNode: unhandled context variable " + sval,
	)
}

// ----------------------------------------------------------------
type FILENAMENode struct {
}

func BuildFILENAMENode() *FILENAMENode {
	return &FILENAMENode{}
}
func (this *FILENAMENode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromString(state.Context.FILENAME)
}

// ----------------------------------------------------------------
type FILENUMNode struct {
}

func BuildFILENUMNode() *FILENUMNode {
	return &FILENUMNode{}
}
func (this *FILENUMNode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromInt64(state.Context.FILENUM)
}

// ----------------------------------------------------------------
type NFNode struct {
}

func BuildNFNode() *NFNode {
	return &NFNode{}
}
func (this *NFNode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromInt64(state.Context.NF)
}

// ----------------------------------------------------------------
type NRNode struct {
}

func BuildNRNode() *NRNode {
	return &NRNode{}
}
func (this *NRNode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromInt64(state.Context.NR)
}

// ----------------------------------------------------------------
type FNRNode struct {
}

func BuildFNRNode() *FNRNode {
	return &FNRNode{}
}
func (this *FNRNode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromInt64(state.Context.FNR)
}

// ----------------------------------------------------------------
type IRSNode struct {
}

func BuildIRSNode() *IRSNode {
	return &IRSNode{}
}
func (this *IRSNode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromString(state.Context.IRS)
}

// ----------------------------------------------------------------
type IFSNode struct {
}

func BuildIFSNode() *IFSNode {
	return &IFSNode{}
}
func (this *IFSNode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromString(state.Context.IFS)
}

// ----------------------------------------------------------------
type IPSNode struct {
}

func BuildIPSNode() *IPSNode {
	return &IPSNode{}
}
func (this *IPSNode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromString(state.Context.IPS)
}

// ----------------------------------------------------------------
type ORSNode struct {
}

func BuildORSNode() *ORSNode {
	return &ORSNode{}
}
func (this *ORSNode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromString(state.Context.ORS)
}

// ----------------------------------------------------------------
type OFSNode struct {
}

func BuildOFSNode() *OFSNode {
	return &OFSNode{}
}
func (this *OFSNode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromString(state.Context.OFS)
}

// ----------------------------------------------------------------
type OPSNode struct {
}

func BuildOPSNode() *OPSNode {
	return &OPSNode{}
}
func (this *OPSNode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromString(state.Context.OPS)
}

// ================================================================
func BuildConstantNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Token == nil)
	sval := string(astNode.Token.Lit)

	switch sval {

	case "M_PI":
		return BuildMathPINode(), nil
		break
	case "M_E":
		return BuildMathENode(), nil
		break

	}

	return nil, errors.New(
		"CST BuildContextVariableNode: unhandled context variable " + sval,
	)
}

// ----------------------------------------------------------------
type MathPINode struct {
}

func BuildMathPINode() *MathPINode {
	return &MathPINode{}
}
func (this *MathPINode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromFloat64(math.Pi)
}

// ----------------------------------------------------------------
type MathENode struct {
}

func BuildMathENode() *MathENode {
	return &MathENode{}
}
func (this *MathENode) Evaluate(state *State) types.Mlrval {
	return types.MlrvalFromFloat64(math.E)
}

// ----------------------------------------------------------------
// The panic token is a special token which causes a panic when evaluated.
// This is for testing that AND/OR short-circuiting is implemented correctly:
// output = input1 || panic should NOT panic the process when input1 is true.

type PanicNode struct {
}

func BuildPanicNode(astNode *dsl.ASTNode) (*PanicNode, error) {
	return &PanicNode{}, nil
}
func (this *PanicNode) Evaluate(state *State) types.Mlrval {
	lib.InternalCodingErrorPanic("Panic token was evaluated, not short-circuited.")
	return types.MlrvalFromError() // not reached
}
