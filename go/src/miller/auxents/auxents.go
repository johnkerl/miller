// ================================================================
// Little side-programs for hex, unhex, LF <-> CR/LF, etc which are delivered
// within the mlr exectuable.
// ================================================================

package auxents

import (
	"fmt"
	"os"
)

// ----------------------------------------------------------------
type tAuxMain func(args []string) int
type tAuxUsage func(verbName string, ostream *os.File, exitCode int)

type tAuxLookupEntry struct {
	name  string
	main  tAuxMain
	usage tAuxUsage
}

// We get a Golang "initialization loop" if this is defined statically. So, we
// use a "package init" function.
var _AUX_LOOKUP_TABLE = []tAuxLookupEntry{}

func init() {
	_AUX_LOOKUP_TABLE = []tAuxLookupEntry{
		{"aux-list", auxListMain, auxListUsage},
		//	{ "lecat",           lecat_main,           lecat_usage           },
		//	{ "termcvt",         termcvt_main,         termcvt_usage         },
		//	{ "hex",             hex_main,             hex_usage             },
		//	{ "unhex",           unhex_main,           unhex_usage           },
	}
}

// ================================================================
func Dispatch(args []string) {
	if len(args) < 2 {
		return
	}
	verb := args[1]

	for _, entry := range _AUX_LOOKUP_TABLE {
		if verb == entry.name {
			os.Exit(entry.main(args))
		}
	}

	// Else, return control to mlr.go for the rest of Miller.
}

// ================================================================
func auxListUsage(verbName string, ostream *os.File, exitCode int) {
	fmt.Fprintf(ostream, "Usage: %s %s [options]\n", os.Args[0], verbName)
	fmt.Fprintf(ostream, "Options:\n")
	fmt.Fprintf(ostream, "-h or --help: print this message\n")
	os.Exit(exitCode)
}

func auxListMain(args []string) int {
	ShowAuxEntries(os.Stdout)
	return 0
}

// This symbol is exported for 'mlr --help'.
func ShowAuxEntries(ostream *os.File) {
	fmt.Fprintf(ostream, "Available subcommands:\n")
	for _, entry := range _AUX_LOOKUP_TABLE {
		fmt.Fprintf(ostream, "  %s\n", entry.name)
	}

	fmt.Fprintf(
		ostream,
		"For more information, please invoke %s {subcommand} --help.\n",
		os.Args[0],
	)
}

// ================================================================
//static void lecat_usage(char* verbName, ostream *os.File, int exitCode) {
//	fmt.Fprintf(ostream, "Usage: %s %s [options] {zero or more file names}\n", os.Args[0], verbName);
//	fmt.Fprintf(ostream, "Simply echoes input, but flags CR characters in red and LF characters in green.\n");
//	fmt.Fprintf(ostream, "If zero file names are supplied, standard input is read.\n");
//	fmt.Fprintf(ostream, "Options:\n");
//	fmt.Fprintf(ostream, "--mono: don't try to colorize the output\n");
//	fmt.Fprintf(ostream, "-h or --help: print this message\n");
//	os.Exit(exitCode);
//}

//static int lecat_main(int argc, char** argv) {
//	int ok = 1;
//	int do_color = TRUE;
//
//	if (argc >= 3) {
//		if (streq(argv[2], "-h") || streq(argv[2], "--help")) {
//			lecat_usage(argv[0], argv[1], stdout, 0);
//		}
//	}
//
//	// 'mlr' and 'lecat' are already argv[0] and argv[1].
//	int argb = 2;
//	if (argc >= 3 && argv[argb][0] == '-') {
//		if (streq(argv[argb], "--mono")) {
//			do_color = FALSE;
//			argb++;
//		} else {
//			fmt.Fprintf(stderr, "%s %s: unrecognized option \"%s\".\n",
//				argv[0], argv[1], argv[argb]);
//			return 1;
//		}
//	}
//
//	if (argb == argc) {
//		ok = ok && lecat_stream(stdin, do_color);
//	} else {
//		for (int argi = argb; argi < argc; argi++) {
//			char* file_name = argv[argi];
//			FILE* input_stream = fopen(file_name, "r");
//			if (input_stream == NULL) {
//				perror(file_name);
//				os.Exit(1);
//			}
//			ok = lecat_stream(input_stream, do_color);
//			fclose(input_stream);
//		}
//	}
//	return ok ? 0 : 1;
//}

