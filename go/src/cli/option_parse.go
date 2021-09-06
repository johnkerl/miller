// ================================================================
// Items which might better belong in miller/cli, but which are placed in a
// deeper package to avoid a package-dependency cycle between miller/cli and
// miller/transforming.
// ================================================================

package cli

import (
	"fmt"
	"os"

	"mlr/src/colorizer"
	"mlr/src/lib"
)

const ASV_FS = "\x1f"
const ASV_RS = "\x1e"
const USV_FS = "\xe2\x90\x9f"
const USV_RS = "\xe2\x90\x9e"

const ASV_FS_FOR_HELP = "0x1f"
const ASV_RS_FOR_HELP = "0x1e"
const USV_FS_FOR_HELP = "U+241F (UTF-8 0xe2909f)"
const USV_RS_FOR_HELP = "U+241E (UTF-8 0xe2909e)"
const DEFAULT_JSON_FLATTEN_SEPARATOR = "."

// ----------------------------------------------------------------
// TODO: move these to their own file

// E.g. if IFS isn't specified, it's space for NIDX and comma for DKVP, etc.

var defaultFSes = map[string]string{
	// "gen" : // TODO
	"csv":      ",",
	"csvlite":  ",",
	"dkvp":     ",",
	"json":     ",", // not honored; not parameterizable in JSON format
	"nidx":     " ",
	"markdown": " ",
	"pprint":   " ",
	"xtab":     "\n", // todo: windows-dependent ...
}

var defaultPSes = map[string]string{
	"csv":      "=",
	"csvlite":  "=",
	"dkvp":     "=",
	"json":     "=", // not honored; not parameterizable in JSON format
	"markdown": "=",
	"nidx":     "=",
	"pprint":   "=",
	"xtab":     " ", // todo: windows-dependent ...
}

var defaultRSes = map[string]string{
	"csv":      "\n",
	"csvlite":  "\n",
	"dkvp":     "\n",
	"json":     "\n", // not honored; not parameterizable in JSON format
	"markdown": "\n",
	"nidx":     "\n",
	"pprint":   "\n",
	"xtab":     "\n\n", // todo: maybe jettison the idea of this being alterable
}

var defaultAllowRepeatIFSes = map[string]bool{
	"csv":      false,
	"csvlite":  false,
	"dkvp":     false,
	"json":     false,
	"markdown": false,
	"nidx":     false,
	"pprint":   true,
	"xtab":     false,
}

var defaultAllowRepeatIPSes = map[string]bool{
	"csv":      false,
	"csvlite":  false,
	"dkvp":     false,
	"json":     false,
	"markdown": false,
	"nidx":     false,
	"pprint":   false,
	"xtab":     true,
}

func ApplyReaderOptionDefaults(readerOptions *TReaderOptions) {
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
		readerOptions.AllowRepeatIFS = defaultAllowRepeatIFSes[readerOptions.InputFileFormat]
	}
	if !readerOptions.AllowRepeatIPSWasSpecified {
		readerOptions.AllowRepeatIPS = defaultAllowRepeatIPSes[readerOptions.InputFileFormat]
	}
}

func ApplyWriterOptionDefaults(writerOptions *TWriterOptions) {
	if !writerOptions.OFSWasSpecified {
		writerOptions.OFS = defaultFSes[writerOptions.OutputFileFormat]
	}
	if !writerOptions.OPSWasSpecified {
		writerOptions.OPS = defaultPSes[writerOptions.OutputFileFormat]
	}
	if !writerOptions.ORSWasSpecified {
		writerOptions.ORS = defaultRSes[writerOptions.OutputFileFormat]
	}
}

// ----------------------------------------------------------------
// Returns true if the current flag was handled. Exported for use by join.
func ParseReaderOptions(
	args []string,
	argc int,
	pargi *int,
	readerOptions *TReaderOptions,
) bool {
	argi := *pargi
	oargi := argi

	if args[argi] == "--ifs" {
		CheckArgCount(args, argi, argc, 2)
		readerOptions.IFS = SeparatorFromArg(args[argi+1])
		readerOptions.IFSWasSpecified = true
		argi += 2

	} else if args[argi] == "--ips" {
		CheckArgCount(args, argi, argc, 2)
		readerOptions.IPS = SeparatorFromArg(args[argi+1])
		readerOptions.IPSWasSpecified = true
		argi += 2

	} else if args[argi] == "--irs" {
		CheckArgCount(args, argi, argc, 2)
		readerOptions.IRS = SeparatorFromArg(args[argi+1])
		readerOptions.IRSWasSpecified = true
		argi += 2

	} else if args[argi] == "--repifs" {
		readerOptions.AllowRepeatIFS = true
		readerOptions.AllowRepeatIFSWasSpecified = true
		argi += 1

	} else if args[argi] == "--repips" {
		readerOptions.AllowRepeatIPS = true
		readerOptions.AllowRepeatIPSWasSpecified = true
		argi += 1

	} else if args[argi] == "--json-fatal-arrays-on-input" {
		// No-op pass-through for backward compatibility with Miller 5
		argi += 1
	} else if args[argi] == "--json-skip-arrays-on-input" {
		// No-op pass-through for backward compatibility with Miller 5
		argi += 1
	} else if args[argi] == "--json-map-arrays-on-input" {
		// No-op pass-through for backward compatibility with Miller 5
		argi += 1
	} else if args[argi] == "--implicit-csv-header" {
		readerOptions.UseImplicitCSVHeader = true
		argi += 1

	} else if args[argi] == "--no-implicit-csv-header" {
		readerOptions.UseImplicitCSVHeader = false
		argi += 1

	} else if args[argi] == "--allow-ragged-csv-input" || args[argi] == "--ragged" {
		readerOptions.AllowRaggedCSVInput = true
		argi += 1

	} else if args[argi] == "-i" {
		CheckArgCount(args, argi, argc, 2)
		readerOptions.InputFileFormat = args[argi+1]
		argi += 2

		//	} else if args[argi] == "--igen" {
		//		readerOptions.InputFileFormat = "gen";
		//		argi += 1;
		//	} else if args[argi] == "--gen-start" {
		//		readerOptions.InputFileFormat = "gen";
		//		CheckArgCount(args, argi, argc, 2);
		//		if (sscanf(args[argi+1], "%lld", &readerOptions.generator_opts.start) != 1) {
		//			fmt.Fprintf(os.Stderr, "%s: could not scan \"%s\".\n",
		//				"mlr", args[argi+1]);
		//		}
		//		argi += 2;
		//	} else if args[argi] == "--gen-stop" {
		//		readerOptions.InputFileFormat = "gen";
		//		CheckArgCount(args, argi, argc, 2);
		//		if (sscanf(args[argi+1], "%lld", &readerOptions.generator_opts.stop) != 1) {
		//			fmt.Fprintf(os.Stderr, "%s: could not scan \"%s\".\n",
		//				"mlr", args[argi+1]);
		//		}
		//		argi += 2;
		//	} else if args[argi] == "--gen-step" {
		//		readerOptions.InputFileFormat = "gen";
		//		CheckArgCount(args, argi, argc, 2);
		//		if (sscanf(args[argi+1], "%lld", &readerOptions.generator_opts.step) != 1) {
		//			fmt.Fprintf(os.Stderr, "%s: could not scan \"%s\".\n",
		//				"mlr", args[argi+1]);
		//		}
		//		argi += 2;

	} else if args[argi] == "--icsv" {
		readerOptions.InputFileFormat = "csv"
		argi += 1

	} else if args[argi] == "--icsvlite" {
		readerOptions.InputFileFormat = "csvlite"
		argi += 1

	} else if args[argi] == "--itsv" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IFS = "\t"
		readerOptions.IFSWasSpecified = true
		argi += 1
	} else if args[argi] == "--itsvlite" {
		readerOptions.InputFileFormat = "csvlite"
		readerOptions.IFS = "\t"
		readerOptions.IFSWasSpecified = true
		argi += 1

	} else if args[argi] == "--iasv" {
		readerOptions.InputFileFormat = "csvlite"
		readerOptions.IFS = ASV_FS
		readerOptions.IRS = ASV_RS
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		argi += 1

	} else if args[argi] == "--iasvlite" {
		readerOptions.InputFileFormat = "csvlite"
		readerOptions.IFS = ASV_FS
		readerOptions.IRS = ASV_RS
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		argi += 1

	} else if args[argi] == "--iusv" {
		readerOptions.InputFileFormat = "csvlite"
		readerOptions.IFS = USV_FS
		readerOptions.IRS = USV_RS
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		argi += 1

	} else if args[argi] == "--iusvlite" {
		readerOptions.InputFileFormat = "csvlite"
		readerOptions.IFS = USV_FS
		readerOptions.IRS = USV_RS
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		argi += 1

	} else if args[argi] == "--idkvp" {
		readerOptions.InputFileFormat = "dkvp"
		argi += 1

	} else if args[argi] == "--ijson" {
		readerOptions.InputFileFormat = "json"
		argi += 1

	} else if args[argi] == "--inidx" {
		readerOptions.InputFileFormat = "nidx"
		argi += 1

	} else if args[argi] == "--ixtab" {
		readerOptions.InputFileFormat = "xtab"
		argi += 1

	} else if args[argi] == "--ipprint" {
		readerOptions.InputFileFormat = "pprint"
		readerOptions.IFS = " "
		readerOptions.IFSWasSpecified = true
		argi += 1

	} else if args[argi] == "--mmap" {
		// No-op as of 5.6.3 (mmap is being abandoned) but don't break
		// the command-line user experience.
		argi += 1

	} else if args[argi] == "--no-mmap" {
		// No-op as of 5.6.3 (mmap is being abandoned) but don't break
		// the command-line user experience.
		argi += 1

	} else if args[argi] == "--prepipe" {
		CheckArgCount(args, argi, argc, 2)
		readerOptions.Prepipe = args[argi+1]
		readerOptions.PrepipeIsRaw = false
		argi += 2

	} else if args[argi] == "--prepipex" {
		CheckArgCount(args, argi, argc, 2)
		readerOptions.Prepipe = args[argi+1]
		readerOptions.PrepipeIsRaw = true
		argi += 2

	} else if args[argi] == "--prepipe-gunzip" {
		readerOptions.Prepipe = "gunzip"
		readerOptions.PrepipeIsRaw = false
		argi += 1

	} else if args[argi] == "--prepipe-zcat" {
		readerOptions.Prepipe = "zcat"
		readerOptions.PrepipeIsRaw = false
		argi += 1

	} else if args[argi] == "--prepipe-bz2" {
		readerOptions.Prepipe = "bz2"
		readerOptions.PrepipeIsRaw = false
		argi += 1

	} else if args[argi] == "--gzin" {
		readerOptions.FileInputEncoding = lib.FileInputEncodingGzip
		argi += 1

	} else if args[argi] == "--zin" {
		readerOptions.FileInputEncoding = lib.FileInputEncodingZlib
		argi += 1

	} else if args[argi] == "--bz2in" {
		readerOptions.FileInputEncoding = lib.FileInputEncodingBzip2
		argi += 1

	} else if args[argi] == "--skip-comments" {
		readerOptions.CommentString = DEFAULT_COMMENT_STRING
		readerOptions.CommentHandling = SkipComments
		argi += 1

	} else if args[argi] == "--skip-comments-with" {
		CheckArgCount(args, argi, argc, 2)
		readerOptions.CommentString = args[argi+1]
		readerOptions.CommentHandling = SkipComments
		argi += 2

	} else if args[argi] == "--pass-comments" {
		readerOptions.CommentString = DEFAULT_COMMENT_STRING
		readerOptions.CommentHandling = PassComments
		argi += 1

	} else if args[argi] == "--pass-comments-with" {
		CheckArgCount(args, argi, argc, 2)
		readerOptions.CommentString = args[argi+1]
		readerOptions.CommentHandling = PassComments
		argi += 2

	}
	*pargi = argi
	return argi != oargi
}

