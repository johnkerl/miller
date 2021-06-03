package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"miller/src/cliutil"
	"miller/src/dsl/cst"
	"miller/src/lib"
	"miller/src/transformers"
	"miller/src/transforming"
	"miller/src/types"
	"miller/src/version"
)

// ----------------------------------------------------------------
func ParseCommandLine(args []string) (
	options cliutil.TOptions,
	recordTransformers []transforming.IRecordTransformer,
	err error,
) {
	options = cliutil.DefaultOptions()
	argc := len(args)
	argi := 1

	// Try .mlrrc overrides (then command-line on top of that).
	// A --norc flag (if provided) must come before all other options.
	// Or, they can set the environment variable MLRRC="__none__".
	if argc >= 2 && args[argi] == "--norc" {
		argi++
	} else {
		loadMlrrcOrDie(&options)
	}

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if args[argi][0] != '-' {
			break // No more flag options to process
		} else if args[argi] == "--cpuprofile" {
			// Already handled in main(); ignore here.
			cliutil.CheckArgCount(args, argi, argc, 1)
			argi += 2
		} else if parseTerminalUsage(args, argc, argi) {
			os.Exit(0)
		} else if cliutil.ParseReaderOptions(args, argc, &argi, &options.ReaderOptions) {
			// handled
		} else if cliutil.ParseWriterOptions(args, argc, &argi, &options.WriterOptions) {
			// handled
		} else if cliutil.ParseReaderWriterOptions(args, argc, &argi,
			&options.ReaderOptions, &options.WriterOptions) {
			// handled
		} else if cliutil.ParseMiscOptions(args, argc, &argi, &options) {
			// handled
		} else {
			// unhandled
			usageUnrecognizedVerb(lib.MlrExeName(), args[argi])
		}
	}

	cliutil.ApplyReaderOptionDefaults(&options.ReaderOptions)
	cliutil.ApplyWriterOptionDefaults(&options.WriterOptions)

	// Set an optional global formatter for floating-point values
	if options.WriterOptions.FPOFMT != "" {
		err = types.SetMlrvalFloatOutputFormat(options.WriterOptions.FPOFMT)
		if err != nil {
			return options, recordTransformers, err
		}
	}

	recordTransformers, ignoresInput, err := parseTransformers(args, &argi, argc, &options)
	if err != nil {
		return options, recordTransformers, err
	}
	if ignoresInput {
		options.NoInput = true // e.g. then-chain begins with seqgen
	}

	if DecideFinalFlatten(&options) {
		// E.g. '{"req": {"method": "GET", "path": "/api/check"}}' becomes
		// req.method=GET,req.path=/api/check.
		transformer, err := transformers.NewTransformerFlatten(options.WriterOptions.OFLATSEP, nil)
		lib.InternalCodingErrorIf(err != nil)
		lib.InternalCodingErrorIf(transformer == nil)
		recordTransformers = append(recordTransformers, transformer)
	}

	if DecideFinalUnflatten(&options) {
		// E.g.  req.method=GET,req.path=/api/check becomes
		// '{"req": {"method": "GET", "path": "/api/check"}}'
		transformer, err := transformers.NewTransformerUnflatten(options.WriterOptions.OFLATSEP, nil)
		lib.InternalCodingErrorIf(err != nil)
		lib.InternalCodingErrorIf(transformer == nil)
		recordTransformers = append(recordTransformers, transformer)
	}

	// There may already be one or more because of --from on the command line,
	// so append.
	for _, arg := range args[argi:] {
		options.FileNames = append(options.FileNames, arg)
	}

	// E.g. mlr -n put -v '...'
	if options.NoInput {
		options.FileNames = nil
	}

	if options.DoInPlace && (options.FileNames == nil || len(options.FileNames) == 0) {
		fmt.Fprintf(os.Stderr, "%s: -I option (in-place operation) requires input files.\n", lib.MlrExeName())
		os.Exit(1)
	}

	if options.HaveRandSeed {
		lib.SeedRandom(int64(options.RandSeed))
	}

	return options, recordTransformers, nil
}

