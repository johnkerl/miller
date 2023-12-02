// ================================================================
// CST build/execute for AST leaf nodes
// ================================================================

package cst

import (
	"fmt"
	"math"

	"github.com/johnkerl/miller/pkg/dsl"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/runtime"
)

// ----------------------------------------------------------------
func (root *RootNode) BuildLeafNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Children != nil)
	sval := string(astNode.Token.Lit)

	switch astNode.Type {

	case dsl.NodeTypeDirectFieldValue:
		return root.BuildDirectFieldRvalueNode(sval), nil
	case dsl.NodeTypeFullSrec:
		return root.BuildFullSrecRvalueNode(), nil

	case dsl.NodeTypeDirectOosvarValue:
		return root.BuildDirectOosvarRvalueNode(sval), nil
	case dsl.NodeTypeFullOosvar:
		return root.BuildFullOosvarRvalueNode(), nil

	case dsl.NodeTypeLocalVariable:
		return root.BuildLocalVariableNode(sval), nil

	case dsl.NodeTypeStringLiteral:
		return root.BuildStringLiteralNode(sval), nil

	case dsl.NodeTypeRegex:
		// During the BNF parse all string literals -- "foo" or "(..)_(...)"
		// regexes etc -- are marked as dsl.NodeTypeStringLiteral. However, a
		// CST-build pre-pass relabels second argument to sub/gsub etc -- any
		// known regex positions -- from dsl.NodeTypeStringLiteral to
		// dsl.NodeTypeRegex. The RegexLiteralNode is responsible for
		// handling backslash sequences for regex literals differently from
		// those for non-regex string literals.
		return root.BuildRegexLiteralNode(sval), nil

	case dsl.NodeTypeRegexCaseInsensitive:
		// StringLiteral nodes like '"abc"' entered by the user come in from
		// the AST as 'abc', with double quotes removed. Case-insensitive
		// regexes like '"a.*b"i' come in with initial '"' and final '"i'
		// intact. We let the sub/regextract/etc functions deal with this.
		// (The alternative would be to make a separate Mlrval type separate
		// from string.)
		return root.BuildRegexLiteralNode(sval), nil

	case dsl.NodeTypeIntLiteral:
		return root.BuildIntLiteralNode(sval), nil
	case dsl.NodeTypeFloatLiteral:
		return root.BuildFloatLiteralNode(sval), nil
	case dsl.NodeTypeBoolLiteral:
		return root.BuildBoolLiteralNode(sval), nil
	case dsl.NodeTypeNullLiteral:
		return root.BuildNullLiteralNode(), nil
	case dsl.NodeTypeContextVariable:
		return root.BuildContextVariableNode(astNode)
	case dsl.NodeTypeConstant:
		return root.BuildConstantNode(astNode)

	case dsl.NodeTypeArraySliceEmptyLowerIndex:
		return root.BuildArraySliceEmptyLowerIndexNode(astNode)
	case dsl.NodeTypeArraySliceEmptyUpperIndex:
		return root.BuildArraySliceEmptyUpperIndexNode(astNode)

	case dsl.NodeTypePanic:
		return root.BuildPanicNode(astNode)
	}

	return nil, fmt.Errorf("at CST BuildLeafNode: unhandled AST node %s", string(astNode.Type))
}

// ----------------------------------------------------------------
type DirectFieldRvalueNode struct {
	fieldName string
}

func (root *RootNode) BuildDirectFieldRvalueNode(fieldName string) *DirectFieldRvalueNode {
	return &DirectFieldRvalueNode{
		fieldName: fieldName,
	}
}
func (node *DirectFieldRvalueNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "$*")
	}
	value := state.Inrec.Get(node.fieldName)
	if value == nil {
		return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "$"+node.fieldName)
	} else {
		return value
	}
}

// ----------------------------------------------------------------
type FullSrecRvalueNode struct {
}

func (root *RootNode) BuildFullSrecRvalueNode() *FullSrecRvalueNode {
	return &FullSrecRvalueNode{}
}
func (node *FullSrecRvalueNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	// For normal DSL use the CST validator will prohibit this from being
	// called in places the current record is undefined (begin and end blocks).
	// However in the REPL people can read past end of stream and still try to
	// print inrec attributes. Also, a UDF/UDS invoked from begin/end could try
	// to access the inrec, and that would get past the validator.
	if state.Inrec == nil {
		return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "$*")
	} else {
		return mlrval.FromMap(state.Inrec)
	}
}

// ----------------------------------------------------------------
type DirectOosvarRvalueNode struct {
	variableName string
}