// Returns true if the current flag was handled.
func ParseWriterOptions(
	args []string,
	argc int,
	pargi *int,
	writerOptions *TWriterOptions,
) bool {
	argi := *pargi
	oargi := argi

	if args[argi] == "--ors" {
		CheckArgCount(args, argi, argc, 2)
		writerOptions.ORS = SeparatorFromArg(args[argi+1])
		writerOptions.ORSWasSpecified = true
		argi += 2

	} else if args[argi] == "--ofs" {
		CheckArgCount(args, argi, argc, 2)
		writerOptions.OFS = SeparatorFromArg(args[argi+1])
		writerOptions.OFSWasSpecified = true
		argi += 2

	} else if args[argi] == "--headerless-csv-output" {
		writerOptions.HeaderlessCSVOutput = true
		argi += 1
	} else if args[argi] == "--ops" {
		CheckArgCount(args, argi, argc, 2)
		writerOptions.OPS = SeparatorFromArg(args[argi+1])
		writerOptions.OPSWasSpecified = true
		argi += 2

	} else if args[argi] == "--flatsep" || args[argi] == "--jflatsep" || args[argi] == "--oflatsep" {
		CheckArgCount(args, argi, argc, 2)
		writerOptions.FLATSEP = SeparatorFromArg(args[argi+1])
		argi += 2

		//	} else if args[argi] == "--xvright" {
		//		writerOptions.right_justify_xtab_value = true;
		//		argi += 1;
		//

	} else if args[argi] == "--jvstack" {
		writerOptions.JSONOutputMultiline = true
		argi += 1

	} else if args[argi] == "--no-jvstack" {
		writerOptions.JSONOutputMultiline = false
		argi += 1

	} else if args[argi] == "--jlistwrap" {
		writerOptions.WrapJSONOutputInOuterList = true
		argi += 1

	} else if args[argi] == "--no-auto-flatten" {
		writerOptions.AutoFlatten = false
		argi += 1

	} else if args[argi] == "--no-auto-unflatten" {
		writerOptions.AutoUnflatten = false
		argi += 1

	} else if args[argi] == "--ofmt" {
		CheckArgCount(args, argi, argc, 2)
		writerOptions.FPOFMT = args[argi+1]
		argi += 2

	} else if args[argi] == "--jknquoteint" {
		// No-op pass-through for backward compatibility with Miller 5
		argi += 1
	} else if args[argi] == "--jquoteall" {
		// No-op pass-through for backward compatibility with Miller 5
		argi += 1
	} else if args[argi] == "--jvquoteall" {
		// No-op pass-through for backward compatibility with Miller 5
		argi += 1

	} else if args[argi] == "--vflatsep" {
		CheckArgCount(args, argi, argc, 2)
		// No-op pass-through for backward compatibility with Miller 5
		argi += 2

	} else if args[argi] == "-o" {
		CheckArgCount(args, argi, argc, 2)
		writerOptions.OutputFileFormat = args[argi+1]
		argi += 2

	} else if args[argi] == "--ocsv" {
		writerOptions.OutputFileFormat = "csv"
		argi += 1

	} else if args[argi] == "--ocsvlite" {
		writerOptions.OutputFileFormat = "csvlite"
		argi += 1

	} else if args[argi] == "--otsv" {
		writerOptions.OutputFileFormat = "csv"
		writerOptions.OFS = "\t"
		writerOptions.OFSWasSpecified = true
		argi += 1

	} else if args[argi] == "--otsvlite" {
		writerOptions.OutputFileFormat = "csvlite"
		writerOptions.OFS = "\t"
		writerOptions.OFSWasSpecified = true
		argi += 1

	} else if args[argi] == "--oasv" {
		writerOptions.OutputFileFormat = "csvlite"
		writerOptions.OFS = ASV_FS
		writerOptions.ORS = ASV_RS
		writerOptions.OFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1

	} else if args[argi] == "--oasvlite" {
		writerOptions.OutputFileFormat = "csvlite"
		writerOptions.OFS = ASV_FS
		writerOptions.ORS = ASV_RS
		writerOptions.OFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1

	} else if args[argi] == "--ousv" {
		writerOptions.OutputFileFormat = "csvlite"
		writerOptions.OFS = USV_FS
		writerOptions.ORS = USV_RS
		writerOptions.OFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1

	} else if args[argi] == "--ousvlite" {
		writerOptions.OutputFileFormat = "csvlite"
		writerOptions.OFS = USV_FS
		writerOptions.ORS = USV_RS
		writerOptions.OFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1

	} else if args[argi] == "--omd" {
		writerOptions.OutputFileFormat = "markdown"
		argi += 1

	} else if args[argi] == "--odkvp" {
		writerOptions.OutputFileFormat = "dkvp"
		argi += 1

	} else if args[argi] == "--ojson" {
		writerOptions.OutputFileFormat = "json"
		argi += 1
	} else if args[argi] == "--ojsonx" {
		// --jvstack is now the default in Miller 6 so this is just for backward compatibility
		writerOptions.OutputFileFormat = "json"
		argi += 1

	} else if args[argi] == "--onidx" {
		writerOptions.OutputFileFormat = "nidx"
		writerOptions.OFS = " "
		writerOptions.OFSWasSpecified = true
		argi += 1

	} else if args[argi] == "--oxtab" {
		writerOptions.OutputFileFormat = "xtab"
		argi += 1

	} else if args[argi] == "--opprint" {
		writerOptions.OutputFileFormat = "pprint"
		argi += 1

	} else if args[argi] == "--right" {
		writerOptions.RightAlignedPprintOutput = true
		argi += 1

	} else if args[argi] == "--barred" {
		writerOptions.BarredPprintOutput = true
		argi += 1
		//
		//	} else if args[argi] == "--quote-all" {
		//		writerOptions.oquoting = QUOTE_ALL;
		//		argi += 1;
		//
		//	} else if args[argi] == "--quote-none" {
		//		writerOptions.oquoting = QUOTE_NONE;
		//		argi += 1;
		//
		//	} else if args[argi] == "--quote-minimal" {
		//		writerOptions.oquoting = QUOTE_MINIMAL;
		//		argi += 1;
		//
		//	} else if args[argi] == "--quote-numeric" {
		//		writerOptions.oquoting = QUOTE_NUMERIC;
		//		argi += 1;
		//
		//	} else if args[argi] == "--quote-original" {
		//		writerOptions.oquoting = QUOTE_ORIGINAL;
		//		argi += 1;
		//

	} else if args[argi] == "--no-fflush" {
		// No-op for Miller 6; accepted at the command line for backward compatibility.
		argi += 1

	} else if args[argi] == "--list-color-codes" {
		colorizer.ListColorCodes()
		os.Exit(0)
		argi += 1

	} else if args[argi] == "--list-color-names" {
		colorizer.ListColorNames()
		os.Exit(0)
		argi += 1

	} else if args[argi] == "--no-color" || args[argi] == "-M" {
		colorizer.SetColorization(colorizer.ColorizeOutputNever)
		argi += 1

	} else if args[argi] == "--always-color" || args[argi] == "-C" {
		colorizer.SetColorization(colorizer.ColorizeOutputAlways)
		argi += 1

	} else if args[argi] == "--key-color" {
		CheckArgCount(args, argi, argc, 2)
		ok := colorizer.SetKeyColor(args[argi+1])
		if !ok {
			fmt.Fprintf(os.Stderr,
				"%s: --key-color argument unrecognized; got \"%s\".\n",
				"mlr", args[argi+1])
			os.Exit(1)
		}
		argi += 2

	} else if args[argi] == "--value-color" {
		CheckArgCount(args, argi, argc, 2)
		ok := colorizer.SetValueColor(args[argi+1])
		if !ok {
			fmt.Fprintf(os.Stderr,
				"%s: --value-color argument unrecognized; got \"%s\".\n",
				"mlr", args[argi+1])
			os.Exit(1)
		}
		argi += 2

	} else if args[argi] == "--pass-color" {
		CheckArgCount(args, argi, argc, 2)
		ok := colorizer.SetPassColor(args[argi+1])
		if !ok {
			fmt.Fprintf(os.Stderr,
				"%s: --pass-color argument unrecognized; got \"%s\".\n",
				"mlr", args[argi+1])
			os.Exit(1)
		}
		argi += 2

	} else if args[argi] == "--fail-color" {
		CheckArgCount(args, argi, argc, 2)
		ok := colorizer.SetFailColor(args[argi+1])
		if !ok {
			fmt.Fprintf(os.Stderr,
				"%s: --fail-color argument unrecognized; got \"%s\".\n",
				"mlr", args[argi+1])
			os.Exit(1)
		}
		argi += 2

	} else if args[argi] == "--help-color" {
		CheckArgCount(args, argi, argc, 2)
		ok := colorizer.SetHelpColor(args[argi+1])
		if !ok {
			fmt.Fprintf(os.Stderr,
				"%s: --help-color argument unrecognized; got \"%s\".\n",
				"mlr", args[argi+1])
			os.Exit(1)
		}
		argi += 2

	}
	*pargi = argi
	return argi != oargi
}

