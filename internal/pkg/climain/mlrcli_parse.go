// ================================================================
// Miller main command-line parsing.
//
// Before Miller 6 the ordering was:
// * mlr
// * main flags like --icsv --ojson
// * verbs and their flags like cat -n
// * data-file names
// and the command-line parser was one-pass.
//
// In Miller 6 we have as keystroke-reducers 'mlr -s', for '#!mlr -s',
// or simply better support for mlr inside of '#!/bin/sh' scripts:
//
//   mlr {flags} {verbs}             -- [more flags] [more verbs] {data file names}
//   [the part inside a script file]    [the part outside]
//
// For example, suppose someone wants to reuse the following:
//   mlr --icsv --json head -n 10
// either via a #!mlr -s script, maybe "peek.mlr":
//   #!/usr/bin/env mlr -s
//   --icsv --json head -n 10
// or a #!/bin/bash script, maybe "peek.sh"
//   #!/bin/bash
//   mlr --icsv --json head -n 10 -- "$@"
// Then they can do 'peek.mlr myfile.csv' or 'peek.sh myfile.csv' which is great.
//
// But suppose they want to do
//   peek.sh --jlistwrap myfile.csv
// Then the Miller command line received here is
//   mlr --icsv --json head -n 10 -- --jlistwrap myfile.csv
// Or, maybe their part inside the '#!mlr' or '#!/bin/sh' file is all verbs,
// and they want to specify format-flags like '--icsv --ojson' outside of that
// script.  It's very reasonable for them to want to put the --jlistwrap,
// --icsv, --ojson, etc. after their keystroke-saver script. But this now means
// that there can be main-flags (and/or 'then someotherverb') *after* the verb
// chain from inside the keystroke-saver.
//
// Also, verbs/transformers must be constructed *after* all main-flags are
// parsed -- since some of them depend on main-flags, e.g. join, put/filter,
// and tee which use things like --csv for their I/O options.
//
// Therefore the command-line parsing is now two-pass.
// * Pass 1:
//   o 'mlr' is first
//   o Split the []args into "sequences" of main-flags, verbs and their flags,
//     and data-file names.
//   o For example in the above 'mlr --icsv --json head -n 10 -- --jlistwrap myfile.csv'
//     we have
//     main-flag sequences ['--icsv'] ['--json'] [--jlistwrap],
//     verb-seqeunce ['head' '-n' '10']
//     data-file names ['myfile.csv'].
//   o Any exiting flags like --version or --help are dispatched here.
//   o To do that splitting we invoke the flag-table parser with throwaway options struct,
//     and we invoke the transformers' ParseCLI functions with doConstruct = false.
// * Pass 2:
//   o Process the flag-sequences in the order they were encountered, into a
//     for-real-use options struct.
//   o Process the verb-sequences in the order they were encountered, and construct
//     transformers.
//   o Some jargon from programming languages we can use here for illustration
//     is that we are "hoisting" the main-flags as if they had been written on
//     the command line before the verbs.
//
// We need to require a '--' between a verb and a main-flag so the main-flag
// doesn't look like a verb flag. For example, in 'mlr head -n 10 --csv
// foo.csv' the '--csv' looks like it belongs to the 'head' verb. When people
// use '#!/bin/sh' scripts they need to insert the '--' in 'mlr head -n 10 --
// --csv foo.csv'; for 'mlr -s' we insert the '--' for them.
// ================================================================

package climain

import (
	"fmt"
	"os"

	"github.com/johnkerl/miller/internal/pkg/auxents/help"
	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/transformers"
	"github.com/johnkerl/miller/internal/pkg/version"
)

// ParseCommandLine is the entrypoint for handling the Miller command line:
// flags, verbs and their flags, and input file name(s).
func ParseCommandLine(
	args []string,
) (
	options *cli.TOptions,
	recordTransformers []transformers.IRecordTransformer,
	err error,
) {
	// mlr -s scriptfile {data-file names ...} means take the contents of
	// scriptfile as if it were command-line items.
	args = maybeInterpolateDashS(args)

	// Pass one as described at the top of this file.
	flagSequences, verbSequences, dataFileNames := parseCommandLinePassOne(args)

	// Pass two as described at the top of this file.
	return parseCommandLinePassTwo(flagSequences, verbSequences, dataFileNames)
}