func (root *RootNode) BuildDirectOosvarRvalueNode(variableName string) *DirectOosvarRvalueNode {
	return &DirectOosvarRvalueNode{
		variableName: variableName,
	}
}
func (node *DirectOosvarRvalueNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	value := state.Oosvars.Get(node.variableName)
	if value == nil {
		return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "@"+node.variableName)
	} else {
		return value
	}
}

// ----------------------------------------------------------------
type FullOosvarRvalueNode struct {
}

func (root *RootNode) BuildFullOosvarRvalueNode() *FullOosvarRvalueNode {
	return &FullOosvarRvalueNode{}
}
func (node *FullOosvarRvalueNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromMap(state.Oosvars)
}

// ----------------------------------------------------------------
// TODO: rename this node-type.
// RHSes can be:
// * Local variables, like `x` in `x = 3`
// * Named user-defined functions, like `f` in `func f(a,b) { return b - a }`
// * Unnamed user-defined functions bound to local variables, like `f` in `f = func (a,b) { return b - a }`
// * Not yet but TODO: built-in functions like `f` in `f = sinh`.

type LocalVariableNode struct {
	stackVariable *runtime.StackVariable
	udfManager    *UDFManager
}

func (root *RootNode) BuildLocalVariableNode(variableName string) *LocalVariableNode {
	return &LocalVariableNode{
		stackVariable: runtime.NewStackVariable(variableName),
		udfManager:    root.udfManager,
	}
}
func (node *LocalVariableNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	value := state.Stack.Get(node.stackVariable)
	if value != nil {
		return value
	}

	functionName := node.stackVariable.GetName()

	udf := node.udfManager.LookUpDisregardingArity(functionName)
	if udf != nil {
		return mlrval.FromFunction(udf, functionName)
	}

	// TODO: allow built-in functions as well. Needs some API-merging as a
	// prerequisite since UDFs and BIFs are managed in quite different
	// structures.

	return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "local variable "+node.stackVariable.GetName())
}

// ----------------------------------------------------------------
// During the BNF parse all string literals -- "foo" or "(..)_(...)" regexes
// etc -- are marked as dsl.NodeTypeStringLiteral. However, a CST-build
// pre-pass relabels second argument to sub/gsub etc -- any known regex
// positions -- from dsl.NodeTypeStringLiteral to dsl.NodeTypeRegex. The
// RegexLiteralNode is responsible for handling backslash sequences for
// regex literals differently from those for non-regex string literals.

type RegexLiteralNode struct {
	literal *mlrval.Mlrval
}

func (root *RootNode) BuildRegexLiteralNode(literal string) IEvaluable {
	return &RegexLiteralNode{
		literal: mlrval.FromString(literal),
	}
}

func (node *RegexLiteralNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.literal
}

// ----------------------------------------------------------------
// StringLiteralNode is for any string literal that doesn't have any "\0" ..
// "\9" in it.
type StringLiteralNode struct {
	literal *mlrval.Mlrval
}

// RegexCaptureReplacementNode is for any string literal that has any "\0" ..
// "\9" in it.  As of the original design of Miller, submatches are captured
// in one place and interpolated in another. For example:
//
//	if ($x =~ "(..)_(...)" {
//	  ... other lines of code ...
//	  $y = "\2:\1";
//	}
//
// This node type is for things like "\2:\1". They can occur quite far from the
// =~ callsite so we need to check all string literals to see if they have "\0"
// .. "\9" anywhere within them. If they do, we precompute a
// replacementCaptureMatrix which is basically compiled information about the
// replacement string -- the start/end offsets of the "\1", "\2", etc
// substrings.
type RegexCaptureReplacementNode struct {
	replacementString        string
	replacementCaptureMatrix [][]int
}

func (root *RootNode) BuildStringLiteralNode(literal string) IEvaluable {
	// Convert "\t" to tab character, "\"" to double-quote character, etc.
	// This is intentionally done for StringLiteralNode but not for
	// RegexLiteralNode.  See also https://github.com/johnkerl/miller/issues/297.
	literal = lib.UnbackslashStringLiteral(literal)

	hasCaptures, replacementCaptureMatrix := lib.RegexReplacementHasCaptures(literal)
	if !hasCaptures {
		return &StringLiteralNode{
			literal: mlrval.FromString(literal),
		}
	} else {
		return &RegexCaptureReplacementNode{
			replacementString:        literal,
			replacementCaptureMatrix: replacementCaptureMatrix,
		}
	}
}