// Returns true if the current flag was handled.
func ParseReaderWriterOptions(
	args []string,
	argc int,
	pargi *int,
	readerOptions *TReaderOptions,
	writerOptions *TWriterOptions,
) bool {
	argi := *pargi
	oargi := argi

	if args[argi] == "--rs" {
		CheckArgCount(args, argi, argc, 2)
		readerOptions.IRS = SeparatorFromArg(args[argi+1])
		writerOptions.ORS = SeparatorFromArg(args[argi+1])
		readerOptions.IRSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 2

	} else if args[argi] == "--fs" {
		CheckArgCount(args, argi, argc, 2)
		readerOptions.IFS = SeparatorFromArg(args[argi+1])
		writerOptions.OFS = SeparatorFromArg(args[argi+1])
		readerOptions.IFSWasSpecified = true
		writerOptions.OFSWasSpecified = true
		argi += 2

	} else if args[argi] == "-p" {
		readerOptions.InputFileFormat = "nidx"
		writerOptions.OutputFileFormat = "nidx"
		readerOptions.IFS = " "
		writerOptions.OFS = " "
		readerOptions.IFSWasSpecified = true
		writerOptions.OFSWasSpecified = true
		readerOptions.AllowRepeatIFS = true
		readerOptions.AllowRepeatIFSWasSpecified = true
		argi += 1

	} else if args[argi] == "--ps" {
		CheckArgCount(args, argi, argc, 2)
		readerOptions.IPS = SeparatorFromArg(args[argi+1])
		writerOptions.OPS = SeparatorFromArg(args[argi+1])
		readerOptions.IPSWasSpecified = true
		writerOptions.OPSWasSpecified = true
		argi += 2

	} else if args[argi] == "--io" {
		CheckArgCount(args, argi, argc, 2)
		if defaultFSes[args[argi+1]] == "" {
			fmt.Fprintf(os.Stderr, "%s: unrecognized I/O format \"%s\".\n",
				"mlr", args[argi+1])
			os.Exit(1)
		}
		readerOptions.InputFileFormat = args[argi+1]
		writerOptions.OutputFileFormat = args[argi+1]
		argi += 2

	} else if args[argi] == "--csv" || args[argi] == "-c" {
		readerOptions.InputFileFormat = "csv"
		writerOptions.OutputFileFormat = "csv"
		argi += 1

	} else if args[argi] == "--csvlite" {
		readerOptions.InputFileFormat = "csvlite"
		writerOptions.OutputFileFormat = "csv"
		argi += 1

	} else if args[argi] == "--tsv" {
		readerOptions.InputFileFormat = "csv"
		writerOptions.OutputFileFormat = "csv"
		readerOptions.IFS = "\t"
		writerOptions.OFS = "\t"
		readerOptions.IFSWasSpecified = true
		writerOptions.OFSWasSpecified = true
		argi += 1

	} else if args[argi] == "--tsvlite" || args[argi] == "-t" {
		readerOptions.InputFileFormat = "csvlite"
		writerOptions.OutputFileFormat = "csvlite"
		readerOptions.IFS = "\t"
		writerOptions.OFS = "\t"
		readerOptions.IFSWasSpecified = true
		writerOptions.OFSWasSpecified = true
		argi += 1

	} else if args[argi] == "--asv" {
		readerOptions.InputFileFormat = "csvlite"
		writerOptions.OutputFileFormat = "csvlite"
		readerOptions.IFS = ASV_FS
		writerOptions.OFS = ASV_FS
		readerOptions.IRS = ASV_RS
		writerOptions.ORS = ASV_RS
		readerOptions.IFSWasSpecified = true

		readerOptions.IRSWasSpecified = true
		writerOptions.OFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1

	} else if args[argi] == "--asvlite" {
		readerOptions.InputFileFormat = "csvlite"
		writerOptions.OutputFileFormat = "csvlite"
		readerOptions.IFS = ASV_FS
		writerOptions.OFS = ASV_FS
		readerOptions.IRS = ASV_RS
		writerOptions.ORS = ASV_RS
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		writerOptions.OFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1

	} else if args[argi] == "--usv" {
		readerOptions.InputFileFormat = "csvlite"
		writerOptions.OutputFileFormat = "csvlite"
		readerOptions.IFS = USV_FS
		writerOptions.OFS = USV_FS
		readerOptions.IRS = USV_RS
		writerOptions.ORS = USV_RS
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		writerOptions.OFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1

	} else if args[argi] == "--usvlite" {
		readerOptions.InputFileFormat = "csvlite"
		writerOptions.OutputFileFormat = "csvlite"
		readerOptions.IFS = USV_FS
		writerOptions.OFS = USV_FS
		readerOptions.IRS = USV_RS
		writerOptions.ORS = USV_RS
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		writerOptions.OFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1

	} else if args[argi] == "--dkvp" {
		readerOptions.InputFileFormat = "dkvp"
		writerOptions.OutputFileFormat = "dkvp"
		argi += 1

	} else if args[argi] == "--json" || args[argi] == "-j" {
		readerOptions.InputFileFormat = "json"
		writerOptions.OutputFileFormat = "json"
		argi += 1
	} else if args[argi] == "--jsonx" {
		// --jvstack is now the default in Miller 6 so this is just for backward compatibility
		readerOptions.InputFileFormat = "json"
		writerOptions.OutputFileFormat = "json"
		argi += 1

	} else if args[argi] == "--nidx" {
		readerOptions.InputFileFormat = "nidx"
		writerOptions.OutputFileFormat = "nidx"
		readerOptions.IFS = " "
		writerOptions.OFS = " "
		readerOptions.IFSWasSpecified = true
		writerOptions.OFSWasSpecified = true
		argi += 1

	} else if args[argi] == "-T" {
		readerOptions.InputFileFormat = "nidx"
		writerOptions.OutputFileFormat = "nidx"
		readerOptions.IFS = "\t"
		writerOptions.OFS = "\t"
		readerOptions.IFSWasSpecified = true
		writerOptions.OFSWasSpecified = true
		argi += 1

	} else if args[argi] == "--xtab" {
		readerOptions.InputFileFormat = "xtab"
		writerOptions.OutputFileFormat = "xtab"
		argi += 1

	} else if args[argi] == "--pprint" {
		readerOptions.InputFileFormat = "pprint"
		readerOptions.IFS = " "
		readerOptions.IFSWasSpecified = true
		writerOptions.OutputFileFormat = "pprint"
		argi += 1
	} else if args[argi] == "--c2t" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "csv"
		writerOptions.ORS = "auto"
		writerOptions.OFS = "\t"
		readerOptions.IRSWasSpecified = true
		writerOptions.OFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1
	} else if args[argi] == "--c2d" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "dkvp"
		readerOptions.IRSWasSpecified = true
		argi += 1
	} else if args[argi] == "--c2n" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "nidx"
		writerOptions.OFS = " "
		readerOptions.IRSWasSpecified = true
		writerOptions.OFSWasSpecified = true
		argi += 1
	} else if args[argi] == "--c2j" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "json"
		readerOptions.IRSWasSpecified = true
		argi += 1
	} else if args[argi] == "--c2p" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "pprint"
		readerOptions.IRSWasSpecified = true
		argi += 1
	} else if args[argi] == "--c2b" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "pprint"
		writerOptions.BarredPprintOutput = true
		readerOptions.IRSWasSpecified = true
		argi += 1
	} else if args[argi] == "--c2x" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "xtab"
		readerOptions.IRSWasSpecified = true
		argi += 1
	} else if args[argi] == "--c2m" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "markdown"
		readerOptions.IRSWasSpecified = true
		argi += 1

	} else if args[argi] == "--t2c" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IFS = "\t"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "csv"
		writerOptions.ORS = "auto"
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1
	} else if args[argi] == "--t2d" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IFS = "\t"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "dkvp"
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		argi += 1
	} else if args[argi] == "--t2n" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IFS = "\t"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "nidx"
		writerOptions.OFS = " "
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		writerOptions.OFSWasSpecified = true
		argi += 1
	} else if args[argi] == "--t2j" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IFS = "\t"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "json"
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		argi += 1
	} else if args[argi] == "--t2p" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IFS = "\t"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "pprint"
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		argi += 1
	} else if args[argi] == "--t2b" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IFS = "\t"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "pprint"
		writerOptions.BarredPprintOutput = true
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		argi += 1
	} else if args[argi] == "--t2x" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IFS = "\t"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "xtab"
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		argi += 1
	} else if args[argi] == "--t2m" {
		readerOptions.InputFileFormat = "csv"
		readerOptions.IFS = "\t"
		readerOptions.IRS = "auto"
		writerOptions.OutputFileFormat = "markdown"
		readerOptions.IFSWasSpecified = true
		readerOptions.IRSWasSpecified = true
		argi += 1

	} else if args[argi] == "--d2c" {
		readerOptions.InputFileFormat = "dkvp"
		writerOptions.OutputFileFormat = "csv"
		writerOptions.ORS = "auto"
		argi += 1
	} else if args[argi] == "--d2t" {
		readerOptions.InputFileFormat = "dkvp"
		writerOptions.OutputFileFormat = "csv"
		writerOptions.ORS = "auto"
		writerOptions.OFS = "\t"
		writerOptions.OFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1
	} else if args[argi] == "--d2n" {
		readerOptions.InputFileFormat = "dkvp"
		writerOptions.OutputFileFormat = "nidx"
		writerOptions.OFS = " "
		writerOptions.OFSWasSpecified = true
		argi += 1
	} else if args[argi] == "--d2j" {
		readerOptions.InputFileFormat = "dkvp"
		writerOptions.OutputFileFormat = "json"
		argi += 1
	} else if args[argi] == "--d2p" {
		readerOptions.InputFileFormat = "dkvp"
		writerOptions.OutputFileFormat = "pprint"
		argi += 1
	} else if args[argi] == "--d2b" {
		readerOptions.InputFileFormat = "dkvp"
		writerOptions.OutputFileFormat = "pprint"
		writerOptions.BarredPprintOutput = true
		argi += 1
	} else if args[argi] == "--d2x" {
		readerOptions.InputFileFormat = "dkvp"
		writerOptions.OutputFileFormat = "xtab"
		argi += 1
	} else if args[argi] == "--d2m" {
		readerOptions.InputFileFormat = "dkvp"
		writerOptions.OutputFileFormat = "markdown"
		argi += 1

	} else if args[argi] == "--n2c" {
		readerOptions.InputFileFormat = "nidx"
		writerOptions.OutputFileFormat = "csv"
		writerOptions.ORS = "auto"
		writerOptions.ORSWasSpecified = true
		argi += 1
	} else if args[argi] == "--n2t" {
		readerOptions.InputFileFormat = "nidx"
		writerOptions.OutputFileFormat = "csv"
		writerOptions.ORS = "auto"
		writerOptions.OFS = "\t"
		writerOptions.OFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1
	} else if args[argi] == "--n2d" {
		readerOptions.InputFileFormat = "nidx"
		writerOptions.OutputFileFormat = "dkvp"
		argi += 1
	} else if args[argi] == "--n2j" {
		readerOptions.InputFileFormat = "nidx"
		writerOptions.OutputFileFormat = "json"
		argi += 1
	} else if args[argi] == "--n2p" {
		readerOptions.InputFileFormat = "nidx"
		writerOptions.OutputFileFormat = "pprint"
		argi += 1
	} else if args[argi] == "--n2b" {
		readerOptions.InputFileFormat = "nidx"
		writerOptions.OutputFileFormat = "pprint"
		writerOptions.BarredPprintOutput = true
		argi += 1
	} else if args[argi] == "--n2x" {
		readerOptions.InputFileFormat = "nidx"
		writerOptions.OutputFileFormat = "xtab"
		argi += 1
	} else if args[argi] == "--n2m" {
		readerOptions.InputFileFormat = "nidx"
		writerOptions.OutputFileFormat = "markdown"
		argi += 1

	} else if args[argi] == "--j2c" {
		readerOptions.InputFileFormat = "json"
		writerOptions.OutputFileFormat = "csv"
		writerOptions.ORS = "auto"
		writerOptions.ORSWasSpecified = true
		argi += 1
	} else if args[argi] == "--j2t" {
		readerOptions.InputFileFormat = "json"
		writerOptions.OutputFileFormat = "csv"
		writerOptions.ORS = "auto"
		writerOptions.OFS = "\t"
		writerOptions.OFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1
	} else if args[argi] == "--j2d" {
		readerOptions.InputFileFormat = "json"
		writerOptions.OutputFileFormat = "dkvp"
		argi += 1
	} else if args[argi] == "--j2n" {
		readerOptions.InputFileFormat = "json"
		writerOptions.OutputFileFormat = "nidx"
		argi += 1
	} else if args[argi] == "--j2p" {
		readerOptions.InputFileFormat = "json"
		writerOptions.OutputFileFormat = "pprint"
		argi += 1
	} else if args[argi] == "--j2b" {
		readerOptions.InputFileFormat = "json"
		writerOptions.OutputFileFormat = "pprint"
		writerOptions.BarredPprintOutput = true
		argi += 1
	} else if args[argi] == "--j2x" {
		readerOptions.InputFileFormat = "json"
		writerOptions.OutputFileFormat = "xtab"
		argi += 1
	} else if args[argi] == "--j2m" {
		readerOptions.InputFileFormat = "json"
		writerOptions.OutputFileFormat = "markdown"
		argi += 1

	} else if args[argi] == "--p2c" {
		readerOptions.InputFileFormat = "pprint"
		readerOptions.IFS = " "
		writerOptions.OutputFileFormat = "csv"
		writerOptions.ORS = "auto"
		readerOptions.IFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1
	} else if args[argi] == "--p2t" {
		readerOptions.InputFileFormat = "pprint"
		readerOptions.IFS = " "
		writerOptions.OutputFileFormat = "csv"
		writerOptions.ORS = "auto"
		writerOptions.OFS = "\t"
		readerOptions.IFSWasSpecified = true
		writerOptions.OFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1
	} else if args[argi] == "--p2d" {
		readerOptions.InputFileFormat = "pprint"
		readerOptions.IFS = " "
		writerOptions.OutputFileFormat = "dkvp"
		readerOptions.IFSWasSpecified = true
		argi += 1
	} else if args[argi] == "--p2n" {
		readerOptions.InputFileFormat = "pprint"
		readerOptions.IFS = " "
		writerOptions.OutputFileFormat = "nidx"
		readerOptions.IFSWasSpecified = true
		argi += 1
	} else if args[argi] == "--p2j" {
		readerOptions.InputFileFormat = "pprint"
		readerOptions.IFS = " "
		writerOptions.OutputFileFormat = "json"
		readerOptions.IFSWasSpecified = true
		argi += 1
	} else if args[argi] == "--p2x" {
		readerOptions.InputFileFormat = "pprint"
		readerOptions.IFS = " "
		writerOptions.OutputFileFormat = "xtab"
		readerOptions.IFSWasSpecified = true
		argi += 1
	} else if args[argi] == "--p2m" {
		readerOptions.InputFileFormat = "pprint"
		readerOptions.IFS = " "
		writerOptions.OutputFileFormat = "markdown"
		readerOptions.IFSWasSpecified = true
		argi += 1

	} else if args[argi] == "--x2c" {
		readerOptions.InputFileFormat = "xtab"
		writerOptions.OutputFileFormat = "csv"
		writerOptions.ORS = "auto"
		writerOptions.ORSWasSpecified = true
		argi += 1
	} else if args[argi] == "--x2t" {
		readerOptions.InputFileFormat = "xtab"
		writerOptions.OutputFileFormat = "csv"
		writerOptions.ORS = "auto"
		writerOptions.OFS = "\t"
		writerOptions.OFSWasSpecified = true
		writerOptions.ORSWasSpecified = true
		argi += 1
	} else if args[argi] == "--x2d" {
		readerOptions.InputFileFormat = "xtab"
		writerOptions.OutputFileFormat = "dkvp"
		argi += 1
	} else if args[argi] == "--x2n" {
		readerOptions.InputFileFormat = "xtab"
		writerOptions.OutputFileFormat = "nidx"
		argi += 1
	} else if args[argi] == "--x2j" {
		readerOptions.InputFileFormat = "xtab"
		writerOptions.OutputFileFormat = "json"
		argi += 1
	} else if args[argi] == "--x2p" {
		readerOptions.InputFileFormat = "xtab"
		writerOptions.OutputFileFormat = "pprint"
		argi += 1
	} else if args[argi] == "--x2b" {
		readerOptions.InputFileFormat = "xtab"
		writerOptions.OutputFileFormat = "pprint"
		writerOptions.BarredPprintOutput = true
		argi += 1
	} else if args[argi] == "--x2m" {
		readerOptions.InputFileFormat = "xtab"
		writerOptions.OutputFileFormat = "markdown"
		argi += 1

	} else if args[argi] == "-N" {
		readerOptions.UseImplicitCSVHeader = true
		writerOptions.HeaderlessCSVOutput = true
		argi += 1
	}
	*pargi = argi
	return argi != oargi
}

