package cst

//// ----------------------------------------------------------------
//typedef void keyword_usage_func_t(FILE* o);
//
//static keyword_usage_func_t mlr_dsl_ENV_keyword_usage;
//static keyword_usage_func_t mlr_dsl_FILENAME_keyword_usage;
//static keyword_usage_func_t mlr_dsl_FILENUM_keyword_usage;
//static keyword_usage_func_t mlr_dsl_FNR_keyword_usage;
//static keyword_usage_func_t mlr_dsl_IFS_keyword_usage;
//static keyword_usage_func_t mlr_dsl_IPS_keyword_usage;
//static keyword_usage_func_t mlr_dsl_IRS_keyword_usage;
//static keyword_usage_func_t mlr_dsl_M_E_keyword_usage;
//static keyword_usage_func_t mlr_dsl_M_PI_keyword_usage;
//static keyword_usage_func_t mlr_dsl_NF_keyword_usage;
//static keyword_usage_func_t mlr_dsl_NR_keyword_usage;
//static keyword_usage_func_t mlr_dsl_OFS_keyword_usage;
//static keyword_usage_func_t mlr_dsl_OPS_keyword_usage;
//static keyword_usage_func_t mlr_dsl_ORS_keyword_usage;
//static keyword_usage_func_t mlr_dsl_all_keyword_usage;
//static keyword_usage_func_t mlr_dsl_begin_keyword_usage;
//static keyword_usage_func_t mlr_dsl_bool_keyword_usage;
//static keyword_usage_func_t mlr_dsl_break_keyword_usage;
//static keyword_usage_func_t mlr_dsl_call_keyword_usage;
//static keyword_usage_func_t mlr_dsl_continue_keyword_usage;
//static keyword_usage_func_t mlr_dsl_do_keyword_usage;
//static keyword_usage_func_t mlr_dsl_dump_keyword_usage;
//static keyword_usage_func_t mlr_dsl_edump_keyword_usage;
//static keyword_usage_func_t mlr_dsl_elif_keyword_usage;
//static keyword_usage_func_t mlr_dsl_else_keyword_usage;
//static keyword_usage_func_t mlr_dsl_emit_keyword_usage;
//static keyword_usage_func_t mlr_dsl_emitf_keyword_usage;
//static keyword_usage_func_t mlr_dsl_emitp_keyword_usage;
//static keyword_usage_func_t mlr_dsl_end_keyword_usage;
//static keyword_usage_func_t mlr_dsl_eprint_keyword_usage;
//static keyword_usage_func_t mlr_dsl_eprintn_keyword_usage;
//static keyword_usage_func_t mlr_dsl_false_keyword_usage;
//static keyword_usage_func_t mlr_dsl_filter_keyword_usage;
//static keyword_usage_func_t mlr_dsl_float_keyword_usage;
//static keyword_usage_func_t mlr_dsl_for_keyword_usage;
//static keyword_usage_func_t mlr_dsl_func_keyword_usage;
//static keyword_usage_func_t mlr_dsl_if_keyword_usage;
//static keyword_usage_func_t mlr_dsl_in_keyword_usage;
//static keyword_usage_func_t mlr_dsl_int_keyword_usage;
//static keyword_usage_func_t mlr_dsl_map_keyword_usage;
//static keyword_usage_func_t mlr_dsl_num_keyword_usage;
//static keyword_usage_func_t mlr_dsl_print_keyword_usage;
//static keyword_usage_func_t mlr_dsl_printn_keyword_usage;
//static keyword_usage_func_t mlr_dsl_return_keyword_usage;
//static keyword_usage_func_t mlr_dsl_stderr_keyword_usage;
//static keyword_usage_func_t mlr_dsl_stdout_keyword_usage;
//static keyword_usage_func_t mlr_dsl_str_keyword_usage;
//static keyword_usage_func_t mlr_dsl_subr_keyword_usage;
//static keyword_usage_func_t mlr_dsl_tee_keyword_usage;
//static keyword_usage_func_t mlr_dsl_true_keyword_usage;
//static keyword_usage_func_t mlr_dsl_unset_keyword_usage;
//static keyword_usage_func_t mlr_dsl_var_keyword_usage;
//static keyword_usage_func_t mlr_dsl_while_keyword_usage;

