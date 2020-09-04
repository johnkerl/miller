package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
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

	case dsl.NodeTypeDirectFieldName:
		return BuildSrecDirectFieldReadNode(sval), nil
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

	case dsl.NodeTypePanic:
		return BuildPanicNode(), nil
		break

		// xxx more
		//	case NodeTypeIndirectFieldName:
		//		return lib.MlrvalFromError(), errors.New("unhandled1")
		//		break

	}

	return nil, errors.New("CST builder: unhandled AST leaf node " + string(astNode.Type))
}

// ----------------------------------------------------------------
type SrecDirectFieldReadNode struct {
	fieldName string
}

func BuildSrecDirectFieldReadNode(fieldName string) *SrecDirectFieldReadNode {
	return &SrecDirectFieldReadNode{
		fieldName: fieldName,
	}
}
func (this *SrecDirectFieldReadNode) Evaluate(state *State) lib.Mlrval {
	value := state.Inrec.Get(&this.fieldName)
	if value == nil {
		return lib.MlrvalFromAbsent()
	} else {
		return *value
	}
}

// ----------------------------------------------------------------
type StringLiteralNode struct {
	literal lib.Mlrval
}

func BuildStringLiteralNode(literal string) *StringLiteralNode {
	return &StringLiteralNode{
		literal: lib.MlrvalFromString(literal),
	}
}
func (this *StringLiteralNode) Evaluate(state *State) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type IntLiteralNode struct {
	literal lib.Mlrval
}

func BuildIntLiteralNode(literal string) *IntLiteralNode {
	return &IntLiteralNode{
		literal: lib.MlrvalFromInt64String(literal),
	}
}
func (this *IntLiteralNode) Evaluate(state *State) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type FloatLiteralNode struct {
	literal lib.Mlrval
}

func BuildFloatLiteralNode(literal string) *FloatLiteralNode {
	return &FloatLiteralNode{
		literal: lib.MlrvalFromFloat64String(literal),
	}
}
func (this *FloatLiteralNode) Evaluate(state *State) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type BoolLiteralNode struct {
	literal lib.Mlrval
}

func BuildBoolLiteralNode(literal string) *BoolLiteralNode {
	return &BoolLiteralNode{
		literal: lib.MlrvalFromBoolString(literal),
	}
}
func (this *BoolLiteralNode) Evaluate(state *State) lib.Mlrval {
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

	return nil, errors.New("CST builder: unhandled context variable " + sval)
}

// ----------------------------------------------------------------
type FILENAMENode struct {
}

func BuildFILENAMENode() *FILENAMENode {
	return &FILENAMENode{}
}
func (this *FILENAMENode) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.FILENAME)
}

// ----------------------------------------------------------------
type FILENUMNode struct {
}

func BuildFILENUMNode() *FILENUMNode {
	return &FILENUMNode{}
}
func (this *FILENUMNode) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromInt64(state.Context.FILENUM)
}

// ----------------------------------------------------------------
type NFNode struct {
}

func BuildNFNode() *NFNode {
	return &NFNode{}
}
func (this *NFNode) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromInt64(state.Context.NF)
}

// ----------------------------------------------------------------
type NRNode struct {
}

func BuildNRNode() *NRNode {
	return &NRNode{}
}
func (this *NRNode) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromInt64(state.Context.NR)
}

// ----------------------------------------------------------------
type FNRNode struct {
}

func BuildFNRNode() *FNRNode {
	return &FNRNode{}
}
func (this *FNRNode) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromInt64(state.Context.FNR)
}

// ----------------------------------------------------------------
type IRSNode struct {
}

func BuildIRSNode() *IRSNode {
	return &IRSNode{}
}
func (this *IRSNode) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.IRS)
}

// ----------------------------------------------------------------
type IFSNode struct {
}

func BuildIFSNode() *IFSNode {
	return &IFSNode{}
}
func (this *IFSNode) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.IFS)
}

// ----------------------------------------------------------------
type IPSNode struct {
}

func BuildIPSNode() *IPSNode {
	return &IPSNode{}
}
func (this *IPSNode) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.IPS)
}

// ----------------------------------------------------------------
type ORSNode struct {
}

func BuildORSNode() *ORSNode {
	return &ORSNode{}
}
func (this *ORSNode) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.ORS)
}

// ----------------------------------------------------------------
type OFSNode struct {
}

func BuildOFSNode() *OFSNode {
	return &OFSNode{}
}
func (this *OFSNode) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.OFS)
}

// ----------------------------------------------------------------
type OPSNode struct {
}

func BuildOPSNode() *OPSNode {
	return &OPSNode{}
}
func (this *OPSNode) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.OPS)
}

// ----------------------------------------------------------------
// The panic token is a special token which causes a panic when evaluated.
// This is for testing that AND/OR short-circuiting is implemented correctly:
// output = input1 || panic should NOT panic the process when input1 is true.

type PanicNode struct {
}

func BuildPanicNode() *PanicNode {
	return &PanicNode{}
}
func (this *PanicNode) Evaluate(state *State) lib.Mlrval {
	lib.InternalCodingErrorPanic("Panic token was evaluated, not short-circuited.")
	return lib.MlrvalFromError() // not reached
}
