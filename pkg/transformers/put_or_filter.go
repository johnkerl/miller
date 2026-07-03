package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/dsl/cst"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/runtime"
	"github.com/johnkerl/miller/v6/pkg/types"
	"github.com/johnkerl/pgpg/go/lib/pkg/asts"
)

const verbNamePut = "put"

var putOptions = []OptionSpec{
	{Flag: "-f", Arg: "{file name}", Type: "filename", Desc: "File containing a DSL expression (see examples below). If the filename is a directory, all *.mlr files in that directory are loaded.", Repeatable: true},
	{Flag: "-e", Arg: "{expression}", Type: "string", Desc: "DSL expression to evaluate. You can use this after -f to add an expression. Example use case: define functions/subroutines in a file you specify with -f, then call them with an expression you specify with -e.", Repeatable: true},
	{Flag: "-s", Arg: "{name=value}", Type: "string", Desc: "Predefines out-of-stream variable @name to have the given value. Thus mlr put -s foo=97 '$column += @foo' is like mlr put 'begin {@foo = 97} $column += @foo'. The value part is subject to type-inferencing. May be specified more than once, e.g. -s name1=value1 -s name2=value2. Note: the value may be an environment variable, e.g. -s sequence=$SEQUENCE.", Repeatable: true},
	{Flag: "-x", Type: "bool", Desc: "Prints records for which {expression} evaluates to false, not true, i.e. invert the sense of the filter expression. Default false."},
	{Flag: "-q", Type: "bool", Desc: "Does not include the modified record in the output stream. Useful for when all desired output is in begin and/or end blocks."},
	{Flag: "-S", Type: "bool", Desc: "No-op in Miller 6 and above, since type-inferencing is now done by the record-readers before filter/put is executed. Supported as a no-op pass-through flag for backward compatibility."},
	{Flag: "-F", Type: "bool", Desc: "No-op in Miller 6 and above, since type-inferencing is now done by the record-readers before filter/put is executed. Supported as a no-op pass-through flag for backward compatibility."},
	{Flag: "-w", Type: "bool", Desc: "Print warnings about things like uninitialized variables."},
	{Flag: "-W", Type: "bool", Desc: "Same as -w, but exit the process if there are any warnings."},
	{Flag: "-p", Type: "bool", Desc: "Prints the expression's AST (abstract syntax tree), which gives full transparency on the precedence and associativity rules of Miller's grammar, to stdout."},
	{Flag: "-d", Type: "bool", Desc: "Like -p but uses a parenthesized-expression format for the AST."},
	{Flag: "-D", Type: "bool", Desc: "Like -d but with output all on one line."},
	{Flag: "-E", Type: "bool", Desc: "Echo DSL expression before printing parse-tree."},
	{Flag: "-v", Type: "bool", Desc: "Same as -E -p."},
	{Flag: "-X", Type: "bool", Desc: "Exit after parsing but before stream-processing. Useful with -v/-d/-D, if you only want to look at parser information."},
	{Flag: "--explain", Type: "bool", Desc: "Parse and type-check the DSL expression, report whether it is valid, and exit without reading the input stream. Exit status is 0 if the expression is valid and non-zero otherwise; combine with --errors-json for a machine-readable error."},
}

var PutSetup = TransformerSetup{
	Verb:         verbNamePut,
	UsageFunc:    transformerPutUsage,
	ParseCLIFunc: transformerPutOrFilterParseCLI,
	IgnoresInput: false,
	Options:      putOptions,
}

const verbNameFilter = "filter"