//// ----------------------------------------------------------------
//typedef struct _keyword_usage_entry_t {
//	char* name;
//	keyword_usage_func_t* pusage_func;
//} keyword_usage_entry_t;

//static keyword_usage_entry_t KEYWORD_USAGE_TABLE[] = {
//
//	{ "all",      mlr_dsl_all_keyword_usage      },
//	{ "begin",    mlr_dsl_begin_keyword_usage    },
//	{ "bool",     mlr_dsl_bool_keyword_usage     },
//	{ "break",    mlr_dsl_break_keyword_usage    },
//	{ "call",     mlr_dsl_call_keyword_usage     },
//	{ "continue", mlr_dsl_continue_keyword_usage },
//	{ "do",       mlr_dsl_do_keyword_usage       },
//	{ "dump",     mlr_dsl_dump_keyword_usage     },
//	{ "edump",    mlr_dsl_edump_keyword_usage    },
//	{ "elif",     mlr_dsl_elif_keyword_usage     },
//	{ "else",     mlr_dsl_else_keyword_usage     },
//	{ "emit",     mlr_dsl_emit_keyword_usage     },
//	{ "emitf",    mlr_dsl_emitf_keyword_usage    },
//	{ "emitp",    mlr_dsl_emitp_keyword_usage    },
//	{ "end",      mlr_dsl_end_keyword_usage      },
//	{ "eprint",   mlr_dsl_eprint_keyword_usage   },
//	{ "eprintn",  mlr_dsl_eprintn_keyword_usage  },
//	{ "false",    mlr_dsl_false_keyword_usage    },
//	{ "filter",   mlr_dsl_filter_keyword_usage   },
//	{ "float",    mlr_dsl_float_keyword_usage    },
//	{ "for",      mlr_dsl_for_keyword_usage      },
//	{ "func",     mlr_dsl_func_keyword_usage     },
//	{ "if",       mlr_dsl_if_keyword_usage       },
//	{ "in",       mlr_dsl_in_keyword_usage       },
//	{ "int",      mlr_dsl_int_keyword_usage      },
//	{ "map",      mlr_dsl_map_keyword_usage      },
//	{ "num",      mlr_dsl_num_keyword_usage      },
//	{ "print",    mlr_dsl_print_keyword_usage    },
//	{ "printn",   mlr_dsl_printn_keyword_usage   },
//	{ "return",   mlr_dsl_return_keyword_usage   },
//	{ "stderr",   mlr_dsl_stderr_keyword_usage   },
//	{ "stdout",   mlr_dsl_stdout_keyword_usage   },
//	{ "str",      mlr_dsl_str_keyword_usage      },
//	{ "subr",     mlr_dsl_subr_keyword_usage     },
//	{ "tee",      mlr_dsl_tee_keyword_usage      },
//	{ "true",     mlr_dsl_true_keyword_usage     },
//	{ "unset",    mlr_dsl_unset_keyword_usage    },
//	{ "var",      mlr_dsl_var_keyword_usage      },
//	{ "while",    mlr_dsl_while_keyword_usage    },
//	{ "ENV",      mlr_dsl_ENV_keyword_usage      },
//	{ "FILENAME", mlr_dsl_FILENAME_keyword_usage },
//	{ "FILENUM",  mlr_dsl_FILENUM_keyword_usage  },
//	{ "FNR",      mlr_dsl_FNR_keyword_usage      },
//	{ "IFS",      mlr_dsl_IFS_keyword_usage      },
//	{ "IPS",      mlr_dsl_IPS_keyword_usage      },
//	{ "IRS",      mlr_dsl_IRS_keyword_usage      },
//	{ "M_E",      mlr_dsl_M_E_keyword_usage      },
//	{ "M_PI",     mlr_dsl_M_PI_keyword_usage     },
//	{ "NF",       mlr_dsl_NF_keyword_usage       },
//	{ "NR",       mlr_dsl_NR_keyword_usage       },
//	{ "OFS",      mlr_dsl_OFS_keyword_usage      },
//	{ "OPS",      mlr_dsl_OPS_keyword_usage      },
//	{ "ORS",      mlr_dsl_ORS_keyword_usage      },
//
//};

