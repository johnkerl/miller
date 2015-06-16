#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "lib/mlrutil.h"
#include "lib/mtrand.h"
#include "containers/slls.h"
#include "input/lrec_readers.h"
#include "mapping/mappers.h"
#include "mapping/lrec_evaluators.h"
#include "output/lrec_writers.h"
#include "cli/mlrcli.h"
#include "cli/argparse.h"

// ----------------------------------------------------------------
static mapper_setup_t* mapper_lookup_table[] = {
	&mapper_cat_setup,
	&mapper_check_setup,
	&mapper_count_distinct_setup,
	&mapper_cut_setup,
	&mapper_filter_setup,
	&mapper_group_by_setup,
	&mapper_group_like_setup,
	&mapper_having_fields_setup,
	&mapper_head_setup,
	&mapper_histogram_setup,
	&mapper_join_setup,
	&mapper_label_setup,
	&mapper_put_setup,
	&mapper_rename_setup,
	&mapper_reorder_setup,
	&mapper_sort_setup,
	&mapper_stats1_setup,
	&mapper_stats2_setup,
	&mapper_step_setup,
	&mapper_tac_setup,
	&mapper_tail_setup,
	&mapper_top_setup,
	&mapper_uniq_setup,
};
static int mapper_lookup_table_length = sizeof(mapper_lookup_table) / sizeof(mapper_lookup_table[0]);

// ----------------------------------------------------------------
#define DEFAULT_RS   '\n'
#define DEFAULT_FS   ','
#define DEFAULT_PS   '='
#define DEFAULT_OFMT "%lf"

// ----------------------------------------------------------------
// xxx cmt stdout/err & 0/1
static void main_usage(char* argv0, int exit_code) {
	FILE* o = exit_code == 0 ? stdout : stderr;
	fprintf(o, "Usage: %s [I/O options] {verb} [verb-dependent options ...] {file names}\n", argv0);
	fprintf(o, "Verbs:\n");
	int linelen = 0;
	for (int i = 0; i < mapper_lookup_table_length; i++) {
		linelen += 1 + strlen(mapper_lookup_table[i]->verb);
		if (linelen > 80) {
			fprintf(o, "\n");
			linelen = 0;
		}
		if ((i > 0) && (linelen > 0))
			fprintf(o, " ");
		else
			fprintf(o, "   ");
		fprintf(o, "%s", mapper_lookup_table[i]->verb);
	}
	fprintf(o, "\n");
	fprintf(o, "Please use \"%s {verb name} --help\" for verb-specific help.\n", argv0);
	fprintf(o, "Please use \"%s --help-all-verbs\" for help on all verbs.\n", argv0);

	fprintf(o, "\n");
	lrec_evaluator_list_functions(o);
	fprintf(o, "Please use \"%s --help-function {function name}\" for function-specific help.\n", argv0);
	fprintf(o, "Please use \"%s --help-all-functions\" or \"%s -f\" for help on all functions.\n", argv0, argv0);
	fprintf(o, "\n");

	fprintf(o, "Separator options, for input, output, or both:\n");
	fprintf(o, "  --rs      --irs     --ors                  Record separators, defaulting to newline\n");
	fprintf(o, "  --fs      --ifs     --ofs    --repifs      Field  separators, defaulting to \"%c\"\n", DEFAULT_FS);
	fprintf(o, "  --ps      --ips     --ops                  Pair   separators, defaulting to \"%c\"\n", DEFAULT_PS);
	fprintf(o, "Data-format options, for input, output, or both:\n");
	fprintf(o, "  --dkvp    --idkvp   --odkvp                Delimited key-value pairs, e.g \"a=1,b=2\" (default)\n");
	fprintf(o, "  --nidx    --inidx   --onidx                Implicitly-integer-indexed fields (Unix-toolkit style)\n");
	fprintf(o, "  --csv     --icsv    --ocsv                 Comma-separated value (or tab-separated with --fs tab, etc.)\n");
	fprintf(o, "  --pprint  --ipprint --opprint --right      Pretty-printed tabular (produces no output until all input is in)\n");
	fprintf(o, "  --xtab    --ixtab   --oxtab                Transposed-tabular (useful for highly multi-column data)\n");
	fprintf(o, "Numerical format:\n");
	fprintf(o, "  --ofmt {format}                            E.g. %%.18lf, %%.0lf. Please use sprintf-style codes for double-precision.\n");
	fprintf(o, "                                             Applies to verbs which compute new values, e.g. put, stats1, stats2.\n");
	fprintf(o, "Other options:\n");
	fprintf(o, "  --seed {n} with n of the form 12345678 or 0xcafefeed. For put/filter urand().\n");
	fprintf(o, "Output of one verb may be chained as input to another using \"then\", e.g.\n");
	fprintf(o, "  %s stats1 -a min,mean,max -f flag,u,v -g color then sort -f color\n", argv0);
	fprintf(o, "Please see http://johnkerl.org/miller/doc and/or http://github.com/johnkerl/miller for more information.\n");

	exit(exit_code);
}