func (node *StringLiteralNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.literal
}

// As noted above, in things like
//
//	if ($x =~ "(..)_(...)" {
//	  ... other lines of code ...
//	  $y = "\2:\1";
//	}
//
// the captures can be set (by =~ or !=~) quite far from where they are used.
// This is why we consult the state.RegexCaptures here, to see if they've been
// set on some previous invocation of =~ or !=~.
func (node *RegexCaptureReplacementNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromString(
		lib.InterpolateCaptures(
			node.replacementString,
			node.replacementCaptureMatrix,
			state.RegexCaptures,
		),
	)
}

// ----------------------------------------------------------------
type IntLiteralNode struct {
	literal *mlrval.Mlrval
}

func (root *RootNode) BuildIntLiteralNode(literal string) *IntLiteralNode {
	ival, ok := lib.TryIntFromString(literal)
	lib.InternalCodingErrorIf(!ok)
	return &IntLiteralNode{
		literal: mlrval.FromPrevalidatedIntString(literal, ival),
	}
}
func (node *IntLiteralNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.literal
}

// ----------------------------------------------------------------
type FloatLiteralNode struct {
	literal *mlrval.Mlrval
}

func (root *RootNode) BuildFloatLiteralNode(literal string) *FloatLiteralNode {
	fval, ok := lib.TryFloatFromString(literal)
	lib.InternalCodingErrorIf(!ok)
	return &FloatLiteralNode{
		literal: mlrval.FromPrevalidatedFloatString(literal, fval),
	}
}
func (node *FloatLiteralNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.literal
}

// ----------------------------------------------------------------
type BoolLiteralNode struct {
	literal *mlrval.Mlrval
}

func (root *RootNode) BuildBoolLiteralNode(literal string) *BoolLiteralNode {
	return &BoolLiteralNode{
		literal: mlrval.FromBoolString(literal),
	}
}
func (node *BoolLiteralNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.literal
}

// ----------------------------------------------------------------
type NullLiteralNode struct {
	literal *mlrval.Mlrval
}

func (root *RootNode) BuildNullLiteralNode() *NullLiteralNode {
	return &NullLiteralNode{
		literal: mlrval.NULL,
	}
}
func (node *NullLiteralNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.literal
}

// ----------------------------------------------------------------
// Used for testing purposes; not used by the main DSL.

type MlrvalLiteralNode struct {
	literal *mlrval.Mlrval
}

func BuildMlrvalLiteralNode(literal *mlrval.Mlrval) *MlrvalLiteralNode {
	return &MlrvalLiteralNode{
		literal: literal.Copy(),
	}
}
func (node *MlrvalLiteralNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.literal
}

// ================================================================
func (root *RootNode) BuildContextVariableNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Token == nil)
	sval := string(astNode.Token.Lit)

	switch sval {

	case "FILENAME":
		return root.BuildFILENAMENode(), nil
	case "FILENUM":
		return root.BuildFILENUMNode(), nil

	case "NF":
		return root.BuildNFNode(), nil
	case "NR":
		return root.BuildNRNode(), nil
	case "FNR":
		return root.BuildFNRNode(), nil

	case "IRS":
		return root.BuildIRSNode(), nil
	case "IFS":
		return root.BuildIFSNode(), nil
	case "IPS":
		return root.BuildIPSNode(), nil

	case "ORS":
		return root.BuildORSNode(), nil
	case "OFS":
		return root.BuildOFSNode(), nil
	case "OPS":
		return root.BuildOPSNode(), nil
	case "FLATSEP":
		return root.BuildFLATSEPNode(), nil

	}

	return nil, fmt.Errorf("at CST BuildContextVariableNode: unhandled context variable %s", sval)
}

// ----------------------------------------------------------------
type FILENAMENode struct {
}

func (root *RootNode) BuildFILENAMENode() *FILENAMENode {
	return &FILENAMENode{}
}
func (node *FILENAMENode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromString(state.Context.FILENAME)
}

// ----------------------------------------------------------------
type FILENUMNode struct {
}

func (root *RootNode) BuildFILENUMNode() *FILENUMNode {
	return &FILENUMNode{}
}
func (node *FILENUMNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromInt(state.Context.FILENUM)
}

// ----------------------------------------------------------------
type NFNode struct {
}

func (root *RootNode) BuildNFNode() *NFNode {
	return &NFNode{}
}
func (node *NFNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromInt(state.Inrec.FieldCount)
}

// ----------------------------------------------------------------
type NRNode struct {
}

