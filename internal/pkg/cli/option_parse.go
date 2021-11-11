// ================================================================
// Items which might better belong in miller/cli, but which are placed in a
// deeper package to avoid a package-dependency cycle between miller/cli and
// miller/transforming.
// ================================================================

package cli

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/colorizer"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// FinalizeReaderOptions does a few things.
// * If a file format was specified but one or more separators were not, a
//   default specific to that file format is applied.
// * Computing regexes for IPS and IFS, and unbackslashing IRS.  This is
//   because the '\n' at the command line which is Go "\\n" (a backslash and an
//   n) needs to become the single newline character, and likewise for "\t", etc.
// * IFS/IPS can have escapes like "\x1f" which aren't valid regex literals
//   so we unhex them. For example, from "\x1f" -- the four bytes '\', 'x', '1', 'f'
//   -- to the single byte with hex code 0x1f.
func FinalizeReaderOptions(readerOptions *TReaderOptions) {

	readerOptions.IFS = lib.UnhexStringLiteral(readerOptions.IFS)
	readerOptions.IPS = lib.UnhexStringLiteral(readerOptions.IPS)

	if !readerOptions.IFSWasSpecified {
		readerOptions.IFS = defaultFSes[readerOptions.InputFileFormat]
	}
	if !readerOptions.IPSWasSpecified {
		readerOptions.IPS = defaultPSes[readerOptions.InputFileFormat]
	}
	if !readerOptions.IRSWasSpecified {
		readerOptions.IRS = defaultRSes[readerOptions.InputFileFormat]
	}
	if !readerOptions.AllowRepeatIFSWasSpecified {
		// Special case for Miller 6 upgrade -- now that we have regexing for mixes of tabs
		// and spaces, that should now be the default for NIDX. But *only* for NIDX format,
		// and if IFS wasn't specified.
		if readerOptions.InputFileFormat == "nidx" && !readerOptions.IFSWasSpecified {
			readerOptions.IFS = WHITESPACE
		} else {
			readerOptions.AllowRepeatIFS = defaultAllowRepeatIFSes[readerOptions.InputFileFormat]
		}
	}
	if !readerOptions.AllowRepeatIPSWasSpecified {
		readerOptions.AllowRepeatIPS = defaultAllowRepeatIPSes[readerOptions.InputFileFormat]
	}

	if readerOptions.SuppressIFSRegexing {
		readerOptions.IFSRegex = nil
	} else if readerOptions.AllowRepeatIFS {
		readerOptions.IFSRegex = lib.CompileMillerRegexOrDie("(" + readerOptions.IFS + ")+")
	} else {
		readerOptions.IFSRegex = lib.CompileMillerRegexOrDie(readerOptions.IFS)
	}

	if readerOptions.SuppressIPSRegexing {
		readerOptions.IPSRegex = nil
	} else if readerOptions.AllowRepeatIPS {
		readerOptions.IPSRegex = lib.CompileMillerRegexOrDie("(" + readerOptions.IPS + ")+")
	} else {
		readerOptions.IPSRegex = lib.CompileMillerRegexOrDie(readerOptions.IPS)
	}

	readerOptions.IRS = lib.UnbackslashStringLiteral(readerOptions.IRS)
}

// FinalizeWriterOptions unbackslashes OPS, OFS, and ORS.  This is because
// because the '\n' at the command line which is Go "\\n" (a backslash and an
// n) needs to become the single newline character., and likewise for "\t", etc.
func FinalizeWriterOptions(writerOptions *TWriterOptions) {
	if !writerOptions.OFSWasSpecified {
		writerOptions.OFS = defaultFSes[writerOptions.OutputFileFormat]
	}
	if !writerOptions.OPSWasSpecified {
		writerOptions.OPS = defaultPSes[writerOptions.OutputFileFormat]
	}
	if !writerOptions.ORSWasSpecified {
		writerOptions.ORS = defaultRSes[writerOptions.OutputFileFormat]
	}

	writerOptions.OFS = lib.UnbackslashStringLiteral(writerOptions.OFS)
	writerOptions.OPS = lib.UnbackslashStringLiteral(writerOptions.OPS)
	writerOptions.ORS = lib.UnbackslashStringLiteral(writerOptions.ORS)
}

// ================================================================
var FLAG_TABLE = FlagTable{
	sections: []*FlagSection{
		&LegacyFlagSection,
		&SeparatorFlagSection,
		&FileFormatFlagSection,
		&FormatConversionKeystrokeSaverFlagSection,
		&JSONOnlyFlagSection,
		&CSVOnlyFlagSection,
		&PPRINTOnlyFlagSection,
		&CompressedDataFlagSection,
		&CommentsInDataFlagSection,
		&OutputColorizationFlagSection,
		&FlattenUnflattenFlagSection,
		&MiscFlagSection,
	},
}

func init() {
	FLAG_TABLE.Sort()
}

// ================================================================
// SEPARATOR FLAGS

func SeparatorPrintInfo() {
	fmt.Println(`See the Separators doc page for more about record separators, field
separators, and pair separators. Also see the File formats doc page, or
` + "`mlr help file-formats`" + `, for more about the file formats Miller supports.

In brief:

* For DKVP records like ` + "`x=1,y=2,z=3`" + `, the fields are separated by a comma,
  the key-value pairs are separated by a comma, and each record is separated
  from the next by a newline.
* Each file format has its own default separators.
* Most formats, such as CSV, don't support pair-separators: keys are on the CSV
  header line and values are on each CSV data line; keys and values are not
  placed next to one another.
* Some separators are not programmable: for example JSON uses a colon as a
  pair separator but this is non-modifiable in the JSON spec.
* You can set separators differently between Miller's input and output --
  hence ` + "`--ifs`" + ` and ` + "`--ofs`" + `, etc.

Notes about line endings:

* Default line endings (` + "`--irs`" + ` and ` + "`--ors`" + `) are newline
  which is interpreted to accept carriage-return/newline files (e.g. on Windows)
  for input, and to produce platform-appropriate line endings on output.

Notes about all other separators:

* IPS/OPS are only used for DKVP and XTAB formats, since only in these formats
  do key-value pairs appear juxtaposed.
* IRS/ORS are ignored for XTAB format. Nominally IFS and OFS are newlines;
  XTAB records are separated by two or more consecutive IFS/OFS -- i.e.
  a blank line. Everything above about ` + "`--irs/--ors/--rs auto`" + ` becomes ` + "`--ifs/--ofs/--fs`" + `
  auto for XTAB format. (XTAB's default IFS/OFS are "auto".)
* OFS must be single-character for PPRINT format. This is because it is used
  with repetition for alignment; multi-character separators would make
  alignment impossible.
* OPS may be multi-character for XTAB format, in which case alignment is
  disabled.
* TSV is simply CSV using tab as field separator (` + "`--fs tab`" + `).
* FS/PS are ignored for markdown format; RS is used.
* All FS and PS options are ignored for JSON format, since they are not relevant
  to the JSON format.
* You can specify separators in any of the following ways, shown by example:
  - Type them out, quoting as necessary for shell escapes, e.g.
    ` + "`--fs '|' --ips :`" + `
  - C-style escape sequences, e.g. ` + "`--rs '\\r\\n' --fs '\\t'`" + `.
  - To avoid backslashing, you can use any of the following names:`)
	fmt.Println()

	// Go doesn't preserve insertion order in its arrays so here we are inlining a sort.
	aliases := lib.GetArrayKeysSorted(SEPARATOR_NAMES_TO_VALUES)
	for _, alias := range aliases {
		// Really absurd level of indent needed to get fixed-with font in mkdocs here,
		// I don't know why. Usually it only takes 4, not 10.
		fmt.Printf("          %-10s = \"%s\"\n", alias, SEPARATOR_NAMES_TO_VALUES[alias])
	}
	fmt.Println()

	fmt.Println("* Default separators by format:")
	fmt.Println()

	formats := lib.GetArrayKeysSorted(defaultFSes)
	// Really absurd level of indent needed to get fixed-with font in mkdocs here,
	// I don't know why. Usually it only takes 4, not 8.
	fmt.Printf("        %-8s %-6s %-6s %-s\n", "Format", "FS", "PS", "RS")
	for _, format := range formats {
		defaultFS := "\"" + strings.ReplaceAll(defaultFSes[format], "\n", "\\n") + "\""
		defaultPS := "\"" + strings.ReplaceAll(defaultPSes[format], "\n", "\\n") + "\""
		defaultRS := "\"" + strings.ReplaceAll(defaultRSes[format], "\n", "\\n") + "\""
		if defaultFS == "\"N/A\"" {
			defaultFS = "N/A"
		}
		if defaultPS == "\"N/A\"" {
			defaultPS = "N/A"
		}
		if defaultRS == "\"N/A\"" {
			defaultRS = "N/A"
		}
		fmt.Printf("        %-8s %-6s %-6s %-s\n", format, defaultFS, defaultPS, defaultRS)
	}
}