// parseCommandLinePassOne is as described at the top of this file.
func parseCommandLinePassOne(
	args []string,
) (
	flagSequences [][]string,
	verbSequences [][]string,
	dataFileNames []string,
) {
	flagSequences = make([][]string, 0)
	verbSequences = make([][]string, 0)
	dataFileNames = make([]string, 0)

	// All verbs after the first must be preceded with "then"
	onFirst := true

	// Throwaway options as described above: passed into the flag-table parser
	// but we'll use for-real-use options in pass two.
	options := cli.DefaultOptions()

	argi := 1
	argc := len(args)

	for argi < argc /* variable increment within loop body */ {

		// Old argi is at start of sequence; argi will be after.
		oargi := argi

		if args[argi][0] == '-' {
			if args[argi] == "--version" {
				// Exiting flag: handle it immediately.
				fmt.Printf("mlr %s\n", version.STRING)
				os.Exit(0)
			} else if args[argi] == "--bare-version" {
				// Exiting flag: handle it immediately.
				fmt.Printf("%s\n", version.STRING)
				os.Exit(0)
			} else if help.ParseTerminalUsage(args[argi]) {
				// Exiting flag: handle it immediately.
				// Most help is in the 'mlr help' auxent but there are a few
				// shorthands like 'mlr -h' and 'mlr -F'.
				os.Exit(0)

			} else if args[argi] == "--norc" {
				flagSequences = append(flagSequences, args[oargi:argi])
				argi += 1

			} else if cli.FLAG_TABLE.Parse(args, argc, &argi, options) {
				flagSequences = append(flagSequences, args[oargi:argi])

			} else if args[argi] == "--" {
				// This separates a main-flag from the verb/verb-flags before it
				argi += 1

			} else {
				// Unrecognized main-flag. Fatal it here, and don't send it to pass two.
				fmt.Fprintf(os.Stderr, "%s: option \"%s\" not recognized.\n", "mlr", args[argi])
				fmt.Fprintf(os.Stderr, "Please run \"%s --help\" for usage information.\n", "mlr")
				os.Exit(1)
			}

		} else if onFirst || args[argi] == "then" {
			// The first verb in the then-chain can *optionally* be preceded by
			// 'then'.  The others one *must* be.
			if args[argi] == "then" {
				cli.CheckArgCount(args, argi, argc, 1)
				oargi++
				argi++
			}
			verb := args[argi]
			onFirst = false

			transformerSetup := transformers.LookUp(verb)
			if transformerSetup == nil {
				fmt.Fprintf(os.Stderr,
					"%s: verb \"%s\" not found. Please use \"%s --help\" for a list.\n",
					"mlr", verb, "mlr")
				os.Exit(1)
			}

			// It's up to the parse func to print its usage, and exit 1, on
			// CLI-parse failure.  Also note: this assumes main reader/writer opts
			// are all parsed *before* transformer parse-CLI methods are invoked.
			transformer := transformerSetup.ParseCLIFunc(
				&argi,
				argc,
				args,
				options,
				false, // false for first pass of CLI-parse, true for second pass -- this is the first pass
			)
			// For pass one we want the verbs to identify the arg-sequences
			// they own within the command line, but not construct
			// transformers.
			lib.InternalCodingErrorIf(transformer != nil)

			verbSequences = append(verbSequences, args[oargi:argi])

		} else {
			// After main-flag sequences and verb sequences, data-file names
			// still come last on the command line.
			break
		}
	}

	for ; argi < argc; argi++ {
		dataFileNames = append(dataFileNames, args[argi])
	}

	if len(verbSequences) == 0 {
		fmt.Fprintf(os.Stderr, "%s: no verb supplied.\n", "mlr")
		help.MainUsage(os.Stderr)
		os.Exit(1)
	}

	return flagSequences, verbSequences, dataFileNames
}

