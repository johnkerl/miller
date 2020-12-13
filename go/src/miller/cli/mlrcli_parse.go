package cli

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/dsl/cst"
	"miller/lib"
	"miller/transforming"
	"miller/version"
)

// ----------------------------------------------------------------
func ParseCommandLine(args []string) (
	options clitypes.TOptions,
	recordTransformers []transforming.IRecordTransformer,
	err error,
) {
	options = clitypes.DefaultOptions()
	argc := len(args)
	argi := 1

	// Try .mlrrc overrides (then command-line on top of that).
	// A --norc flag (if provided) must come before all other options.
	// Or, they can set the environment variable MLRRC="__none__".
	//	if (argc >= 2 && args[1] == "--norc" {
	//		argi++;
	//	} else {
	//		loadMlrrcOrDie(popts);
	//	}

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if args[argi][0] != '-' {
			break // No more flag options to process
		} else if args[argi] == "--cpuprofile" {
			// Already handled in main(); ignore here.
			clitypes.CheckArgCount(args, argi, argc, 1)
			argi += 2
		} else if parseTerminalUsage(args, argc, argi) {
			os.Exit(0)
		} else if clitypes.ParseReaderOptions(args, argc, &argi, &options.ReaderOptions) {
			// handled
		} else if clitypes.ParseWriterOptions(args, argc, &argi, &options.WriterOptions) {
			// handled
		} else if clitypes.ParseReaderWriterOptions(args, argc, &argi,
			&options.ReaderOptions, &options.WriterOptions) {
			// handled
		} else if clitypes.ParseMiscOptions(args, argc, &argi, &options) {
			// handled
		} else {
			// unhandled
			usageUnrecognizedVerb(args[0], args[argi])
		}
	}

	// xxx to do
	//	cli_apply_defaults(popts);

	//	lhmss_t* default_rses = get_default_rses();
	//	lhmss_t* default_fses = get_default_fses();
	//	lhmss_t* default_pses = get_default_pses();
	//	lhmsll_t* default_repeat_ifses = get_default_repeat_ifses();
	//	lhmsll_t* default_repeat_ipses = get_default_repeat_ipses();
	//
	//	if (options.ReaderOptions.IRS == nil)
	//		options.ReaderOptions.IRS = lhmss_get_or_die(default_rses, options.ReaderOptions.InputFileFormat);
	//	if (options.ReaderOptions.IFS == nil)
	//		options.ReaderOptions.IFS = lhmss_get_or_die(default_fses, options.ReaderOptions.InputFileFormat);
	//	if (options.ReaderOptions.ips == nil)
	//		options.ReaderOptions.ips = lhmss_get_or_die(default_pses, options.ReaderOptions.InputFileFormat);
	//
	//	if (options.ReaderOptions.allow_repeat_ifs == NEITHER_TRUE_NOR_FALSE)
	//		options.ReaderOptions.allow_repeat_ifs = lhmsll_get_or_die(default_repeat_ifses, options.ReaderOptions.InputFileFormat);
	//	if (options.ReaderOptions.allow_repeat_ips == NEITHER_TRUE_NOR_FALSE)
	//		options.ReaderOptions.allow_repeat_ips = lhmsll_get_or_die(default_repeat_ipses, options.ReaderOptions.InputFileFormat);
	//
	//	if (options.WriterOptions.ORS == nil)
	//		options.WriterOptions.ORS = lhmss_get_or_die(default_rses, options.WriterOptions.OutputFileFormat);
	//	if (options.WriterOptions.OFS == nil)
	//		options.WriterOptions.OFS = lhmss_get_or_die(default_fses, options.WriterOptions.OutputFileFormat);
	//	if (options.WriterOptions.ops == nil)
	//		options.WriterOptions.ops = lhmss_get_or_die(default_pses, options.WriterOptions.OutputFileFormat);
	//
	//	if options.WriterOptions.OutputFileFormat == "pprint") && len(options.WriterOptions.OFS) != 1) {
	//		fmt.Fprintf(os.Stderr, "%s: OFS for PPRINT format must be single-character; got \"%s\".\n",
	//			os.Args[0], options.WriterOptions.OFS);
	//		return nil;
	//	}

	//	// Construct the transformer list for single use, e.g. the normal streaming case wherein the
	//	// transformers operate on all input files. Also retain information needed to construct them
	//	// for each input file, for in-place mode.
	//	options.transformer_argb = argi;
	//	options.original_argv = args;
	//	options.non_in_place_argv = copy_argv(args);
	//	options.argc = argc;
	//	*pptransformer_list = cli_parse_transformers(options.non_in_place_argv, &argi, argc, popts);

	recordTransformers, ignoresInput, err := parseTransformers(args, &argi, argc, &options)
	if err != nil {
		return options, recordTransformers, err
	}
	if ignoresInput {
		options.NoInput = true // e.g. then-chain begins with seqgen
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

	//	if (options.do_in_place && (options.FileNames == nil || options.FileNames.length == 0)) {
	//		fmt.Fprintf(os.Stderr, "%s: -I option (in-place operation) requires input files.\n", os.Args[0]);
	//		os.Exit(1);
	//	}

	if options.HaveRandSeed {
		lib.SeedRandom(options.RandSeed)
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
	options *clitypes.TOptions,
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
		fmt.Fprintf(os.Stderr, "%s: no verb supplied.\n", os.Args[0])
		mainUsageShort()
		os.Exit(1)
	}

	onFirst := true

	for {
		clitypes.CheckArgCount(args, argi, argc, 1)
		verb := args[argi]

		transformerSetup := lookUpTransformerSetup(verb)
		if transformerSetup == nil {
			fmt.Fprintf(os.Stderr,
				"%s: verb \"%s\" not found. Please use \"%s --help\" for a list.\n",
				os.Args[0], verb, os.Args[0])
			os.Exit(1)
		}

		// E.g. then-chain begins with seqgen
		if onFirst && transformerSetup.IgnoresInput {
			ignoresInput = true
		}
		onFirst = false

		// It's up to the parse func to print its usage on CLI-parse failure.
		// Also note: this assumes main reader/writer opts are all parsed
		// *before* transformer parse-CLI methods are invoked.
		transformer := transformerSetup.ParseCLIFunc(
			&argi,
			argc,
			args,
			flag.ExitOnError,
			&options.ReaderOptions,
			&options.WriterOptions,
		)

		if transformer == nil {
			// Error message already printed out
			os.Exit(1)
		}

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
		mainUsageLong(os.Stdout, os.Args[0])
		return true
	} else if args[argi] == "--help" {
		mainUsageLong(os.Stdout, os.Args[0])
		return true
		//	} else if args[argi] == "--print-type-arithmetic-info" {
		//		printTypeArithmeticInfo(os.Stdout, os.Args[0]);
		//		return true;
		//
	} else if args[argi] == "--help-all-verbs" || args[argi] == "--usage-all-verbs" {
		usageAllVerbs(os.Args[0])
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
		//	} else if args[argi] == "--help-function" || args[argi] == "--hf" {
		//		clitypes.CheckArgCount(args, argi, argc, 2);
		//		fmgr_t* pfmgr = fmgr_alloc();
		//		fmgr_function_usage(pfmgr, os.Stdout, args[argi+1]);
		//		fmgr_free(pfmgr, nil);
		//		return true;
		//
		//	} else if args[argi] == "--list-all-keywords-raw" || args[argi] == "-K" {
		//		mlr_dsl_list_all_keywords_raw(os.Stdout);
		//		return true;
		//	} else if args[argi] == "--help-all-keywords" || args[argi] == "-k" {
		//		mlr_dsl_keyword_usage(os.Stdout, nil);
		//		return true;
		//	} else if args[argi] == "--help-keyword" || args[argi] == "--hk" {
		//		clitypes.CheckArgCount(args, argi, argc, 2);
		//		mlr_dsl_keyword_usage(os.Stdout, args[argi+1]);
		//		return true;
		//
		//	// main-usage subsections, individually accessible for the benefit of
		//	// the manpage-autogenerator
	} else if args[argi] == "--usage-synopsis" {
		mainUsageSynopsis(os.Stdout, os.Args[0])
		return true
	} else if args[argi] == "--usage-examples" {
		mainUsageExamples(os.Stdout, os.Args[0], "")
		return true
	} else if args[argi] == "--usage-list-all-verbs" {
		listAllVerbs(os.Stdout, "")
		return true
	} else if args[argi] == "--usage-help-options" {
		mainUsageHelpOptions(os.Stdout, os.Args[0])
		return true
		//	} else if args[argi] == "--usage-mlrrc" {
		//		mainUsageMlrrc(os.Stdout, os.Args[0]);
		//		return true;
		//	} else if args[argi] == "--usage-functions" {
		//		mainUsageFunctions(os.Stdout, os.Args[0], "");
		//		return true;
		//	} else if args[argi] == "--usage-data-format-examples" {
		mainUsageDataFormatExamples(os.Stdout, os.Args[0])
		//		return true;
	} else if args[argi] == "--usage-data-format-options" {
		mainUsageDataFormatOptions(os.Stdout, os.Args[0])
		return true
		//	} else if args[argi] == "--usage-comments-in-data" {
		//		mainUsageCommentsInData(os.Stdout, os.Args[0]);
		//		return true;
	} else if args[argi] == "--usage-format-conversion-keystroke-saver-options" {
		mainUsageFormatConversionKeystrokeSaverOptions(os.Stdout, os.Args[0])
		return true
		//	} else if args[argi] == "--usage-compressed-data-options" {
		//		mainUsageCompressedDataOptions(os.Stdout, os.Args[0]);
		//		return true;
		//	} else if args[argi] == "--usage-separator-options" {
		//		mainUsageSeparatorOptions(os.Stdout, os.Args[0]);
		//		return true;
		//	} else if args[argi] == "--usage-csv-options" {
		//		mainUsageCsvOptions(os.Stdout, os.Args[0]);
		//		return true;
		//	} else if args[argi] == "--usage-double-quoting" {
		//		mainUsageDoubleQuoting(os.Stdout, os.Args[0]);
		//		return true;
		//	} else if args[argi] == "--usage-numerical-formatting" {
		//		mainUsageNumericalFormatting(os.Stdout, os.Args[0]);
		//		return true;
		//	} else if args[argi] == "--usage-other-options" {
		//		mainUsageOtherOptions(os.Stdout, os.Args[0]);
		//		return true;
	} else if args[argi] == "--usage-then-chaining" {
		mainUsageThenChaining(os.Stdout, os.Args[0])
		return true
		//	} else if args[argi] == "--usage-auxents" {
		//		mainUsageAuxents(os.Stdout, os.Args[0]);
		//		return true;
	} else if args[argi] == "--usage-see-also" {
		mainUsageSeeAlso(os.Stdout, os.Args[0])
		return true
	}
	return false
}