func ListSeparatorAliasesForOnlineHelp() {
	// Go doesn't preserve insertion order in its arrays so here we are inlining a sort.
	aliases := lib.GetArrayKeysSorted(SEPARATOR_NAMES_TO_VALUES)
	for _, alias := range aliases {
		// Really absurd level of indent needed to get fixed-with font in mkdocs here,
		// I don't know why. Usually it only takes 4, not 10.
		fmt.Printf("%-10s = \"%s\"\n", alias, SEPARATOR_NAMES_TO_VALUES[alias])
	}
}

func init() { SeparatorFlagSection.Sort() }

var SeparatorFlagSection = FlagSection{
	name:        "Separator flags",
	infoPrinter: SeparatorPrintInfo,
	flags: []Flag{

		{
			name: "--ifs",
			arg:  "{string}",
			help: "Specify FS for input.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				// Backward compatibility with Miller <= 5. Auto-inference of
				// LF vs CR/LF line endings is handled within Go libraries so
				// we needn't do anything ourselves.
				if args[*pargi+1] != "auto" {
					options.ReaderOptions.IFS = SeparatorFromArg(args[*pargi+1])
					options.ReaderOptions.IFSWasSpecified = true
				}
				*pargi += 2
			},
		},

		{
			name: "--ips",
			arg:  "{string}",
			help: "Specify PS for input.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.IPS = SeparatorFromArg(args[*pargi+1])
				options.ReaderOptions.IPSWasSpecified = true
				*pargi += 2
			},
		},

		{
			name: "--irs",
			arg:  "{string}",
			help: "Specify RS for input.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				// Backward compatibility with Miller <= 5. Auto-inference of
				// LF vs CR/LF line endings is handled within Go libraries so
				// we needn't do anything ourselves.
				if args[*pargi+1] != "auto" {
					options.ReaderOptions.IRS = SeparatorFromArg(args[*pargi+1])
					options.ReaderOptions.IRSWasSpecified = true
				}
				*pargi += 2
			},
		},

		{
			name: "--repifs",
			help: "Let IFS be repeated: e.g. for splitting on multiple spaces.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.AllowRepeatIFS = true
				options.ReaderOptions.AllowRepeatIFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--ors",
			arg:  "{string}",
			help: "Specify RS for output.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				// Backward compatibility with Miller <= 5. Auto-inference of
				// LF vs CR/LF line endings is handled within Go libraries so
				// we needn't do anything ourselves.
				if args[*pargi+1] != "auto" {
					options.WriterOptions.ORS = SeparatorFromArg(args[*pargi+1])
					options.WriterOptions.ORSWasSpecified = true
				}
				*pargi += 2
			},
		},

		{
			name: "--ofs",
			arg:  "{string}",
			help: "Specify FS for output.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.WriterOptions.OFS = SeparatorFromArg(args[*pargi+1])
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 2
			},
		},

		{
			name: "--ops",
			arg:  "{string}",
			help: "Specify PS for output.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.WriterOptions.OPS = SeparatorFromArg(args[*pargi+1])
				options.WriterOptions.OPSWasSpecified = true
				*pargi += 2
			},
		},

		{
			name: "--rs",
			arg:  "{string}",
			help: "Specify RS for input and output.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				// Backward compatibility with Miller <= 5. Auto-inference of
				// LF vs CR/LF line endings is handled within Go libraries so
				// we needn't do anything ourselves.
				if args[*pargi+1] != "auto" {
					options.ReaderOptions.IRS = SeparatorFromArg(args[*pargi+1])
					options.WriterOptions.ORS = SeparatorFromArg(args[*pargi+1])
					options.ReaderOptions.IRSWasSpecified = true
					options.WriterOptions.ORSWasSpecified = true
				}
				*pargi += 2
			},
		},

		{
			name: "--fs",
			arg:  "{string}",
			help: "Specify FS for input and output.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				// Backward compatibility with Miller <= 5. Auto-inference of
				// LF vs CR/LF line endings is handled within Go libraries so
				// we needn't do anything ourselves.
				if args[*pargi+1] != "auto" {
					options.ReaderOptions.IFS = SeparatorFromArg(args[*pargi+1])
					options.WriterOptions.OFS = SeparatorFromArg(args[*pargi+1])
					options.ReaderOptions.IFSWasSpecified = true
					options.WriterOptions.OFSWasSpecified = true
				}
				*pargi += 2
			},
		},

		{
			name: "--ps",
			arg:  "{string}",
			help: "Specify PS for input and output.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.IPS = SeparatorFromArg(args[*pargi+1])
				options.WriterOptions.OPS = SeparatorFromArg(args[*pargi+1])
				options.ReaderOptions.IPSWasSpecified = true
				options.WriterOptions.OPSWasSpecified = true
				*pargi += 2
			},
		},

		{
			name: "--no-ifs-regex",
			help: `Don't treat IFS value as a regular expression. Useful if your IFS is ".".`,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.SuppressIFSRegexing = true
				*pargi += 1
			},
		},

		{
			name: "--no-ips-regex",
			help: `Don't treat IPS value as a regular expression. Useful if your IPS is ".".`,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.SuppressIPSRegexing = true
				*pargi += 1
			},
		},
	},
}

// ================================================================
// JSON-ONLY FLAGS

func JSONOnlyPrintInfo() {
	fmt.Println("These are flags which are applicable to JSON format.")
}

func init() { JSONOnlyFlagSection.Sort() }

var JSONOnlyFlagSection = FlagSection{
	name:        "JSON-only flags",
	infoPrinter: JSONOnlyPrintInfo,
	flags: []Flag{

		{
			name: "--jvstack",
			help: "Put one key-value pair per line for JSON output (multi-line output).",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.JSONOutputMultiline = true
				*pargi += 1
			},
		},

		{
			name: "--no-jvstack",
			help: "Put objects/arrays all on one line for JSON output.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.JSONOutputMultiline = false
				*pargi += 1
			},
		},

		{
			name:     "--jlistwrap",
			altNames: []string{"--jl"},
			help:     "Wrap JSON output in outermost `[ ]`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.WrapJSONOutputInOuterList = true
				*pargi += 1
			},
		},
	},
}

