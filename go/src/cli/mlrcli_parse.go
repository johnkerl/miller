package cli

import (
	"fmt"
	"os"

	"miller/src/auxents/help"
	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transformers"
	"miller/src/types"
	"miller/src/version"
)

// ----------------------------------------------------------------
func ParseCommandLine(args []string) (
	options cliutil.TOptions,
	recordTransformers []transformers.IRecordTransformer,
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
	transformerList []transformers.IRecordTransformer,
	ignoresInput bool,
	err error,
) {

	transformerList = make([]transformers.IRecordTransformer, 0)
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

		transformerSetup := transformers.LookUp(verb)
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
func parseTerminalUsage(args []string, argc int, argi int) bool {
	if args[argi] == "--version" {
		fmt.Printf("Miller %s\n", version.STRING)
		return true

	} else if args[argi] == "-h" || args[argi] == "--help" {
		help.MainUsage(os.Stdout)
		os.Exit(0)
		return true

	} else if args[argi] == "-l" {
		// TODO: move to help?
		help.ListAllVerbNamesAsParagraph()
		return true
	} else if args[argi] == "-L" {
		help.ListAllVerbNames()
		return true

	} else if args[argi] == "-f" {
		// TODO: mlr help function-details
		// all functions with usage-strings
		return true
	} else if args[argi] == "-F" {
		// TODO: mlr help function-names
		// all functions, names only
		return true

	} else if args[argi] == "-k" {
		help.HelpKeyword([]string{})
		// TODO: all keywords, long version
		return true
	} else if args[argi] == "-K" {
		// TODO: refacctor
		// TODO: all keywords, names only
		help.ListKeywords([]string{})
		return true

	} else {
		return false
	}
}