// Returns true if the current flag was handled.
func ParseMiscOptions(
	args []string,
	argc int,
	pargi *int,
	options *TOptions,
) bool {
	argi := *pargi
	oargi := argi

	if args[argi] == "-n" {
		options.NoInput = true
		argi += 1

	} else if args[argi] == "-I" {
		options.DoInPlace = true
		argi += 1

	} else if args[argi] == "--from" {
		CheckArgCount(args, argi, argc, 2)
		options.FileNames = append(options.FileNames, args[argi+1])
		argi += 2

	} else if args[argi] == "--mfrom" {
		CheckArgCount(args, argi, argc, 2)
		argi += 1
		for argi < argc && args[argi] != "--" {
			options.FileNames = append(options.FileNames, args[argi])
			argi += 1
		}
		if args[argi] == "--" {
			argi += 1
		}

	} else if args[argi] == "--load" {
		CheckArgCount(args, argi, argc, 2)
		options.DSLPreloadFileNames = append(options.DSLPreloadFileNames, args[argi+1])
		argi += 2

	} else if args[argi] == "--mload" {
		CheckArgCount(args, argi, argc, 2)
		argi += 1
		for argi < argc && args[argi] != "--" {
			options.DSLPreloadFileNames = append(options.DSLPreloadFileNames, args[argi])
			argi += 1
		}
		if args[argi] == "--" {
			argi += 1
		}

		//	} else if args[argi] == "--nr-progress-mod" {
		//		CheckArgCount(args, argi, argc, 2);
		//		if (sscanf(args[argi+1], "%lld", &options.nr_progress_mod) != 1) {
		//			fmt.Fprintf(os.Stderr,
		//				"%s: --nr-progress-mod argument must be a positive integer; got \"%s\".\n",
		//				"mlr", args[argi+1]);
		//			mainUsageShort()
		//			os.Exit(1);
		//		}
		//		if (options.nr_progress_mod <= 0) {
		//			fmt.Fprintf(os.Stderr,
		//				"%s: --nr-progress-mod argument must be a positive integer; got \"%s\".\n",
		//				"mlr", args[argi+1]);
		//			mainUsageShort()
		//			os.Exit(1);
		//		}
		//		argi += 2;
		//
	} else if args[argi] == "--seed" {
		CheckArgCount(args, argi, argc, 2)
		randSeed, ok := lib.TryIntFromString(args[argi+1])
		if ok {
			options.RandSeed = randSeed
			options.HaveRandSeed = true
		} else {
			fmt.Fprintf(os.Stderr,
				"%s: --seed argument must be a decimal or hexadecimal integer; got \"%s\".\n",
				"mlr", args[argi+1])
			fmt.Fprintf(os.Stderr, "Please run \"%s --help\" for detailed usage information.\n", "mlr")
			os.Exit(1)
		}
		argi += 2

	} else if args[argi] == "--list" {
		argi += 1
		FLAG_TABLE.ListTemp()
		os.Exit(0)

	}
	*pargi = argi
	return argi != oargi
}

