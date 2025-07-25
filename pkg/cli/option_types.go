// ================================================================
// Items which might better belong in miller/cli, but which are placed in a
// deeper package to avoid a package-dependency cycle between miller/cli and
// miller/transforming.
// ================================================================

package cli

import (
	"regexp"

	"github.com/johnkerl/miller/v6/pkg/lib"
)

type TCommentHandling int

const (
	CommentsAreData TCommentHandling = iota
	SkipComments
	PassComments
)
const DEFAULT_COMMENT_STRING = "#"

const DEFAULT_GEN_FIELD_NAME = "i"
const DEFAULT_GEN_START_AS_STRING = "1"
const DEFAULT_GEN_STEP_AS_STRING = "1"
const DEFAULT_GEN_STOP_AS_STRING = "100"

const DEFAULT_RECORDS_PER_BATCH = 500

type TGeneratorOptions struct {
	FieldName     string
	StartAsString string
	StepAsString  string
	StopAsString  string
}

type TReaderOptions struct {
	InputFileFormat  string
	IFS              string
	IPS              string
	IRS              string
	AllowRepeatIFS   bool
	IFSRegex         *regexp.Regexp
	IPSRegex         *regexp.Regexp
	DedupeFieldNames bool

	// If unspecified on the command line, these take input-format-dependent
	// defaults.  E.g. default FS is comma for DKVP but space for NIDX;
	// default AllowRepeatIFS is false for CSV but true for PPRINT.
	ifsWasSpecified            bool
	ipsWasSpecified            bool
	irsWasSpecified            bool
	allowRepeatIFSWasSpecified bool

	UseImplicitHeader    bool
	AllowRaggedCSVInput  bool
	CSVLazyQuotes        bool
	CSVTrimLeadingSpace  bool
	BarredPprintInput    bool
	IncrementImplicitKey bool

	CommentHandling TCommentHandling
	CommentString   string

	// Fake internal-data-generator 'reader'
	GeneratorOptions TGeneratorOptions

	// For out-of-process handling of compressed data, via popen
	Prepipe string
	// For most things like gunzip we do 'gunzip < filename | mlr ...' if
	// filename is present, else 'gunzip | mlr ...' if reading from stdin.
	// However some commands like 'unzip -qc' are weird so this option lets
	// people give the command and we won't insert the '<'.
	PrepipeIsRaw bool
	// For in-process gunzip/bunzip2/zcat (distinct from prepipe)
	FileInputEncoding lib.TFileInputEncoding

	// TODO: comment
	RecordsPerBatch int64
}

// ----------------------------------------------------------------
type TWriterOptions struct {
	OutputFileFormat string
	ORS              string
	OFS              string
	OPS              string
	FLATSEP          string

	FlushOnEveryRecord             bool
	flushOnEveryRecordWasSpecified bool

	// If unspecified on the command line, these take input-format-dependent
	// defaults.  E.g. default FS is comma for DKVP but space for NIDX.
	ofsWasSpecified bool
	opsWasSpecified bool
	orsWasSpecified bool

	HeaderlessOutput         bool
	BarredPprintOutput       bool
	RightAlignedPPRINTOutput bool
	RightAlignedXTABOutput   bool

	// JSON output: --jlistwrap on, --jvstack on
	// JSON Lines output: --jlistwrap off, --jvstack off
	WrapJSONOutputInOuterList bool // --jlistwrap
	JSONOutputMultiline       bool // --jvstack
	JVQuoteAll                bool // --jvquoteall
	// Not using miller/types enum to avoid package cycle

	CSVQuoteAll bool // --quote-all

	// When we read things like
	//
	//   x:a=1,x:b=2
	//
	// which is how we write out nested data structures for non-nested formats
	// (all but JSON), the default behavior is to unflatten them back to
	//
	//   {"x": {"a": 1}, {"b": 2}}
	//
	// unless the user explicitly asks to suppress that.
	AutoUnflatten bool

	// The default behavior is to flatten nested data structures like
	//
	//   {"x": {"a": 1}, {"b": 2}}
	//
	// down to
	//
	//   x:a=1,x:b=2
	//
	// which is how we write out nested data structures for non-nested formats
	// (all but JSON) -- unless the user explicitly asks to suppress that.
	AutoFlatten bool

	// Default CSV/TSV:
	//   a=1,b=2,c=3
	//   a=4,b=5
	// leads to
	//   a,b,c
	//   1 2,3
	//   4,5, <-- note trailing empty field
	// and
	//   a=1,b=2,c=3
	//   d=4,e=5
	// leads to
	//   fatal error
	//
	// With this flag:
	//   a=1,b=2,c=3
	//   a=4,b=5
	// leads to
	//   a,b,c
	//   1 2,3
	//
	//   a,b
	//   4,5
	//
	// and
	//   a=1,b=2,c=3
	//   d=4,e=5
	// leads to
	//   a,b,c
	//   1,2,3
	//
	//   d,e
	//   4,5
	NoAutoUnsparsify bool

	// For floating-point numbers: "" means use the Go default.
	FPOFMT string

	// Fatal the process when error data in a given record is about to be output.
	FailOnDataError bool
}