// ================================================================
// PPRINT-ONLY FLAGS

func PPRINTOnlyPrintInfo() {
	fmt.Println("These are flags which are applicable to PPRINT output format.")
}

func init() { PPRINTOnlyFlagSection.Sort() }

var PPRINTOnlyFlagSection = FlagSection{
	name:        "PPRINT-only flags",
	infoPrinter: PPRINTOnlyPrintInfo,
	flags: []Flag{

		{
			name: "--right",
			help: "Right-justifies all fields for PPRINT output.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.RightAlignedPPRINTOutput = true
				*pargi += 1
			},
		},

		{
			name: "--barred",
			help: "Prints a border around PPRINT output (not available for input).",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.BarredPprintOutput = true
				*pargi += 1
			},
		},
	},
}

// ================================================================
// LEGACY FLAGS

func LegacyFlagInfoPrint() {
	fmt.Println(`These are flags which don't do anything in the current Miller version.
They are accepted as no-op flags in order to keep old scripts from breaking.`)
}

func init() { LegacyFlagSection.Sort() }

var LegacyFlagSection = FlagSection{
	name:        "Legacy flags",
	infoPrinter: LegacyFlagInfoPrint,
	flags: []Flag{

		{
			name:   "--mmap",
			help:   "Miller no longer uses memory-mapping to access data files.",
			parser: NoOpParse1,
		},

		{
			name:   "--no-mmap",
			help:   "Miller no longer uses memory-mapping to access data files.",
			parser: NoOpParse1,
		},

		{
			name:   "--no-fflush",
			help:   "The current implementation of Miller does not use buffered output, so there is no longer anything to suppress here.",
			parser: NoOpParse1,
		},

		{
			name:   "--jsonx",
			help:   "The `--jvstack` flag is now default true in Miller 6.",
			parser: NoOpParse1,
		},
		{
			name:   "--ojsonx",
			help:   "The `--jvstack` flag is now default true in Miller 6.",
			parser: NoOpParse1,
		},

		{
			name:   "--jknquoteint",
			help:   "Type information from JSON input files is now preserved throughout the processing stream.",
			parser: NoOpParse1,
		},

		{
			name:   "--jquoteall",
			help:   "Type information from JSON input files is now preserved throughout the processing stream.",
			parser: NoOpParse1,
		},

		{
			name:   "--jvquoteall",
			help:   "Type information from JSON input files is now preserved throughout the processing stream.",
			parser: NoOpParse1,
		},

		{
			name:   "--json-fatal-arrays-on-input",
			help:   "Miller now supports arrays as of version 6.",
			parser: NoOpParse1,
		},

		{
			name:   "--json-map-arrays-on-input",
			help:   "Miller now supports arrays as of version 6.",
			parser: NoOpParse1,
		},

		{
			name:   "--json-skip-arrays-on-input",
			help:   "Miller now supports arrays as of version 6.",
			parser: NoOpParse1,
		},

		{
			name:   "--vflatsep",
			help:   "Ignored as of version 6. This functionality is subsumed into JSON formatting.",
			parser: NoOpParse1,
		},

		{
			name:   "--quote-all",
			help:   "Ignored as of version 6. Types are inferred/retained through the processing flow now.",
			parser: NoOpParse1,
		},
		{
			name:   "--quote-none",
			help:   "Ignored as of version 6. Types are inferred/retained through the processing flow now.",
			parser: NoOpParse1,
		},
		{
			name:   "--quote-minimal",
			help:   "Ignored as of version 6. Types are inferred/retained through the processing flow now.",
			parser: NoOpParse1,
		},
		{
			name:   "--quote-numeric",
			help:   "Ignored as of version 6. Types are inferred/retained through the processing flow now.",
			parser: NoOpParse1,
		},
		{
			name:   "--quote-original",
			help:   "Ignored as of version 6. Types are inferred/retained through the processing flow now.",
			parser: NoOpParse1,
		},
	},
}

// ================================================================
// FILE-FORMAT FLAGS

func FileFormatPrintInfo() {
	fmt.Println(`See the File formats doc page, and or ` + "`mlr help file-formats`" + `, for more
about file formats Miller supports.

Examples: ` + "`--csv`" + ` for CSV-formatted input and output; ` + "`--icsv --opprint`" + ` for
CSV-formatted input and pretty-printed output.

Please use ` + "`--iformat1 --oformat2`" + ` rather than ` + "`--format1 --oformat2`" + `.
The latter sets up input and output flags for ` + "`format1`" + `, not all of which
are overridden in all cases by setting output format to ` + "`format2`" + `.`)
}

func init() { FileFormatFlagSection.Sort() }