// parseCommandLinePassTwo is as described at the top of this file.
func parseCommandLinePassTwo(
	flagSequences [][]string,
	verbSequences [][]string,
	dataFileNames []string,
) (
	options *cli.TOptions,
	recordTransformers []transformers.IRecordTransformer,
	err error,
) {
	// Options take in-code defaults, then overridden by .mlrrc (if any and if
	// desired), then those in turn overridden by command-line flags.
	options = cli.DefaultOptions()
	recordTransformers = make([]transformers.IRecordTransformer, 0)
	err = nil
	ignoresInput := false

	// Load a .mlrrc file unless --norc was a main-flag on the command line.
	loadMlrrc := true
	for _, flagSequence := range flagSequences {
		lib.InternalCodingErrorIf(len(flagSequence) < 1)
		if flagSequence[0] == "--norc" {
			loadMlrrc = false
			break
		}
	}
	if loadMlrrc {
		loadMlrrcOrDie(options)
	}

	// Process the flag-sequences in order from pass one. We assume all the
	// exiting flags like --help and --version were already processed, so all
	// main-flags making it here to pass two are for the flag-table parser.
	for _, flagSequence := range flagSequences {
		argi := 0
		args := flagSequence
		argc := len(args)
		lib.InternalCodingErrorIf(argc == 0)

		// Parse the main-flag into the options struct.
		rc := cli.FLAG_TABLE.Parse(args, argc, &argi, options)

		// Should have been parsed OK in pass one.
		lib.InternalCodingErrorIf(rc != true)
		// Make sure we consumed the entire flag sequence as parsed by pass one.
		lib.InternalCodingErrorIf(argi != argc)
	}

	// Check now to avoid confusing timezone-library behavior later on
	lib.SetTZFromEnv()

	cli.FinalizeReaderOptions(&options.ReaderOptions)
	cli.FinalizeWriterOptions(&options.WriterOptions)

	// Set an optional global formatter for floating-point values
	if options.WriterOptions.FPOFMT != "" {
		err = mlrval.SetFloatOutputFormat(options.WriterOptions.FPOFMT)
		if err != nil {
			return options, recordTransformers, err
		}
	}

	// Now process the verb-sequences from pass one, with options-struct set up
	// and finalized.
	for i, verbSequence := range verbSequences {
		argi := 0 // xxx needed?
		args := verbSequence
		argc := len(args)
		lib.InternalCodingErrorIf(argc == 0)

		// Non-existent verbs should have been fatalled in pass one.
		transformerSetup := transformers.LookUp(args[0])
		lib.InternalCodingErrorIf(transformerSetup == nil)

		// It's up to the parse func to print its usage, and exit 1, on
		// CLI-parse failure.
		transformer := transformerSetup.ParseCLIFunc(
			&argi,
			argc,
			args,
			options,
			true, // false for first pass of CLI-parse, true for second pass -- this is the first pass
		)
		// Unparsable verb-setups should have been found in pass one.
		lib.InternalCodingErrorIf(transformer == nil)
		// Make sure we consumed the entire verb sequence as parsed by pass one.
		lib.InternalCodingErrorIf(argi != argc)

		// E.g. then-chain begins with seqgen
		if i == 0 && transformerSetup.IgnoresInput {
			ignoresInput = true
		}

		recordTransformers = append(recordTransformers, transformer)
	}

	if ignoresInput {
		options.NoInput = true // e.g. then-chain begins with seqgen
	}

	if cli.DecideFinalFlatten(&options.WriterOptions) {
		// E.g. '{"req": {"method": "GET", "path": "/api/check"}}' becomes
		// req.method=GET,req.path=/api/check.
		transformer, err := transformers.NewTransformerFlatten(options.WriterOptions.FLATSEP, options, nil)
		lib.InternalCodingErrorIf(err != nil)
		lib.InternalCodingErrorIf(transformer == nil)
		recordTransformers = append(recordTransformers, transformer)
	}

	if cli.DecideFinalUnflatten(options) {
		// E.g.  req.method=GET,req.path=/api/check becomes
		// '{"req": {"method": "GET", "path": "/api/check"}}'
		transformer, err := transformers.NewTransformerUnflatten(options.WriterOptions.FLATSEP, options, nil)
		lib.InternalCodingErrorIf(err != nil)
		lib.InternalCodingErrorIf(transformer == nil)
		recordTransformers = append(recordTransformers, transformer)
	}

	// There may already be one or more because of --from on the command line,
	// so append.
	for _, dataFileName := range dataFileNames {
		options.FileNames = append(options.FileNames, dataFileName)
	}

	// E.g. mlr -n put -v '...'
	if options.NoInput {
		options.FileNames = nil
	}

	if options.DoInPlace && (options.FileNames == nil || len(options.FileNames) == 0) {
		fmt.Fprintf(os.Stderr, "%s: -I option (in-place operation) requires input files.\n", "mlr")
		os.Exit(1)
	}

	if options.HaveRandSeed {
		lib.SeedRandom(int64(options.RandSeed))
	}

	return options, recordTransformers, nil
}
