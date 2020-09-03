package cli

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/mapping"
	"miller/version"
)

// ----------------------------------------------------------------
func ParseCommandLine(args []string) (
	options clitypes.TOptions,
	recordMappers []mapping.IRecordMapper,
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
			checkArgCount(args, argi, argc, 1)
			argi += 2
		} else if parseTerminalUsage(args, argc, argi) {
			os.Exit(0)
		} else if parseReaderOptions(args, argc, &argi, &options.ReaderOptions) {
			// handled
		} else if parseWriterOptions(args, argc, &argi, &options.WriterOptions) {
			// handled
		} else if parseReaderWriterOptions(args, argc, &argi,
			&options.ReaderOptions, &options.WriterOptions) {
			// handled
		} else if parseMiscOptions(args, argc, &argi, &options) {
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

	//	// Construct the mapper list for single use, e.g. the normal streaming case wherein the
	//	// mappers operate on all input files. Also retain information needed to construct them
	//	// for each input file, for in-place mode.
	//	options.mapper_argb = argi;
	//	options.original_argv = args;
	//	options.non_in_place_argv = copy_argv(args);
	//	options.argc = argc;
	//	*ppmapper_list = cli_parse_mappers(options.non_in_place_argv, &argi, argc, popts);
	recordMappers, err = parseMappers(args, &argi, argc, &options)
	if err != nil {
		return options, recordMappers, err
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

	//	if (options.have_rand_seed) {
	//		mtrand_init(options.rand_seed);
	//	} else {
	//		mtrand_init_default();
	//	}

	return options, recordMappers, nil
}

// ----------------------------------------------------------------
// Returns a list of mappers, from the starting point in args given by *pargi.
// Bumps *pargi to point to remaining post-mapper-setup args, i.e. filenames.

func parseMappers(
	args []string,
	pargi *int,
	argc int,
	options *clitypes.TOptions,
) ([]mapping.IRecordMapper, error) {

	mapperList := make([]mapping.IRecordMapper, 0)
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

	for {
		checkArgCount(args, argi, argc, 1)
		verb := args[argi]

		mapperSetup := lookUpMapperSetup(verb)
		if mapperSetup == nil {
			fmt.Fprintf(os.Stderr,
				"%s: verb \"%s\" not found. Please use \"%s --help\" for a list.\n",
				os.Args[0], verb, os.Args[0])
			os.Exit(1)
		}

		// It's up to the parse func to print its usage on CLI-parse failure.
		// Also note: this assumes main reader/writer opts are all parsed
		// *before* mapper parse-CLI methods are invoked.

		mapper := mapperSetup.ParseCLIFunc(
			&argi,
			argc,
			args,
			flag.ExitOnError,
			&options.ReaderOptions,
			&options.WriterOptions,
		)

		if mapper == nil {
			// Error message already printed out
			os.Exit(1)
		}

		//		if (mapperSetup.IgnoresInput && len(mapperList) == 0) {
		//			// e.g. then-chain starts with seqgen
		//			options.no_input = true;
		//		}

		mapperList = append(mapperList, mapper)

		if argi >= argc || args[argi] != "then" {
			break
		}
		argi++
	}

	*pargi = argi
	return mapperList, nil
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

		//	} else if args[argi] == "--list-all-functions-raw" || args[argi] == "-F" {
		//		fmgr_t* pfmgr = fmgr_alloc();
		//		fmgr_list_all_functions_raw(pfmgr, os.Stdout);
		//		fmgr_free(pfmgr, nil);
		//		return true;
		//	} else if args[argi] == "--list-all-functions-as-table" {
		//		fmgr_t* pfmgr = fmgr_alloc();
		//		fmgr_list_all_functions_as_table(pfmgr, os.Stdout);
		//		fmgr_free(pfmgr, nil);
		//		return true;
		//	} else if args[argi] == "--help-all-functions" || args[argi] == "-f" {
		//		fmgr_t* pfmgr = fmgr_alloc();
		//		fmgr_function_usage(pfmgr, os.Stdout, nil);
		//		fmgr_free(pfmgr, nil);
		//		return true;
		//	} else if args[argi] == "--help-function" || args[argi] == "--hf" {
		//		checkArgCount(args, argi, argc, 2);
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
		//		checkArgCount(args, argi, argc, 2);
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

// Returns true if the current flag was handled.
func parseReaderOptions(args []string, argc int, pargi *int, readerOptions *clitypes.TReaderOptions) bool {
	argi := *pargi
	oargi := argi

	if args[argi] == "--irs" {
		checkArgCount(args, argi, argc, 2)
		readerOptions.IRS = SeparatorFromArg(args[argi+1])
		argi += 2

	} else if args[argi] == "--ifs" {
		checkArgCount(args, argi, argc, 2)
		readerOptions.IFS = SeparatorFromArg(args[argi+1])
		argi += 2

	} else if args[argi] == "--ips" {
		checkArgCount(args, argi, argc, 2)
		readerOptions.IPS = SeparatorFromArg(args[argi+1])
		argi += 2

		//	} else if args[argi] == "--repifs" {
		//		readerOptions.allow_repeat_ifs = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--json-fatal-arrays-on-input" {
		//		readerOptions.json_array_ingest = JSON_ARRAY_INGEST_FATAL;
		//		argi += 1;
		//	} else if args[argi] == "--json-skip-arrays-on-input" {
		//		readerOptions.json_array_ingest = JSON_ARRAY_INGEST_SKIP;
		//		argi += 1;
		//	} else if args[argi] == "--json-map-arrays-on-input" {
		//		readerOptions.json_array_ingest = JSON_ARRAY_INGEST_AS_MAP;
		//		argi += 1;
		//
		//	} else if args[argi] == "--implicit-csv-header" {
		//		readerOptions.use_implicit_csv_header = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--no-implicit-csv-header" {
		//		readerOptions.use_implicit_csv_header = false;
		//		argi += 1;
		//
		//	} else if args[argi] == "--allow-ragged-csv-input" || args[argi] == "--ragged" {
		//		readerOptions.allow_ragged_csv_input = true;
		//		argi += 1;

		//	} else if args[argi] == "-i" {
		//		checkArgCount(args, argi, argc, 2);
		//		if (!lhmss_has_key(get_default_rses(), args[argi+1])) {
		//			fmt.Fprintf(os.Stderr, "%s: unrecognized input format \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//			os.Exit(1);
		//		}
		//		readerOptions.InputFileFormat = args[argi+1];
		//		argi += 2;
		//
		//	} else if args[argi] == "--igen" {
		//		readerOptions.InputFileFormat = "gen";
		//		argi += 1;
		//	} else if args[argi] == "--gen-start" {
		//		readerOptions.InputFileFormat = "gen";
		//		checkArgCount(args, argi, argc, 2);
		//		if (sscanf(args[argi+1], "%lld", &readerOptions.generator_opts.start) != 1) {
		//			fmt.Fprintf(os.Stderr, "%s: could not scan \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//		}
		//		argi += 2;
		//	} else if args[argi] == "--gen-stop" {
		//		readerOptions.InputFileFormat = "gen";
		//		checkArgCount(args, argi, argc, 2);
		//		if (sscanf(args[argi+1], "%lld", &readerOptions.generator_opts.stop) != 1) {
		//			fmt.Fprintf(os.Stderr, "%s: could not scan \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//		}
		//		argi += 2;
		//	} else if args[argi] == "--gen-step" {
		//		readerOptions.InputFileFormat = "gen";
		//		checkArgCount(args, argi, argc, 2);
		//		if (sscanf(args[argi+1], "%lld", &readerOptions.generator_opts.step) != 1) {
		//			fmt.Fprintf(os.Stderr, "%s: could not scan \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//		}
		//		argi += 2;

	} else if args[argi] == "--icsv" {
		readerOptions.InputFileFormat = "csv"
		argi += 1

		//	} else if args[argi] == "--icsvlite" {
		//		readerOptions.InputFileFormat = "csvlite";
		//		argi += 1;
		//
		//	} else if args[argi] == "--itsv" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS = "\t";
		//		argi += 1;
		//
		//	} else if args[argi] == "--itsvlite" {
		//		readerOptions.InputFileFormat = "csvlite";
		//		readerOptions.IFS = "\t";
		//		argi += 1;
		//
		//	} else if args[argi] == "--iasv" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS = ASV_FS;
		//		readerOptions.IRS = ASV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--iasvlite" {
		//		readerOptions.InputFileFormat = "csvlite";
		//		readerOptions.IFS = ASV_FS;
		//		readerOptions.IRS = ASV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--iusv" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS = USV_FS;
		//		readerOptions.IRS = USV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--iusvlite" {
		//		readerOptions.InputFileFormat = "csvlite";
		//		readerOptions.IFS = USV_FS;
		//		readerOptions.IRS = USV_RS;
		//		argi += 1;
		//
	} else if args[argi] == "--idkvp" {
		readerOptions.InputFileFormat = "dkvp"
		argi += 1

	} else if args[argi] == "--ijson" {
		readerOptions.InputFileFormat = "json"
		argi += 1

	} else if args[argi] == "--inidx" {
		readerOptions.InputFileFormat = "nidx"
		argi += 1

		//	} else if args[argi] == "--ixtab" {
		//		readerOptions.InputFileFormat = "xtab";
		//		argi += 1;
		//
		//	} else if args[argi] == "--ipprint" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--mmap" {
		//		// No-op as of 5.6.3 (mmap is being abandoned) but don't break
		//		// the command-line user experience.
		//		argi += 1;
		//
		//	} else if args[argi] == "--no-mmap" {
		//		// No-op as of 5.6.3 (mmap is being abandoned) but don't break
		//		// the command-line user experience.
		//		argi += 1;
		//
		//	} else if args[argi] == "--prepipe" {
		//		checkArgCount(args, argi, argc, 2);
		//		readerOptions.prepipe = args[argi+1];
		//		argi += 2;
		//
		//	} else if args[argi] == "--skip-comments" {
		//		readerOptions.comment_string = DEFAULT_COMMENT_STRING;
		//		readerOptions.comment_handling = SKIP_COMMENTS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--skip-comments-with" {
		//		checkArgCount(args, argi, argc, 2);
		//		readerOptions.comment_string = args[argi+1];
		//		readerOptions.comment_handling = SKIP_COMMENTS;
		//		argi += 2;
		//
		//	} else if args[argi] == "--pass-comments" {
		//		readerOptions.comment_string = DEFAULT_COMMENT_STRING;
		//		readerOptions.comment_handling = PASS_COMMENTS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--pass-comments-with" {
		//		checkArgCount(args, argi, argc, 2);
		//		readerOptions.comment_string = args[argi+1];
		//		readerOptions.comment_handling = PASS_COMMENTS;
		//		argi += 2;
		//
	}
	*pargi = argi
	return argi != oargi
}

// Returns true if the current flag was handled.
func parseWriterOptions(args []string, argc int, pargi *int, writerOptions *clitypes.TWriterOptions) bool {
	argi := *pargi
	oargi := argi

	if args[argi] == "--ors" {
		checkArgCount(args, argi, argc, 2)
		writerOptions.ORS = SeparatorFromArg(args[argi+1])
		argi += 2

	} else if args[argi] == "--ofs" {
		checkArgCount(args, argi, argc, 2)
		writerOptions.OFS = SeparatorFromArg(args[argi+1])
		argi += 2

		//	} else if args[argi] == "--headerless-csv-output" {
		//		writerOptions.headerless_csv_output = true;
		//		argi += 1;
		//
	} else if args[argi] == "--ops" {
		checkArgCount(args, argi, argc, 2)
		writerOptions.OPS = SeparatorFromArg(args[argi+1])
		argi += 2

		//	} else if args[argi] == "--xvright" {
		//		writerOptions.right_justify_xtab_value = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--jvstack" {
		//		writerOptions.stack_json_output_vertically = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--jlistwrap" {
		//		writerOptions.wrap_json_output_in_outer_list = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--jknquoteint" {
		//		writerOptions.json_quote_int_keys = false;
		//		argi += 1;
		//	} else if args[argi] == "--jquoteall" {
		//		writerOptions.json_quote_non_string_values = true;
		//		argi += 1;
		//	} else if args[argi] == "--jvquoteall" {
		//		writerOptions.json_quote_non_string_values = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--vflatsep" {
		//		checkArgCount(args, argi, argc, 2);
		//		writerOptions.oosvar_flatten_separator = SeparatorFromArg(args[argi+1]);
		//		argi += 2;
		//
		//	} else if args[argi] == "-o" {
		//		checkArgCount(args, argi, argc, 2);
		//		if (!lhmss_has_key(get_default_rses(), args[argi+1])) {
		//			fmt.Fprintf(os.Stderr, "%s: unrecognized output format \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//			os.Exit(1);
		//		}
		//		writerOptions.OutputFileFormat = args[argi+1];
		//		argi += 2;
		//
	} else if args[argi] == "--ocsv" {
		writerOptions.OutputFileFormat = "csv"
		argi += 1

		//	} else if args[argi] == "--ocsvlite" {
		//		writerOptions.OutputFileFormat = "csvlite";
		//		argi += 1;
		//
		//	} else if args[argi] == "--otsv" {
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.OFS = "\t";
		//		argi += 1;
		//
		//	} else if args[argi] == "--otsvlite" {
		//		writerOptions.OutputFileFormat = "csvlite";
		//		writerOptions.OFS = "\t";
		//		argi += 1;
		//
		//	} else if args[argi] == "--oasv" {
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.OFS = ASV_FS;
		//		writerOptions.ORS = ASV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--oasvlite" {
		//		writerOptions.OutputFileFormat = "csvlite";
		//		writerOptions.OFS = ASV_FS;
		//		writerOptions.ORS = ASV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--ousv" {
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.OFS = USV_FS;
		//		writerOptions.ORS = USV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--ousvlite" {
		//		writerOptions.OutputFileFormat = "csvlite";
		//		writerOptions.OFS = USV_FS;
		//		writerOptions.ORS = USV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--omd" {
		//		writerOptions.OutputFileFormat = "markdown";
		//		argi += 1;
		//
	} else if args[argi] == "--odkvp" {
		writerOptions.OutputFileFormat = "dkvp"
		argi += 1

	} else if args[argi] == "--ojson" {
		writerOptions.OutputFileFormat = "json"
		argi += 1
		//	} else if args[argi] == "--ojsonx" {
		//		writerOptions.OutputFileFormat = "json";
		//		writerOptions.stack_json_output_vertically = true;
		//		argi += 1;

	} else if args[argi] == "--onidx" {
		writerOptions.OutputFileFormat = "nidx"
		argi += 1

	} else if args[argi] == "--oxtab" {
		writerOptions.OutputFileFormat = "xtab"
		argi += 1

	} else if args[argi] == "--opprint" {
		writerOptions.OutputFileFormat = "pprint"
		argi += 1

		//	} else if args[argi] == "--right" {
		//		writerOptions.right_align_pprint = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--barred" {
		//		writerOptions.pprint_barred = true;
		//		argi += 1;
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
	}
	*pargi = argi
	return argi != oargi
}

// Returns true if the current flag was handled.
func parseReaderWriterOptions(
	args []string,
	argc int,
	pargi *int,
	readerOptions *clitypes.TReaderOptions,
	writerOptions *clitypes.TWriterOptions,
) bool {

	argi := *pargi
	oargi := argi

	if args[argi] == "--rs" {
		checkArgCount(args, argi, argc, 2)
		readerOptions.IRS = SeparatorFromArg(args[argi+1])
		writerOptions.ORS = SeparatorFromArg(args[argi+1])
		argi += 2

	} else if args[argi] == "--fs" {
		checkArgCount(args, argi, argc, 2)
		readerOptions.IFS = SeparatorFromArg(args[argi+1])
		writerOptions.OFS = SeparatorFromArg(args[argi+1])
		argi += 2

		//	} else if args[argi] == "-p" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "nidx";
		//		readerOptions.IFS = " ";
		//		writerOptions.OFS = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		argi += 1;
		//
	} else if args[argi] == "--ps" {
		checkArgCount(args, argi, argc, 2)
		readerOptions.IPS = SeparatorFromArg(args[argi+1])
		writerOptions.OPS = SeparatorFromArg(args[argi+1])
		argi += 2

		//	} else if args[argi] == "--jflatsep" {
		//		checkArgCount(args, argi, argc, 2);
		//		readerOptions.input_json_flatten_separator  = SeparatorFromArg(args[argi+1]);
		//		writerOptions.output_json_flatten_separator = SeparatorFromArg(args[argi+1]);
		//		argi += 2;
		//
		//	} else if args[argi] == "--io" {
		//		checkArgCount(args, argi, argc, 2);
		//		if (!lhmss_has_key(get_default_rses(), args[argi+1])) {
		//			fmt.Fprintf(os.Stderr, "%s: unrecognized I/O format \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//			os.Exit(1);
		//		}
		//		readerOptions.InputFileFormat = args[argi+1];
		//		writerOptions.OutputFileFormat = args[argi+1];
		//		argi += 2;
		//
	} else if args[argi] == "--csv" {
		readerOptions.InputFileFormat = "csv"
		writerOptions.OutputFileFormat = "csv"
		argi += 1

		//	} else if args[argi] == "--csvlite" {
		//		readerOptions.InputFileFormat = "csvlite";
		//		writerOptions.OutputFileFormat = "csvlite";
		//		argi += 1;
		//
		//	} else if args[argi] == "--tsv" {
		//		readerOptions.InputFileFormat = writerOptions.OutputFileFormat = "csv";
		//		readerOptions.IFS = "\t";
		//		writerOptions.OFS = "\t";
		//		argi += 1;
		//
		//	} else if args[argi] == "--tsvlite" || args[argi] == "-t" {
		//		readerOptions.InputFileFormat = writerOptions.OutputFileFormat = "csvlite";
		//		readerOptions.IFS = "\t";
		//		writerOptions.OFS = "\t";
		//		argi += 1;
		//
		//	} else if args[argi] == "--asv" {
		//		readerOptions.InputFileFormat = writerOptions.OutputFileFormat = "csv";
		//		readerOptions.IFS = ASV_FS;
		//		writerOptions.OFS = ASV_FS;
		//		readerOptions.IRS = ASV_RS;
		//		writerOptions.ORS = ASV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--asvlite" {
		//		readerOptions.InputFileFormat = writerOptions.OutputFileFormat = "csvlite";
		//		readerOptions.IFS = ASV_FS;
		//		writerOptions.OFS = ASV_FS;
		//		readerOptions.IRS = ASV_RS;
		//		writerOptions.ORS = ASV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--usv" {
		//		readerOptions.InputFileFormat = writerOptions.OutputFileFormat = "csv";
		//		readerOptions.IFS = USV_FS;
		//		writerOptions.OFS = USV_FS;
		//		readerOptions.IRS = USV_RS;
		//		writerOptions.ORS = USV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--usvlite" {
		//		readerOptions.InputFileFormat = writerOptions.OutputFileFormat = "csvlite";
		//		readerOptions.IFS = USV_FS;
		//		writerOptions.OFS = USV_FS;
		//		readerOptions.IRS = USV_RS;
		//		writerOptions.ORS = USV_RS;
		//		argi += 1;
		//
	} else if args[argi] == "--dkvp" {
		readerOptions.InputFileFormat = "dkvp"
		writerOptions.OutputFileFormat = "dkvp"
		argi += 1

	} else if args[argi] == "--json" {
		readerOptions.InputFileFormat = "json"
		writerOptions.OutputFileFormat = "json"
		argi += 1
		//	} else if args[argi] == "--jsonx" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "json";
		//		writerOptions.stack_json_output_vertically = true;
		//		argi += 1;
		//
	} else if args[argi] == "--nidx" {
		readerOptions.InputFileFormat = "nidx"
		writerOptions.OutputFileFormat = "nidx"
		argi += 1

		//	} else if args[argi] == "-T" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "nidx";
		//		readerOptions.IFS = "\t";
		//		writerOptions.OFS = "\t";
		//		argi += 1;
		//
		//	} else if args[argi] == "--xtab" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "xtab";
		//		argi += 1;
		//
		//	} else if args[argi] == "--pprint" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat        = "pprint";
		//		argi += 1;
		//
		//	} else if args[argi] == "--c2t" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		writerOptions.OFS       = "\t";
		//		argi += 1;
		//	} else if args[argi] == "--c2d" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "dkvp";
		//		argi += 1;
		//	} else if args[argi] == "--c2n" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "nidx";
		//		argi += 1;
		//	} else if args[argi] == "--c2j" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "json";
		//		argi += 1;
		//	} else if args[argi] == "--c2p" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "pprint";
		//		argi += 1;
		//	} else if args[argi] == "--c2x" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "xtab";
		//		argi += 1;
		//	} else if args[argi] == "--c2m" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "markdown";
		//		argi += 1;
		//
		//	} else if args[argi] == "--t2c" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS       = "\t";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		argi += 1;
		//	} else if args[argi] == "--t2d" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS       = "\t";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "dkvp";
		//		argi += 1;
		//	} else if args[argi] == "--t2n" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS       = "\t";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "nidx";
		//		argi += 1;
		//	} else if args[argi] == "--t2j" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS       = "\t";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "json";
		//		argi += 1;
		//	} else if args[argi] == "--t2p" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS       = "\t";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "pprint";
		//		argi += 1;
		//	} else if args[argi] == "--t2x" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS       = "\t";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "xtab";
		//		argi += 1;
		//	} else if args[argi] == "--t2m" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS       = "\t";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "markdown";
		//		argi += 1;
		//
		//	} else if args[argi] == "--d2c" {
		//		readerOptions.InputFileFormat = "dkvp";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		argi += 1;
		//	} else if args[argi] == "--d2t" {
		//		readerOptions.InputFileFormat = "dkvp";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		writerOptions.OFS       = "\t";
		//		argi += 1;
		//	} else if args[argi] == "--d2n" {
		//		readerOptions.InputFileFormat = "dkvp";
		//		writerOptions.OutputFileFormat = "nidx";
		//		argi += 1;
		//	} else if args[argi] == "--d2j" {
		//		readerOptions.InputFileFormat = "dkvp";
		//		writerOptions.OutputFileFormat = "json";
		//		argi += 1;
		//	} else if args[argi] == "--d2p" {
		//		readerOptions.InputFileFormat = "dkvp";
		//		writerOptions.OutputFileFormat = "pprint";
		//		argi += 1;
		//	} else if args[argi] == "--d2x" {
		//		readerOptions.InputFileFormat = "dkvp";
		//		writerOptions.OutputFileFormat = "xtab";
		//		argi += 1;
		//	} else if args[argi] == "--d2m" {
		//		readerOptions.InputFileFormat = "dkvp";
		//		writerOptions.OutputFileFormat = "markdown";
		//		argi += 1;
		//
		//	} else if args[argi] == "--n2c" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		argi += 1;
		//	} else if args[argi] == "--n2t" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		writerOptions.OFS       = "\t";
		//		argi += 1;
		//	} else if args[argi] == "--n2d" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "dkvp";
		//		argi += 1;
		//	} else if args[argi] == "--n2j" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "json";
		//		argi += 1;
		//	} else if args[argi] == "--n2p" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "pprint";
		//		argi += 1;
		//	} else if args[argi] == "--n2x" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "xtab";
		//		argi += 1;
		//	} else if args[argi] == "--n2m" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "markdown";
		//		argi += 1;
		//
		//	} else if args[argi] == "--j2c" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		argi += 1;
		//	} else if args[argi] == "--j2t" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		writerOptions.OFS       = "\t";
		//		argi += 1;
		//	} else if args[argi] == "--j2d" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "dkvp";
		//		argi += 1;
		//	} else if args[argi] == "--j2n" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "nidx";
		//		argi += 1;
		//	} else if args[argi] == "--j2p" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "pprint";
		//		argi += 1;
		//	} else if args[argi] == "--j2x" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "xtab";
		//		argi += 1;
		//	} else if args[argi] == "--j2m" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "markdown";
		//		argi += 1;
		//
		//	} else if args[argi] == "--p2c" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat        = "csv";
		//		writerOptions.ORS              = "auto";
		//		argi += 1;
		//	} else if args[argi] == "--p2t" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat        = "csv";
		//		writerOptions.ORS              = "auto";
		//		writerOptions.OFS              = "\t";
		//		argi += 1;
		//	} else if args[argi] == "--p2d" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat        = "dkvp";
		//		argi += 1;
		//	} else if args[argi] == "--p2n" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat        = "nidx";
		//		argi += 1;
		//	} else if args[argi] == "--p2j" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat = "json";
		//		argi += 1;
		//	} else if args[argi] == "--p2x" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat        = "xtab";
		//		argi += 1;
		//	} else if args[argi] == "--p2m" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat        = "markdown";
		//		argi += 1;
		//
		//	} else if args[argi] == "--x2c" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		argi += 1;
		//	} else if args[argi] == "--x2t" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		writerOptions.OFS       = "\t";
		//		argi += 1;
		//	} else if args[argi] == "--x2d" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "dkvp";
		//		argi += 1;
		//	} else if args[argi] == "--x2n" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "nidx";
		//		argi += 1;
		//	} else if args[argi] == "--x2j" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "json";
		//		argi += 1;
		//	} else if args[argi] == "--x2p" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "pprint";
		//		argi += 1;
		//	} else if args[argi] == "--x2m" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "markdown";
		//		argi += 1;
		//
		//	} else if args[argi] == "-N" {
		//		readerOptions.use_implicit_csv_header = true;
		//		writerOptions.headerless_csv_output = true;
		//		argi += 1;
	}
	*pargi = argi
	return argi != oargi
}

// Returns true if the current flag was handled.
func parseMiscOptions(
	args []string,
	argc int,
	pargi *int,
	options *clitypes.TOptions,
) bool {

	argi := *pargi
	oargi := argi

	if args[argi] == "-n" {
		options.NoInput = true
		argi += 1

		//	} else if args[argi] == "-I" {
		//		options.do_in_place = true;
		//		argi += 1;
		//
	} else if args[argi] == "--from" {
		checkArgCount(args, argi, argc, 2)
		options.FileNames = append(options.FileNames, args[argi+1])
		argi += 2

		//	} else if args[argi] == "--ofmt" {
		//		checkArgCount(args, argi, argc, 2);
		//		options.ofmt = args[argi+1];
		//		argi += 2;
		//
		//	} else if args[argi] == "--nr-progress-mod" {
		//		checkArgCount(args, argi, argc, 2);
		//		if (sscanf(args[argi+1], "%lld", &options.nr_progress_mod) != 1) {
		//			fmt.Fprintf(os.Stderr,
		//				"%s: --nr-progress-mod argument must be a positive integer; got \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//			mainUsageShort()
		//			os.Exit(1);
		//		}
		//		if (options.nr_progress_mod <= 0) {
		//			fmt.Fprintf(os.Stderr,
		//				"%s: --nr-progress-mod argument must be a positive integer; got \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//			mainUsageShort()
		//			os.Exit(1);
		//		}
		//		argi += 2;
		//
		//	} else if args[argi] == "--seed" {
		//		checkArgCount(args, argi, argc, 2);
		//		if (sscanf(args[argi+1], "0x%x", &options.rand_seed) == 1) {
		//			options.have_rand_seed = true;
		//		} else if (sscanf(args[argi+1], "%u", &options.rand_seed) == 1) {
		//			options.have_rand_seed = true;
		//		} else {
		//			fmt.Fprintf(os.Stderr,
		//				"%s: --seed argument must be a decimal or hexadecimal integer; got \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//			mainUsageShort()
		//			os.Exit(1);
		//		}
		//		argi += 2;

	}
	*pargi = argi
	return argi != oargi
}