var FileFormatFlagSection = FlagSection{
	name:        "File-format flags",
	infoPrinter: FileFormatPrintInfo,
	flags: []Flag{

		{
			name: "--icsv",
			help: "Use CSV format for input data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				*pargi += 1
			},
		},

		{
			name: "--icsvlite",
			help: "Use CSV-lite format for input data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csvlite"
				*pargi += 1
			},
		},

		{
			name: "--itsv",
			help: "Use TSV format for input data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.ReaderOptions.IFSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--itsvlite",
			help: "Use TSV-lite format for input data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csvlite"
				options.ReaderOptions.IFS = "\t"
				options.ReaderOptions.IFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name:     "--iasv",
			altNames: []string{"--iasvlite"},
			help:     "Use ASV format for input data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csvlite"
				options.ReaderOptions.IFS = ASV_FS
				options.ReaderOptions.IRS = ASV_RS
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name:     "--iusv",
			altNames: []string{"--iusvlite"},
			help:     "Use USV format for input data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csvlite"
				options.ReaderOptions.IFS = USV_FS
				options.ReaderOptions.IRS = USV_RS
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--idkvp",
			help: "Use DKVP format for input data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				*pargi += 1
			},
		},

		{
			name: "--ijson",
			help: "Use JSON format for input data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				*pargi += 1
			},
		},

		{
			name: "--inidx",
			help: "Use NIDX format for input data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				*pargi += 1
			},
		},

		{
			name: "--ixtab",
			help: "Use XTAB format for input data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				*pargi += 1
			},
		},

		{
			name: "--ipprint",
			help: "Use PPRINT format for input data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "pprint"
				options.ReaderOptions.IFS = " "
				options.ReaderOptions.IFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "-i",
			arg:  "{format name}",
			help: "Use format name for input data. For example: `-i csv` is the same as `--icsv`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.InputFileFormat = args[*pargi+1]
				*pargi += 2
			},
		},

		{
			name: "--igen",
			help: `Ignore input files and instead generate sequential numeric input using --gen-field-name,
--gen-start, --gen-step, and --gen-stop values. See also the seqgen verb, which is more useful/intuitive.`,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "gen"
				*pargi += 1
			},
		},
		{
			name: "--gen-field-name",
			help: `Specify field name for --igen. Defaults to "` + DEFAULT_GEN_FIELD_NAME + `".`,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "gen"
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.GeneratorOptions.FieldName = args[*pargi+1]
				*pargi += 2
			},
		},
		{
			name: "--gen-start",
			help: "Specify start value for --igen. Defaults to " + DEFAULT_GEN_START_AS_STRING + ".",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "gen"
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.GeneratorOptions.StartAsString = args[*pargi+1]
				*pargi += 2
			},
		},
		{
			name: "--gen-step",
			help: "Specify step value for --igen. Defaults to " + DEFAULT_GEN_STEP_AS_STRING + ".",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "gen"
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.GeneratorOptions.StepAsString = args[*pargi+1]
				*pargi += 2
			},
		},
		{
			name: "--gen-stop",
			help: "Specify stop value for --igen. Defaults to " + DEFAULT_GEN_STOP_AS_STRING + ".",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "gen"
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.GeneratorOptions.StopAsString = args[*pargi+1]
				*pargi += 2
			},
		},

		{
			name: "-o",
			arg:  "{format name}",
			help: "Use format name for output data.  For example: `-o csv` is the same as `--ocsv`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.WriterOptions.OutputFileFormat = args[*pargi+1]
				*pargi += 2
			},
		},

		{
			name: "--ocsv",
			help: "Use CSV format for output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "csv"
				*pargi += 1
			},
		},

		{
			name: "--ocsvlite",
			help: "Use CSV-lite format for output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "csvlite"
				*pargi += 1
			},
		},

		{
			name: "--otsv",
			help: "Use TSV format for output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.OFS = "\t"
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--otsvlite",
			help: "Use TSV-lite format for output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "csvlite"
				options.WriterOptions.OFS = "\t"
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name:     "--oasv",
			altNames: []string{"--oasvlite"},
			help:     "Use ASV format for output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "csvlite"
				options.WriterOptions.OFS = ASV_FS
				options.WriterOptions.ORS = ASV_RS
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name:     "--ousv",
			altNames: []string{"--ousvlite"},
			help:     "Use USV format for output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "csvlite"
				options.WriterOptions.OFS = USV_FS
				options.WriterOptions.ORS = USV_RS
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--omd",
			help: "Use markdown-tabular format for output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "markdown"
				*pargi += 1
			},
		},

		{
			name: "--odkvp",
			help: "Use DKVP format for output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "dkvp"
				*pargi += 1
			},
		},

		{
			name: "--ojson",
			help: "Use JSON format for output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "json"
				*pargi += 1
			},
		},

		{
			name: "--onidx",
			help: "Use NIDX format for output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "nidx"
				options.WriterOptions.OFS = " "
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--oxtab",
			help: "Use XTAB format for output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "xtab"
				*pargi += 1
			},
		},

		{
			name: "--opprint",
			help: "Use PPRINT format for output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "pprint"
				*pargi += 1
			},
		},

		{
			name: "--io",
			arg:  "{format name}",
			help: "Use format name for input and output data. For example: `--io csv` is the same as `--csv`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				if defaultFSes[args[*pargi+1]] == "" {
					fmt.Fprintf(os.Stderr, "mlr: unrecognized I/O format \"%s\".\n",
						args[*pargi+1])
					os.Exit(1)
				}
				options.ReaderOptions.InputFileFormat = args[*pargi+1]
				options.WriterOptions.OutputFileFormat = args[*pargi+1]
				*pargi += 2
			},
		},

		{
			name:     "--csv",
			help:     "Use CSV format for input and output data.",
			altNames: []string{"-c"},
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.WriterOptions.OutputFileFormat = "csv"
				*pargi += 1
			},
		},

		{
			name: "--csvlite",
			help: "Use CSV-lite format for input and output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csvlite"
				options.WriterOptions.OutputFileFormat = "csvlite"
				*pargi += 1
			},
		},

		{
			name: "--tsv",
			help: "Use TSV format for input and output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.WriterOptions.OutputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.WriterOptions.OFS = "\t"
				options.ReaderOptions.IFSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name:     "--tsvlite",
			help:     "Use TSV-lite format for input and output data.",
			altNames: []string{"-t"},
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csvlite"
				options.WriterOptions.OutputFileFormat = "csvlite"
				options.ReaderOptions.IFS = "\t"
				options.WriterOptions.OFS = "\t"
				options.ReaderOptions.IFSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name:     "--asv",
			altNames: []string{"--asvlite"},
			help:     "Use ASV format for input and output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csvlite"
				options.WriterOptions.OutputFileFormat = "csvlite"
				options.ReaderOptions.IFS = ASV_FS
				options.WriterOptions.OFS = ASV_FS
				options.ReaderOptions.IRS = ASV_RS
				options.WriterOptions.ORS = ASV_RS
				options.ReaderOptions.IFSWasSpecified = true

				options.ReaderOptions.IRSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name:     "--usv",
			altNames: []string{"--usvlite"},
			help:     "Use USV format for input and output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csvlite"
				options.WriterOptions.OutputFileFormat = "csvlite"
				options.ReaderOptions.IFS = USV_FS
				options.WriterOptions.OFS = USV_FS
				options.ReaderOptions.IRS = USV_RS
				options.WriterOptions.ORS = USV_RS
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--dkvp",
			help: "Use DKVP format for input and output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "dkvp"
				*pargi += 1
			},
		},

		{
			name:     "--json",
			help:     "Use JSON format for input and output data.",
			altNames: []string{"-j"},
			parser: func(args []string, argc int, pargi *int, options *TOptions) {

				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "json"
				*pargi += 1
			},
		},

		{
			name: "--nidx",
			help: "Use NIDX format for input and output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "nidx"
				options.ReaderOptions.IFS = " "
				options.WriterOptions.OFS = " "
				options.ReaderOptions.IFSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--xtab",
			help: "Use XTAB format for input and output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "xtab"
				*pargi += 1
			},
		},

		{
			name: "--pprint",
			help: "Use PPRINT format for input and output data.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "pprint"
				options.ReaderOptions.IFS = " "
				options.ReaderOptions.IFSWasSpecified = true
				options.WriterOptions.OutputFileFormat = "pprint"
				*pargi += 1
			},
		},
	},
}

// ================================================================
// FORMAT-CONVERSION KEYSTROKE-SAVER FLAGS

func FormatConversionKeystrokeSaverPrintInfo() {
	fmt.Println(`As keystroke-savers for format-conversion you may use the following.
The letters c, t, j, d, n, x, p, and m refer to formats CSV, TSV, DKVP, NIDX,
JSON, XTAB, PPRINT, and markdown, respectively. Note that markdown format is
available for output only.

| In\out | CSV   | TSV   | JSON   | DKVP   | NIDX   | XTAB   | PPRINT | Markdown |
+--------+-------+-------+--------+--------+--------+--------+--------+----------+
| CSV    |       | --c2t | --c2j  | --c2d  | --c2n  | --c2x  | --c2p  | --c2m    |
| TSV    | --t2c |       | --t2j  | --t2d  | --t2n  | --t2x  | --t2p  | --t2m    |
| JSON   | --j2c | --j2t |        | --j2d  | --j2n  | --j2x  | --j2p  | --j2m    |
| DKVP   | --d2c | --d2t | --d2j  |        | --d2n  | --d2x  | --d2p  | --d2m    |
| NIDX   | --n2c | --n2t | --n2j  | --n2d  |        | --n2x  | --n2p  | --n2m    |
| XTAB   | --x2c | --x2t | --x2j  | --x2d  | --x2n  |        | --x2p  | --x2m    |
| PPRINT | --p2c | --p2t | --p2j  | --p2d  | --p2n  | --p2x  |        | --p2m    |`)
}

func init() { FormatConversionKeystrokeSaverFlagSection.Sort() }