// ================================================================
// Pass function_name == NULL to get usage for all keywords.
// Note keywords are defined in parsing/mlr_dsl_lexer.l.
//void mlr_dsl_keyword_usage(FILE* o, char* keyword) {
//	if (keyword == NULL) {
//		for (int i = 0; i < KEYWORD_USAGE_TABLE_SIZE; i++) {
//			if (i > 0) {
//				fmt.Fprintf(o, "\n");
//			}
//			KEYWORD_USAGE_TABLE[i].pusage_func(o);
//		}
//	} else {
//
//		int found = FALSE;
//		for (int i = 0; i < KEYWORD_USAGE_TABLE_SIZE; i++) {
//			if (streq(keyword, KEYWORD_USAGE_TABLE[i].name)) {
//				KEYWORD_USAGE_TABLE[i].pusage_func(o);
//				found = TRUE;
//				break;
//			}
//		}
//		if (!found) {
//			fmt.Fprintf(o, "%s: unrecognized keyword \"%s\".\n", MLR_GLOBALS.bargv0, keyword);
//		}
//	}
//}

//void mlr_dsl_list_all_keywords_raw(FILE* o) {
//	for (int i = 0; i < KEYWORD_USAGE_TABLE_SIZE; i++) {
//		printf("%s\n", KEYWORD_USAGE_TABLE[i].name);
//	}
//}

