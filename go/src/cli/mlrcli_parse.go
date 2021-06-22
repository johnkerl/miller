package cli

import (
	"fmt"
	"os"

	"miller/src/auxents/help"
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
			fmt.Fprintf(os.Stderr, "%s: option \"%s\" not recognized.\n", lib.MlrExeName(), args[argi])
			fmt.Fprintf(os.Stderr, "Please run \"%s --help\" for usage information.\n", lib.MlrExeName())
			os.Exit(1)
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

	if cliutil.DecideFinalFlatten(&options) {
		// E.g. '{"req": {"method": "GET", "path": "/api/check"}}' becomes
		// req.method=GET,req.path=/api/check.
		transformer, err := transformers.NewTransformerFlatten(options.WriterOptions.OFLATSEP, nil)
		lib.InternalCodingErrorIf(err != nil)
		lib.InternalCodingErrorIf(transformer == nil)
		recordTransformers = append(recordTransformers, transformer)
	}

	if cliutil.DecideFinalUnflatten(&options) {
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
		help.MainUsage(os.Stderr)
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

		// argi now points to:
		// * A "then", followed by the start of the next verb in the chain, if any
		// * Filenames after the verb (there are no more verbs listed)
		// * None of the above; argi == argc
		if argi >= argc {
			break
		} else if args[argi] != "then" {
			break
		} else {
			if argi == argc-1 {
				fmt.Fprintf(os.Stderr, "%s: missing next verb after \"then\".\n", lib.MlrExeName())
				os.Exit(1)
			} else {
				argi++
			}
		}
	}

	*pargi = argi
	return transformerList, ignoresInput, nil
}

// ----------------------------------------------------------------
// TODO: move to src/auxents/help -- ?
func parseTerminalUsage(args []string, argc int, argi int) bool {
	if args[argi] == "--version" {
		fmt.Printf("Miller %s\n", version.STRING)
		return true
	} else if args[argi] == "-h" || args[argi] == "--help" {
		help.MainUsage(os.Stdout)
		os.Exit(0)
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
		fmt.Println("TODO: port mlr_dsl_list_all_keywords_raw")
		//		mlr_dsl_list_all_keywords_raw(os.Stdout);
		return true
	} else if args[argi] == "--help-all-keywords" || args[argi] == "-k" {
		fmt.Println("TODO: port mlr_dsl_list_all_keywords")
		//		mlr_dsl_keyword_usage(os.Stdout, nil);
		return true
	} else if args[argi] == "--help-keyword" || args[argi] == "--hk" {
		cliutil.CheckArgCount(args, argi, argc, 2)
		fmt.Println("TODO: port mlr_dsl_keyword_usage")
		//		mlr_dsl_keyword_usage(os.Stdout, args[argi+1]);
		return true

		// main-usage subsections, individually accessible for the benefit of
		// the manpage-autogenerator
	} else if args[argi] == "--usage-list-all-verbs" {
		listAllVerbs(os.Stdout, "")
		return true
	} else if args[argi] == "--usage-help-options" {
		help.MainUsage(os.Stdout)
		return true
	} else if args[argi] == "--usage-functions" {
		help.ListBuiltinFunctions(os.Stdout)
		return true
	}
	return false
}
