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
func NewEvaluableLeafNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Children != nil)
	sval := string(astNode.Token.Lit)

	switch astNode.Type {

	case dsl.NodeTypeDirectFieldName:
		return NewSrecDirectFieldRead(sval), nil
		break

	case dsl.NodeTypeStringLiteral:
		return NewStringLiteral(sval), nil
		break
	case dsl.NodeTypeIntLiteral:
		return NewIntLiteral(sval), nil
		break
	case dsl.NodeTypeFloatLiteral:
		return NewFloatLiteral(sval), nil
		break
	case dsl.NodeTypeBoolLiteral:
		return NewBoolLiteral(sval), nil
		break
	case dsl.NodeTypeContextVariable:
		return NewContextVariable(astNode)
		break

	case dsl.NodeTypePanic:
		return NewPanic(), nil
		break

		// xxx more
		//	case NodeTypeIndirectFieldName:
		//		return lib.MlrvalFromError(), errors.New("unhandled1")
		//		break

	}

	return nil, errors.New("CST builder: unhandled AST leaf node " + string(astNode.Type))
}

// ----------------------------------------------------------------
type SrecDirectFieldRead struct {
	fieldName string
}

func NewSrecDirectFieldRead(fieldName string) *SrecDirectFieldRead {
	return &SrecDirectFieldRead{
		fieldName: fieldName,
	}
}
func (this *SrecDirectFieldRead) Evaluate(state *State) lib.Mlrval {
	value := state.Inrec.Get(&this.fieldName)
	if value == nil {
		return lib.MlrvalFromAbsent()
	} else {
		return *value
	}
}

// ----------------------------------------------------------------
type StringLiteral struct {
	literal lib.Mlrval
}

func NewStringLiteral(literal string) *StringLiteral {
	return &StringLiteral{
		literal: lib.MlrvalFromString(literal),
	}
}
func (this *StringLiteral) Evaluate(state *State) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type IntLiteral struct {
	literal lib.Mlrval
}

func NewIntLiteral(literal string) *IntLiteral {
	return &IntLiteral{
		literal: lib.MlrvalFromInt64String(literal),
	}
}
func (this *IntLiteral) Evaluate(state *State) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type FloatLiteral struct {
	literal lib.Mlrval
}

func NewFloatLiteral(literal string) *FloatLiteral {
	return &FloatLiteral{
		literal: lib.MlrvalFromFloat64String(literal),
	}
}
func (this *FloatLiteral) Evaluate(state *State) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type BoolLiteral struct {
	literal lib.Mlrval
}

func NewBoolLiteral(literal string) *BoolLiteral {
	return &BoolLiteral{
		literal: lib.MlrvalFromBoolString(literal),
	}
}
func (this *BoolLiteral) Evaluate(state *State) lib.Mlrval {
	return this.literal
}

// ================================================================
func NewContextVariable(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Token == nil)
	sval := string(astNode.Token.Lit)

	switch sval {

	case "FILENAME":
		return NewFILENAME(), nil
		break
	case "FILENUM":
		return NewFILENUM(), nil
		break

	case "NF":
		return NewNF(), nil
		break
	case "NR":
		return NewNR(), nil
		break
	case "FNR":
		return NewFNR(), nil
		break

	case "IRS":
		return NewIRS(), nil
		break
	case "IFS":
		return NewIFS(), nil
		break
	case "IPS":
		return NewIPS(), nil
		break

	case "ORS":
		return NewORS(), nil
		break
	case "OFS":
		return NewOFS(), nil
		break
	case "OPS":
		return NewOPS(), nil
		break

	}

	return nil, errors.New("CST builder: unhandled context variable " + sval)
}

// ----------------------------------------------------------------
type FILENAME struct {
}

func NewFILENAME() *FILENAME {
	return &FILENAME{}
}
func (this *FILENAME) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.FILENAME)
}

// ----------------------------------------------------------------
type FILENUM struct {
}

func NewFILENUM() *FILENUM {
	return &FILENUM{}
}
func (this *FILENUM) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromInt64(state.Context.FILENUM)
}

// ----------------------------------------------------------------
type NF struct {
}

func NewNF() *NF {
	return &NF{}
}
func (this *NF) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromInt64(state.Context.NF)
}

// ----------------------------------------------------------------
type NR struct {
}

func NewNR() *NR {
	return &NR{}
}
func (this *NR) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromInt64(state.Context.NR)
}

// ----------------------------------------------------------------
type FNR struct {
}

func NewFNR() *FNR {
	return &FNR{}
}
func (this *FNR) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromInt64(state.Context.FNR)
}

// ----------------------------------------------------------------
type IRS struct {
}

func NewIRS() *IRS {
	return &IRS{}
}
func (this *IRS) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.IRS)
}

// ----------------------------------------------------------------
type IFS struct {
}

func NewIFS() *IFS {
	return &IFS{}
}
func (this *IFS) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.IFS)
}

// ----------------------------------------------------------------
type IPS struct {
}

func NewIPS() *IPS {
	return &IPS{}
}
func (this *IPS) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.IPS)
}

// ----------------------------------------------------------------
type ORS struct {
}

func NewORS() *ORS {
	return &ORS{}
}
func (this *ORS) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.ORS)
}

// ----------------------------------------------------------------
type OFS struct {
}

func NewOFS() *OFS {
	return &OFS{}
}
func (this *OFS) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.OFS)
}

// ----------------------------------------------------------------
type OPS struct {
}

func NewOPS() *OPS {
	return &OPS{}
}
func (this *OPS) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.OPS)
}

// ----------------------------------------------------------------
// The panic token is a special token which causes a panic when evaluated.
// This is for testing that AND/OR short-circuiting is implemented correctly:
// output = input1 || panic should NOT panic the process when input1 is true.

type Panic struct {
}

func NewPanic() *Panic {
	return &Panic{}
}
func (this *Panic) Evaluate(state *State) lib.Mlrval {
	lib.InternalCodingErrorPanic("Panic token was evaluated, not short-circuited.")
	return lib.MlrvalFromError() // not reached
}