var filterOptions = []OptionSpec{
	{Flag: "-f", Arg: "{file name}", Type: "filename", Desc: "File containing a DSL expression (see examples below). If the filename is a directory, all *.mlr files in that directory are loaded.", Repeatable: true},
	{Flag: "-e", Arg: "{expression}", Type: "string", Desc: "DSL expression to evaluate. You can use this after -f to add an expression. Example use case: define functions/subroutines in a file you specify with -f, then call them with an expression you specify with -e.", Repeatable: true},
	{Flag: "-s", Arg: "{name=value}", Type: "string", Desc: "Predefines out-of-stream variable @name to have the given value. Thus mlr put -s foo=97 '$column += @foo' is like mlr put 'begin {@foo = 97} $column += @foo'. The value part is subject to type-inferencing. May be specified more than once, e.g. -s name1=value1 -s name2=value2. Note: the value may be an environment variable, e.g. -s sequence=$SEQUENCE.", Repeatable: true},
	{Flag: "-x", Type: "bool", Desc: "Prints records for which {expression} evaluates to false, not true, i.e. invert the sense of the filter expression. Default false."},
	{Flag: "-q", Type: "bool", Desc: "Does not include the modified record in the output stream. Useful for when all desired output is in begin and/or end blocks."},
	{Flag: "-S", Type: "bool", Desc: "No-op in Miller 6 and above, since type-inferencing is now done by the record-readers before filter/put is executed. Supported as a no-op pass-through flag for backward compatibility."},
	{Flag: "-F", Type: "bool", Desc: "No-op in Miller 6 and above, since type-inferencing is now done by the record-readers before filter/put is executed. Supported as a no-op pass-through flag for backward compatibility."},
	{Flag: "-w", Type: "bool", Desc: "Print warnings about things like uninitialized variables."},
	{Flag: "-W", Type: "bool", Desc: "Same as -w, but exit the process if there are any warnings."},
	{Flag: "-p", Type: "bool", Desc: "Prints the expression's AST (abstract syntax tree), which gives full transparency on the precedence and associativity rules of Miller's grammar, to stdout."},
	{Flag: "-d", Type: "bool", Desc: "Like -p but uses a parenthesized-expression format for the AST."},
	{Flag: "-D", Type: "bool", Desc: "Like -d but with output all on one line."},
	{Flag: "-E", Type: "bool", Desc: "Echo DSL expression before printing parse-tree."},
	{Flag: "-v", Type: "bool", Desc: "Same as -E -p."},
	{Flag: "-X", Type: "bool", Desc: "Exit after parsing but before stream-processing. Useful with -v/-d/-D, if you only want to look at parser information."},
	{Flag: "--explain", Type: "bool", Desc: "Parse and type-check the DSL expression, report whether it is valid, and exit without reading the input stream. Exit status is 0 if the expression is valid and non-zero otherwise; combine with --errors-json for a machine-readable error."},
}

var FilterSetup = TransformerSetup{
	Verb:         verbNameFilter,
	UsageFunc:    transformerFilterUsage,
	ParseCLIFunc: transformerPutOrFilterParseCLI,
	IgnoresInput: false,
	Options:      filterOptions,
}

func transformerPutUsage(
	o *os.File,
) {
	transformerPutOrFilterUsage(o, "put")
}

func transformerFilterUsage(
	o *os.File,
) {
	transformerPutOrFilterUsage(o, "filter")
}

func transformerPutOrFilterUsage(
	o *os.File,
	verb string,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {DSL expression}\n", "mlr", verb)
	switch verb {
	case "put":
		fmt.Fprintf(o, "Lets you use a domain-specific language to programmatically alter stream records.\n")
	case "filter":
		fmt.Fprintf(o, "Lets you use a domain-specific language to programmatically filter which\n")
		fmt.Fprintf(o, "stream records will be output.\n")
	}
	fmt.Fprintf(o, "See also: https://miller.readthedocs.io/en/latest/reference-verbs\n")
	fmt.Fprintf(o, "\n")
	if verb == "put" {
		WriteVerbOptions(o, putOptions)
	} else {
		WriteVerbOptions(o, filterOptions)
	}
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "If you mix -e and -f then the expressions are evaluated in the order encountered.\n")
	fmt.Fprintf(o, "Since the expression pieces are simply concatenated, please be sure to use intervening\n")
	fmt.Fprintf(o, "semicolons to separate expressions.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Parser-info options are -w, -W, -p, -d, -D, -E, -v, and -X.\n")

	if verb == "put" {
		fmt.Fprintln(o)
		fmt.Fprint(o, `Examples:
  mlr --from example.csv put '$qr = $quantity * $rate'
More example put expressions:
  If-statements:
    'if ($flag == true) { $quantity *= 10}'
    'if ($x > 0.0) { $y=log10($x); $z=sqrt($y) } else {$y = 0.0; $z = 0.0}'
  Newly created fields can be read after being written:
    '$new_field = $index**2; $qn = $quantity * $new_field'
  Regex-replacement:
    '$name = sub($name, "http.*com"i, "")'
  Regex-capture:
	'if ($a =~ "([a-z]+)_([0-9]+)") { $b = "left_\1"; $c = "right_\2" }'
  Built-in variables:
    '$filename = FILENAME'
  Aggregations (use mlr put -q):
    '@sum += $x; end {emit @sum}'
    '@sum[$shape] += $quantity; end {emit @sum, "shape"}'
    '@sum[$shape][$color] += $x; end {emit @sum, "shape", "color"}'
    '
      @min = min(@min,$x);
      @max=max(@max,$x);
      end{emitf @min, @max}
    '
`)
	}

	if verb == "filter" {
		fmt.Fprintln(o)
		fmt.Fprintf(o, `Records will pass the filter depending on the last bare-boolean statement in
the DSL expression. That can be the result of <, ==, >, etc., the return value of a function call
which returns boolean, etc.
`)
		fmt.Fprintln(o)
		fmt.Fprint(o, `Examples:
  mlr --csv --from example.csv filter '$color == "red"'
  mlr --csv --from example.csv filter '$color == "red" && flag == true'
More example filter expressions:
  First record in each file:
    'FNR == 1'
  Subsampling:
    'urand() < 0.001'
  Compound booleans:
    '$color != "blue" && $value > 4.2'
    '($x < 0.5 && $y < 0.5) || ($x > 0.5 && $y > 0.5)'
  Regexes with case-insensitive flag
    '($name =~ "^sys.*east$") || ($name =~ "^dev.[0-9]+"i)'
  Assignments, then bare-boolean filter statement:
    '$ab = $a+$b; $cd = $c+$d; $ab != $cd'
  Bare-boolean filter statement within a conditional:
    'if (NR < 100) {
      $x > 0.3;
    } else {
      $x > 0.002;
    }
    '
  Using 'any' higher-order function to see if $index is 10, 20, or 30:
    'any([10,20,30], func(e) {return $index == e})'
`)

	}

	fmt.Fprintln(o)
	fmt.Fprintf(o, "See also %s/reference-dsl for more context.\n", lib.DOC_URL)
}