// ================================================================

var FLAG_TABLE = FlagTable{
	sections: []*FlagSection{
		&LegacyFlagSection,
		&SeparatorFlagSection,
		&JSONOnlyFlagSection,
		&FileFormatFlagSection,
		&CSVOnlyFlagSection,
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

func init() { SeparatorFlagSection.Sort() }

var SeparatorFlagSection = FlagSection{
	name: "separator flags",
	flags: []Flag{

		{
			name: "--ifs",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.IFS = SeparatorFromArg(args[*pargi+1])
				options.ReaderOptions.IFSWasSpecified = true
				*pargi += 2
			},
			forReader: true,
			forWriter: false,
		},

		{
			name: "--ips",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.IPS = SeparatorFromArg(args[*pargi+1])
				options.ReaderOptions.IPSWasSpecified = true
				*pargi += 2
			},
			forReader: true,
			forWriter: false,
		},

		{
			name: "--irs",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.IRS = SeparatorFromArg(args[*pargi+1])
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 2
			},
			forReader: true,
			forWriter: false,
		},

		{
			name: "--repifs",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.AllowRepeatIFS = true
				options.ReaderOptions.AllowRepeatIFSWasSpecified = true
				*pargi += 1
			},
			forReader: true,
		},

		{
			name: "--rs",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.IRS = SeparatorFromArg(args[*pargi+1])
				options.WriterOptions.ORS = SeparatorFromArg(args[*pargi+1])
				options.ReaderOptions.IRSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 2
			},
			forReader: true,
			forWriter: true,
		},

		{
			name: "--fs",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.IFS = SeparatorFromArg(args[*pargi+1])
				options.WriterOptions.OFS = SeparatorFromArg(args[*pargi+1])
				options.ReaderOptions.IFSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 2
			},
			forReader: true,
			forWriter: true,
		},

		{
			name: "--ps",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.IPS = SeparatorFromArg(args[*pargi+1])
				options.WriterOptions.OPS = SeparatorFromArg(args[*pargi+1])
				options.ReaderOptions.IPSWasSpecified = true
				options.WriterOptions.OPSWasSpecified = true
				*pargi += 2
			},
			forReader: true,
			forWriter: true,
		},
	},
}

