package cli

// ================================================================
// Stuff needing to be ported from C
// ================================================================

//// ----------------------------------------------------------------
//#define DEFAULT_OFMT                     "%lf"
//#define DEFAULT_OQUOTING                 QUOTE_MINIMAL
//#define DEFAULT_OOSVAR_FLATTEN_SEPARATOR ":"
//#define DEFAULT_COMMENT_STRING           "#"
//
//// ASCII 1f and 1e
//#define ASV_FS "\x1f"
//#define ASV_RS "\x1e"
//
//#define ASV_FS_FOR_HELP "0x1f"
//#define ASV_RS_FOR_HELP "0x1e"
//
//// Unicode code points U+241F and U+241E, encoded as UTF-8.
//#define USV_FS "\xe2\x90\x9f"
//#define USV_RS "\xe2\x90\x9e"
//
//#define USV_FS_FOR_HELP "U+241F (UTF-8 0xe2909f)"
//#define USV_RS_FOR_HELP "U+241E (UTF-8 0xe2909e)"

// ----------------------------------------------------------------
//static lhmss_t* singleton_pdesc_to_chars_map = nil;
//static lhmss_t* get_desc_to_chars_map() {
//	if (singleton_pdesc_to_chars_map == nil) {
//		singleton_pdesc_to_chars_map = lhmss_alloc();
//		lhmss_put(singleton_pdesc_to_chars_map, "cr",        "\r",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "crcr",      "\r\r",     NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "newline",   "\n",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "lf",        "\n",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "lflf",      "\n\n",     NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "crlf",      "\r\n",     NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "crlfcrlf",  "\r\n\r\n", NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "tab",       "\t",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "space",     " ",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "comma",     ",",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "newline",   "\n",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "pipe",      "|",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "slash",     "/",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "colon",     ":",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "semicolon", ";",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "equals",    "=",        NO_FREE);
//	}
//	return singleton_pdesc_to_chars_map;
//}