//static int lecat_stream(inputStream *os.File, doColor bool) {
//	while (1) {
//		int c = fgetc(input_stream);
//		if (c == EOF)
//			break;
//		if (c == '\r') {
//			if (do_color)
//				printf("\033[31;01m"); // xterm red
//			printf("[CR]");
//			if (do_color)
//				printf("\033[0m");
//		} else if (c == '\n') {
//			if (do_color)
//				printf("\033[32;01m"); // xterm green
//			printf("[LF]\n");
//			if (do_color)
//				printf("\033[0m");
//		} else {
//			putchar(c);
//		}
//	}
//	return 1;
//}

//// ================================================================
//static int termcvt_stream(inputStream *os.File, outputStream *os.File, char* inend, char* outend) {
//	size_t line_length = MLR_ALLOC_READ_LINE_INITIAL_SIZE;
//	int inend_length = strlen(inend);
//	while (1) {
//		char* line = mlr_alloc_read_line_multiple_delimiter(input_stream, inend, inend_length, &line_length);
//		if (line == NULL) {
//			break;
//		}
//		fputs(line, output_stream);
//		fputs(outend, output_stream);
//		free(line);
//	}
//	return 1;
//}

//// ----------------------------------------------------------------
//static void termcvt_usage(char* verbName, ostream *os.File, int exitCode) {
//	fmt.Fprintf(ostream, "Usage: %s %s [option] {zero or more file names}\n", os.Args[0], verbName);
//	fmt.Fprintf(ostream, "Option (exactly one is required):\n");
//	fmt.Fprintf(ostream, "--cr2crlf\n");
//	fmt.Fprintf(ostream, "--lf2crlf\n");
//	fmt.Fprintf(ostream, "--crlf2cr\n");
//	fmt.Fprintf(ostream, "--crlf2lf\n");
//	fmt.Fprintf(ostream, "--cr2lf\n");
//	fmt.Fprintf(ostream, "--lf2cr\n");
//	fmt.Fprintf(ostream, "-I in-place processing (default is to write to stdout)\n");
//	fmt.Fprintf(ostream, "-h or --help: print this message\n");
//	fmt.Fprintf(ostream, "Zero file names means read from standard input.\n");
//	fmt.Fprintf(ostream, "Output is always to standard output; files are not written in-place.\n");
//	os.Exit(exitCode);
//}

//// ----------------------------------------------------------------
//static int termcvt_main(int argc, char** argv) {
//	int ok = 1;
//	char* inend  = "\n";
//	char* outend = "\n";
//	int do_in_place = FALSE;
//
//	// argv[0] is 'mlr'
//	// argv[1] is 'termcvt'
//	// argv[2] is '--some-option'
//	// argv[3] and above are filenames
//	if (argc < 2)
//		termcvt_usage(argv[0], argv[1], stderr, 1);
//	int argi;
//	for (argi = 2; argi < argc; argi++) {
//		char* opt = argv[argi];
//
//		if (opt[0] != '-')
//			break;
//
//		if (streq(opt, "-h") || streq(opt, "--help")) {
//			termcvt_usage(argv[0], argv[1], stdout, 0);
//		} else if (streq(opt, "-I")) {
//			do_in_place = TRUE;
//		} else if (streq(opt, "--cr2crlf")) {
//			inend  = "\r";
//			outend = "\r\n";
//		} else if (streq(opt, "--lf2crlf")) {
//			inend  = "\n";
//			outend = "\r\n";
//		} else if (streq(opt, "--crlf2cr")) {
//			inend  = "\r\n";
//			outend = "\r";
//		} else if (streq(opt, "--lf2cr")) {
//			inend  = "\n";
//			outend = "\r";
//		} else if (streq(opt, "--crlf2lf")) {
//			inend  = "\r\n";
//			outend = "\n";
//		} else if (streq(opt, "--cr2lf")) {
//			inend  = "\r";
//			outend = "\n";
//		} else {
//			termcvt_usage(argv[0], argv[1], stdout, 0);
//		}
//	}
//
//	int nfiles = argc - argi;
//	if (nfiles == 0) {
//		ok = ok && termcvt_stream(stdin, stdout, inend, outend);
//
//	} else if (do_in_place) {
//		for (; argi < argc; argi++) {
//			char* file_name = argv[argi];
//			char* temp_name = alloc_suffixed_temp_file_name(file_name);
//			FILE* input_stream = fopen(file_name, "r");
//			FILE* output_stream = fopen(temp_name, "wb");
//
//			if (input_stream == NULL) {
//				perror("fopen");
//				fmt.Fprintf(stderr, "%s: Could not open \"%s\" for read.\n",
//					os.Args[0], file_name);
//				os.Exit(1);
//			}
//			if (output_stream == NULL) {
//				perror("fopen");
//				fmt.Fprintf(stderr, "%s: Could not open \"%s\" for write.\n",
//					os.Args[0], temp_name);
//				os.Exit(1);
//			}
//
//			ok = termcvt_stream(input_stream, output_stream, inend, outend);
//
//			fclose(input_stream);
//			fclose(output_stream);
//
//			int rc = rename(temp_name, file_name);
//			if (rc != 0) {
//				perror("rename");
//				fmt.Fprintf(stderr, "%s: Could not rename \"%s\" to \"%s\".\n",
//					os.Args[0], temp_name, file_name);
//				os.Exit(1);
//			}
//			free(temp_name);
//		}
//
//	} else {
//		for (; argi < argc; argi++) {
//			char* file_name = argv[argi];
//			FILE* input_stream = fopen(file_name, "r");
//			if (input_stream == NULL) {
//				perror(file_name);
//				os.Exit(1);
//			}
//			ok = termcvt_stream(input_stream, stdout, inend, outend);
//			fclose(input_stream);
//		}
//
//	}
//	return ok ? 0 : 1;
//}