var FormatConversionKeystrokeSaverFlagSection = FlagSection{
	name:        "Format-conversion keystroke-saver flags",
	infoPrinter: FormatConversionKeystrokeSaverPrintInfo,

	flags: []Flag{

		{
			name: "-p",
			help: "Keystroke-saver for `--nidx --fs space --repifs`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "nidx"
				options.ReaderOptions.IFS = " "
				options.WriterOptions.OFS = " "
				options.ReaderOptions.IFSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				options.ReaderOptions.AllowRepeatIFS = true
				options.ReaderOptions.AllowRepeatIFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "-T",
			help: "Keystroke-saver for `--nidx --fs tab`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "nidx"
				options.ReaderOptions.IFS = "\t"
				options.WriterOptions.OFS = "\t"
				options.ReaderOptions.IFSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--c2t",
			help: "Use CSV for input, TSV for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.OFS = "\t"
				options.ReaderOptions.IRSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--c2d",
			help: "Use CSV for input, DKVP for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.WriterOptions.OutputFileFormat = "dkvp"
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--c2n",
			help: "Use CSV for input, NIDX for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.WriterOptions.OutputFileFormat = "nidx"
				options.WriterOptions.OFS = " "
				options.ReaderOptions.IRSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--c2j",
			help: "Use CSV for input, JSON for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.WriterOptions.OutputFileFormat = "json"
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--c2p",
			help: "Use CSV for input, PPRINT for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--c2b",
			help: "Use CSV for input, PPRINT with `--barred` for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.WriterOptions.BarredPprintOutput = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--c2x",
			help: "Use CSV for input, XTAB for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.WriterOptions.OutputFileFormat = "xtab"
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--c2m",
			help: "Use CSV for input, markdown-tabular for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.WriterOptions.OutputFileFormat = "markdown"
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--t2c",
			help: "Use TSV for input, CSV for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.WriterOptions.OutputFileFormat = "csv"
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--t2d",
			help: "Use TSV for input, DKVP for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.WriterOptions.OutputFileFormat = "dkvp"
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--t2n",
			help: "Use TSV for input, NIDX for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.WriterOptions.OutputFileFormat = "nidx"
				options.WriterOptions.OFS = " "
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--t2j",
			help: "Use TSV for input, JSON for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.WriterOptions.OutputFileFormat = "json"
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--t2p",
			help: "Use TSV for input, PPRINT for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--t2b",
			help: "Use TSV for input, PPRINT with `--barred` for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.WriterOptions.BarredPprintOutput = true
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--t2x",
			help: "Use TSV for input, XTAB for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.WriterOptions.OutputFileFormat = "xtab"
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--t2m",
			help: "Use TSV for input, markdown-tabular for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.WriterOptions.OutputFileFormat = "markdown"
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--d2c",
			help: "Use DKVP for input, CSV for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "csv"
				*pargi += 1
			},
		},
		{
			name: "--d2t",
			help: "Use DKVP for input, TSV for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.OFS = "\t"
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--d2n",
			help: "Use DKVP for input, NIDX for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "nidx"
				options.WriterOptions.OFS = " "
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--d2j",
			help: "Use DKVP for input, JSON for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "json"
				*pargi += 1
			},
		},
		{
			name: "--d2p",
			help: "Use DKVP for input, PPRINT for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "pprint"
				*pargi += 1
			},
		},
		{
			name: "--d2b",
			help: "Use DKVP for input, PPRINT with `--barred` for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.WriterOptions.BarredPprintOutput = true
				*pargi += 1
			},
		},
		{
			name: "--d2x",
			help: "Use DKVP for input, XTAB for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "xtab"
				*pargi += 1
			},
		},
		{
			name: "--d2m",
			help: "Use DKVP for input, markdown-tabular for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "markdown"
				*pargi += 1
			},
		},

		{
			name: "--n2c",
			help: "Use NIDX for input, CSV for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--n2t",
			help: "Use NIDX for input, TSV for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.OFS = "\t"
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--n2d",
			help: "Use NIDX for input, DKVP for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "dkvp"
				*pargi += 1
			},
		},
		{
			name: "--n2j",
			help: "Use NIDX for input, JSON for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "json"
				*pargi += 1
			},
		},
		{
			name: "--n2p",
			help: "Use NIDX for input, PPRINT for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "pprint"
				*pargi += 1
			},
		},
		{
			name: "--n2b",
			help: "Use NIDX for input, PPRINT with `--barred` for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.WriterOptions.BarredPprintOutput = true
				*pargi += 1
			},
		},
		{
			name: "--n2x",
			help: "Use NIDX for input, XTAB for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "xtab"
				*pargi += 1
			},
		},
		{
			name: "--n2m",
			help: "Use NIDX for input, markdown-tabular for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "markdown"
				*pargi += 1
			},
		},

		{
			name: "--j2c",
			help: "Use JSON for input, CSV for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--j2t",
			help: "Use JSON for input, TSV for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.OFS = "\t"
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--j2d",
			help: "Use JSON for input, DKVP for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "dkvp"
				*pargi += 1
			},
		},
		{
			name: "--j2n",
			help: "Use JSON for input, NIDX for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "nidx"
				*pargi += 1
			},
		},
		{
			name: "--j2p",
			help: "Use JSON for input, PPRINT for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "pprint"
				*pargi += 1
			},
		},
		{
			name: "--j2b",
			help: "Use JSON for input, PPRINT with --barred for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.WriterOptions.BarredPprintOutput = true
				*pargi += 1
			},
		},
		{
			name: "--j2x",
			help: "Use JSON for input, XTAB for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "xtab"
				*pargi += 1
			},
		},
		{
			name: "--j2m",
			help: "Use JSON for input, markdown-tabular for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "markdown"
				*pargi += 1
			},
		},

		{
			name: "--p2c",
			help: "Use PPRINT for input, CSV for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "pprint"
				options.ReaderOptions.IFS = " "
				options.WriterOptions.OutputFileFormat = "csv"
				options.ReaderOptions.IFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--p2t",
			help: "Use PPRINT for input, TSV for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "pprint"
				options.ReaderOptions.IFS = " "
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.OFS = "\t"
				options.ReaderOptions.IFSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--p2d",
			help: "Use PPRINT for input, DKVP for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "pprint"
				options.ReaderOptions.IFS = " "
				options.WriterOptions.OutputFileFormat = "dkvp"
				options.ReaderOptions.IFSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--p2n",
			help: "Use PPRINT for input, NIDX for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "pprint"
				options.ReaderOptions.IFS = " "
				options.WriterOptions.OutputFileFormat = "nidx"
				options.ReaderOptions.IFSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--p2j",
			help: "Use PPRINT for input, JSON for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "pprint"
				options.ReaderOptions.IFS = " "
				options.WriterOptions.OutputFileFormat = "json"
				options.ReaderOptions.IFSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--p2x",
			help: "Use PPRINT for input, XTAB for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "pprint"
				options.ReaderOptions.IFS = " "
				options.WriterOptions.OutputFileFormat = "xtab"
				options.ReaderOptions.IFSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--p2m",
			help: "Use PPRINT for input, markdown-tabular for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "pprint"
				options.ReaderOptions.IFS = " "
				options.WriterOptions.OutputFileFormat = "markdown"
				options.ReaderOptions.IFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--x2c",
			help: "Use XTAB for input, CSV for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--x2t",
			help: "Use XTAB for input, TSV for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.OFS = "\t"
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--x2d",
			help: "Use XTAB for input, DKVP for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "dkvp"
				*pargi += 1
			},
		},
		{
			name: "--x2n",
			help: "Use XTAB for input, NIDX for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "nidx"
				*pargi += 1
			},
		},
		{
			name: "--x2j",
			help: "Use XTAB for input, JSON for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "json"
				*pargi += 1
			},
		},
		{
			name: "--x2p",
			help: "Use XTAB for input, PPRINT for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "pprint"
				*pargi += 1
			},
		},
		{
			name: "--x2b",
			help: "Use XTAB for input, PPRINT with `--barred` for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.WriterOptions.BarredPprintOutput = true
				*pargi += 1
			},
		},
		{
			name: "--x2m",
			help: "Use XTAB for input, markdown-tabular for output.",
			// For format-conversion keystroke-savers, a matrix is plenty -- we don't
			// need to print a tedious 60-line list.
			suppressFlagEnumeration: true,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "markdown"
				*pargi += 1
			},
		},
	},
}

