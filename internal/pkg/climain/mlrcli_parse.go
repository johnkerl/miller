package climain

import (
	"fmt"
	"os"

	"mlr/internal/pkg/auxents/help"
	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/transformers"
	"mlr/internal/pkg/types"
	"mlr/internal/pkg/version"
)

// ParseCommandLine is the entrypoint for handling the Miller command line:
// flags, verbs and their flags, and input file name(s).
func ParseCommandLine(args []string) (
	options cli.TOptions,
	recordTransformers []transformers.IRecordTransformer,
	err error,
) {
	options = cli.DefaultOptions()
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
			cli.CheckArgCount(args, argi, argc, 1)
			argi += 2
		} else if args[argi] == "--version" {
			fmt.Printf("mlr %s\n", version.STRING)
			os.Exit(0)
		} else if args[argi] == "--bare-version" {
			fmt.Printf("%s\n", version.STRING)
			os.Exit(0)

		} else if help.ParseTerminalUsage(args[argi]) {
			// Most help is in the 'mlr help' auxent but there are a few shorthands
			// like 'mlr -h' and 'mlr -F'.
			os.Exit(0)

		} else if cli.FLAG_TABLE.Parse(args, argc, &argi, &options) {
			// handled

		} else {
			// unhandled
			fmt.Fprintf(os.Stderr, "%s: option \"%s\" not recognized.\n", "mlr", args[argi])
			fmt.Fprintf(os.Stderr, "Please run \"%s --help\" for usage information.\n", "mlr")
			os.Exit(1)
		}
	}

	// Check now to avoid confusing timezone-library behavior later on
	lib.SetTZFromEnv()

	cli.FinalizeReaderOptions(&options.ReaderOptions)
	cli.FinalizeWriterOptions(&options.WriterOptions)

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

	if cli.DecideFinalFlatten(&options.WriterOptions) {
		// E.g. '{"req": {"method": "GET", "path": "/api/check"}}' becomes
		// req.method=GET,req.path=/api/check.
		transformer, err := transformers.NewTransformerFlatten(options.WriterOptions.FLATSEP, nil)
		lib.InternalCodingErrorIf(err != nil)
		lib.InternalCodingErrorIf(transformer == nil)
		recordTransformers = append(recordTransformers, transformer)
	}

	if cli.DecideFinalUnflatten(&options) {
		// E.g.  req.method=GET,req.path=/api/check becomes
		// '{"req": {"method": "GET", "path": "/api/check"}}'
		transformer, err := transformers.NewTransformerUnflatten(options.WriterOptions.FLATSEP, nil)
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
		fmt.Fprintf(os.Stderr, "%s: -I option (in-place operation) requires input files.\n", "mlr")
		os.Exit(1)
	}

	if options.HaveRandSeed {
		lib.SeedRandom(int64(options.RandSeed))
	}

	return options, recordTransformers, nil
}

// parseTransformers returns a list of transformers, from the starting point in
// args given by *pargi.  Bumps *pargi to point to remaining
// post-transformer-setup args, i.e. filenames.
func parseTransformers(
	args []string,
	pargi *int,
	argc int,
	options *cli.TOptions,
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
		fmt.Fprintf(os.Stderr, "%s: no verb supplied.\n", "mlr")
		help.MainUsage(os.Stderr)
		os.Exit(1)
	}

	onFirst := true

	for {
		cli.CheckArgCount(args, argi, argc, 1)
		verb := args[argi]

		transformerSetup := transformers.LookUp(verb)
		if transformerSetup == nil {
			fmt.Fprintf(os.Stderr,
				"%s: verb \"%s\" not found. Please use \"%s --help\" for a list.\n",
				"mlr", verb, "mlr")
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
				fmt.Fprintf(os.Stderr, "%s: missing next verb after \"then\".\n", "mlr")
				os.Exit(1)
			} else {
				argi++
			}
		}
	}

	*pargi = argi
	return transformerList, ignoresInput, nil
}
