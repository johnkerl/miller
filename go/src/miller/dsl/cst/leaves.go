package cst

import (
	"miller/containers"
	"miller/lib"
)

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

func NewIntLiteral(literal int64) *IntLiteral {
	return &IntLiteral{
		literal: lib.MlrvalFromInt64(literal),
	}
}
func (this *IntLiteral) Evaluate(state *State) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type FloatLiteral struct {
	literal lib.Mlrval
}

func NewFloatLiteral(literal float64) *FloatLiteral {
	return &FloatLiteral{
		literal: lib.MlrvalFromFloat64(literal),
	}
}
func (this *FloatLiteral) Evaluate(state *State) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type BoolLiteral struct {
	literal lib.Mlrval
}

func NewBoolLiteral(literal bool) *BoolLiteral {
	return &BoolLiteral{
		literal: lib.MlrvalFromBool(literal),
	}
}
func (this *BoolLiteral) Evaluate(state *State) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type FILENAME struct {
}

func NewFILENAME() *FILENAME {
	return &FILENAME{}
}
func (this *BoolLiteral) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromString(state.Context.FILENAME)
}

// ----------------------------------------------------------------
type FILENUM struct {
}

func NewFILENUM() *FILENUM {
	return &FILENUM{}
}
func (this *BoolLiteral) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromInt64(state.Context.FILENUM)
}

// ----------------------------------------------------------------
type NF struct {
}

func NewNF() *NF {
	return &NF{}
}
func (this *BoolLiteral) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromInt64(state.Context.NF)
}

// ----------------------------------------------------------------
type NR struct {
}

func NewNR() *NR {
	return &NR{}
}
func (this *BoolLiteral) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromInt64(state.Context.NR)
}

// ----------------------------------------------------------------
type FNR struct {
}

func NewFNR() *FNR {
	return &FNR{}
}
func (this *BoolLiteral) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromInt64(state.Context.FNR)
}
