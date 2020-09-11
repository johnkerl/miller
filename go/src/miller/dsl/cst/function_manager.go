package cst

import (
	"fmt"
	"miller/lib"
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
	zaryFunc           lib.ZaryFunc
	unaryFunc          lib.UnaryFunc
	binaryFunc         lib.BinaryFunc
	ternaryFunc        lib.TernaryFunc
	variadicFunc       lib.VariadicFunc
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
		zaryFunc: lib.MlrvalSystime,
	},
	{
		name:     "systimeint",
		help:     "help string will go here",
		zaryFunc: lib.MlrvalSystimeInt,
	},
	{
		name:     "urand",
		zaryFunc: lib.MlrvalUrand,
	},
	{
		name:     "urand32",
		zaryFunc: lib.MlrvalUrand32,
	},

	// ----------------------------------------------------------------
	// Multiple-arity built-in functions
	{
		name:               "+",
		unaryFunc:          lib.MlrvalUnaryPlus,
		binaryFunc:         lib.MlrvalBinaryPlus,
		hasMultipleArities: true,
	},
	{
		name:               "-",
		unaryFunc:          lib.MlrvalUnaryMinus,
		binaryFunc:         lib.MlrvalBinaryMinus,
		hasMultipleArities: true,
	},

	// ----------------------------------------------------------------
	// Unary built-in functions
	{
		name:      "~",
		unaryFunc: lib.MlrvalBitwiseNOT,
	},
	{
		name:      "!",
		unaryFunc: lib.MlrvalLogicalNOT,
	},

	{
		name:      "abs",
		help:      "Absolute value.",
		unaryFunc: lib.MlrvalAbs,
	},
	{
		name:      "acos",
		help:      "Inverse trigonometric cosine.",
		unaryFunc: lib.MlrvalAcos,
	},
	{
		name:      "acosh",
		help:      "Inverse hyperbolic cosine.",
		unaryFunc: lib.MlrvalAcosh,
	},
	{
		name:      "asin",
		help:      "Inverse trigonometric sine.",
		unaryFunc: lib.MlrvalAsin,
	},
	{
		name:      "asinh",
		help:      "Inverse hyperbolic sine.",
		unaryFunc: lib.MlrvalAsinh,
	},
	{
		name:      "atan",
		help:      "One-argument arctangent.",
		unaryFunc: lib.MlrvalAtan,
	},
	{
		name:      "atanh",
		help:      "Inverse hyperbolic tangent.",
		unaryFunc: lib.MlrvalAtanh,
	},
	{
		name:      "cbrt",
		help:      "Cube root.",
		unaryFunc: lib.MlrvalCbrt,
	},
	{
		name:      "ceil",
		help:      "Ceiling: nearest integer at or above.",
		unaryFunc: lib.MlrvalCeil,
	},
	{
		name:      "cos",
		help:      "Trigonometric cosine.",
		unaryFunc: lib.MlrvalCos,
	},
	{
		name:      "cosh",
		help:      "Hyperbolic cosine.",
		unaryFunc: lib.MlrvalCosh,
	},
	{
		name:      "erf",
		help:      "Error function.",
		unaryFunc: lib.MlrvalErf,
	},
	{
		name:      "erfc",
		help:      "Complementary error function.",
		unaryFunc: lib.MlrvalErfc,
	},
	{
		name:      "exp",
		help:      "Exponential function e**x.",
		unaryFunc: lib.MlrvalExp,
	},
	{
		name:      "expm1",
		help:      "e**x - 1.",
		unaryFunc: lib.MlrvalExpm1,
	},
	{
		name:      "floor",
		help:      "Floor: nearest integer at or below.",
		unaryFunc: lib.MlrvalFloor,
	},
	{
		name:      "log",
		help:      "Natural (base-e) logarithm.",
		unaryFunc: lib.MlrvalLog,
	},
	{
		name:      "log10",
		help:      "Base-10 logarithm.",
		unaryFunc: lib.MlrvalLog10,
	},
	{
		name:      "log1p",
		help:      "log(1-x).",
		unaryFunc: lib.MlrvalLog1p,
	},
	{
		name:      "round",
		help:      "Round to nearest integer.",
		unaryFunc: lib.MlrvalRound,
	},
	{
		name:      "sin",
		help:      "Trigonometric sine.",
		unaryFunc: lib.MlrvalSin,
	},
	{
		name:      "sinh",
		help:      "Hyperbolic sine.",
		unaryFunc: lib.MlrvalSinh,
	},
	{
		name:      "sqrt",
		help:      "Square root.",
		unaryFunc: lib.MlrvalSqrt,
	},
	{
		name:      "tan",
		help:      "Trigonometric tangent.",
		unaryFunc: lib.MlrvalTan,
	},
	{
		name:      "tanh",
		help:      "Hyperbolic tangent.",
		unaryFunc: lib.MlrvalTanh,
	},

	// ----------------------------------------------------------------
	// Binary built-in functions
	{
		name:       ".",
		binaryFunc: lib.MlrvalDot,
	},
	{
		name:       "*",
		binaryFunc: lib.MlrvalTimes,
	},
	{
		name:       "/",
		binaryFunc: lib.MlrvalDivide,
	},
	{
		name:       "//",
		binaryFunc: lib.MlrvalIntDivide,
	},
	{
		name:       "**",
		binaryFunc: lib.MlrvalPow,
	},
	{
		name:       ".+",
		binaryFunc: lib.MlrvalDotPlus,
	},
	{
		name:       ".-",
		binaryFunc: lib.MlrvalDotMinus,
	},
	{
		name:       ".*",
		binaryFunc: lib.MlrvalDotTimes,
	},
	{
		name:       "./",
		binaryFunc: lib.MlrvalDotDivide,
	},
	{
		name:       "%",
		binaryFunc: lib.MlrvalModulus,
	},

	{
		name:       "==",
		binaryFunc: lib.MlrvalEquals,
	},
	{
		name:       "!=",
		binaryFunc: lib.MlrvalNotEquals,
	},
	{
		name:       ">",
		binaryFunc: lib.MlrvalGreaterThan,
	},
	{
		name:       ">=",
		binaryFunc: lib.MlrvalGreaterThanOrEquals,
	},
	{
		name:       "<",
		binaryFunc: lib.MlrvalLessThan,
	},
	{
		name:       "<=",
		binaryFunc: lib.MlrvalLessThanOrEquals,
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
		binaryFunc: lib.MlrvalLogicalXOR,
	},
	{
		name:       "&",
		binaryFunc: lib.MlrvalBitwiseAND,
	},
	{
		name:       "|",
		binaryFunc: lib.MlrvalBitwiseOR,
	},
	{
		name:       "^",
		binaryFunc: lib.MlrvalBitwiseXOR,
	},
	{
		name:       "<<",
		binaryFunc: lib.MlrvalLeftShift,
	},
	{
		name:       ">>",
		binaryFunc: lib.MlrvalSignedRightShift,
	},
	{
		name:       ">>>",
		binaryFunc: lib.MlrvalUnsignedRightShift,
	},

	{
		name: "urandint",
	},
	{
		name: "urandrange",
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

	// Variadic built-in functions
	{
		name:         "max",
		variadicFunc: lib.MlrvalVariadicMax,
	},
	{
		name:         "min",
		variadicFunc: lib.MlrvalVariadicMin,
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

// ================================================================
// Standard singleton. UDFs are still to come. :)
var BuiltinFunctionManager *FunctionManager = NewFunctionManager()