//// ----------------------------------------------------------------
//static lhmss_t* singleton_default_rses = nil;
//static lhmss_t* singleton_default_fses = nil;
//static lhmss_t* singleton_default_pses = nil;
//static lhmsll_t* singleton_default_repeat_ifses = nil;
//static lhmsll_t* singleton_default_repeat_ipses = nil;
//
//static lhmss_t* get_default_rses() {
//	if (singleton_default_rses == nil) {
//		singleton_default_rses = lhmss_alloc();
//
//		lhmss_put(singleton_default_rses, "gen",      "N/A",  NO_FREE);
//		lhmss_put(singleton_default_rses, "dkvp",     "auto",  NO_FREE);
//		lhmss_put(singleton_default_rses, "json",     "auto",  NO_FREE);
//		lhmss_put(singleton_default_rses, "nidx",     "auto",  NO_FREE);
//		lhmss_put(singleton_default_rses, "csv",      "auto",  NO_FREE);
//		lhmss_put(singleton_default_rses, "csvlite",  "auto",  NO_FREE);
//		lhmss_put(singleton_default_rses, "markdown", "auto",  NO_FREE);
//		lhmss_put(singleton_default_rses, "pprint",   "auto",  NO_FREE);
//		lhmss_put(singleton_default_rses, "xtab",     "(N/A)", NO_FREE);
//	}
//	return singleton_default_rses;
//}
//
//static lhmss_t* get_default_fses() {
//	if (singleton_default_fses == nil) {
//		singleton_default_fses = lhmss_alloc();
//		lhmss_put(singleton_default_fses, "gen",      "(N/A)",  NO_FREE);
//		lhmss_put(singleton_default_fses, "dkvp",     ",",      NO_FREE);
//		lhmss_put(singleton_default_fses, "json",     "(N/A)",  NO_FREE);
//		lhmss_put(singleton_default_fses, "nidx",     " ",      NO_FREE);
//		lhmss_put(singleton_default_fses, "csv",      ",",      NO_FREE);
//		lhmss_put(singleton_default_fses, "csvlite",  ",",      NO_FREE);
//		lhmss_put(singleton_default_fses, "markdown", "(N/A)",  NO_FREE);
//		lhmss_put(singleton_default_fses, "pprint",   " ",      NO_FREE);
//		lhmss_put(singleton_default_fses, "xtab",     "auto",   NO_FREE);
//	}
//	return singleton_default_fses;
//}
//
//static lhmss_t* get_default_pses() {
//	if (singleton_default_pses == nil) {
//		singleton_default_pses = lhmss_alloc();
//		lhmss_put(singleton_default_pses, "gen",      "(N/A)", NO_FREE);
//		lhmss_put(singleton_default_pses, "dkvp",     "=",     NO_FREE);
//		lhmss_put(singleton_default_pses, "json",     "(N/A)", NO_FREE);
//		lhmss_put(singleton_default_pses, "nidx",     "(N/A)", NO_FREE);
//		lhmss_put(singleton_default_pses, "csv",      "(N/A)", NO_FREE);
//		lhmss_put(singleton_default_pses, "csvlite",  "(N/A)", NO_FREE);
//		lhmss_put(singleton_default_pses, "markdown", "(N/A)", NO_FREE);
//		lhmss_put(singleton_default_pses, "pprint",   "(N/A)", NO_FREE);
//		lhmss_put(singleton_default_pses, "xtab",     " ",     NO_FREE);
//	}
//	return singleton_default_pses;
//}
//
//static lhmsll_t* get_default_repeat_ifses() {
//	if (singleton_default_repeat_ifses == nil) {
//		singleton_default_repeat_ifses = lhmsll_alloc();
//		lhmsll_put(singleton_default_repeat_ifses, "gen",      false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "dkvp",     false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "json",     false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "csv",      false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "csvlite",  false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "markdown", false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "nidx",     false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "xtab",     false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "pprint",   true,  NO_FREE);
//	}
//	return singleton_default_repeat_ifses;
//}
//
//static lhmsll_t* get_default_repeat_ipses() {
//	if (singleton_default_repeat_ipses == nil) {
//		singleton_default_repeat_ipses = lhmsll_alloc();
//		lhmsll_put(singleton_default_repeat_ipses, "gen",      false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "dkvp",     false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "json",     false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "csv",      false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "csvlite",  false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "markdown", false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "nidx",     false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "xtab",     true,  NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "pprint",   false, NO_FREE);
//	}
//	return singleton_default_repeat_ipses;
//}
//
//func free_opt_singletons() {
//	lhmss_free(singleton_pdesc_to_chars_map);
//	lhmss_free(singleton_default_rses);
//	lhmss_free(singleton_default_fses);
//	lhmss_free(singleton_default_pses);
//	lhmsll_free(singleton_default_repeat_ifses);
//	lhmsll_free(singleton_default_repeat_ipses);
//}
//
//// For displaying the default separators in on-line help
//static char* rebackslash(char* sep) {
//	if sep == "\r"))
//		return "\\r";
//	else if sep == "\n"))
//		return "\\n";
//	else if sep == "\r\n"))
//		return "\\r\\n";
//	else if sep == "\t"))
//		return "\\t";
//	else if sep == " "))
//		return "space";
//	else
//		return sep;
//}

//func printTypeArithmeticInfo(o *os.File, argv0 string) {
//	for (int i = -2; i < MT_DIM; i++) {
//		mv_t a = (mv_t) {.type = i, .free_flags = NO_FREE, .u.intv = 0};
//		if (i == -2)
//			fmt.Printf("%-6s |", "(+)");
//		else if (i == -1)
//			fmt.Printf("%-6s +", "------");
//		else
//			fmt.Printf("%-6s |", mt_describe_type_simple(a.type));
//
//		for (int j = 0; j < MT_DIM; j++) {
//			mv_t b = (mv_t) {.type = j, .free_flags = NO_FREE, .u.intv = 0};
//			if (i == -2) {
//				fmt.Printf(" %-6s", mt_describe_type_simple(b.type));
//			} else if (i == -1) {
//				fmt.Printf(" %-6s", "------");
//			} else {
//				mv_t c = x_xx_plus_func(&a, &b);
//				fmt.Printf(" %-6s", mt_describe_type_simple(c.type));
//			}
//		}
//
//		fmt.Fprintf(o, "\n");
//	}
//}

//	cli_reader_opts_init(&options.ReaderOptions);
//	cli_writer_opts_init(&options.WriterOptions);
//
//	options.transformer_argb     = 0;
//	options.filenames       = slls_alloc();
//
//	options.ofmt            = nil;
//	options.nr_progress_mod = 0LL;
//
//	options.DoInPlace     = false;
//
//	options.no_input        = false;
//	options.have_rand_seed  = false;
//	options.rand_seed       = 0;
//}

