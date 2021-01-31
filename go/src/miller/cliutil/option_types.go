// ================================================================
// Items which might better belong in miller/cli, but which are placed in a
// deeper package to avoid a package-dependency cycle between miller/cli and
// miller/transforming.
// ================================================================

package cliutil

// ----------------------------------------------------------------
//typedef struct _generator_opts_t {
//	char* field_name;
//	// xxx to do: convert to mv_t
//	long long start;
//	long long stop;
//	long long step;
//} generator_opts_t;

// ----------------------------------------------------------------
type TReaderOptions struct {
	InputFileFormat string
	IRS             string
	IFS             string
	IPS             string
	IFLATSEP        string

	//	char* input_json_flatten_separator;
	//	json_array_ingest_t  json_array_ingest;

	//	allow_repeat_ifs bool;
	//	allow_repeat_ips bool;
	UseImplicitCSVHeader bool
	AllowRaggedCSVInput  bool
	//
	//	// Command for popen on input, e.g. "zcat -cf <". Can be null in which case
	//	// files are read directly rather than through a pipe.
	//	prepipe string;
	//
	//	comment_handling_t comment_handling;
	//	comment_string string;
	//
	//	// Fake internal-data-generator 'reader'
	//	generator_opts_t generator_opts;

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
}

// ----------------------------------------------------------------
type TWriterOptions struct {
	OutputFileFormat string
	ORS              string
	OFS              string
	OPS              string
	OFLATSEP         string

	HeaderlessCSVOutput bool
	BarredPprintOutput  bool

	//	right_justify_xtab_value bool;
	//	right_align_pprint bool;
	WrapJSONOutputInOuterList bool
	JSONOutputMultiline       bool // Not using miller/types enum to avoid package cycle
	//	json_quote_int_keys bool;
	//	json_quote_non_string_values bool;
	//	output_json_flatten_separator string;
	//	oosvar_flatten_separator string;
	//
	//	quoting_t oquoting;

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
}

// ----------------------------------------------------------------
type TOptions struct {
	ReaderOptions TReaderOptions
	WriterOptions TWriterOptions

	FileNames []string

	// These are used to construct the transformer list. In particular, for in-place mode
	// they're reconstructed for each file.  We make copies since each pass through a
	// CLI-parser operates destructively, principally by running strtok over
	// comma-delimited field-name lists.
	//
	//	char**  original_argv;
	//	char**  non_in_place_argv;
	//	int     argc;
	//	int     transformer_argb;
	//
	//
	//	char* ofmt;
	//	nr_progress_mod int64u;
	//
	//	do_in_place bool;
	//
	NoInput bool // mlr -n

	HaveRandSeed bool
	RandSeed     int
}

// ----------------------------------------------------------------
func DefaultOptions() TOptions {
	return TOptions{
		ReaderOptions: DefaultReaderOptions(),
		WriterOptions: DefaultWriterOptions(),

		FileNames: make([]string, 0),
		NoInput:   false,
	}
}

func DefaultReaderOptions() TReaderOptions {
	return TReaderOptions{
		InputFileFormat: "dkvp", // xxx constify at top
		IRS:             "\n",
		IFS:             ",",
		IPS:             "=",
		IFLATSEP:        ".",
		AutoUnflatten:   true,
	}
}

func DefaultWriterOptions() TWriterOptions {
	return TWriterOptions{
		OutputFileFormat: "dkvp",
		ORS:              "\n",
		OFS:              ",",
		OPS:              "=",
		OFLATSEP:         ".",

		HeaderlessCSVOutput:       false,
		WrapJSONOutputInOuterList: false,
		JSONOutputMultiline:       true,
		AutoFlatten:               true,
	}
}