// ================================================================
// JSON-ONLY FLAGS

func init() { JSONOnlyFlagSection.Sort() }

var JSONOnlyFlagSection = FlagSection{
	name: "json-only flags",
	flags: []Flag{

		{
			name: "--jvstack",
			help: "Help goes here",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.JSONOutputMultiline = true
				*pargi += 1
			},
			forWriter: true,
		},

		{
			name: "--no-jvstack",
			help: "Help goes here",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.JSONOutputMultiline = false
				*pargi += 1
			},
			forWriter: true,
		},

		{
			name: "--jlistwrap",
			help: "Help goes here",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.WrapJSONOutputInOuterList = true
				*pargi += 1
			},
			forWriter: true,
		},

		{
			name:   "--jknquoteint",
			help:   NoOpHelp,
			parser: NoOpParse1,
		},

		{
			name:   "--jquoteall",
			help:   NoOpHelp,
			parser: NoOpParse1,
		},

		{
			name:   "--jvquoteall",
			help:   NoOpHelp,
			parser: NoOpParse1,
		},

		{
			name:   "--json-fatal-arrays-on-input",
			help:   NoOpHelp,
			parser: NoOpParse1,
		},

		{
			name:   "--json-skip-arrays-on-input",
			help:   NoOpHelp,
			parser: NoOpParse1,
		},

		{
			name:   "--json-skip-arrays-on-input",
			help:   NoOpHelp,
			parser: NoOpParse1,
		},
	},
}

// ================================================================
// LEGACY FLAGS

func init() { LegacyFlagSection.Sort() }

var LegacyFlagSection = FlagSection{
	name: "legacy flags",
	flags: []Flag{

		{
			name:      "--mmap",
			help:      NoOpHelp,
			parser:    NoOpParse1,
			forReader: true,
		},

		{
			name:      "--no-mmap",
			help:      NoOpHelp,
			parser:    NoOpParse1,
			forReader: true,
		},

		{
			name:      "--no-fflush",
			help:      NoOpHelp,
			parser:    NoOpParse1,
			forWriter: true,
		},
	},
}

// ================================================================
// FILE-FORMAT FLAGS

func init() { FileFormatFlagSection.Sort() }

