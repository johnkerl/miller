// ================================================================
// Items which might better belong in miller/cli, but which are placed in a
// deeper package to avoid a package-dependency cycle between miller/cli and
// miller/transforming.
// ================================================================

package cliutil

import (
	"fmt"
	"os"

	"miller/src/colorizer"
	"miller/src/lib"
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

	} else if args[argi] == "--oflatsep" {
		CheckArgCount(args, argi, argc, 2)
		writerOptions.OFLATSEP = SeparatorFromArg(args[argi+1])
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

	} else if args[argi] == "--list-colors" {
		colorizer.ListColors(os.Stdout)
		os.Exit(0)
		argi += 1

	} else if args[argi] == "--no-color" || args[argi] == "-M" {
		colorizer.SetColorization(colorizer.ColorizeOutputNever)
		argi += 1

	} else if args[argi] == "--always-color" || args[argi] == "-C" {
		colorizer.SetColorization(colorizer.ColorizeOutputAlways)
		argi += 1

	} else if args[argi] == "--key-color" {
		colorizer.SetColorization(colorizer.ColorizeOutputAlways)
		CheckArgCount(args, argi, argc, 2)
		code, ok := lib.TryIntFromString(args[argi+1])
		if !ok {
			fmt.Fprintf(os.Stderr,
				"%s: --key-color argument must be a decimal integer; got \"%s\".\n",
				"mlr", args[argi+1])
			os.Exit(1)
		}
		colorizer.SetKeyColor(code)
		argi += 2

	} else if args[argi] == "--value-color" {
		colorizer.SetColorization(colorizer.ColorizeOutputAlways)
		CheckArgCount(args, argi, argc, 2)
		code, ok := lib.TryIntFromString(args[argi+1])
		if !ok {
			fmt.Fprintf(os.Stderr,
				"%s: --value-color argument must be a decimal integer; got \"%s\".\n",
				"mlr", args[argi+1])
			os.Exit(1)
		}
		colorizer.SetValueColor(code)
		argi += 2

	} else if args[argi] == "--pass-color" {
		colorizer.SetColorization(colorizer.ColorizeOutputAlways)
		CheckArgCount(args, argi, argc, 2)
		code, ok := lib.TryIntFromString(args[argi+1])
		if !ok {
			fmt.Fprintf(os.Stderr,
				"%s: --pass-color argument must be a decimal integer; got \"%s\".\n",
				"mlr", args[argi+1])
			os.Exit(1)
		}
		colorizer.SetPassColor(code)
		argi += 2

	} else if args[argi] == "--fail-color" {
		colorizer.SetColorization(colorizer.ColorizeOutputAlways)
		CheckArgCount(args, argi, argc, 2)
		code, ok := lib.TryIntFromString(args[argi+1])
		if !ok {
			fmt.Fprintf(os.Stderr,
				"%s: --fail-color argument must be a decimal integer; got \"%s\".\n",
				"mlr", args[argi+1])
			os.Exit(1)
		}
		colorizer.SetFailColor(code)
		argi += 2

	} else if args[argi] == "--help-color" {
		colorizer.SetColorization(colorizer.ColorizeOutputAlways)
		CheckArgCount(args, argi, argc, 2)
		code, ok := lib.TryIntFromString(args[argi+1])
		if !ok {
			fmt.Fprintf(os.Stderr,
				"%s: --help-color argument must be a decimal integer; got \"%s\".\n",
				"mlr", args[argi+1])
			os.Exit(1)
		}
		colorizer.SetHelpColor(code)
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

	} else if args[argi] == "--csv" {
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

	} else if args[argi] == "--json" {
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

	}
	*pargi = argi
	return argi != oargi
}

// ----------------------------------------------------------------
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