// ----------------------------------------------------------------
//void cli_reader_opts_init(cliutil.TReaderOptions* readerOptions) {
//	readerOptions.InputFileFormat                      = nil;
//	readerOptions.IRS                            = nil;
//	readerOptions.IFS                            = nil;
//	readerOptions.ips                            = nil;
//	readerOptions.input_json_flatten_separator   = nil;
//	readerOptions.json_array_ingest              = JSON_ARRAY_INGEST_UNSPECIFIED;
//
//	readerOptions.allow_repeat_ifs               = NEITHER_TRUE_NOR_FALSE;
//	readerOptions.allow_repeat_ips               = NEITHER_TRUE_NOR_FALSE;
//	readerOptions.use_implicit_csv_header        = NEITHER_TRUE_NOR_FALSE;
//	readerOptions.allow_ragged_csv_input         = NEITHER_TRUE_NOR_FALSE;
//
//	readerOptions.prepipe                        = nil;
//	readerOptions.comment_handling               = COMMENTS_ARE_DATA;
//	readerOptions.comment_string                 = nil;
//
//	readerOptions.generator_opts.field_name     = "i";
//	readerOptions.generator_opts.start          = 0LL;
//	readerOptions.generator_opts.stop           = 100LL;
//	readerOptions.generator_opts.step           = 1LL;
//}

//void cli_writer_opts_init(cliutil.TWriterOptions* writerOptions) {
//	writerOptions.OutputFileFormat                      = nil;
//	writerOptions.ORS                            = nil;
//	writerOptions.OFS                            = nil;
//	writerOptions.ops                            = nil;
//
//	writerOptions.headerless_csv_output          = NEITHER_TRUE_NOR_FALSE;
//	writerOptions.right_justify_xtab_value       = NEITHER_TRUE_NOR_FALSE;
//	writerOptions.right_align_pprint             = NEITHER_TRUE_NOR_FALSE;
//	writerOptions.pprint_barred                  = NEITHER_TRUE_NOR_FALSE;
//	writerOptions.stack_json_output_vertically   = NEITHER_TRUE_NOR_FALSE;
//	writerOptions.wrap_json_output_in_outer_list = NEITHER_TRUE_NOR_FALSE;
//	writerOptions.json_quote_int_keys            = NEITHER_TRUE_NOR_FALSE;
//	writerOptions.json_quote_non_string_values   = NEITHER_TRUE_NOR_FALSE;
//
//	writerOptions.output_json_flatten_separator  = nil;
//	writerOptions.oosvar_flatten_separator       = nil;
//
//	writerOptions.oquoting                       = QUOTE_UNSPECIFIED;
//}

//void cli_apply_defaults(cli_opts_t* popts) {
//
//	cli_apply_reader_defaults(&options.ReaderOptions);
//
//	cli_apply_writer_defaults(&options.WriterOptions);
//
//	if (options.ofmt == nil)
//		options.ofmt = DEFAULT_OFMT;
//}

//void cli_apply_reader_defaults(cliutil.TReaderOptions* readerOptions) {
//	if (readerOptions.InputFileFormat == nil)
//		readerOptions.InputFileFormat = "dkvp";
//
//	if (readerOptions.json_array_ingest == JSON_ARRAY_INGEST_UNSPECIFIED)
//		readerOptions.json_array_ingest = JSON_ARRAY_INGEST_AS_MAP;
//
//	if (readerOptions.use_implicit_csv_header == NEITHER_TRUE_NOR_FALSE)
//		readerOptions.use_implicit_csv_header = false;
//
//	if (readerOptions.allow_ragged_csv_input == NEITHER_TRUE_NOR_FALSE)
//		readerOptions.allow_ragged_csv_input = false;
//
//	if (readerOptions.input_json_flatten_separator == nil)
//		readerOptions.input_json_flatten_separator = DEFAULT_JSON_FLATTEN_SEPARATOR;
//}

