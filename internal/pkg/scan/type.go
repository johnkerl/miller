package scan

// TODO: comment re context

type ScanType int

const (
	scanTypeString                ScanType = 0
	scanTypeDecimalInt                     = 1
	scanTypeLeadingZeroDecimalInt          = 2
	scanTypeOctalInt                       = 3
	scanTypeLeadingZeroOctalInt            = 4
	scanTypeHexInt                         = 5
	scanTypeBinaryInt                      = 6
	scanTypeMaybeFloat                     = 7
)

const typeNameString = "string"
const typeNameDecimalInt = "decint"              // e.g.       123
const typeNameLeadingZeroDecimalInt = "lzdecint" // e.g.       0899
const typeNameOctalInt = "octint"                // e.g.       0o377
const typeNameLeadingZeroOctalInt = "lzoctint"   // e.g.       0377
const typeNameHexInt = "hexint"                  // e.g.       0xcafe
const typeNameBinaryInt = "binint"               // e.g.       0b1011
const typeNameMaybeFloat = "float?"              // characters in     [0-9\.-+eE] but needs parse to be sure

var TypeNames = []string{
	typeNameString,
	typeNameDecimalInt,
	typeNameLeadingZeroDecimalInt,
	typeNameOctalInt,
	typeNameLeadingZeroOctalInt,
	typeNameHexInt,
	typeNameBinaryInt,
	typeNameMaybeFloat,
}
