package lib

// The `lib.Mlrval` structure includes **string, int, float, boolean, void,
// absent, and error** types (not unlike PHP's `zval`) as well as
// type-conversion logic for various operators.
//
// Whenever I say "int" and "float" with regard to mlrvals I always mean
// "int64" and "float64". If I ever miss a spot and use Go int/float types then
// that is a bug. It would be great to be able to somehow lint for this.

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
	intval        int64
	floatval      float64
	boolval       bool
	arrayval      []Mlrval
	mapval        *Mlrmap
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
	// E.g. error encountered in one eval & it propagates up the AST at
	// evaluation time.  Various runtime errors, such as file-not-found, result
	// in a message to stderr and os.Exit(1). But errors in user-provided data
	// are intended to result in "(error)"-valued output rather than a crash.
	// This is analogous to the way that IEEE-754 arithmetic carries around
	// Inf and NaN through computation chains.
	MT_ERROR MVType = 0

	// Key not present in input record, e.g. 'foo = $nosuchkey'
	MT_ABSENT = 1

	// Key present in input record with empty value, e.g. input data '$x=,$y=2'
	MT_VOID = 2

	MT_STRING = 3

	MT_INT = 4

	MT_FLOAT = 5

	MT_BOOL = 6

	MT_ARRAY = 7

	MT_MAP = 8

	// Not a type -- this is a dimension for disposition vectors and
	// disposition matrices. For example, when we want to add two mlrvals,
	// instead of if/elsing or switching on the types of both operands, we
	// instead jump directly to a type-specific function in a matrix of
	// function pointers which is MT_DIM x MT_DIM.
	MT_DIM = 9
)
