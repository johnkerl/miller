package scan

// TODO: comment re context

type ScanType int

const (
	scanTypeString     ScanType = 0
	scanTypeDecimalInt          = 1
	scanTypeOctalInt            = 2
	scanTypeHexInt              = 3
	scanTypeBinaryInt           = 4
	scanTypeMaybeFloat          = 5
	scanTypeBool                = 6
)

const typeNameString = "string"
const typeNameDecimalInt = "decint"
const typeNameOctalInt = "octint"
const typeNameHexInt = "hexint"
const typeNameBinaryInt = "binint"
const typeNameMaybeFloat = "float?"
const typeNameBool = "bool"

var TypeNames = []string{
	typeNameString,
	typeNameDecimalInt,
	typeNameOctalInt,
	typeNameHexInt,
	typeNameBinaryInt,
	typeNameMaybeFloat,
	typeNameBool,
}
