// ================================================================
// The `types.Mlrval` structure includes **string, int, float, boolean, void,
// absent, and error** types (not unlike PHP's `zval`) as well as
// type-conversion logic for various operators.
//
// Whenever I say "int" and "float" with regard to mlrvals I always mean
// "int" and "float64". If I ever miss a spot and use Go int/float types then
// that is a bug. It would be great to be able to somehow lint for this.
// ================================================================

package types

type Mlrval struct {

	// Enumeration for string / int / float / boolean / etc.
	// I would call this "type" not "mvtype" but "type" is a keyword in Go.
	mvtype MVType

	// An int/float always starts from a string -- be it from record data from
	// a file, or a literal within a DSL expression. The printrep is exactly
	// that string, however the user formatted it, and the intval/floatval is
	// computed from that -- and in sync with it -- at construction time.
	//
	// When a mlrval is computed from one or more others -- e.g. '$z = $x + 4'
	// -- the printrep is not updated. That would be wasted CPU, since the
	// string representation is not needed until when/if the value is printed
	// as output. For computation methods the printrep is neglected and the
	// printrepValid is set to false.
	//
	// In the String() method the printrep is computed from the intval/floatval
	// and printrepValid is set back to true.
	//
	// Note that for MT_STRING, the printrep is always valid since it is the
	// only payload for the mlrval.
	//
	// Thus we (a) keep user-specific input-formatting when possible, for the
	// principle of least surprise; (b) avoid reformatting strings during
	// intermediate arithmetic expressions; (c) resync arithmetic results to
	// string formatting on a just-in-time basis when output is printed.
	printrep      string
	printrepValid bool
	intval        int
	floatval      float64
	boolval       bool
	arrayval      []Mlrval
	mapval        *Mlrmap
	// These are first-class-function literals. Stored here as interface{} to
	// avoid what would otherwise be a package-dependency cycle with the
	// mlr/internal/pkg/dsl/cst package.
	funcval interface{}
}

// Enumeration for mlrval types
//
// There are two kinds of null: ABSENT (key not present in a record) and VOID
// (key present with empty value).  Note void is an acceptable string (empty
// string) but not an acceptable number. (In Javascript, similarly, there are
// undefined and null, respectively.)

type MVType int

// Important: the values of these enums are used to index into disposition
// matrices. If they are changed, it will break the disposition matrices, or
// they will all need manual re-indexing.
const (
	// Type not yet determined, during JSON decode.
	MT_PENDING MVType = -1

	// E.g. error encountered in one eval & it propagates up the AST at
	// evaluation time.  Various runtime errors, such as file-not-found, result
	// in a message to stderr and os.Exit(1). But errors in user-provided data
	// are intended to result in "(error)"-valued output rather than a crash.
	// This is analogous to the way that IEEE-754 arithmetic carries around
	// Inf and NaN through computation chains.
	MT_ERROR = 0

	// Key not present in input record, e.g. 'foo = $nosuchkey'
	MT_ABSENT = 1

	// Used only for JSON null, and for 'empty' slots when an array is
	// auto-extended by assigning to an index having a gap from the last index.
	// E.g. x=[1,2,3] then x[5]=5; now x[4] is null
	MT_NULL = 2

	// Key present in input record with empty value, e.g. input data '$x=,$y=2'
	MT_VOID = 3

	MT_STRING = 4

	MT_INT = 5

	MT_FLOAT = 6

	MT_BOOL = 7

	MT_ARRAY = 8

	MT_MAP = 9

	MT_FUNC = 10

	// Not a type -- this is a dimension for disposition vectors and
	// disposition matrices. For example, when we want to add two mlrvals,
	// instead of if/elsing or switching on the types of both operands, we
	// instead jump directly to a type-specific function in a matrix of
	// function pointers which is MT_DIM x MT_DIM.
	MT_DIM = 11
)

var TYPE_NAMES = [MT_DIM]string{
	"error",
	"absent",
	"null",
	"empty", // For backward compatiblity with the C impl: this is user-visible
	"string",
	"int",
	"float",
	"bool",
	"array",
	"map",
	"funct",
}

// ----------------------------------------------------------------
// For typed assignments in the DSL

// TODO: comment more re typedecls
const MT_TYPE_MASK_STRING = (1 << MT_STRING) | (1 << MT_VOID)
const MT_TYPE_MASK_INT = 1 << MT_INT
const MT_TYPE_MASK_FLOAT = 1 << MT_FLOAT
const MT_TYPE_MASK_NUM = (1 << MT_INT) | (1 << MT_FLOAT)
const MT_TYPE_MASK_BOOL = 1 << MT_BOOL
const MT_TYPE_MASK_ARRAY = 1 << MT_ARRAY
const MT_TYPE_MASK_MAP = 1 << MT_MAP
const MT_TYPE_MASK_VAR = (1 << MT_VOID) |
	(1 << MT_STRING) |
	(1 << MT_INT) |
	(1 << MT_FLOAT) |
	(1 << MT_BOOL) |
	(1 << MT_ARRAY) |
	(1 << MT_MAP)
const MT_TYPE_MASK_FUNC = 1 << MT_FUNC

// Not exposed in userspace
const MT_TYPE_MASK_ANY = (1 << MT_ERROR) | (1 << MT_ABSENT) | MT_TYPE_MASK_VAR | MT_TYPE_MASK_FUNC

// TODO: const these throughout
var typeNameToMaskMap = map[string]int{
	"var":   MT_TYPE_MASK_VAR,
	"str":   MT_TYPE_MASK_STRING,
	"int":   MT_TYPE_MASK_INT,
	"float": MT_TYPE_MASK_FLOAT,
	"num":   MT_TYPE_MASK_NUM,
	"bool":  MT_TYPE_MASK_BOOL,
	"arr":   MT_TYPE_MASK_ARRAY,
	"map":   MT_TYPE_MASK_MAP,
	"any":   MT_TYPE_MASK_ANY,
	"funct": MT_TYPE_MASK_FUNC,
}

func TypeNameToMask(typeName string) (mask int, present bool) {
	retval := typeNameToMaskMap[typeName]
	if retval != 0 {
		return retval, true
	} else {
		return 0, false
	}
}

func (mv *Mlrval) GetTypeBit() int {
	return 1 << mv.mvtype
}