// ================================================================
// Copyright (c) 1998 John Kerl.
// ================================================================
// This is a simple hex dump with hex offsets to the left, hex data in the
// middle, and ASCII at the right.  This is a subset of the functionality of
// Unix od; I wrote it in my NT days.
//
// Example:
//
// $ d2h $(jot 0 128) | unhex | hex
// 00000000: 00 01 02 03  04 05 06 07  08 09 0a 0b  0c 0d 0e 0f |................|
// 00000010: 10 11 12 13  14 15 16 17  18 19 1a 1b  1c 1d 1e 1f |................|
// 00000020: 20 21 22 23  24 25 26 27  28 29 2a 2b  2c 2d 2e 2f | !"#$%&'()*+,-./|
// 00000030: 30 31 32 33  34 35 36 37  38 39 3a 3b  3c 3d 3e 3f |0123456789:;<=>?|
// 00000040: 40 41 42 43  44 45 46 47  48 49 4a 4b  4c 4d 4e 4f |@ABCDEFGHIJKLMNO|
// 00000050: 50 51 52 53  54 55 56 57  58 59 5a 5b  5c 5d 5e 5f |PQRSTUVWXYZ[\]^_|
// 00000060: 60 61 62 63  64 65 66 67  68 69 6a 6b  6c 6d 6e 6f |`abcdefghijklmno|
// 00000070: 70 71 72 73  74 75 76 77  78 79 7a 7b  7c 7d 7e 7f |pqrstuvwxyz{|}~.|
// ================================================================

//#define LINE_LENGTH_MAX 8192

// ----------------------------------------------------------------
//static void hex_usage(char* verbName, ostream *os.File, int exitCode) {
//	fmt.Fprintf(ostream, "Usage: %s %s [options] {zero or more file names}\n", os.Args[0], verbName);
//	fmt.Fprintf(ostream, "Simple hex-dump.\n");
//	fmt.Fprintf(ostream, "If zero file names are supplied, standard input is read.\n");
//	fmt.Fprintf(ostream, "Options:\n");
//	fmt.Fprintf(ostream, "-r: print only raw hex without leading offset indicators or trailing ASCII dump.\n");
//	fmt.Fprintf(ostream, "-h or --help: print this message\n");
//	os.Exit(exitCode);
//}

//----------------------------------------------------------------------
// 'mlr' and 'hex' are already argv[0] and argv[1].