// ================================================================
// CSV FLAGS

func CSVOnlyPrintInfo() {
	fmt.Println("These are flags which are applicable to CSV format.")
}

func init() { CSVOnlyFlagSection.Sort() }

var CSVOnlyFlagSection = FlagSection{
	name:        "CSV-only flags",
	infoPrinter: CSVOnlyPrintInfo,
	flags: []Flag{

		{
			name: "--no-implicit-csv-header",
			help: "Opposite of `--implicit-csv-header`. This is the default anyway -- the main use is for the flags to `mlr join` if you have main file(s) which are headerless but you want to join in on a file which does have a CSV header. Then you could use `mlr --csv --implicit-csv-header join --no-implicit-csv-header -l your-join-in-with-header.csv ... your-headerless.csv`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.UseImplicitCSVHeader = false
				*pargi += 1
			},
		},

		{
			name:     "--allow-ragged-csv-input",
			altNames: []string{"--ragged"},
			help:     "If a data line has fewer fields than the header line, fill remaining keys with empty string. If a data line has more fields than the header line, use integer field labels as in the implicit-header case.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.AllowRaggedCSVInput = true
				*pargi += 1
			},
		},

		{
			name:     "--implicit-csv-header",
			altNames: []string{"--headerless-csv-input", "--hi"},
			help:     "Use 1,2,3,... as field labels, rather than from line 1 of input files. Tip: combine with `label` to recreate missing headers.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.UseImplicitCSVHeader = true
				*pargi += 1
			},
		},

		{
			name:     "--headerless-csv-output",
			altNames: []string{"--ho"},
			help:     "Print only CSV data lines; do not print CSV header lines.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.HeaderlessCSVOutput = true
				*pargi += 1
			},
		},

		{
			name: "-N",
			help: "Keystroke-saver for `--implicit-csv-header --headerless-csv-output`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.UseImplicitCSVHeader = true
				options.WriterOptions.HeaderlessCSVOutput = true
				*pargi += 1
			},
		},
	},
}

// ================================================================
// COMPRESSED-DATA FLAGS

func CompressedDataPrintInfo() {
	fmt.Print(`Miller offers a few different ways to handle reading data files
	which have been compressed.

* Decompression done within the Miller process itself: ` + "`--bz2in`" + ` ` + "`--gzin`" + ` ` + "`--zin`" + `
* Decompression done outside the Miller process: ` + "`--prepipe`" + ` ` + "`--prepipex`" + `

Using ` + "`--prepipe`" + ` and ` + "`--prepipex`" + ` you can specify an action to be
taken on each input file.  The prepipe command must be able to read from
standard input; it will be invoked with ` + "`{command} < {filename}`" + `.  The
prepipex command must take a filename as argument; it will be invoked with
` + "`{command} {filename}`" + `.

Examples:

    mlr --prepipe gunzip
    mlr --prepipe zcat -cf
    mlr --prepipe xz -cd
    mlr --prepipe cat

Note that this feature is quite general and is not limited to decompression
utilities. You can use it to apply per-file filters of your choice.  For output
compression (or other) utilities, simply pipe the output:
` + "`mlr ... | {your compression command} > outputfilenamegoeshere`" + `

Lastly, note that if ` + "`--prepipe`" + ` or ` + "`--prepipex`" + ` is specified, it replaces any
decisions that might have been made based on the file suffix. Likewise,
` + "`--gzin`" + `/` + "`--bz2in`" + `/` + "`--zin`" + ` are ignored if ` + "`--prepipe`" + ` is also specified.
`)
}

func init() { CompressedDataFlagSection.Sort() }

var CompressedDataFlagSection = FlagSection{
	name:        "Compressed-data flags",
	infoPrinter: CompressedDataPrintInfo,
	flags: []Flag{

		{
			name: "--prepipe",
			arg:  "{decompression command}",
			help: "You can, of course, already do without this for single input files, e.g. `gunzip < myfile.csv.gz | mlr ...`.  Allowed at the command line, but not in `.mlrrc` to avoid unexpected code execution.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.Prepipe = args[*pargi+1]
				options.ReaderOptions.PrepipeIsRaw = false
				*pargi += 2
			},
		},

		{
			name: "--prepipex",
			arg:  "{decompression command}",
			help: "Like `--prepipe` with one exception: doesn't insert `<` between command and filename at runtime. Useful for some commands like `unzip -qc` which don't read standard input.  Allowed at the command line, but not in `.mlrrc` to avoid unexpected code execution.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.Prepipe = args[*pargi+1]
				options.ReaderOptions.PrepipeIsRaw = true
				*pargi += 2
			},
		},

		{
			name: "--prepipe-gunzip",
			help: "Same as  `--prepipe gunzip`, except this is allowed in `.mlrrc`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.Prepipe = "gunzip"
				options.ReaderOptions.PrepipeIsRaw = false
				*pargi += 1
			},
		},

		{
			name: "--prepipe-zcat",
			help: "Same as  `--prepipe zcat`, except this is allowed in `.mlrrc`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.Prepipe = "zcat"
				options.ReaderOptions.PrepipeIsRaw = false
				*pargi += 1
			},
		},

		{
			name: "--prepipe-bz2",
			help: "Same as  `--prepipe bz2`, except this is allowed in `.mlrrc`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.Prepipe = "bz2"
				options.ReaderOptions.PrepipeIsRaw = false
				*pargi += 1
			},
		},

		{
			name: "--gzin",
			help: "Uncompress gzip within the Miller process. Done by default if file ends in `.gz`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.FileInputEncoding = lib.FileInputEncodingGzip
				*pargi += 1
			},
		},

		{
			name: "--zin",
			help: "Uncompress zlib within the Miller process. Done by default if file ends in `.z`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.FileInputEncoding = lib.FileInputEncodingZlib
				*pargi += 1
			},
		},

		{
			name: "--bz2in",
			help: "Uncompress bzip2 within the Miller process. Done by default if file ends in `.bz2`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.FileInputEncoding = lib.FileInputEncodingBzip2
				*pargi += 1
			},
		},
	},
}

// ================================================================
// COMMENTS-IN-DATA FLAGS

func CommentsInDataPrintInfo() {
	fmt.Printf(`Miller lets you put comments in your data, such as

    # This is a comment for a CSV file
    a,b,c
    1,2,3
    4,5,6

Notes:

* Comments are only honored at the start of a line.
* In the absence of any of the below four options, comments are data like
  any other text. (The comments-in-data feature is opt-in.)
* When ` + "`--pass-comments`" + ` is used, comment lines are written to standard output
  immediately upon being read; they are not part of the record stream.  Results
  may be counterintuitive. A suggestion is to place comments at the start of
  data files.
`)
}

func init() { CommentsInDataFlagSection.Sort() }