// ================================================================
// Decide whether to insert a flatten or unflatten verb at the end of the
// chain.  See also repl/verbs.go which handles the same issue in the REPL.
//
// ----------------------------------------------------------------
// PROBLEM TO BE SOLVED:
//
// JSON has nested structures and CSV et al. do not. For example:
// {
//   "req" : {
//     "method": "GET",
//     "path":   "api/check",
//   }
// }
//
// For CSV we flatten this down to
//
// {
//   "req.method": "GET",
//   "req.path":   "api/check"
// }
//
// ----------------------------------------------------------------
// APPROACH:
//
// Use the Principle of Least Surprise (POLS).
//
// * If input is JSON and output is JSON:
//   o Records can be nested from record-read
//   o They remain that way through the Miller record-processing stream
//   o They are nested on record-write
//   o No action needs to be taken
//
// * If input is JSON and output is non-JSON:
//   o Records can be nested from record-read
//   o They remain that way through the Miller record-processing stream
//   o On record-write, nested structures will be converted to string (carriage
//     returns and all) using json_stringify. People *might* want this but
//     (using POLS) we will (by default) AUTO-FLATTEN for them. There is a
//     --no-auto-unflatten CLI flag for those who want it.
//
// * If input is non-JSON and output is non-JSON:
//   o If there is a "req.method" field, people should be able to do
//     'mlr sort -f req.method' with no surprises. (Again, POLS.) Therefore
//     no auto-unflatten on input.  People can insert an unflatten verb
//     into their verb chain if they really want unflatten for non-JSON
//     files.
//   o The DSL can make nested data, so AUTO-FLATTEN at output.
//
// * If input is non-JSON and output is JSON:
//   o Default is to auto-unflatten at output.
//   o There is a --no-auto-unflatten for those who want it.
// ================================================================

func DecideFinalFlatten(options *cliutil.TOptions) bool {
	ofmt := options.WriterOptions.OutputFileFormat
	if options.WriterOptions.AutoFlatten {
		if ofmt != "json" {
			return true
		}
	}
	return false
}

func DecideFinalUnflatten(options *cliutil.TOptions) bool {
	ifmt := options.ReaderOptions.InputFileFormat
	ofmt := options.WriterOptions.OutputFileFormat

	if options.WriterOptions.AutoUnflatten {
		if ifmt != "json" {
			if ofmt == "json" {
				return true
			}
		}
	}
	return false
}

// ----------------------------------------------------------------
// Returns a list of transformers, from the starting point in args given by *pargi.
// Bumps *pargi to point to remaining post-transformer-setup args, i.e. filenames.

func parseTransformers(
	args []string,
	pargi *int,
	argc int,
	options *cliutil.TOptions,
) (
	transformerList []transforming.IRecordTransformer,
	ignoresInput bool,
	err error,
) {

	transformerList = make([]transforming.IRecordTransformer, 0)
	ignoresInput = false

	argi := *pargi

	// Allow then-chains to start with an initial 'then': 'mlr verb1 then verb2 then verb3' or
	// 'mlr then verb1 then verb2 then verb3'. Particuarly useful in backslashy scripting contexts.
	if (argc-argi) >= 1 && args[argi] == "then" {
		argi++
	}

	if (argc - argi) < 1 {
		fmt.Fprintf(os.Stderr, "%s: no verb supplied.\n", lib.MlrExeName())
		mainUsageShort()
		os.Exit(1)
	}

	onFirst := true

	for {
		cliutil.CheckArgCount(args, argi, argc, 1)
		verb := args[argi]

		transformerSetup := lookUpTransformerSetup(verb)
		if transformerSetup == nil {
			fmt.Fprintf(os.Stderr,
				"%s: verb \"%s\" not found. Please use \"%s --help\" for a list.\n",
				lib.MlrExeName(), verb, lib.MlrExeName())
			os.Exit(1)
		}

		// E.g. then-chain begins with seqgen
		if onFirst && transformerSetup.IgnoresInput {
			ignoresInput = true
		}
		onFirst = false

		// It's up to the parse func to print its usage, and exit 1, on
		// CLI-parse failure.  Also note: this assumes main reader/writer opts
		// are all parsed *before* transformer parse-CLI methods are invoked.
		transformer := transformerSetup.ParseCLIFunc(
			&argi,
			argc,
			args,
			options,
		)
		lib.InternalCodingErrorIf(transformer == nil)

		transformerList = append(transformerList, transformer)

		if argi >= argc || args[argi] != "then" {
			break
		}
		argi++
	}

	*pargi = argi
	return transformerList, ignoresInput, nil
}