var FileFormatFlagSection = FlagSection{
	name: "file-format flags",
	flags: []Flag{

		{
			name: "--icsv",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				*pargi += 1
			},
		},

		{
			name: "--icsvlite",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csvlite"
				*pargi += 1
			},
		},

		{
			name: "--itsv",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.ReaderOptions.IFSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--itsvlite",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csvlite"
				options.ReaderOptions.IFS = "\t"
				options.ReaderOptions.IFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--iasv",
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
			name: "--iasvlite",
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
			name: "--iusv",
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
			name: "--iusvlite",
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
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				*pargi += 1
			},
		},

		{
			name: "--ijson",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				*pargi += 1
			},
		},

		{
			name: "--inidx",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				*pargi += 1
			},
		},

		{
			name: "--ixtab",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				*pargi += 1
			},
		},

		{
			name: "--ipprint",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "pprint"
				options.ReaderOptions.IFS = " "
				options.ReaderOptions.IFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "-i",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.InputFileFormat = args[*pargi+1]
				*pargi += 2
			},
		},

		//{
		//	name: "--igen",
		//	parser: func(args []string, argc int, pargi *int, options *TOptions) {
		//		options.ReaderOptions.InputFileFormat = "gen"
		//		*pargi += 1
		//	},
		//},
		//{
		//	name: "--gen-start",
		//	parser: func(args []string, argc int, pargi *int, options *TOptions) {
		//		options.ReaderOptions.InputFileFormat = "gen"
		//		CheckArgCount(args, *pargi, argc, 2)
		//		if sscanf(args[*pargi+1], "%lld", &options.ReaderOptions.generator_opts.start) != 1 {
		//			fmt.Fprintf(os.Stderr, "%s: could not scan \"%s\".\n",
		//				"mlr", args[*pargi+1])
		//		}
		//		*pargi += 2
		//	},
		//},
		//{
		//	name: "--gen-stop",
		//	parser: func(args []string, argc int, pargi *int, options *TOptions) {
		//		options.ReaderOptions.InputFileFormat = "gen"
		//		CheckArgCount(args, *pargi, argc, 2)
		//		if sscanf(args[*pargi+1], "%lld", &options.ReaderOptions.generator_opts.stop) != 1 {
		//			fmt.Fprintf(os.Stderr, "%s: could not scan \"%s\".\n",
		//				"mlr", args[*pargi+1])
		//		}
		//		*pargi += 2
		//	},
		//},
		//{
		//	name: "--gen-step",
		//	parser: func(args []string, argc int, pargi *int, options *TOptions) {
		//		options.ReaderOptions.InputFileFormat = "gen"
		//		CheckArgCount(args, *pargi, argc, 2)
		//		if sscanf(args[*pargi+1], "%lld", &options.ReaderOptions.generator_opts.step) != 1 {
		//			fmt.Fprintf(os.Stderr, "%s: could not scan \"%s\".\n",
		//				"mlr", args[*pargi+1])
		//		}
		//		*pargi += 2
		//	},
		//},

		{
			name: "--ors",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.WriterOptions.ORS = SeparatorFromArg(args[*pargi+1])
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 2
			},
		},

		{
			name: "--ofs",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.WriterOptions.OFS = SeparatorFromArg(args[*pargi+1])
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 2
			},
		},

		{
			name: "--ops",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.WriterOptions.OPS = SeparatorFromArg(args[*pargi+1])
				options.WriterOptions.OPSWasSpecified = true
				*pargi += 2
			},
		},

		{
			name: "-o",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.WriterOptions.OutputFileFormat = args[*pargi+1]
				*pargi += 2
			},
		},

		{
			name: "--ocsv",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "csv"
				*pargi += 1
			},
		},

		{
			name: "--ocsvlite",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "csvlite"
				*pargi += 1
			},
		},

		{
			name: "--otsv",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.OFS = "\t"
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--otsvlite",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "csvlite"
				options.WriterOptions.OFS = "\t"
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--oasv",
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
			name: "--oasvlite",
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
			name: "--ousv",
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
			name: "--ousvlite",
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
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "markdown"
				*pargi += 1
			},
		},

		{
			name: "--odkvp",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "dkvp"
				*pargi += 1
			},
		},

		{
			name: "--ojson",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "json"
				*pargi += 1
			},
		},
		{
			name: "--ojsonx",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				// --jvstack is now the default in Miller 6 so this is just for backward compatibility
				options.WriterOptions.OutputFileFormat = "json"
				*pargi += 1
			},
		},

		{
			name: "--onidx",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "nidx"
				options.WriterOptions.OFS = " "
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--oxtab",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "xtab"
				*pargi += 1
			},
		},

		{
			name: "--opprint",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.OutputFileFormat = "pprint"
				*pargi += 1
			},
		},

		{
			name: "--right",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.RightAlignedPprintOutput = true
				*pargi += 1
			},
		},

		{
			name: "--barred",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.BarredPprintOutput = true
				*pargi += 1
			},
		},

		{
			name: "--io",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				if defaultFSes[args[*pargi+1]] == "" {
					fmt.Fprintf(os.Stderr, "%s: unrecognized I/O format \"%s\".\n",
						"mlr", args[*pargi+1])
					os.Exit(1)
				}
				options.ReaderOptions.InputFileFormat = args[*pargi+1]
				options.WriterOptions.OutputFileFormat = args[*pargi+1]
				*pargi += 2
			},
		},

		{
			name:     "--csv",
			altNames: []string{"-c"},
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.WriterOptions.OutputFileFormat = "csv"
				*pargi += 1
			},
		},

		{
			name: "--csvlite",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csvlite"
				options.WriterOptions.OutputFileFormat = "csv"
				*pargi += 1
			},
		},

		{
			name: "--tsv",
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
			name: "--asv",
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
			name: "--asvlite",
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
			name: "--usv",
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
			name: "--usvlite",
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
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "dkvp"
				*pargi += 1
			},
		},

		{
			name:     "--json",
			altNames: []string{"-j"},
			parser: func(args []string, argc int, pargi *int, options *TOptions) {

				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "json"
				*pargi += 1
			},
		},
		{
			name: "--jsonx",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				// --jvstack is now the default in Miller 6 so this is just for backward compatibility
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "json"
				*pargi += 1
			},
		},

		{
			name: "--nidx",
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
			name: "-T",
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
			name: "--xtab",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "xtab"
				*pargi += 1
			},
		},

		{
			name: "--pprint",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "pprint"
				options.ReaderOptions.IFS = " "
				options.ReaderOptions.IFSWasSpecified = true
				options.WriterOptions.OutputFileFormat = "pprint"
				*pargi += 1
			},
		},
		{
			name: "--c2t",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORS = "auto"
				options.WriterOptions.OFS = "\t"
				options.ReaderOptions.IRSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--c2d",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "dkvp"
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--c2n",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "nidx"
				options.WriterOptions.OFS = " "
				options.ReaderOptions.IRSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--c2j",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "json"
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--c2p",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--c2b",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.WriterOptions.BarredPprintOutput = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--c2x",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "xtab"
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--c2m",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "markdown"
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--t2c",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORS = "auto"
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--t2d",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "dkvp"
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--t2n",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.ReaderOptions.IRS = "auto"
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
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "json"
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--t2p",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--t2b",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.WriterOptions.BarredPprintOutput = true
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--t2x",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "xtab"
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--t2m",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "csv"
				options.ReaderOptions.IFS = "\t"
				options.ReaderOptions.IRS = "auto"
				options.WriterOptions.OutputFileFormat = "markdown"
				options.ReaderOptions.IFSWasSpecified = true
				options.ReaderOptions.IRSWasSpecified = true
				*pargi += 1
			},
		},

		{
			name: "--d2c",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORS = "auto"
				*pargi += 1
			},
		},
		{
			name: "--d2t",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORS = "auto"
				options.WriterOptions.OFS = "\t"
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--d2n",
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
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "json"
				*pargi += 1
			},
		},
		{
			name: "--d2p",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "pprint"
				*pargi += 1
			},
		},
		{
			name: "--d2b",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.WriterOptions.BarredPprintOutput = true
				*pargi += 1
			},
		},
		{
			name: "--d2x",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "xtab"
				*pargi += 1
			},
		},
		{
			name: "--d2m",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "dkvp"
				options.WriterOptions.OutputFileFormat = "markdown"
				*pargi += 1
			},
		},

		{
			name: "--n2c",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORS = "auto"
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--n2t",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORS = "auto"
				options.WriterOptions.OFS = "\t"
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--n2d",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "dkvp"
				*pargi += 1
			},
		},
		{
			name: "--n2j",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "json"
				*pargi += 1
			},
		},
		{
			name: "--n2p",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "pprint"
				*pargi += 1
			},
		},
		{
			name: "--n2b",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.WriterOptions.BarredPprintOutput = true
				*pargi += 1
			},
		},
		{
			name: "--n2x",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "xtab"
				*pargi += 1
			},
		},
		{
			name: "--n2m",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "nidx"
				options.WriterOptions.OutputFileFormat = "markdown"
				*pargi += 1
			},
		},

		{
			name: "--j2c",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORS = "auto"
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--j2t",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORS = "auto"
				options.WriterOptions.OFS = "\t"
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--j2d",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "dkvp"
				*pargi += 1
			},
		},
		{
			name: "--j2n",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "nidx"
				*pargi += 1
			},
		},
		{
			name: "--j2p",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "pprint"
				*pargi += 1
			},
		},
		{
			name: "--j2b",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.WriterOptions.BarredPprintOutput = true
				*pargi += 1
			},
		},
		{
			name: "--j2x",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "xtab"
				*pargi += 1
			},
		},
		{
			name: "--j2m",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "json"
				options.WriterOptions.OutputFileFormat = "markdown"
				*pargi += 1
			},
		},

		{
			name: "--p2c",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "pprint"
				options.ReaderOptions.IFS = " "
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORS = "auto"
				options.ReaderOptions.IFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--p2t",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "pprint"
				options.ReaderOptions.IFS = " "
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORS = "auto"
				options.WriterOptions.OFS = "\t"
				options.ReaderOptions.IFSWasSpecified = true
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--p2d",
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
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORS = "auto"
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--x2t",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "csv"
				options.WriterOptions.ORS = "auto"
				options.WriterOptions.OFS = "\t"
				options.WriterOptions.OFSWasSpecified = true
				options.WriterOptions.ORSWasSpecified = true
				*pargi += 1
			},
		},
		{
			name: "--x2d",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "dkvp"
				*pargi += 1
			},
		},
		{
			name: "--x2n",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "nidx"
				*pargi += 1
			},
		},
		{
			name: "--x2j",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "json"
				*pargi += 1
			},
		},
		{
			name: "--x2p",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "pprint"
				*pargi += 1
			},
		},
		{
			name: "--x2b",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "pprint"
				options.WriterOptions.BarredPprintOutput = true
				*pargi += 1
			},
		},
		{
			name: "--x2m",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.InputFileFormat = "xtab"
				options.WriterOptions.OutputFileFormat = "markdown"
				*pargi += 1
			},
		},

		{
			name: "-p",
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
	},
}

// ================================================================
// CSV FLAGS

func init() { CSVOnlyFlagSection.Sort() }

var CSVOnlyFlagSection = FlagSection{
	name: "CSV-only flags",
	flags: []Flag{

		{
			name: "--no-implicit-csv-header",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.UseImplicitCSVHeader = false
				*pargi += 1
			},
		},

		{
			name:     "--allow-ragged-csv-input",
			altNames: []string{"--ragged"},
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.AllowRaggedCSVInput = true
				*pargi += 1
			},
		},

		{
			name: "--implicit-csv-header",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.UseImplicitCSVHeader = true
				*pargi += 1
			},
		},

		{
			name: "--headerless-csv-output",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.HeaderlessCSVOutput = true
				*pargi += 1
			},
		},

		{
			name: "-N",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.UseImplicitCSVHeader = true
				options.WriterOptions.HeaderlessCSVOutput = true
				*pargi += 1
			},
		},

		//{
		//		name: "--quote-all",
		//	parser: func(args []string, argc int, pargi *int, options *TOptions) {
		//		options.WriterOptions.oquoting = QUOTE_ALL
		//		*pargi += 1
		//	},
		//},
		//{
		//	name: "--quote-none",
		//	parser: func(args []string, argc int, pargi *int, options *TOptions) {
		//		options.WriterOptions.oquoting = QUOTE_NONE
		//		*pargi += 1
		//	},
		//},
		//{
		//	name: "--quote-minimal",
		//	parser: func(args []string, argc int, pargi *int, options *TOptions) {
		//		options.WriterOptions.oquoting = QUOTE_MINIMAL
		//		*pargi += 1
		//	},
		//},
		//{
		//	name: "--quote-numeric",
		//	parser: func(args []string, argc int, pargi *int, options *TOptions) {
		//		options.WriterOptions.oquoting = QUOTE_NUMERIC
		//		*pargi += 1
		//	},
		//},
		//{
		//	name: "--quote-original",
		//	parser: func(args []string, argc int, pargi *int, options *TOptions) {
		//		options.WriterOptions.oquoting = QUOTE_ORIGINAL
		//		*pargi += 1
		//	},
		//},

	},
}

// ================================================================
// COMPRESSED-DATA FLAGS

func init() { CompressedDataFlagSection.Sort() }