//void cli_apply_writer_defaults(cliutil.TWriterOptions* writerOptions) {
//	if (writerOptions.OutputFileFormat == nil)
//		writerOptions.OutputFileFormat = "dkvp";
//
//	if (writerOptions.headerless_csv_output == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.headerless_csv_output = false;
//
//	if (writerOptions.right_justify_xtab_value == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.right_justify_xtab_value = false;
//
//	if (writerOptions.right_align_pprint == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.right_align_pprint = false;
//
//	if (writerOptions.pprint_barred == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.pprint_barred = false;
//
//	if (writerOptions.stack_json_output_vertically == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.stack_json_output_vertically = false;
//
//	if (writerOptions.wrap_json_output_in_outer_list == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.wrap_json_output_in_outer_list = false;
//
//	if (writerOptions.json_quote_int_keys == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.json_quote_int_keys = true;
//
//	if (writerOptions.json_quote_non_string_values == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.json_quote_non_string_values = false;
//
//	if (writerOptions.output_json_flatten_separator == nil)
//		writerOptions.output_json_flatten_separator = DEFAULT_JSON_FLATTEN_SEPARATOR;
//
//	if (writerOptions.oosvar_flatten_separator == nil)
//		writerOptions.oosvar_flatten_separator = DEFAULT_OOSVAR_FLATTEN_SEPARATOR;
//
//	if (writerOptions.oquoting == QUOTE_UNSPECIFIED)
//		writerOptions.oquoting = DEFAULT_OQUOTING;
//}

// ----------------------------------------------------------------
// For transformer join which has its own input-format overrides.
//
// Mainly this just takes the main-opts flag whenever the join-opts flag was not
// specified by the user. But it's a bit more complex when main and join input
// formats are different. Example: main input format is CSV, for which IPS is
// "(N/A)", and join input format is DKVP. Then we should not use "(N/A)"
// for DKVP IPS. However if main input format were DKVP with IPS set to ":",
// then we should take that.
//
// The logic is:
//
// * If the join input format was unspecified, take all unspecified values from
//   main opts.
//
// * If the join input format was specified and is the same as main input
//   format, take unspecified values from main opts.
//
// * If the join input format was specified and is not the same as main input
//   format, take unspecified values from defaults for the join input format.

//void cli_merge_reader_opts(cliutil.TReaderOptions* pfunc_opts, TReaderOptions* pmain_opts) {
//
//	if (pfunc_opts.InputFileFormat == nil) {
//		pfunc_opts.InputFileFormat = pmain_opts.InputFileFormat;
//	}
//
//	if pfunc_opts.InputFileFormat == pmain_opts.InputFileFormat {
//
//		if (pfunc_opts.IRS == nil)
//			pfunc_opts.IRS = pmain_opts.IRS;
//		if (pfunc_opts.IFS == nil)
//			pfunc_opts.IFS = pmain_opts.IFS;
//		if (pfunc_opts.ips == nil)
//			pfunc_opts.ips = pmain_opts.ips;
//		if (pfunc_opts.allow_repeat_ifs  == NEITHER_TRUE_NOR_FALSE)
//			pfunc_opts.allow_repeat_ifs = pmain_opts.allow_repeat_ifs;
//		if (pfunc_opts.allow_repeat_ips  == NEITHER_TRUE_NOR_FALSE)
//			pfunc_opts.allow_repeat_ips = pmain_opts.allow_repeat_ips;
//
//	} else {
//
//		if (pfunc_opts.IRS == nil)
//			pfunc_opts.IRS = lhmss_get_or_die(get_default_rses(), pfunc_opts.InputFileFormat);
//		if (pfunc_opts.IFS == nil)
//			pfunc_opts.IFS = lhmss_get_or_die(get_default_fses(), pfunc_opts.InputFileFormat);
//		if (pfunc_opts.ips == nil)
//			pfunc_opts.ips = lhmss_get_or_die(get_default_pses(), pfunc_opts.InputFileFormat);
//		if (pfunc_opts.allow_repeat_ifs  == NEITHER_TRUE_NOR_FALSE)
//			pfunc_opts.allow_repeat_ifs = lhmsll_get_or_die(get_default_repeat_ifses(), pfunc_opts.InputFileFormat);
//		if (pfunc_opts.allow_repeat_ips  == NEITHER_TRUE_NOR_FALSE)
//			pfunc_opts.allow_repeat_ips = lhmsll_get_or_die(get_default_repeat_ipses(), pfunc_opts.InputFileFormat);
//
//	}
//
//	if (pfunc_opts.json_array_ingest == JSON_ARRAY_INGEST_UNSPECIFIED)
//		pfunc_opts.json_array_ingest = pmain_opts.json_array_ingest;
//
//	if (pfunc_opts.use_implicit_csv_header == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.use_implicit_csv_header = pmain_opts.use_implicit_csv_header;
//
//	if (pfunc_opts.allow_ragged_csv_input == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.allow_ragged_csv_input = pmain_opts.allow_ragged_csv_input;
//
//	if (pfunc_opts.input_json_flatten_separator == nil)
//		pfunc_opts.input_json_flatten_separator = pmain_opts.input_json_flatten_separator;
//}