var CommentsInDataFlagSection = FlagSection{
	name:        "Comments-in-data flags",
	infoPrinter: CommentsInDataPrintInfo,
	flags: []Flag{

		{
			name: "--skip-comments",
			help: "Ignore commented lines (prefixed by `" + DEFAULT_COMMENT_STRING + "`) within the input.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.CommentString = DEFAULT_COMMENT_STRING
				options.ReaderOptions.CommentHandling = SkipComments
				*pargi += 1
			},
		},

		{
			name: "--skip-comments-with",
			arg:  "{string}",
			help: "Ignore commented lines within input, with specified prefix.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.CommentString = args[*pargi+1]
				options.ReaderOptions.CommentHandling = SkipComments
				*pargi += 2
			},
		},

		{
			name: "--pass-comments",
			help: "Immediately print commented lines (prefixed by `" + DEFAULT_COMMENT_STRING + "`) within the input.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.CommentString = DEFAULT_COMMENT_STRING
				options.ReaderOptions.CommentHandling = PassComments
				*pargi += 1
			},
		},

		{
			name: "--pass-comments-with",
			arg:  "{string}",
			help: "Immediately print commented lines within input, with specified prefix.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.CommentString = args[*pargi+1]
				options.ReaderOptions.CommentHandling = PassComments
				*pargi += 2
			},
		},
	},
}

// ================================================================
// OUTPUT-COLORIZATION FLAGS

func OutputColorizationPrintInfo() {
	fmt.Print(`Miller uses colors to highlight outputs. You can specify color preferences.
Note: output colorization does not work on Windows.

Things having colors:

* Keys in CSV header lines, JSON keys, etc
* Values in CSV data lines, JSON scalar values, etc in regression-test output
* Some online-help strings

Rules for coloring:

* By default, colorize output only if writing to stdout and stdout is a TTY.
    * Example: color: ` + "`mlr --csv cat foo.csv`" + `
    * Example: no color: ` + "`mlr --csv cat foo.csv > bar.csv`" + `
    * Example: no color: ` + "`mlr --csv cat foo.csv | less`" + `
* The default colors were chosen since they look OK with white or black
  terminal background, and are differentiable with common varieties of human
  color vision.

Mechanisms for coloring:

* Miller uses ANSI escape sequences only. This does not work on Windows
  except within Cygwin.
* Requires ` + "`TERM`" + ` environment variable to be set to non-empty string.
* Doesn't try to check to see whether the terminal is capable of 256-color
  ANSI vs 16-color ANSI. Note that if colors are in the range 0..15
  then 16-color ANSI escapes are used, so this is in the user's control.

How you can control colorization:

* Suppression/unsuppression:
    * Environment variable ` + "`export MLR_NO_COLOR=true`" + ` means don't color
      even if stdout+TTY.
    * Environment variable ` + "`export MLR_ALWAYS_COLOR=true`" + ` means do color
      even if not stdout+TTY.
      For example, you might want to use this when piping mlr output to ` + "`less -r`" + `.
    * Command-line flags ` + "`--no-color`" + ` or ` + "`-M`" + `, ` + "`--always-color`" + ` or ` + "`-C`" + `.

* Color choices can be specified by using environment variables, or command-line
  flags, with values 0..255:
    * ` + "`export MLR_KEY_COLOR=208`" + `, ` + "`MLR_VALUE_COLOR=33`" + `, etc.:
        ` + "`MLR_KEY_COLOR`" + ` ` + "`MLR_VALUE_COLOR`" + ` ` + "`MLR_PASS_COLOR`" + ` ` + "`MLR_FAIL_COLOR`" + `
        ` + "`MLR_REPL_PS1_COLOR`" + ` ` + "`MLR_REPL_PS2_COLOR`" + ` ` + "`MLR_HELP_COLOR`" + `
    * Command-line flags ` + "`--key-color 208`" + `, ` + "`--value-color 33`" + `, etc.:
        ` + "`--key-color`" + ` ` + "`--value-color`" + ` ` + "`--pass-color`" + ` ` + "`--fail-color`" + `
        ` + "`--repl-ps1-color`" + ` ` + "`--repl-ps2-color`" + ` ` + "`--help-color`" + `
    * This is particularly useful if your terminal's background color clashes
      with current settings.

If environment-variable settings and command-line flags are both provided, the
latter take precedence.

Please do mlr ` + "`--list-color-codes`" + ` to see the available color codes (like 170),
and ` + "`mlr --list-color-names`" + ` to see available names (like ` + "`orchid`" + `).
`)
}

func init() { OutputColorizationFlagSection.Sort() }

var OutputColorizationFlagSection = FlagSection{
	name:        "Output-colorization flags",
	infoPrinter: OutputColorizationPrintInfo,
	flags: []Flag{

		{
			name: "--list-color-codes",
			help: "Show the available color codes in the range 0..255, such as 170 for example.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				colorizer.ListColorCodes()
				os.Exit(0)
				*pargi += 1
			},
		},

		{
			name: "--list-color-names",
			help: "Show the names for the available color codes, such as `orchid` for example.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				colorizer.ListColorNames()
				os.Exit(0)
				*pargi += 1
			},
		},

		{
			name:     "--no-color",
			altNames: []string{"-M"},
			help:     "Instructs Miller to not colorize any output.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				colorizer.SetColorization(colorizer.ColorizeOutputNever)
				*pargi += 1
			},
		},

		{
			name:     "--always-color",
			altNames: []string{"-C"},
			help:     "Instructs Miller to colorize output even when it normally would not. Useful for piping output to `less -r`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				colorizer.SetColorization(colorizer.ColorizeOutputAlways)
				*pargi += 1
			},
		},

		{
			name: "--key-color",
			help: "Specify the color (see `--list-color-codes` and `--list-color-names`) for record keys.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				ok := colorizer.SetKeyColor(args[*pargi+1])
				if !ok {
					fmt.Fprintf(os.Stderr,
						"mlr: --key-color argument unrecognized; got \"%s\".\n",
						args[*pargi+1])
					os.Exit(1)
				}
				*pargi += 2
			},
		},

		{
			name: "--value-color",
			help: "Specify the color (see `--list-color-codes` and `--list-color-names`) for record values.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				ok := colorizer.SetValueColor(args[*pargi+1])
				if !ok {
					fmt.Fprintf(os.Stderr,
						"mlr: --value-color argument unrecognized; got \"%s\".\n",
						args[*pargi+1])
					os.Exit(1)
				}
				*pargi += 2
			},
		},

		{
			name: "--pass-color",
			help: "Specify the color (see `--list-color-codes` and `--list-color-names`) for passing cases in `mlr regtest`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				ok := colorizer.SetPassColor(args[*pargi+1])
				if !ok {
					fmt.Fprintf(os.Stderr,
						"mlr: --pass-color argument unrecognized; got \"%s\".\n",
						args[*pargi+1])
					os.Exit(1)
				}
				*pargi += 2
			},
		},

		{
			name: "--fail-color",
			help: "Specify the color (see `--list-color-codes` and `--list-color-names`) for failing cases in `mlr regtest`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				ok := colorizer.SetFailColor(args[*pargi+1])
				if !ok {
					fmt.Fprintf(os.Stderr,
						"mlr: --fail-color argument unrecognized; got \"%s\".\n",
						args[*pargi+1])
					os.Exit(1)
				}
				*pargi += 2
			},
		},

		{
			name: "--help-color",
			help: "Specify the color (see `--list-color-codes` and `--list-color-names`) for highlights in `mlr help` output.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				ok := colorizer.SetHelpColor(args[*pargi+1])
				if !ok {
					fmt.Fprintf(os.Stderr,
						"mlr: --help-color argument unrecognized; got \"%s\".\n",
						args[*pargi+1])
					os.Exit(1)
				}
				*pargi += 2
			},
		},
	},
}