// ----------------------------------------------------------------
//func mlr_dsl_all_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"all: used in \"emit\", \"emitp\", and \"unset\" as a synonym for @*\n"
//	);
//}
//
//func mlr_dsl_begin_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"begin: defines a block of statements to be executed before input records\n"
//		"are ingested. The body statements must be wrapped in curly braces.\n"
//		"Example: 'begin { @count = 0 }'\n"
//	);
//}
//
//func mlr_dsl_bool_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"bool: declares a boolean local variable in the current curly-braced scope.\n"
//		"Type-checking happens at assignment: 'bool b = 1' is an error.\n"
//	);
//}
//
//func mlr_dsl_break_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"break: causes execution to continue after the body of the current\n"
//		"for/while/do-while loop.\n"
//	);
//}
//
//func mlr_dsl_call_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"call: used for invoking a user-defined subroutine.\n"
//		"Example: 'subr s(k,v) { print k . \" is \" . v} call s(\"a\", $a)'\n"
//	);
//}
//
//func mlr_dsl_continue_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"continue: causes execution to skip the remaining statements in the body of\n"
//		"the current for/while/do-while loop. For-loop increments are still applied.\n"
//	);
//}
//
//func mlr_dsl_do_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"do: with \"while\", introduces a do-while loop. The body statements must be wrapped\n"
//		"in curly braces.\n"
//	);
//}
//
//func mlr_dsl_dump_keyword_usage(FILE* o) {
//    fmt.Fprintf(o,
//		"dump: prints all currently defined out-of-stream variables immediately\n"
//		"  to stdout as JSON.\n"
//		"\n"
//		"  With >, >>, or |, the data do not become part of the output record stream but\n"
//		"  are instead redirected.\n"
//		"\n"
//		"  The > and >> are for write and append, as in the shell, but (as with awk) the\n"
//		"  file-overwrite for > is on first write, not per record. The | is for piping to\n"
//		"  a process which will process the data. There will be one open file for each\n"
//		"  distinct file name (for > and >>) or one subordinate process for each distinct\n"
//		"  value of the piped-to command (for |). Output-formatting flags are taken from\n"
//		"  the main command line.\n"
//		"\n"
//		"  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump }'\n"
//		"  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump >  \"mytap.dat\"}'\n"
//		"  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump >> \"mytap.dat\"}'\n"
//		"  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump | \"jq .[]\"}'\n");
//}
//
//func mlr_dsl_edump_keyword_usage(FILE* o) {
//    fmt.Fprintf(o,
//		"edump: prints all currently defined out-of-stream variables immediately\n"
//		"  to stderr as JSON.\n"
//		"\n"
//		"  Example: mlr --from f.dat put -q '@v[NR]=$*; end { edump }'\n");
//}
//
//func mlr_dsl_elif_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"elif: the way Miller spells \"else if\". The body statements must be wrapped\n"
//		"in curly braces.\n"
//	);
//}
//
//func mlr_dsl_else_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"else: terminates an if/elif/elif chain. The body statements must be wrapped\n"
//		"in curly braces.\n"
//	);
//}
//
//func mlr_dsl_emit_keyword_usage(FILE* o) {
//    fmt.Fprintf(o,
//		"emit: inserts an out-of-stream variable into the output record stream. Hashmap\n"
//		"  indices present in the data but not slotted by emit arguments are not output.\n"
//		"\n"
//		"  With >, >>, or |, the data do not become part of the output record stream but\n"
//		"  are instead redirected.\n"
//		"\n"
//		"  The > and >> are for write and append, as in the shell, but (as with awk) the\n"
//		"  file-overwrite for > is on first write, not per record. The | is for piping to\n"
//		"  a process which will process the data. There will be one open file for each\n"
//		"  distinct file name (for > and >>) or one subordinate process for each distinct\n"
//		"  value of the piped-to command (for |). Output-formatting flags are taken from\n"
//		"  the main command line.\n"
//		"\n"
//		"  You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,\n"
//		"  etc., to control the format of the output if the output is redirected. See also %s -h.\n"
//		"\n"
//		"  Example: mlr --from f.dat put 'emit >  \"/tmp/data-\".$a, $*'\n"
//		"  Example: mlr --from f.dat put 'emit >  \"/tmp/data-\".$a, mapexcept($*, \"a\")'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @sums'\n"
//		"  Example: mlr --from f.dat put --ojson '@sums[$a][$b]+=$x; emit > \"tap-\".$a.$b.\".dat\", @sums'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @sums, \"index1\", \"index2\"'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @*, \"index1\", \"index2\"'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit >  \"mytap.dat\", @*, \"index1\", \"index2\"'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit >> \"mytap.dat\", @*, \"index1\", \"index2\"'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit | \"gzip > mytap.dat.gz\", @*, \"index1\", \"index2\"'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit > stderr, @*, \"index1\", \"index2\"'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit | \"grep somepattern\", @*, \"index1\", \"index2\"'\n"
//		"\n"
//		"  Please see http://johnkerl.org/miller/doc for more information.\n",
//		MLR_GLOBALS.bargv0);
//}
//
//func mlr_dsl_emitf_keyword_usage(FILE* o) {
//    fmt.Fprintf(o,
//		"emitf: inserts non-indexed out-of-stream variable(s) side-by-side into the\n"
//		"  output record stream.\n"
//		"\n"
//		"  With >, >>, or |, the data do not become part of the output record stream but\n"
//		"  are instead redirected.\n"
//		"\n"
//		"  The > and >> are for write and append, as in the shell, but (as with awk) the\n"
//		"  file-overwrite for > is on first write, not per record. The | is for piping to\n"
//		"  a process which will process the data. There will be one open file for each\n"
//		"  distinct file name (for > and >>) or one subordinate process for each distinct\n"
//		"  value of the piped-to command (for |). Output-formatting flags are taken from\n"
//		"  the main command line.\n"
//		"\n"
//		"  You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,\n"
//		"  etc., to control the format of the output if the output is redirected. See also %s -h.\n"
//		"\n"
//		"  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf @a'\n"
//		"  Example: mlr --from f.dat put --oxtab '@a=$i;@b+=$x;@c+=$y; emitf > \"tap-\".$i.\".dat\", @a'\n"
//		"  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf @a, @b, @c'\n"
//		"  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf > \"mytap.dat\", @a, @b, @c'\n"
//		"  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf >> \"mytap.dat\", @a, @b, @c'\n"
//		"  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf > stderr, @a, @b, @c'\n"
//		"  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf | \"grep somepattern\", @a, @b, @c'\n"
//		"  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf | \"grep somepattern > mytap.dat\", @a, @b, @c'\n"
//		"\n"
//		"  Please see http://johnkerl.org/miller/doc for more information.\n",
//		MLR_GLOBALS.bargv0);
//}
//
//func mlr_dsl_emitp_keyword_usage(FILE* o) {
//    fmt.Fprintf(o,
//		"emitp: inserts an out-of-stream variable into the output record stream.\n"
//		"  Hashmap indices present in the data but not slotted by emitp arguments are\n"
//		"  output concatenated with \":\".\n"
//		"\n"
//		"  With >, >>, or |, the data do not become part of the output record stream but\n"
//		"  are instead redirected.\n"
//		"\n"
//		"  The > and >> are for write and append, as in the shell, but (as with awk) the\n"
//		"  file-overwrite for > is on first write, not per record. The | is for piping to\n"
//		"  a process which will process the data. There will be one open file for each\n"
//		"  distinct file name (for > and >>) or one subordinate process for each distinct\n"
//		"  value of the piped-to command (for |). Output-formatting flags are taken from\n"
//		"  the main command line.\n"
//		"\n"
//		"  You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,\n"
//		"  etc., to control the format of the output if the output is redirected. See also %s -h.\n"
//		"\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @sums'\n"
//		"  Example: mlr --from f.dat put --opprint '@sums[$a][$b]+=$x; emitp > \"tap-\".$a.$b.\".dat\", @sums'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @sums, \"index1\", \"index2\"'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @*, \"index1\", \"index2\"'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp >  \"mytap.dat\", @*, \"index1\", \"index2\"'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp >> \"mytap.dat\", @*, \"index1\", \"index2\"'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp | \"gzip > mytap.dat.gz\", @*, \"index1\", \"index2\"'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp > stderr, @*, \"index1\", \"index2\"'\n"
//		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp | \"grep somepattern\", @*, \"index1\", \"index2\"'\n"
//		"\n"
//		"  Please see http://johnkerl.org/miller/doc for more information.\n",
//		MLR_GLOBALS.bargv0);
//}
//
//func mlr_dsl_end_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"end: defines a block of statements to be executed after input records\n"
//		"are ingested. The body statements must be wrapped in curly braces.\n"
//		"Example: 'end { emit @count }'\n"
//		"Example: 'end { eprint \"Final count is \" . @count }'\n"
//	);
//}
//
//func mlr_dsl_eprint_keyword_usage(FILE* o) {
//    fmt.Fprintf(o,
//		"eprint: prints expression immediately to stderr.\n"
//		"  Example: mlr --from f.dat put -q 'eprint \"The sum of x and y is \".($x+$y)'\n"
//		"  Example: mlr --from f.dat put -q 'for (k, v in $*) { eprint k . \" => \" . v }'\n"
//		"  Example: mlr --from f.dat put  '(NR %% 1000 == 0) { eprint \"Checkpoint \".NR}'\n");
//}
//
//func mlr_dsl_eprintn_keyword_usage(FILE* o) {
//    fmt.Fprintf(o,
//		"eprintn: prints expression immediately to stderr, without trailing newline.\n"
//		"  Example: mlr --from f.dat put -q 'eprintn \"The sum of x and y is \".($x+$y); eprint \"\"'\n");
//}
//
//func mlr_dsl_false_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"false: the boolean literal value.\n"
//	);
//}
//
//func mlr_dsl_filter_keyword_usage(FILE* o) {
//    fmt.Fprintf(o,
//		"filter: includes/excludes the record in the output record stream.\n"
//		"\n"
//		"  Example: mlr --from f.dat put 'filter (NR == 2 || $x > 5.4)'\n"
//		"\n"
//		"  Instead of put with 'filter false' you can simply use put -q.  The following\n"
//		"  uses the input record to accumulate data but only prints the running sum\n"
//		"  without printing the input record:\n"
//		"\n"
//		"  Example: mlr --from f.dat put -q '@running_sum += $x * $y; emit @running_sum'\n");
//}
//
//func mlr_dsl_float_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"float: declares a floating-point local variable in the current curly-braced scope.\n"
//		"Type-checking happens at assignment: 'float x = 0' is an error.\n"
//	);
//}
//
//func mlr_dsl_for_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"for: defines a for-loop using one of three styles. The body statements must\n"
//		"be wrapped in curly braces.\n"
//		"For-loop over stream record:\n"
//		"  Example:  'for (k, v in $*) { ... }'\n"
//		"For-loop over out-of-stream variables:\n"
//		"  Example: 'for (k, v in @counts) { ... }'\n"
//		"  Example: 'for ((k1, k2), v in @counts) { ... }'\n"
//		"  Example: 'for ((k1, k2, k3), v in @*) { ... }'\n"
//		"C-style for-loop:\n"
//		"  Example:  'for (var i = 0, var b = 1; i < 10; i += 1, b *= 2) { ... }'\n"
//	);
//}
//
//func mlr_dsl_func_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"func: used for defining a user-defined function.\n"
//		"Example: 'func f(a,b) { return sqrt(a**2+b**2)} $d = f($x, $y)'\n"
//	);
//}
//
//func mlr_dsl_if_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"if: starts an if/elif/elif chain. The body statements must be wrapped\n"
//		"in curly braces.\n"
//	);
//}
//
//func mlr_dsl_in_keyword_usage(FILE* o) {
//	fmt.Fprintf(o, "in: used in for-loops over stream records or out-of-stream variables.\n");
//}
//
//func mlr_dsl_int_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"int: declares an integer local variable in the current curly-braced scope.\n"
//		"Type-checking happens at assignment: 'int x = 0.0' is an error.\n"
//	);
//}
//
//func mlr_dsl_map_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"map: declares an map-valued local variable in the current curly-braced scope.\n"
//		"Type-checking happens at assignment: 'map b = 0' is an error. map b = {} is\n"
//		"always OK. map b = a is OK or not depending on whether a is a map.\n"
//	);
//}
//
//func mlr_dsl_num_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"num: declares an int/float local variable in the current curly-braced scope.\n"
//		"Type-checking happens at assignment: 'num b = true' is an error.\n"
//	);
//}
//
//func mlr_dsl_print_keyword_usage(FILE* o) {
//    fmt.Fprintf(o,
//		"print: prints expression immediately to stdout.\n"
//		"  Example: mlr --from f.dat put -q 'print \"The sum of x and y is \".($x+$y)'\n"
//		"  Example: mlr --from f.dat put -q 'for (k, v in $*) { print k . \" => \" . v }'\n"
//		"  Example: mlr --from f.dat put  '(NR %% 1000 == 0) { print > stderr, \"Checkpoint \".NR}'\n");
//}
//
//func mlr_dsl_printn_keyword_usage(FILE* o) {
//    fmt.Fprintf(o,
//		"printn: prints expression immediately to stdout, without trailing newline.\n"
//		"  Example: mlr --from f.dat put -q 'printn \".\"; end { print \"\" }'\n");
//}
//
//func mlr_dsl_return_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"return: specifies the return value from a user-defined function.\n"
//		"Omitted return statements (including via if-branches) result in an absent-null\n"
//		"return value, which in turns results in a skipped assignment to an LHS.\n"
//	);
//}
//
//func mlr_dsl_stderr_keyword_usage(FILE* o) {
//    fmt.Fprintf(o,
//		"stderr: Used for tee, emit, emitf, emitp, print, and dump in place of filename\n"
//		"  to print to standard error.\n");
//}
//
//func mlr_dsl_stdout_keyword_usage(FILE* o) {
//    fmt.Fprintf(o,
//		"stdout: Used for tee, emit, emitf, emitp, print, and dump in place of filename\n"
//		"  to print to standard output.\n");
//}
//
//func mlr_dsl_str_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"str: declares a string local variable in the current curly-braced scope.\n"
//		"Type-checking happens at assignment.\n"
//	);
//}
//
//func mlr_dsl_subr_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"subr: used for defining a subroutine.\n"
//		"Example: 'subr s(k,v) { print k . \" is \" . v} call s(\"a\", $a)'\n"
//	);
//}
//
//func mlr_dsl_tee_keyword_usage(FILE* o) {
//    fmt.Fprintf(o,
//		"tee: prints the current record to specified file.\n"
//		"  This is an immediate print to the specified file (except for pprint format\n"
//		"  which of course waits until the end of the input stream to format all output).\n"
//		"\n"
//		"  The > and >> are for write and append, as in the shell, but (as with awk) the\n"
//		"  file-overwrite for > is on first write, not per record. The | is for piping to\n"
//		"  a process which will process the data. There will be one open file for each\n"
//		"  distinct file name (for > and >>) or one subordinate process for each distinct\n"
//		"  value of the piped-to command (for |). Output-formatting flags are taken from\n"
//		"  the main command line.\n"
//		"\n"
//		"  You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,\n"
//		"  etc., to control the format of the output. See also %s -h.\n"
//		"\n"
//		"  emit with redirect and tee with redirect are identical, except tee can only\n"
//		"  output $*.\n"
//		"\n"
//		"  Example: mlr --from f.dat put 'tee >  \"/tmp/data-\".$a, $*'\n"
//		"  Example: mlr --from f.dat put 'tee >> \"/tmp/data-\".$a.$b, $*'\n"
//		"  Example: mlr --from f.dat put 'tee >  stderr, $*'\n"
//		"  Example: mlr --from f.dat put -q 'tee | \"tr \[a-z\\] \[A-Z\\]\", $*'\n"
//		"  Example: mlr --from f.dat put -q 'tee | \"tr \[a-z\\] \[A-Z\\] > /tmp/data-\".$a, $*'\n"
//		"  Example: mlr --from f.dat put -q 'tee | \"gzip > /tmp/data-\".$a.\".gz\", $*'\n"
//		"  Example: mlr --from f.dat put -q --ojson 'tee | \"gzip > /tmp/data-\".$a.\".gz\", $*'\n",
//		MLR_GLOBALS.bargv0);
//}
//
//func mlr_dsl_true_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"true: the boolean literal value.\n"
//	);
//}
//
//func mlr_dsl_unset_keyword_usage(FILE* o) {
//    fmt.Fprintf(o,
//		"unset: clears field(s) from the current record, or an out-of-stream or local variable.\n"
//		"\n"
//		"  Example: mlr --from f.dat put 'unset $x'\n"
//		"  Example: mlr --from f.dat put 'unset $*'\n"
//		"  Example: mlr --from f.dat put 'for (k, v in $*) { if (k =~ \"a.*\") { unset $[k] } }'\n"
//		"  Example: mlr --from f.dat put '...; unset @sums'\n"
//		"  Example: mlr --from f.dat put '...; unset @sums[\"green\"]'\n"
//		"  Example: mlr --from f.dat put '...; unset @*'\n");
//}
//
//func mlr_dsl_var_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"var: declares an untyped local variable in the current curly-braced scope.\n"
//		"Examples: 'var a=1', 'var xyz=\"\"'\n"
//	);
//}
//
//func mlr_dsl_while_keyword_usage(FILE* o) {
//	fmt.Fprintf(o,
//		"while: introduces a while loop, or with \"do\", introduces a do-while loop.\n"
//		"The body statements must be wrapped in curly braces.\n"
//	);
//}
//
//func mlr_dsl_ENV_keyword_usage(o *os.File) {
//	fmt.Fprintf(o,
//		"ENV: access to environment variables by name, e.g. '$home = ENV[\"HOME\"]'\n"
//	);
//}
//
//func mlr_dsl_FILENAME_keyword_usage(o *os.File) {
//	fmt.Fprintf(o,
//		"FILENAME: evaluates to the name of the current file being processed.\n"
//	);
//}
//
//func mlr_dsl_FILENUM_keyword_usage(o *os.File) {
//	fmt.Fprintf(o,
//		"FILENUM: evaluates to the number of the current file being processed,\n"
//		"starting with 1.\n"
//	);
//}
//
//func mlr_dsl_FNR_keyword_usage(o *os.File) {
//	fmt.Fprintf(o,
//		"FNR: evaluates to the number of the current record within the current file\n"
//		"being processed, starting with 1. Resets at the start of each file.\n"
//	);
//}
//
//func mlr_dsl_IFS_keyword_usage(o *os.File) {
//	fmt.Fprintf(o,
//		"IFS: evaluates to the input field separator from the command line.\n"
//	);
//}
//
//func mlr_dsl_IPS_keyword_usage(o *os.File) {
//	fmt.Fprintf(o,
//		"IPS: evaluates to the input pair separator from the command line.\n"
//	);
//}
//
//func mlr_dsl_IRS_keyword_usage(o *os.File) {
//	fmt.Fprintf(o,
//		"IRS: evaluates to the input record separator from the command line,\n"
//		"or to LF or CRLF from the input data if in autodetect mode (which is\n"
//		"the default).\n"
//	);
//}
//
//func mlr_dsl_M_E_keyword_usage(o *os.File) {
//	fmt.Fprintf(o,
//		"M_E: the mathematical constant e.\n"
//	);
//}
//
//func mlr_dsl_M_PI_keyword_usage(o *os.File) {
//	fmt.Fprintf(o,
//		"M_PI: the mathematical constant pi.\n"
//	);
//}
//
//func mlr_dsl_NF_keyword_usage(o *os.File) {
//	fmt.Fprintf(o,
//		"NF: evaluates to the number of fields in the current record.\n"
//	);
//}
//
//func mlr_dsl_NR_keyword_usage(o *os.File) {
//	fmt.Fprintf(o,
//		"NR: evaluates to the number of the current record over all files\n"
//		"being processed, starting with 1. Does not reset at the start of each file.\n"
//	);
//}
//
//func mlr_dsl_OFS_keyword_usage(o *os.File) {
//	fmt.Fprintf(o,
//		"OFS: evaluates to the output field separator from the command line.\n"
//	);
//}
//
//func mlr_dsl_OPS_keyword_usage(o *os.File) {
//	fmt.Fprintf(o,
//		"OPS: evaluates to the output pair separator from the command line.\n"
//	);
//}
//
//func mlr_dsl_ORS_keyword_usage(o *os.File) {
//	fmt.Fprintf(o,
//		"ORS: evaluates to the output record separator from the command line,\n"
//		"or to LF or CRLF from the input data if in autodetect mode (which is\n"
//		"the default).\n"
//	);
//}
//