//static int hex_main(int argc, char **argv) {
//	char * filename;
//	FILE * in_fp;
//	FILE * out_fp;
//	int do_raw = 0;
//	int argi = 2;
//
//	if (argc >= 3) {
//		if (streq(argv[2], "-r")) {
//			do_raw = 1;
//			argi++;
//		} else if (streq(argv[2], "-h") || streq(argv[2], "--help")) {
//			hex_usage(argv[0], argv[1], stdout, 0);
//		}
//	}
//
//	int num_file_names = argc - argi;
//
//	if (num_file_names == 0) {
//#ifdef WINDOWS
//		setmode(fileno(stdin), O_BINARY);
//#endif //WINDOWS
//		hex_dump_fp(stdin, stdout, do_raw);
//	} else {
//		for ( ; argi < argc; argi++) {
//			if (!do_raw) {
//				if (num_file_names > 1)
//					printf("%s:\n", argv[argi]);
//			}
//			filename = argv[argi];
//			in_fp    = fopen(filename, "rb");
//			out_fp   = stdout;
//			if (in_fp == NULL) {
//				fmt.Fprintf(stderr, "Couldn't open \"%s\"; skipping.\n",
//					filename);
//			}
//			else {
//				hex_dump_fp(in_fp, out_fp, do_raw);
//				fclose(in_fp);
//				if (out_fp != stdout)
//					fclose(out_fp);
//
//			}
//			if (!do_raw) {
//				if (num_file_names > 1)
//					printf("\n");
//			}
//		}
//	}
//
//	return 0;
//}

//// ----------------------------------------------------------------
//static void hex_dump_fp(FILE *in_fp, out_fp *os.File, int do_raw) {
//	const int bytes_per_clump = 4;
//	const int clumps_per_line = 4;
//	const int buffer_size = bytes_per_clump * clumps_per_line;
//	unsigned char buf[buffer_size];
//	long num_bytes_read;
//	long num_bytes_total = 0;
//	int byteno;
//
//	while ((num_bytes_read=fread(buf, sizeof(unsigned char),
//		buffer_size, in_fp)) > 0)
//	{
//		if (!do_raw) {
//			printf("%08lx: ", num_bytes_total);
//		}
//
//		for (byteno = 0; byteno < num_bytes_read; byteno++) {
//			unsigned int temp = buf[byteno];
//			printf("%02x ", temp);
//			if ((byteno % bytes_per_clump) ==
//			(bytes_per_clump - 1))
//			{
//				if ((byteno > 0) && (byteno < buffer_size-1))
//					printf(" ");
//			}
//		}
//		for (byteno = num_bytes_read; byteno < buffer_size; byteno++) {
//			printf("   ");
//			if ((byteno % bytes_per_clump) ==
//			(bytes_per_clump - 1))
//			{
//				if ((byteno > 0) && (byteno < buffer_size-1))
//					printf(" ");
//			}
//		}
//
//		if (!do_raw) {
//			printf("|");
//			for (byteno = 0; byteno < num_bytes_read; byteno++) {
//				unsigned char temp = buf[byteno];
//				if (!isprint(temp))
//					temp = '.';
//				printf("%c", temp);
//			}
//
//			printf("|");
//		}
//
//		printf("\n");
//		num_bytes_total += num_bytes_read;
//	}
//}

// ================================================================
//static void unhex_usage(char* verbName, ostream *os.File, int exitCode) {
//	fmt.Fprintf(ostream, "Usage: %s %s [option] {zero or more file names}\n", os.Args[0], verbName);
//	fmt.Fprintf(ostream, "Options:\n");
//	fmt.Fprintf(ostream, "-h or --help: print this message\n");
//	fmt.Fprintf(ostream, "Zero file names means read from standard input.\n");
//	fmt.Fprintf(ostream, "Output is always to standard output; files are not written in-place.\n");
//	os.Exit(exitCode);
//}

// ----------------------------------------------------------------
//int unhex_main(int argc, char ** argv) {
//	// 'mlr' and 'unhex' are already argv[0] and argv[1].
//	if (argc >= 3) {
//		if (streq(argv[2], "-h") || streq(argv[2], "--help")) {
//			unhex_usage(argv[0], argv[1], stdout, 0);
//		}
//	}
//
//	int exitCode = 0;
//	if (argc == 2) {
//		unhex_fp(stdin, stdout);
//	} else {
//		for (int argi = 2; argi < argc; argi++) {
//			char* filename = argv[argi];
//			FILE* istream = fopen(filename, "rb");
//			if (istream == NULL) {
//				fmt.Fprintf(stderr, "%s %s: Couldn't open \"%s\"; skipping.\n",
//					argv[0], argv[1], filename);
//				exitCode = 1;
//			} else {
//				unhex_fp(istream, stdout);
//				fclose(istream);
//			}
//		}
//	}
//
//	return exitCode;
//}

// ----------------------------------------------------------------

//static void unhex_fp(FILE *istream, FILE *ostream) {
//	unsigned char byte;
//	unsigned temp;
//	int count;
//	while ((count=fscanf(istream, "%x", &temp)) > 0) {
//		byte = temp;
//		fwrite (&byte, sizeof(byte), 1, ostream);
//	}
//}
