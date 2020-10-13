package cst

import (
	"fmt"
	"miller/types"
	"os"
)

// ================================================================
// Adding a new builtin function:
// * New entry in BUILTIN_FUNCTION_LOOKUP_TABLE
// * Implement the function in mlrval_functions.go
// ================================================================

// ================================================================
type FunctionInfo struct {
	name string
	// class      string -- "math", "time", "typing", "maps", etc
	help               string
	hasMultipleArities bool
	zaryFunc           types.ZaryFunc
	unaryFunc          types.UnaryFunc
	binaryFunc         types.BinaryFunc
	ternaryFunc        types.TernaryFunc
	variadicFunc       types.VariadicFunc
}

//// ----------------------------------------------------------------
//typedef enum _func_class_t {
//	FUNC_CLASS_ARITHMETIC,
//	FUNC_CLASS_MATH,
//	FUNC_CLASS_BOOLEAN,
//	FUNC_CLASS_STRING,
//	FUNC_CLASS_CONVERSION,
//	FUNC_CLASS_TYPING,
//	FUNC_CLASS_MAPS,
//	FUNC_CLASS_TIME
//} func_class_t;

// ================================================================
var BUILTIN_FUNCTION_LOOKUP_TABLE = []FunctionInfo{

	// ----------------------------------------------------------------
	// Zary built-in functions
	{
		name:     "systime",
		help:     "help string will go here",
		zaryFunc: types.MlrvalSystime,
	},
	{
		name:     "systimeint",
		help:     "help string will go here",
		zaryFunc: types.MlrvalSystimeInt,
	},
	{
		name:     "urand",
		zaryFunc: types.MlrvalUrand,
	},
	{
		name:     "urand32",
		zaryFunc: types.MlrvalUrand32,
	},

	// ----------------------------------------------------------------
	// Multiple-arity built-in functions
	{
		name:               "+",
		unaryFunc:          types.MlrvalUnaryPlus,
		binaryFunc:         types.MlrvalBinaryPlus,
		hasMultipleArities: true,
	},
	{
		name:               "-",
		unaryFunc:          types.MlrvalUnaryMinus,
		binaryFunc:         types.MlrvalBinaryMinus,
		hasMultipleArities: true,
	},
	{
		name:      "sec2gmt",
		help:      `Formats seconds since epoch (integer part)
as GMT timestamp, e.g. sec2gmt(1440768801.7) = "2015-08-28T13:33:21Z".
Leaves non-numbers as-is.`,
		unaryFunc: types.MlrvalSec2GMTUnary,
		binaryFunc: types.MlrvalSec2GMTBinary,
		hasMultipleArities: true,
	},

	// ----------------------------------------------------------------
	// Unary built-in functions
	{
		name:      "~",
		unaryFunc: types.MlrvalBitwiseNOT,
	},
	{
		name:      "!",
		unaryFunc: types.MlrvalLogicalNOT,
	},

	{
		name:      "abs",
		help:      "Absolute value.",
		unaryFunc: types.MlrvalAbs,
	},
	{
		name:      "acos",
		help:      "Inverse trigonometric cosine.",
		unaryFunc: types.MlrvalAcos,
	},
	{
		name:      "acosh",
		help:      "Inverse hyperbolic cosine.",
		unaryFunc: types.MlrvalAcosh,
	},
	{
		name:      "asin",
		help:      "Inverse trigonometric sine.",
		unaryFunc: types.MlrvalAsin,
	},
	{
		name:      "asinh",
		help:      "Inverse hyperbolic sine.",
		unaryFunc: types.MlrvalAsinh,
	},
	{
		name:      "atan",
		help:      "One-argument arctangent.",
		unaryFunc: types.MlrvalAtan,
	},
	{
		name:      "atanh",
		help:      "Inverse hyperbolic tangent.",
		unaryFunc: types.MlrvalAtanh,
	},
	{
		name:      "cbrt",
		help:      "Cube root.",
		unaryFunc: types.MlrvalCbrt,
	},
	{
		name:      "ceil",
		help:      "Ceiling: nearest integer at or above.",
		unaryFunc: types.MlrvalCeil,
	},
	{
		name:      "cos",
		help:      "Trigonometric cosine.",
		unaryFunc: types.MlrvalCos,
	},
	{
		name:      "cosh",
		help:      "Hyperbolic cosine.",
		unaryFunc: types.MlrvalCosh,
	},
	{
		name:      "erf",
		help:      "Error function.",
		unaryFunc: types.MlrvalErf,
	},
	{
		name:      "erfc",
		help:      "Complementary error function.",
		unaryFunc: types.MlrvalErfc,
	},
	{
		name:      "exp",
		help:      "Exponential function e**x.",
		unaryFunc: types.MlrvalExp,
	},
	{
		name:      "expm1",
		help:      "e**x - 1.",
		unaryFunc: types.MlrvalExpm1,
	},
	{
		name:      "floor",
		help:      "Floor: nearest integer at or below.",
		unaryFunc: types.MlrvalFloor,
	},
	{
		name:      "log",
		help:      "Natural (base-e) logarithm.",
		unaryFunc: types.MlrvalLog,
	},
	{
		name:      "log10",
		help:      "Base-10 logarithm.",
		unaryFunc: types.MlrvalLog10,
	},
	{
		name:      "log1p",
		help:      "log(1-x).",
		unaryFunc: types.MlrvalLog1p,
	},
	{
		name:      "round",
		help:      "Round to nearest integer.",
		unaryFunc: types.MlrvalRound,
	},
	{
		name:      "sin",
		help:      "Trigonometric sine.",
		unaryFunc: types.MlrvalSin,
	},
	{
		name:      "sinh",
		help:      "Hyperbolic sine.",
		unaryFunc: types.MlrvalSinh,
	},
	{
		name:      "sqrt",
		help:      "Square root.",
		unaryFunc: types.MlrvalSqrt,
	},
	{
		name:      "tan",
		help:      "Trigonometric tangent.",
		unaryFunc: types.MlrvalTan,
	},
	{
		name:      "tanh",
		help:      "Hyperbolic tangent.",
		unaryFunc: types.MlrvalTanh,
	},
	{
		name:      "clean_whitespace",
		help:      "Same as collapse_whitespace and strip.",
		unaryFunc: types.MlrvalCleanWhitespace,
	},
	{
		name:      "collapse_whitespace",
		help:      "Strip repeated whitespace from string.",
		unaryFunc: types.MlrvalCollapseWhitespace,
	},
	{
		name:      "lstrip",
		help:      "Strip leading whitespace from string.",
		unaryFunc: types.MlrvalLStrip,
	},
	{
		name:      "rstrip",
		help:      "Strip trailing whitespace from string.",
		unaryFunc: types.MlrvalRStrip,
	},
	{
		name:      "strip",
		help:      "Strip leading and trailing whitespace from string.",
		unaryFunc: types.MlrvalStrip,
	},

	// ----------------------------------------------------------------
	// Binary built-in functions
	{
		name:       ".",
		binaryFunc: types.MlrvalDot,
	},
	{
		name:       "*",
		binaryFunc: types.MlrvalTimes,
	},
	{
		name:       "/",
		binaryFunc: types.MlrvalDivide,
	},
	{
		name:       "//",
		binaryFunc: types.MlrvalIntDivide,
	},
	{
		name:       "**",
		binaryFunc: types.MlrvalPow,
	},
	{
		name:       ".+",
		binaryFunc: types.MlrvalDotPlus,
	},
	{
		name:       ".-",
		binaryFunc: types.MlrvalDotMinus,
	},
	{
		name:       ".*",
		binaryFunc: types.MlrvalDotTimes,
	},
	{
		name:       "./",
		binaryFunc: types.MlrvalDotDivide,
	},
	{
		name:       "%",
		binaryFunc: types.MlrvalModulus,
	},

	{
		name:       "==",
		binaryFunc: types.MlrvalEquals,
	},
	{
		name:       "!=",
		binaryFunc: types.MlrvalNotEquals,
	},
	{
		name:       ">",
		binaryFunc: types.MlrvalGreaterThan,
	},
	{
		name:       ">=",
		binaryFunc: types.MlrvalGreaterThanOrEquals,
	},
	{
		name:       "<",
		binaryFunc: types.MlrvalLessThan,
	},
	{
		name:       "<=",
		binaryFunc: types.MlrvalLessThanOrEquals,
	},

	{
		name:       "&&",
		binaryFunc: BinaryShortCircuitPlaceholder,
	},
	{
		name:       "||",
		binaryFunc: BinaryShortCircuitPlaceholder,
	},
	{
		name:       "^&",
		binaryFunc: types.MlrvalLogicalXOR,
	},
	{
		name:       "&",
		binaryFunc: types.MlrvalBitwiseAND,
	},
	{
		name:       "|",
		binaryFunc: types.MlrvalBitwiseOR,
	},
	{
		name:       "^",
		binaryFunc: types.MlrvalBitwiseXOR,
	},
	{
		name:       "<<",
		binaryFunc: types.MlrvalLeftShift,
	},
	{
		name:       ">>",
		binaryFunc: types.MlrvalSignedRightShift,
	},
	{
		name:       ">>>",
		binaryFunc: types.MlrvalUnsignedRightShift,
	},

	{
		name: "urandint",
	},
	{
		name: "urandrange",
	},
	{
		name:       "truncate",
		binaryFunc: types.MlrvalTruncate,
	},

	//pow (class=math #args=2): Exponentiation; same as **.
	//roundm (class=math #args=2): Round to nearest multiple of m: roundm($x,$m) is
	//urandrange (class=math #args=2): Floating-point numbers uniformly distributed on the interval [a, b).
	//urandint (class=math #args=2): Integer uniformly distributed between inclusive
	//atan2 (class=math #args=2): Two-argument arctangent.

	// Ternary built-in functions
	//logifit (class=math #args=3): Given m and b from logistic regression, compute
	//madd (class=math #args=3): a + b mod m (integers)
	//mexp (class=math #args=3): a ** b mod m (integers)
	//mmul (class=math #args=3): a * b mod m (integers)
	//msub (class=math #args=3): a - b mod m (integers)
	{
		name:        "?:",
		ternaryFunc: TernaryShortCircuitPlaceholder,
	},
	{
		name:        "ssub",
		ternaryFunc: types.MlrvalSsub,
	},
	{
		name:        "gsub",
		ternaryFunc: types.MlrvalGsub,
	},

	// Variadic built-in functions
	{
		name:         "max",
		variadicFunc: types.MlrvalVariadicMax,
	},
	{
		name:         "min",
		variadicFunc: types.MlrvalVariadicMin,
	},
}