static void usage_all_verbs(char* argv0) {
	char* separator = "================================================================";

	for (int i = 0; i < mapper_lookup_table_length; i++) {
		fprintf(stdout, "%s\n", separator);
		mapper_lookup_table[i]->pusage_func(argv0, mapper_lookup_table[i]->verb);
		fprintf(stdout, "\n");
	}
	fprintf(stdout, "%s\n", separator);
	exit(0);
}

static void nusage(char* argv0, char* arg) {
	fprintf(stdout, "%s: option \"%s\" not recognized.\n", argv0, arg);
	fprintf(stdout, "\n");
	main_usage(argv0, 1);
}

static void check_arg_count(char** argv, int argi, int argc, int n) {
	if ((argc - argi) < n) {
		main_usage(argv[0], 1);
	}
}

static char sep_from_arg(char* arg, char* argv0) {
	if (streq(arg, "tab"))
		return '\t';
	if (streq(arg, "space"))
		return ' ';
	if (streq(arg, "newline"))
		return '\n';
	if (streq(arg, "pipe"))
		return '|';
	if (streq(arg, "semicolon"))
		return '|';
	if (strlen(arg) != 1)
		main_usage(argv0, 1);
	return arg[0];
}

static mapper_setup_t* look_up_mapper_setup(char* verb) {
	mapper_setup_t* pmapper_setup = NULL;
	for (int i = 0; i < mapper_lookup_table_length; i++) {
		if (streq(mapper_lookup_table[i]->verb, verb))
			return mapper_lookup_table[i];
	}

	return pmapper_setup;
}