func (root *RootNode) BuildNRNode() *NRNode {
	return &NRNode{}
}
func (node *NRNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromInt(state.Context.NR)
}

// ----------------------------------------------------------------
type FNRNode struct {
}

func (root *RootNode) BuildFNRNode() *FNRNode {
	return &FNRNode{}
}
func (node *FNRNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromInt(state.Context.FNR)
}

// ----------------------------------------------------------------
type IRSNode struct {
}

func (root *RootNode) BuildIRSNode() *IRSNode {
	return &IRSNode{}
}
func (node *IRSNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromString(state.Options.ReaderOptions.IRS)
}

// ----------------------------------------------------------------
type IFSNode struct {
}

func (root *RootNode) BuildIFSNode() *IFSNode {
	return &IFSNode{}
}
func (node *IFSNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromString(state.Options.ReaderOptions.IFS)
}

// ----------------------------------------------------------------
type IPSNode struct {
}

func (root *RootNode) BuildIPSNode() *IPSNode {
	return &IPSNode{}
}
func (node *IPSNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromString(state.Options.ReaderOptions.IPS)
}

// ----------------------------------------------------------------
type ORSNode struct {
}

func (root *RootNode) BuildORSNode() *ORSNode {
	return &ORSNode{}
}
func (node *ORSNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromString(state.Options.WriterOptions.ORS)
}

// ----------------------------------------------------------------
type OFSNode struct {
}

func (root *RootNode) BuildOFSNode() *OFSNode {
	return &OFSNode{}
}
func (node *OFSNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromString(state.Options.WriterOptions.OFS)
}

// ----------------------------------------------------------------
type OPSNode struct {
}

func (root *RootNode) BuildOPSNode() *OPSNode {
	return &OPSNode{}
}
func (node *OPSNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromString(state.Options.WriterOptions.OPS)
}

// ----------------------------------------------------------------
type FLATSEPNode struct {
}

func (root *RootNode) BuildFLATSEPNode() *FLATSEPNode {
	return &FLATSEPNode{}
}
func (node *FLATSEPNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromString(state.Options.WriterOptions.FLATSEP)
}

// ================================================================
func (root *RootNode) BuildConstantNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Token == nil)
	sval := string(astNode.Token.Lit)

	switch sval {

	case "M_PI":
		return root.BuildMathPINode(), nil
	case "M_E":
		return root.BuildMathENode(), nil

	}

	return nil, fmt.Errorf("at CST BuildContextVariableNode: unhandled context variable %s", sval)
}

// ----------------------------------------------------------------
type MathPINode struct {
}

func (root *RootNode) BuildMathPINode() *MathPINode {
	return &MathPINode{}
}
func (node *MathPINode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromFloat(math.Pi)
}

// ----------------------------------------------------------------
type MathENode struct {
}

func (root *RootNode) BuildMathENode() *MathENode {
	return &MathENode{}
}
func (node *MathENode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromFloat(math.E)
}

// ================================================================
type LiteralOneNode struct {
}

// In array slices like 'myarray[:4]', the lower index is always 1 since Miller
// user-space indices are 1-up.
func (root *RootNode) BuildArraySliceEmptyLowerIndexNode(
	astNode *dsl.ASTNode,
) (*LiteralOneNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArraySliceEmptyLowerIndex)
	return &LiteralOneNode{}, nil
}
func (node *LiteralOneNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromInt(1)
}

// ================================================================
type LiteralEmptyStringNode struct {
}

// In array slices like 'myarray[4:]', the upper index is always n, where n is
// the length of the array, since Miller user-space indices are 1-up. However,
// we don't have access to the array length in this AST node so we return ""
// so the slice-index CST node can compute it.
func (root *RootNode) BuildArraySliceEmptyUpperIndexNode(
	astNode *dsl.ASTNode,
) (*LiteralEmptyStringNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArraySliceEmptyUpperIndex)
	return &LiteralEmptyStringNode{}, nil
}
func (node *LiteralEmptyStringNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return mlrval.FromString("")
}

// ----------------------------------------------------------------
// The panic token is a special token which causes a panic when evaluated.
// This is for testing that AND/OR short-circuiting is implemented correctly:
// output = input1 || panic should NOT panic the process when input1 is true.

type PanicNode struct {
}

func (root *RootNode) BuildPanicNode(astNode *dsl.ASTNode) (*PanicNode, error) {
	return &PanicNode{}, nil
}
func (node *PanicNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	lib.InternalCodingErrorPanic("Panic token was evaluated, not short-circuited.")
	return nil // not reached
}