// ================================================================
type FunctionManager struct {
	// We need both the array and the hashmap since Go maps are not
	// insertion-order-preserving: to produce a sensical help-all-functions
	// list, etc., we want the original ordering.
	lookupTable *[]FunctionInfo
	hashTable   map[string]*FunctionInfo
}

func NewFunctionManager() *FunctionManager {
	// TODO: temp -- one big one -- pending UDFs
	lookupTable := &BUILTIN_FUNCTION_LOOKUP_TABLE
	hashTable := hashifyLookupTable(lookupTable)
	return &FunctionManager{
		lookupTable: lookupTable,
		hashTable:   hashTable,
	}
}

func (this *FunctionManager) LookUp(functionName string) *FunctionInfo {
	return this.hashTable[functionName]
}

func hashifyLookupTable(lookupTable *[]FunctionInfo) map[string]*FunctionInfo {
	hashTable := make(map[string]*FunctionInfo)
	for _, functionInfo := range *lookupTable {
		// Each function name should appear only once in the table.  If it has
		// multiple arities (e.g. unary and binary "-") there should be
		// multiple function-pointers in a single row.
		if hashTable[functionInfo.name] != nil {
			fmt.Fprintf(
				os.Stderr,
				"Internal coding error: function name \"%s\" is non-unique",
				functionInfo.name,
			)
			os.Exit(1)
		}
		clone := functionInfo
		hashTable[functionInfo.name] = &clone
	}
	return hashTable
}

// ----------------------------------------------------------------
func (this *FunctionManager) ListBuiltinFunctionsRaw(o *os.File) {
	for _, functionInfo := range *this.lookupTable {
		fmt.Fprintln(o, functionInfo.name)
	}
}

// ----------------------------------------------------------------
func (this *FunctionManager) ListBuiltinFunctionUsages(o *os.File) {
	for _, functionInfo := range *this.lookupTable {
		fmt.Fprintf(o, "%-20s  %s\n", functionInfo.name, functionInfo.help)
	}
}

// ================================================================
// Standard singleton. UDFs are still to come. :)
var BuiltinFunctionManager *FunctionManager = NewFunctionManager()