var CompressedDataFlagSection = FlagSection{
	name: "compressed-data flags",
	flags: []Flag{

		{
			name: "--prepipe",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.Prepipe = args[*pargi+1]
				options.ReaderOptions.PrepipeIsRaw = false
				*pargi += 2
			},
		},

		{
			name: "--prepipex",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.Prepipe = args[*pargi+1]
				options.ReaderOptions.PrepipeIsRaw = true
				*pargi += 2
			},
		},

		{
			name: "--prepipe-gunzip",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.Prepipe = "gunzip"
				options.ReaderOptions.PrepipeIsRaw = false
				*pargi += 1
			},
		},

		{
			name: "--prepipe-zcat",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.Prepipe = "zcat"
				options.ReaderOptions.PrepipeIsRaw = false
				*pargi += 1
			},
		},

		{
			name: "--prepipe-bz2",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.Prepipe = "bz2"
				options.ReaderOptions.PrepipeIsRaw = false
				*pargi += 1
			},
		},

		{
			name: "--gzin",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.FileInputEncoding = lib.FileInputEncodingGzip
				*pargi += 1
			},
		},

		{
			name: "--zin",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.FileInputEncoding = lib.FileInputEncodingZlib
				*pargi += 1
			},
		},

		{
			name: "--bz2in",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.FileInputEncoding = lib.FileInputEncodingBzip2
				*pargi += 1
			},
		},
	},
}

// ================================================================
// COMMENTS-IN-DATA FLAGS

func init() { CommentsInDataFlagSection.Sort() }

var CommentsInDataFlagSection = FlagSection{
	name: "comments-in-data flags",
	flags: []Flag{

		{
			name: "--skip-comments",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.CommentString = DEFAULT_COMMENT_STRING
				options.ReaderOptions.CommentHandling = SkipComments
				*pargi += 1
			},
		},

		{
			name: "--skip-comments-with",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.ReaderOptions.CommentString = args[*pargi+1]
				options.ReaderOptions.CommentHandling = SkipComments
				*pargi += 2
			},
		},

		{
			name: "--pass-comments",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.ReaderOptions.CommentString = DEFAULT_COMMENT_STRING
				options.ReaderOptions.CommentHandling = PassComments
				*pargi += 1
			},
		},

		{
			name: "--pass-comments-with",
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

func init() { OutputColorizationFlagSection.Sort() }

var OutputColorizationFlagSection = FlagSection{
	name: "output-colorization flags",
	flags: []Flag{

		{
			name: "--list-color-codes",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				colorizer.ListColorCodes()
				os.Exit(0)
				*pargi += 1
			},
		},

		{
			name: "--list-color-names",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				colorizer.ListColorNames()
				os.Exit(0)
				*pargi += 1
			},
		},

		{
			name:     "--no-color",
			altNames: []string{"-M"},
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				colorizer.SetColorization(colorizer.ColorizeOutputNever)
				*pargi += 1
			},
		},

		{
			name:     "--always-color",
			altNames: []string{"-C"},
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				colorizer.SetColorization(colorizer.ColorizeOutputAlways)
				*pargi += 1
			},
		},

		{
			name: "--key-color",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				ok := colorizer.SetKeyColor(args[*pargi+1])
				if !ok {
					fmt.Fprintf(os.Stderr,
						"%s: --key-color argument unrecognized; got \"%s\".\n",
						"mlr", args[*pargi+1])
					os.Exit(1)
				}
				*pargi += 2
			},
		},

		{
			name: "--value-color",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				ok := colorizer.SetValueColor(args[*pargi+1])
				if !ok {
					fmt.Fprintf(os.Stderr,
						"%s: --value-color argument unrecognized; got \"%s\".\n",
						"mlr", args[*pargi+1])
					os.Exit(1)
				}
				*pargi += 2
			},
		},

		{
			name: "--pass-color",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				ok := colorizer.SetPassColor(args[*pargi+1])
				if !ok {
					fmt.Fprintf(os.Stderr,
						"%s: --pass-color argument unrecognized; got \"%s\".\n",
						"mlr", args[*pargi+1])
					os.Exit(1)
				}
				*pargi += 2
			},
		},

		{
			name: "--fail-color",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				ok := colorizer.SetFailColor(args[*pargi+1])
				if !ok {
					fmt.Fprintf(os.Stderr,
						"%s: --fail-color argument unrecognized; got \"%s\".\n",
						"mlr", args[*pargi+1])
					os.Exit(1)
				}
				*pargi += 2
			},
		},

		{
			name: "--help-color",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				ok := colorizer.SetHelpColor(args[*pargi+1])
				if !ok {
					fmt.Fprintf(os.Stderr,
						"%s: --help-color argument unrecognized; got \"%s\".\n",
						"mlr", args[*pargi+1])
					os.Exit(1)
				}
				*pargi += 2
			},
		},
	},
}

// ================================================================
// FLATTEN/UNFLATTEN FLAGS

func init() { FlattenUnflattenFlagSection.Sort() }

var FlattenUnflattenFlagSection = FlagSection{
	name: "flatten-unflatten flags",
	flags: []Flag{

		{
			name:     "--flatsep",
			altNames: []string{"--jflatsep", "--oflatsep"}, // TODO: really need all for miller5 back-compat?
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.WriterOptions.FLATSEP = SeparatorFromArg(args[*pargi+1])
				*pargi += 2
			},
		},

		//{
		//	name: "--xvright",
		//	parser: func(args []string, argc int, pargi *int, options *TOptions) {
		//		options.WriterOptions.right_justify_xtab_value = true
		//		*pargi += 1
		//	},
		//},

		//{
		//	name: "--vflatsep",
		//	parser: func(args []string, argc int, pargi *int, options *TOptions) {
		//		CheckArgCount(args, *pargi, argc, 2)
		//		// No-op pass-through for backward compatibility with Miller 5
		//		*pargi += 2
		// },
		//},

		{
			name: "--no-auto-flatten",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.AutoFlatten = false
				*pargi += 1
			},
		},

		{
			name: "--no-auto-unflatten",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.WriterOptions.AutoUnflatten = false
				*pargi += 1
			},
		},
	},
}

// ================================================================
// MISC FLAGS

func init() { MiscFlagSection.Sort() }

var MiscFlagSection = FlagSection{
	name: "miscellaneous flags",
	flags: []Flag{

		{
			name: "-n",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.NoInput = true
				*pargi += 1
			},
		},

		{
			name: "-I",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				options.DoInPlace = true
				*pargi += 1
			},
		},

		{
			name: "--from",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.FileNames = append(options.FileNames, args[*pargi+1])
				*pargi += 2
			},
		},

		{
			name: "--mfrom",
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
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.WriterOptions.FPOFMT = args[*pargi+1]
				*pargi += 2
			},
		},

		// TODO: move to another (or new) section
		{
			name: "--load",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				options.DSLPreloadFileNames = append(options.DSLPreloadFileNames, args[*pargi+1])
				*pargi += 2
			},
		},

		{
			name: "--mload",
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

		//		name: "--nr-progress-mod",
		//			parser: func(args []string, argc int, pargi *int, options *TOptions) {
		//				CheckArgCount(args, *pargi, argc, 2);
		//				if (sscanf(args[*pargi+1], "%lld", &options.nr_progress_mod) != 1) {
		//					fmt.Fprintf(os.Stderr,
		//						"%s: --nr-progress-mod argument must be a positive integer; got \"%s\".\n",
		//						"mlr", args[*pargi+1]);
		//					mainUsageShort()
		//					os.Exit(1);
		//				}
		//				if (options.nr_progress_mod <= 0) {
		//					fmt.Fprintf(os.Stderr,
		//						"%s: --nr-progress-mod argument must be a positive integer; got \"%s\".\n",
		//						"mlr", args[*pargi+1]);
		//					mainUsageShort()
		//					os.Exit(1);
		//				}
		//				*pargi += 2;
		// },

		{
			name: "--seed",
			parser: func(args []string, argc int, pargi *int, options *TOptions) {
				CheckArgCount(args, *pargi, argc, 2)
				randSeed, ok := lib.TryIntFromString(args[*pargi+1])
				if ok {
					options.RandSeed = randSeed
					options.HaveRandSeed = true
				} else {
					fmt.Fprintf(os.Stderr,
						"%s: --seed argument must be a decimal or hexadecimal integer; got \"%s\".\n",
						"mlr", args[*pargi+1])
					fmt.Fprintf(os.Stderr, "Please run \"%s --help\" for detailed usage information.\n", "mlr")
					os.Exit(1)
				}
				*pargi += 2
			},
		},
	},
}