// Similar to cli_merge_reader_opts but for transformer tee & transformer put which have their
// own output-format overrides.

//void cli_merge_writer_opts(cliutil.TWriterOptions* pfunc_opts, TWriterOptions* pmain_opts) {
//
//	if (pfunc_opts.OutputFileFormat == nil) {
//		pfunc_opts.OutputFileFormat = pmain_opts.OutputFileFormat;
//	}
//
//	if pfunc_opts.OutputFileFormat == pmain_opts.OutputFileFormat {
//		if (pfunc_opts.ORS == nil)
//			pfunc_opts.ORS = pmain_opts.ORS;
//		if (pfunc_opts.OFS == nil)
//			pfunc_opts.OFS = pmain_opts.OFS;
//		if (pfunc_opts.ops == nil)
//			pfunc_opts.ops = pmain_opts.ops;
//	} else {
//		if (pfunc_opts.ORS == nil)
//			pfunc_opts.ORS = lhmss_get_or_die(get_default_rses(), pfunc_opts.OutputFileFormat);
//		if (pfunc_opts.OFS == nil)
//			pfunc_opts.OFS = lhmss_get_or_die(get_default_fses(), pfunc_opts.OutputFileFormat);
//		if (pfunc_opts.ops == nil)
//			pfunc_opts.ops = lhmss_get_or_die(get_default_pses(), pfunc_opts.OutputFileFormat);
//	}
//
//	if (pfunc_opts.headerless_csv_output == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.headerless_csv_output = pmain_opts.headerless_csv_output;
//
//	if (pfunc_opts.right_justify_xtab_value == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.right_justify_xtab_value = pmain_opts.right_justify_xtab_value;
//
//	if (pfunc_opts.right_align_pprint == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.right_align_pprint = pmain_opts.right_align_pprint;
//
//	if (pfunc_opts.pprint_barred == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.pprint_barred = pmain_opts.pprint_barred;
//
//	if (pfunc_opts.stack_json_output_vertically == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.stack_json_output_vertically = pmain_opts.stack_json_output_vertically;
//
//	if (pfunc_opts.wrap_json_output_in_outer_list == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.wrap_json_output_in_outer_list = pmain_opts.wrap_json_output_in_outer_list;
//
//	if (pfunc_opts.json_quote_int_keys == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.json_quote_int_keys = pmain_opts.json_quote_int_keys;
//
//	if (pfunc_opts.json_quote_non_string_values == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.json_quote_non_string_values = pmain_opts.json_quote_non_string_values;
//
//	if (pfunc_opts.output_json_flatten_separator == nil)
//		pfunc_opts.output_json_flatten_separator = pmain_opts.output_json_flatten_separator;
//
//	if (pfunc_opts.oosvar_flatten_separator == nil)
//		pfunc_opts.oosvar_flatten_separator = pmain_opts.oosvar_flatten_separator;
//
//	if (pfunc_opts.oquoting == QUOTE_UNSPECIFIED)
//		pfunc_opts.oquoting = pmain_opts.oquoting;
//}

// ----------------------------------------------------------------
//static char* lhmss_get_or_die(lhmss_t* pmap, char* key) {
//	char* value = lhmss_get(pmap, key);
//	MLR_INTERNAL_CODING_ERROR_IF(value == nil);
//	return value;
//}

// ----------------------------------------------------------------
//static int lhmsll_get_or_die(lhmsll_t* pmap, char* key) {
//	MLR_INTERNAL_CODING_ERROR_UNLESS(lhmsll_has_key(pmap, key));
//	return lhmsll_get(pmap, key);
//}