// ----------------------------------------------------------------
func parseTerminalUsage(args []string, argc int, argi int) bool {
	if args[argi] == "--version" {
		fmt.Printf("Miller %s\n", version.STRING)
		return true
	} else if args[argi] == "-h" {
		mainUsageLong(os.Stdout, lib.MlrExeName())
		return true
	} else if args[argi] == "--help" {
		mainUsageLong(os.Stdout, lib.MlrExeName())
		return true
	} else if args[argi] == "--print-type-arithmetic-info" {
		fmt.Println("TODO: port printTypeArithmeticInfo")
		//		printTypeArithmeticInfo(os.Stdout, lib.MlrExeName());
		return true

	} else if args[argi] == "--help-all-verbs" || args[argi] == "--usage-all-verbs" {
		usageAllVerbs(lib.MlrExeName())
	} else if args[argi] == "--list-all-verbs" || args[argi] == "-l" {
		listAllVerbs(os.Stdout, " ")
		return true
	} else if args[argi] == "--list-all-verbs-raw" || args[argi] == "-L" {
		listAllVerbsRaw(os.Stdout)
		return true

	} else if args[argi] == "--list-all-functions-raw" || args[argi] == "-F" {
		cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionsRaw(os.Stdout)
		return true
		//	} else if args[argi] == "--list-all-functions-as-table" {
		//		fmgr_t* pfmgr = fmgr_alloc();
		//		fmgr_list_all_functions_as_table(pfmgr, os.Stdout);
		//		fmgr_free(pfmgr, nil);
		//		return true;
	} else if args[argi] == "--help-all-functions" || args[argi] == "-f" {
		cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsages(os.Stdout)
		return true
	} else if args[argi] == "--help-function" || args[argi] == "--hf" {
		cliutil.CheckArgCount(args, argi, argc, 2)
		cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsage(args[argi+1], os.Stdout)
		argi++
		return true

	} else if args[argi] == "--list-all-keywords-raw" || args[argi] == "-K" {
		fmt.Println("TOD: port mlr_dsl_list_all_keywords_raw")
		//		mlr_dsl_list_all_keywords_raw(os.Stdout);
		return true
	} else if args[argi] == "--help-all-keywords" || args[argi] == "-k" {
		fmt.Println("TOD: port mlr_dsl_list_all_keywords")
		//		mlr_dsl_keyword_usage(os.Stdout, nil);
		return true
	} else if args[argi] == "--help-keyword" || args[argi] == "--hk" {
		cliutil.CheckArgCount(args, argi, argc, 2)
		fmt.Println("TOD: port mlr_dsl_keyword_usage")
		//		mlr_dsl_keyword_usage(os.Stdout, args[argi+1]);
		return true

		//	// main-usage subsections, individually accessible for the benefit of
		//	// the manpage-autogenerator
	} else if args[argi] == "--usage-synopsis" {
		mainUsageSynopsis(os.Stdout, lib.MlrExeName())
		return true
	} else if args[argi] == "--usage-examples" {
		mainUsageExamples(os.Stdout, lib.MlrExeName(), "")
		return true
	} else if args[argi] == "--usage-list-all-verbs" {
		listAllVerbs(os.Stdout, "")
		return true
	} else if args[argi] == "--usage-help-options" {
		mainUsageHelpOptions(os.Stdout, lib.MlrExeName())
		return true
	} else if args[argi] == "--usage-mlrrc" {
		mainUsageMlrrc(os.Stdout, lib.MlrExeName())
		return true
	} else if args[argi] == "--usage-functions" {
		mainUsageFunctions(os.Stdout)
		return true
	} else if args[argi] == "--usage-data-format-examples" {
		mainUsageDataFormatExamples(os.Stdout, lib.MlrExeName())
		return true
	} else if args[argi] == "--usage-data-format-options" {
		mainUsageDataFormatOptions(os.Stdout, lib.MlrExeName())
		return true
	} else if args[argi] == "--usage-comments-in-data" {
		mainUsageCommentsInData(os.Stdout, lib.MlrExeName())
		return true
	} else if args[argi] == "--usage-format-conversion-keystroke-saver-options" {
		mainUsageFormatConversionKeystrokeSaverOptions(os.Stdout, lib.MlrExeName())
		return true
	} else if args[argi] == "--usage-compressed-data-options" {
		mainUsageCompressedDataOptions(os.Stdout, lib.MlrExeName())
		return true
		//	} else if args[argi] == "--usage-separator-options" {
		//		mainUsageSeparatorOptions(os.Stdout, lib.MlrExeName());
		//		return true;
	} else if args[argi] == "--usage-csv-options" {
		mainUsageCsvOptions(os.Stdout, lib.MlrExeName())
		return true
		//	} else if args[argi] == "--usage-double-quoting" {
		//		mainUsageDoubleQuoting(os.Stdout, lib.MlrExeName());
		//		return true;
		//	} else if args[argi] == "--usage-numerical-formatting" {
		//		mainUsageNumericalFormatting(os.Stdout, lib.MlrExeName());
		//		return true;
	} else if args[argi] == "--usage-other-options" {
		mainUsageOtherOptions(os.Stdout, lib.MlrExeName())
		return true
	} else if args[argi] == "--usage-then-chaining" {
		mainUsageThenChaining(os.Stdout, lib.MlrExeName())
		return true
	} else if args[argi] == "--usage-auxents" {
		mainUsageAuxents(os.Stdout)
		return true
	} else if args[argi] == "--usage-see-also" {
		mainUsageSeeAlso(os.Stdout, lib.MlrExeName())
		return true
	}
	return false
}