// ================================================================
// FLATTEN/UNFLATTEN FLAGS

func FlattenUnflattenPrintInfo() {
	fmt.Println("These flags control how Miller converts record values which are maps or arrays, when input is JSON and ouput is non-JSON (flattening) or input is non-JSON and output is JSON (unflattening).")
	fmt.Println()
	fmt.Println("See the Flatten/unflatten doc page for more information.")
}

func init() { FlattenUnflattenFlagSection.Sort() }

var FlattenUnflattenFlagSection = FlagSection{
	name:        "Flatten-unflatten flags",
	infoPrinter: FlattenUnflattenPrintInfo,
	flags: []Flag{

		{
			name:     "--flatsep",
			altNames: []string{"--jflatsep"},
			arg:      "{string}",
			help:     "Separator for flattening multi-level JSON keys, e.g. `{\"a\":{\"b\":3}}` becomes `a:b => 3` for non-JSON formats. Defaults to `" + DEFAULT_JSON_FLATTEN_SEPARATOR + "`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.WriterOptions.FLATSEP = SeparatorFromArg(args[*pargi+1])
				*pargi += 2
			},
		},

		{
			name: "--xvright",
			help: "Right-justify values for XTAB format.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.RightAlignedXTABOutput = true
				*pargi += 1
			},
		},

		{
			name: "--no-auto-flatten",
			help: "When output is non-JSON, suppress the default auto-flatten behavior. Default: if `$y = [7,8,9]` then this flattens to `y.1=7,y.2=8,y.3=9, and similarly for maps. With `--no-auto-flatten`, instead we get `$y=[1, 2, 3]`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.AutoFlatten = false
				*pargi += 1
			},
		},

		{
			name: "--no-auto-unflatten",
			help: "When input non-JSON and output is JSON, suppress the default auto-unflatten behavior. Default: if the input has `y.1=7,y.2=8,y.3=9` then this unflattens to `$y=[7,8,9]`.  flattens to `y.1=7,y.2=8,y.3=9. With `--no-auto-flatten`, instead we get `${y.1}=7,${y.2}=8,${y.3}=9`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.AutoUnflatten = false
				*pargi += 1
			},
		},
	},
}

// ================================================================
// MISC FLAGS

func MiscPrintInfo() {
	fmt.Print("These are flags which don't fit into any other category.")
}

func init() { MiscFlagSection.Sort() }

var MiscFlagSection = FlagSection{
	name:        "Miscellaneous flags",
	infoPrinter: MiscPrintInfo,
	flags: []Flag{

		{
			name: "-n",
			help: "Process no input files, nor standard input either. Useful for `mlr put` with `begin`/`end` statements only. (Same as `--from /dev/null`.) Also useful in `mlr -n put -v '...'` for analyzing abstract syntax trees (if that's your thing).",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.NoInput = true
				*pargi += 1
			},
		},

		{
			name: "-I",
			help: "Process files in-place. For each file name on the command line, output is written to a temp file in the same directory, which is then renamed over the original. Each file is processed in isolation: if the output format is CSV, CSV headers will be present in each output file, statistics are only over each file's own records; and so on.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.DoInPlace = true
				*pargi += 1
			},
		},

		{
			name: "--from",
			arg:  "{filename}",
			help: "Use this to specify an input file before the verb(s), rather than after. May be used more than once. Example: `mlr --from a.dat --from b.dat cat` is the same as `mlr cat a.dat b.dat`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.FileNames = append(options.FileNames, args[*pargi+1])
				*pargi += 2
			},
		},

		{
			name: "--mfrom",
			arg:  "{filenames}",
			help: "Use this to specify one of more input files before the verb(s), rather than after. May be used more than once.  The list of filename must end with `--`. This is useful for example since `--from *.csv` doesn't do what you might hope but `--mfrom *.csv --` does.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				*pargi += 1
				for *pargi < argc && args[*pargi] != "--" {
					options.FileNames = append(options.FileNames, args[*pargi])
					*pargi += 1
				}
				if args[*pargi] == "--" {
					*pargi += 1
				}
			},
		},

		{
			name: "--ofmt",
			arg:  "{format}",
			help: "E.g. `%.18f`, `%.0f`, `%9.6e`. Please use sprintf-style codes for floating-point nummbers. If not specified, default formatting is used.  See also the `fmtnum` function and the `format-values` verb.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.WriterOptions.FPOFMT = args[*pargi+1]
				*pargi += 2
			},
		},

		{
			name: "--load",
			arg:  "{filename}",
			help: "Load DSL script file for all put/filter operations on the command line.  If the name following `--load` is a directory, load all `*.mlr` files in that directory. This is just like `put -f` and `filter -f` except it's up-front on the command line, so you can do something like `alias mlr='mlr --load ~/myscripts'` if you like.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.DSLPreloadFileNames = append(options.DSLPreloadFileNames, args[*pargi+1])
				*pargi += 2
			},
		},

		{
			name: "--mload",
			arg:  "{filenames}",
			help: "Like `--load` but works with more than one filename, e.g. `--mload *.mlr --`.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				*pargi += 1
				for *pargi < argc && args[*pargi] != "--" {
					options.DSLPreloadFileNames = append(options.DSLPreloadFileNames, args[*pargi])
					*pargi += 1
				}
				if args[*pargi] == "--" {
					*pargi += 1
				}
			},
		},

		{
			name: "--tz",
			arg:  "{timezone}",
			help: "Specify timezone, overriding `$TZ` environment variable (if any).",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				os.Setenv("TZ", args[*pargi+1])
				*pargi += 2
			},
		},

		{
			name: "--nr-progress-mod",
			arg:  "{m}",
			help: "With m a positive integer: print filename and record count to os.Stderr every m input records.",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				nrProgressMod, ok := lib.TryIntFromString(args[*pargi+1])
				if !ok || nrProgressMod <= 0 {
					fmt.Fprintf(os.Stderr,
						"%s: --nr-progress-mod argument must be a positive integer; got \"%s\".\n",

						"mlr", args[*pargi+1])
					os.Exit(1)
				}
				options.NRProgressMod = nrProgressMod
				*pargi += 2
			},
		},

		{
			name: "--seed",
			arg:  "{n}",
			help: "with `n` of the form `12345678` or `0xcafefeed`. For `put`/`filter` `urand`, `urandint`, and `urand32`.",

			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				randSeed, ok := lib.TryIntFromString(args[*pargi+1])
				if ok {
					options.RandSeed = randSeed
					options.HaveRandSeed = true
				} else {
					fmt.Fprintf(os.Stderr,
						"mlr: --seed argument must be a decimal or hexadecimal integer; got \"%s\".\n",
						args[*pargi+1])
					fmt.Fprintf(os.Stderr, "Please run \"mlr --help\" for detailed usage information.\n")
					os.Exit(1)
				}
				*pargi += 2
			},
		},

		{
			name:     "--infer-none",
			altNames: []string{"-S"},
			help:     `Don't treat values like 123 or 456.7 in data files as int/float; leave them as strings.`,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				types.SetInferrerStringOnly()
				*pargi += 1
			},
		},

		{
			name:     "--infer-int-as-float",
			altNames: []string{"-A"},
			help:     `Cast all integers in data files to floats.`,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				types.SetInferrerIntAsFloat()
				*pargi += 1
			},
		},

		{
			name:     "--infer-no-octal",
			altNames: []string{"-O"},
			help:     `Treat numbers like 0123 in data files as string "0123", not octal for decimal 83 etc.`,
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				types.SetInferrerNoOctal()
				*pargi += 1
			},
		},
	},
}