// ----------------------------------------------------------------
cli_opts_t* parse_command_line(int argc, char** argv) {
	cli_opts_t* popts = mlr_malloc_or_die(sizeof(cli_opts_t));
	memset(popts, 0, sizeof(*popts));

	popts->irs  = DEFAULT_RS;
	popts->ifs  = DEFAULT_FS;
	popts->ips  = DEFAULT_PS;
	popts->allow_repeat_ifs = FALSE;
	popts->allow_repeat_ips = FALSE;

	popts->ors  = DEFAULT_RS;
	popts->ofs  = DEFAULT_FS;
	popts->ops  = DEFAULT_PS;
	popts->ofmt = DEFAULT_OFMT;

	popts->plrec_reader_stdio   = NULL;
	popts->plrec_reader_mmap    = NULL;
	popts->plrec_writer         = NULL;
	popts->filenames            = NULL;
	popts->use_file_reader_mmap = TRUE;

	char* rdesc = "dkvp";
	char* wdesc = "dkvp";
	int left_align_pprint = TRUE;

	int have_rand_seed = FALSE;
	unsigned rand_seed = 0;

	int argi = 1;
	for (; argi < argc; argi++) {
		if (argv[argi][0] != '-')
			break;

		else if (streq(argv[argi], "-h"))
			main_usage(argv[0], 0);
		else if (streq(argv[argi], "--help"))
			main_usage(argv[0], 0);
		else if (streq(argv[argi], "--help-all-verbs"))
			usage_all_verbs(argv[0]);
		else if (streq(argv[argi], "--help-all-functions") || streq(argv[argi], "-f")) {
			lrec_evaluator_function_usage(stdout, NULL);
			exit(0);
		}

		else if (streq(argv[argi], "--help-function") || streq(argv[argi], "--hf")) {
			check_arg_count(argv, argi, argc, 2);
			lrec_evaluator_function_usage(stdout, argv[argi+1]);
			exit(0);
		}

		else if (streq(argv[argi], "--rs")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ors = popts->irs = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		}
		else if (streq(argv[argi], "--irs")) {
			check_arg_count(argv, argi, argc, 2);
			popts->irs = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		}
		else if (streq(argv[argi], "--ors")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ors = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		}

		else if (streq(argv[argi], "--fs")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ofs = popts->ifs = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		}
		else if (streq(argv[argi], "--ifs")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ifs = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		}
		else if (streq(argv[argi], "--ofs")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ofs = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		}
		else if (streq(argv[argi], "--repifs")) {
			popts->allow_repeat_ifs = TRUE;
		}

		else if (streq(argv[argi], "--ps")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ops = popts->ips = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		}
		else if (streq(argv[argi], "--ips")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ips = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		}
		else if (streq(argv[argi], "--ops")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ops = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		}

		else if (streq(argv[argi], "--dkvp"))    { rdesc = wdesc = "dkvp"; }
		else if (streq(argv[argi], "--idkvp"))   { rdesc = "dkvp"; }
		else if (streq(argv[argi], "--odkvp"))   { wdesc = "dkvp"; }

		else if (streq(argv[argi], "--csv"))     { rdesc = wdesc = "csv";  }
		else if (streq(argv[argi], "--icsv"))    { rdesc = "csv";  }
		else if (streq(argv[argi], "--ocsv"))    { wdesc = "csv";  }

		else if (streq(argv[argi], "--nidx"))    { rdesc = wdesc = "nidx"; }
		else if (streq(argv[argi], "--inidx"))   { rdesc = "nidx"; }
		else if (streq(argv[argi], "--onidx"))   { wdesc = "nidx"; }

		else if (streq(argv[argi], "--xtab"))    { rdesc = wdesc = "xtab"; }
		else if (streq(argv[argi], "--ixtab"))   { rdesc = "xtab"; }
		else if (streq(argv[argi], "--oxtab"))   { wdesc = "xtab"; }

		else if (streq(argv[argi], "--ipprint")) {
			rdesc                   = "csv";
			popts->ifs              = ' ';
			popts->allow_repeat_ifs = TRUE;

		}
		else if (streq(argv[argi], "--opprint")) {
			wdesc = "pprint";
		}
		else if (streq(argv[argi], "--pprint")) {
			rdesc                   = "csv";
			popts->ifs              = ' ';
			popts->allow_repeat_ifs = TRUE;
			wdesc                   = "pprint";
		}
		else if (streq(argv[argi], "--right"))   {
			left_align_pprint = FALSE;
		}

		else if (streq(argv[argi], "--ofmt")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ofmt = argv[argi+1];
			argi++;
		}

		// xxx negate the default once this is working.
		// xxx put into online help.
		else if (streq(argv[argi], "--mmap")) {
			popts->use_file_reader_mmap = TRUE;
		}
		else if (streq(argv[argi], "--no-mmap")) {
			popts->use_file_reader_mmap = FALSE;
		}
		else if (streq(argv[argi], "--seed")) {
			check_arg_count(argv, argi, argc, 2);
			if (sscanf(argv[argi+1], "0x%x", &rand_seed) == 1) {
				have_rand_seed = TRUE;
			} else if (sscanf(argv[argi+1], "%u", &rand_seed) == 1) {
				have_rand_seed = TRUE;
			} else {
				main_usage(argv[0], 1);
			}
			argi++;
		}

		else
			nusage(argv[0], argv[argi]);
	}

	// xxx final coalesce: file count -> 0 <-> no mmap
	popts->plrec_reader_mmap  = lrec_reader_alloc(rdesc, TRUE,
		popts->irs, popts->ifs, popts->allow_repeat_ifs, popts->ips, popts->allow_repeat_ips);
	popts->plrec_reader_stdio = lrec_reader_alloc(rdesc, FALSE,
		popts->irs, popts->ifs, popts->allow_repeat_ifs, popts->ips, popts->allow_repeat_ips);
	if (popts->plrec_reader_mmap == NULL || popts->plrec_reader_stdio == NULL)
		main_usage(argv[0], 1);

	if      (streq(wdesc, "dkvp"))   popts->plrec_writer = lrec_writer_dkvp_alloc(popts->ors, popts->ofs, popts->ops);
	else if (streq(wdesc, "csv"))    popts->plrec_writer = lrec_writer_csv_alloc(popts->ors, popts->ofs);
	else if (streq(wdesc, "nidx"))   popts->plrec_writer = lrec_writer_nidx_alloc(popts->ors, popts->ofs);
	else if (streq(wdesc, "xtab"))   popts->plrec_writer = lrec_writer_xtab_alloc();
	else if (streq(wdesc, "pprint")) popts->plrec_writer = lrec_writer_pprint_alloc(left_align_pprint);
	else {
		main_usage(argv[0], 1);
	}

	if ((argc - argi) < 1) {
		main_usage(argv[0], 1);
	}

	if (have_rand_seed) {
		mtrand_init(rand_seed);
	} else {
		mtrand_init_default();
	}

	popts->pmapper_list = sllv_alloc();
	while (TRUE) {
		check_arg_count(argv, argi, argc, 1);
		char* verb = argv[argi];

		mapper_setup_t* pmapper_setup = look_up_mapper_setup(verb);
		if (pmapper_setup == NULL) {
			fprintf(stderr, "%s: verb \"%s\" not found. Please use \"%s --help\" for a list.\n",
				argv[0], verb, argv[0]);
			exit(1);
		}

		if ((argc - argi) >= 2) {
			if (streq(argv[argi+1], "-h") || streq(argv[argi+1], "--help")) {
				pmapper_setup->pusage_func(argv[0], verb);
				exit(0);
			}
		}

		// It's up to the parse func to print its usage on CLI-parse failure.
		mapper_t* pmapper = pmapper_setup->pparse_func(&argi, argc, argv);
		if (pmapper == NULL) {
			exit(1);
		}
		sllv_add(popts->pmapper_list, pmapper);

		// xxx cmt
		if (argi >= argc || !streq(argv[argi], "then"))
			break;
		argi++;
	}

	popts->filenames = &argv[argi];

	return popts;
}

// ----------------------------------------------------------------
void cli_opts_free(cli_opts_t* popts) {

	popts->plrec_reader_stdio->pfree_func(popts->plrec_reader_stdio->pvstate);

	for (sllve_t* pe = popts->pmapper_list->phead; pe != NULL; pe = pe->pnext) {
		mapper_t* pmapper = pe->pvdata;
		pmapper->pfree_func(pmapper->pvstate);
	}
	sllv_free(popts->pmapper_list);

	popts->plrec_writer->pfree_func(popts->plrec_writer->pvstate);
	free(popts);
}