// ----------------------------------------------------------------
// * If $MLRRC is set, use it and only it.
// * Otherwise try first $HOME/.mlrrc and then ./.mlrrc but let them
//   stack: e.g. $HOME/.mlrrc is lots of settings and maybe in one
//   subdir you want to override just a setting or two.

// TODO: move to separate file?
func loadMlrrcOrDie(
	options *cliutil.TOptions,
) {
	env_mlrrc := os.Getenv("MLRRC")

	if env_mlrrc != "" {
		if env_mlrrc == "__none__" {
			return
		}
		if tryLoadMlrrc(options, env_mlrrc) {
			return
		}
	}

	env_home := os.Getenv("HOME")
	if env_home != "" {
		path := env_home + "/.mlrrc"
		tryLoadMlrrc(options, path)
	}

	tryLoadMlrrc(options, "./.mlrrc")
}

func tryLoadMlrrc(
	options *cliutil.TOptions,
	path string,
) bool {
	handle, err := os.Open(path)
	if err != nil {
		return false
	}
	defer handle.Close()

	lineReader := bufio.NewReader(handle)

	eof := false
	lineno := 0
	for !eof {
		line, err := lineReader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
			break
		}
		lineno++

		if err != nil {
			fmt.Fprintln(os.Stderr, lib.MlrExeName(), err)
			os.Exit(1)
			return false
		}

		// This is how to do a chomp:
		// TODO: handle \r\n with libified solution.
		line = strings.TrimRight(line, "\n")

		if !handleMlrrcLine(options, line) {
			fmt.Fprintf(os.Stderr, "%s: parse error at file \"%s\" line %d: %s\n",
				lib.MlrExeName(), path, lineno, line,
			)
			os.Exit(1)
		}
	}

	return true
}

func handleMlrrcLine(
	options *cliutil.TOptions,
	line string,
) bool {

	// Comment-strip
	re := regexp.MustCompile("#.*")
	line = re.ReplaceAllString(line, "")

	// Left-trim / right-trim
	line = strings.TrimSpace(line)

	if line == "" { // line was whitespace-only
		return true
	}

	// Prepend initial "--" if it's not already there
	if !strings.HasPrefix(line, "-") {
		line = "--" + line
	}

	// Split line into args array
	args := strings.Fields(line)
	argi := 0
	argc := len(args)

	if args[0] == "--prepipe" || args[0] == "--prepipex" {
		// Don't allow code execution via .mlrrc
		return false
	} else if args[0] == "--load" || args[0] == "--mload" {
		// Don't allow code execution via .mlrrc
		return false
	} else if cliutil.ParseReaderOptions(args, argc, &argi, &options.ReaderOptions) {
		// handled
	} else if cliutil.ParseWriterOptions(args, argc, &argi, &options.WriterOptions) {
		// handled
	} else if cliutil.ParseReaderWriterOptions(args, argc, &argi,
		&options.ReaderOptions, &options.WriterOptions) {
		// handled
	} else if cliutil.ParseMiscOptions(args, argc, &argi, options) {
		// handled
	} else {
		return false
	}

	return true
}