func transformerPutOrFilterParseCLI(
	pargi *int,
	argc int,
	args []string,
	mainOptions *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	dslStrings := []string{}
	haveDSLStringsHere := false
	echoDSLString := false
	printASTAsTree := false
	printASTMultiLine := false
	printASTSingleLine := false
	exitAfterParse := false
	doExplain := false
	doWarnings := false
	warningsAreFatal := false
	strictMode := false
	invertFilter := false
	suppressOutputRecord := false
	presets := []string{}

	// TODO: make sure this is a full nested-struct copy.
	var options *cli.TOptions = nil
	if mainOptions != nil {
		copyThereof := *mainOptions // struct copy
		options = &copyThereof
	}

	// If there was a global --load/--mload, load those DSL strings here (e.g.
	// someone's local function library).
	for _, filename := range options.DSLPreloadFileNames {
		theseDSLStrings, err := lib.LoadStringsFromFileOrDir(filename, ".mlr")
		if err != nil {
			return nil, fmt.Errorf("%s %s: cannot load DSL expression from \"%s\": %w",
				"mlr", verb, filename, err)
		}
		dslStrings = append(dslStrings, theseDSLStrings...)
	}

	// Parse local flags.
	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		if args[argi] == "--" {
			break // All transformers must do this so main-flags can follow verb-flags
		}
		argi++

		switch opt {
		case "-h", "--help":
			transformerPutOrFilterUsage(os.Stdout, verb)
			return nil, cli.ErrHelpRequested

		case "-f":
			// Get a DSL string from the user-specified filename
			filename, err := cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

			// Miller has a two-pass command-line parser. If the user does
			//   `mlr put -f foo.mlr`
			// then that file can be parsed both times. But if the user does
			//   `mlr put -f <( echo 'some expression goes here' )`
			// that will read stdin. (The filename will come in as "dev/fd/63" or what have you.)
			// But this file _cannot_ be read twice. So, if doConstruct==false -- we're
			// on the first pass of the command-line parser -- don't bother to parse
			// the DSL-contents file.
			//
			// See also https://github.com/johnkerl/miller/issues/1515

			if doConstruct {
				theseDSLStrings, err := lib.LoadStringsFromFileOrDir(filename, ".mlr")
				if err != nil {
					return nil, fmt.Errorf("%s %s: cannot load DSL expression from file \"%s\": %w",
						"mlr", verb, filename, err)
				}
				dslStrings = append(dslStrings, theseDSLStrings...)
			}
			haveDSLStringsHere = true

		case "-e":
			dslString, err := cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			dslStrings = append(dslStrings, dslString)
			haveDSLStringsHere = true

		case "-s":
			// E.g.
			//   mlr put -s sum=0
			// is like
			//   mlr put -s 'begin {@sum = 0}'
			preset, err := cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			presets = append(presets, preset)

		case "-x":
			invertFilter = true
		case "-q":
			suppressOutputRecord = true

		case "-E":
			echoDSLString = true
		case "-p":
			printASTAsTree = true
		case "-v":
			echoDSLString = true
			printASTAsTree = true
		case "-d":
			printASTMultiLine = true
		case "-D":
			printASTSingleLine = true
		case "-X":
			exitAfterParse = true
		case "--explain":
			doExplain = true
		case "-w":
			doWarnings = true
			warningsAreFatal = false
		case "-z":
			// TODO: perhaps doWarnings and warningsAreFatal as well.
			// But first I want to see what can be caught at runtime
			// without static analysis.
			strictMode = true
		case "-W":
			doWarnings = true
			warningsAreFatal = true

		case "-S":
			// TODO: this is a no-op in Miller 6 and above.
			// Comment this in more detail.

		case "-F":
			// TODO: this is a no-op in Miller 6 and above.
			// Comment this in more detail.

		default:
			// This is inelegant. For error-proofing we advance argi already in our
			// loop (so individual if-statements don't need to). However,
			// ParseWriterOptions expects it unadvanced.
			largi := argi - 1
			if cli.FLAG_TABLE.Parse(args, argc, &largi, options) {
				// This lets mlr main and mlr put have different output formats.
				// Nothing else to handle here.
				argi = largi
			} else {
				return nil, cli.VerbErrorf(verb, "option not recognized")
			}
		}
	}

	if err := cli.FinalizeWriterOptions(&options.WriterOptions); err != nil {
		return nil, cli.VerbErrorf(verb, "%v", err)
	}

	// If they've used either of 'mlr put -f {filename}' or 'mlr put -e
	// {expression}' then that specifies their DSL expression. But if they've
	// done neither then we expect 'mlr put {expression}'.
	if !haveDSLStringsHere {
		// Get the DSL string from the command line, after the flags
		if argi >= argc {
			return nil, cli.VerbErrorf(verb, "expression or -f/-e is required")
		}
		dslString := args[argi]
		dslStrings = append(dslStrings, dslString)
		argi++
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	var dslInstanceType cst.DSLInstanceType = cst.DSLInstanceTypePut
	if verb == "filter" {
		dslInstanceType = cst.DSLInstanceTypeFilter
	}

	doFilter := (verb == "filter")

	transformer, err := NewTransformerPut(
		doFilter,
		dslStrings,
		dslInstanceType,
		presets,
		echoDSLString,
		printASTAsTree,
		printASTMultiLine,
		printASTSingleLine,
		exitAfterParse,
		doExplain,
		doWarnings,
		warningsAreFatal,
		strictMode,
		invertFilter,
		suppressOutputRecord,
		options,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerPut struct {
	doFilter             bool // false for the put verb, true for the filter verb
	cstRootNode          *cst.RootNode
	runtimeState         *runtime.State
	callCount            int
	invertFilter         bool
	suppressOutputRecord bool
	executedBeginBlocks  bool
}

func NewTransformerPut(
	doFilter bool, // false for the put verb, true for the filter verb
	dslStrings []string,
	dslInstanceType cst.DSLInstanceType,
	presets []string,
	echoDSLString bool,
	printASTAsTree bool,
	printASTMultiLine bool,
	printASTSingleLine bool,
	exitAfterParse bool,
	doExplain bool,
	doWarnings bool,
	warningsAreFatal bool,
	strictMode bool,
	invertFilter bool,
	suppressOutputRecord bool,
	options *cli.TOptions,
) (*TransformerPut, error) {

	cstRootNode := cst.NewEmptyRoot(&options.WriterOptions, dslInstanceType).WithStrictMode(strictMode)

	hadWarnings, err := cstRootNode.Build(
		dslStrings,
		dslInstanceType,
		false, // isReplImmediate
		doWarnings,

		func(dslString string, astNode *asts.AST) {

			if echoDSLString {
				fmt.Println("DSL EXPRESSION:")
				fmt.Println(dslString)
				fmt.Println()
			}
			if printASTAsTree {
				fmt.Println("AST:")
				astNode.Print()
				fmt.Println()
			}
			if printASTMultiLine {
				astNode.PrintParex()
				fmt.Println()
			}
			if printASTSingleLine {
				astNode.PrintParexOneLine()
				fmt.Println()
			}

		},
	)

	// --explain is a validate/dry-run: report whether the DSL parsed and
	// type-checked, then exit without reading the input stream. A parse/build
	// error is returned so it flows through the normal error path (including
	// --errors-json); a valid expression prints a confirmation and exits 0.
	if doExplain {
		if err != nil {
			return nil, err
		}
		verbName := "put"
		if doFilter {
			verbName = "filter"
		}
		if warningsAreFatal && hadWarnings {
			fmt.Fprintf(os.Stderr, "mlr %s: DSL expression has warnings treated as fatal.\n", verbName)
			os.Exit(1)
		}
		fmt.Printf("mlr %s: DSL expression is valid.\n", verbName)
		os.Exit(0)
	}

	if warningsAreFatal && hadWarnings {
		fmt.Printf(
			"%s: Exiting due to warnings treated as fatal.\n",
			"mlr",
		)
		os.Exit(1)
	}

	if exitAfterParse {
		os.Exit(0)
	}

	if err != nil {
		return nil, err
	}

	runtimeState := runtime.NewEmptyState(options, strictMode)

	// E.g.
	//   mlr put -s sum=0
	// is like
	//   mlr put -s 'begin {@sum = 0}'
	if len(presets) > 0 {
		for _, preset := range presets {
			pair := strings.SplitN(preset, "=", 2)
			if len(pair) != 2 {
				return nil, fmt.Errorf(`missing "=" in preset expression "%s"`, preset)
			}
			key := pair[0]
			svalue := pair[1]
			mvalue := mlrval.FromInferredType(svalue)
			runtimeState.Oosvars.PutCopy(key, mvalue)
		}
	}

	return &TransformerPut{
		doFilter:             doFilter,
		cstRootNode:          cstRootNode,
		runtimeState:         runtimeState,
		callCount:            0,
		invertFilter:         invertFilter,
		suppressOutputRecord: suppressOutputRecord,
		executedBeginBlocks:  false,
	}, nil
}

func (tr *TransformerPut) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.runtimeState.OutputRecordsAndContexts = outputRecordsAndContexts

	inrec := inrecAndContext.Record
	context := inrecAndContext.Context
	if !inrecAndContext.EndOfStream {

		// Execute the begin { ... } before the first input record
		tr.callCount++
		if tr.callCount == 1 {
			tr.runtimeState.Update(nil, &context)
			err := tr.cstRootNode.ExecuteBeginBlocks(tr.runtimeState)
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
				os.Exit(1)
			}
			tr.executedBeginBlocks = true
		}

		tr.runtimeState.Update(inrec, &context)

		// Execute the main block on the current input record
		outrec, err := tr.cstRootNode.ExecuteMainBlock(tr.runtimeState)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
			os.Exit(1)
		}

		if !tr.suppressOutputRecord {
			// The tr.runtimeState.FilterExpression defaults to null. It evaluates to null
			// for assignment statements, etc.
			// * If the verb is put, then tr.runtimeState.FilterExpression will get set to
			//   something only when a filter DSL statement is encountered.
			// * If the verb is filter, then we insist that the expression evaluate to either
			//   boolean, or absent. The latter is for Miller's record-heterogeneity guarantees,
			//   e.g. mlr filter '$x > 10' for records not having a $x.

			filterBool, isBool := tr.runtimeState.FilterExpression.GetBoolValue()

			if tr.doFilter {
				// This is mlr filter
				if !isBool {
					if tr.runtimeState.FilterExpression.IsAbsent() {
						filterBool = false
					} else {
						fmt.Fprintf(os.Stderr,
							"Filter expression did not evaluate to boolean: got %s value %s",
							tr.runtimeState.FilterExpression.String(),
							tr.runtimeState.FilterExpression.GetTypeName(),
						)
						os.Exit(1)
					}
				}
			} else {
				// This is mlr put.
				if !isBool {
					filterBool = true
				}
			}

			wantToEmit := lib.BooleanXOR(filterBool, tr.invertFilter)
			if wantToEmit {
				*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(outrec, &context))
			}
		}

	} else {
		tr.runtimeState.Update(nil, &context)

		// If there were no input records then we never executed the
		// begin-blocks. Do so now.
		if !tr.executedBeginBlocks {
			err := tr.cstRootNode.ExecuteBeginBlocks(tr.runtimeState)
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
				os.Exit(1)
			}
		}

		// Execute the end { ... } after the last input record
		err := tr.cstRootNode.ExecuteEndBlocks(tr.runtimeState)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
			os.Exit(1)
		}

		// Send all registered OutputHandlerManager instances the end-of-stream
		// indicator.
		tr.cstRootNode.ProcessEndOfStream()

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewEndOfStreamMarker(&context))
	}
}