// ----------------------------------------------------------------
type TOptions struct {
	ReaderOptions TReaderOptions
	WriterOptions TWriterOptions

	// Data files to be operated on: e.g. given 'mlr cat foo.dat bar.dat', this
	// is ["foo.dat", "bar.dat"].
	FileNames []string

	// DSL files to be loaded for every put/filter operation -- like 'put -f'
	// or 'filter -f' but specified up front on the command line, suitable for
	// .mlrrc. Use-case is someone has DSL functions they always want to be
	// defined.
	//
	// Risk of CVE if this is in .mlrrc so --load and --mload are explicitly
	// denied in the .mlrrc reader.
	DSLPreloadFileNames []string

	NRProgressMod int64
	DoInPlace     bool // mlr -I
	NoInput       bool // mlr -n

	HaveRandSeed bool
	RandSeed     int64

	PrintElapsedTime bool // mlr --time
}

// Not usable until FinalizeReaderOptions and FinalizeWriterOptions are called.
func DefaultOptions() *TOptions {
	return &TOptions{
		ReaderOptions: DefaultReaderOptions(),
		WriterOptions: DefaultWriterOptions(),

		FileNames:           make([]string, 0),
		DSLPreloadFileNames: make([]string, 0),
		NoInput:             false,
	}
}

// Not usable until FinalizeReaderOptions is called on it.
func DefaultReaderOptions() TReaderOptions {
	return TReaderOptions{
		InputFileFormat: "dkvp", // TODO: constify at top, or maybe formats.DKVP in package
		// FinalizeReaderOptions will compute IFSRegex and IPSRegex.
		IRS:               "\n",
		IFS:               ",",
		IPS:               "=",
		CommentHandling:   CommentsAreData,
		FileInputEncoding: lib.FileInputEncodingDefault,
		GeneratorOptions: TGeneratorOptions{
			FieldName:     DEFAULT_GEN_FIELD_NAME,
			StartAsString: DEFAULT_GEN_START_AS_STRING,
			StepAsString:  DEFAULT_GEN_STEP_AS_STRING,
			StopAsString:  DEFAULT_GEN_STOP_AS_STRING,
		},
		DedupeFieldNames: true,

		// TODO: comment
		RecordsPerBatch: DEFAULT_RECORDS_PER_BATCH,
	}
}

// Not usable until FinalizeWriterOptions is called on it.
func DefaultWriterOptions() TWriterOptions {
	return TWriterOptions{
		OutputFileFormat:   "dkvp",
		ORS:                "\n",
		OFS:                ",",
		OPS:                "=",
		FLATSEP:            ".",
		FlushOnEveryRecord: true,

		HeaderlessOutput: false,

		WrapJSONOutputInOuterList: true,
		JSONOutputMultiline:       true,

		AutoUnflatten: true,
		AutoFlatten:   true,

		FPOFMT: "",
	}
}
