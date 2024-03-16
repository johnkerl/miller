package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/dsl"
	"github.com/johnkerl/miller/pkg/dsl/cst"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/runtime"
	"github.com/johnkerl/miller/pkg/types"
)

// ----------------------------------------------------------------
const verbNamePut = "put"

var PutSetup = TransformerSetup{
	Verb:         verbNamePut,
	UsageFunc:    transformerPutUsage,
	ParseCLIFunc: transformerPutOrFilterParseCLI,
	IgnoresInput: false,
}

const verbNameFilter = "filter"

var FilterSetup = TransformerSetup{
	Verb:         verbNameFilter,
	UsageFunc:    transformerFilterUsage,
	ParseCLIFunc: transformerPutOrFilterParseCLI,
	IgnoresInput: false,
}

// ----------------------------------------------------------------
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
	if verb == "put" {
		fmt.Fprintf(o, "Lets you use a domain-specific language to progamatically alter stream records.\n")
	} else if verb == "filter" {
		fmt.Fprintf(o, "Lets you use a domain-specific language to progamatically filter which\n")
		fmt.Fprintf(o, "stream records will be output.\n")
	}
	fmt.Fprintf(o, "See also: https://miller.readthedocs.io/en/latest/reference-verbs\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o,
		`-f {file name} File containing a DSL expression (see examples below). If the filename
   is a directory, all *.mlr files in that directory are loaded.

-e {expression} You can use this after -f to add an expression. Example use
   case: define functions/subroutines in a file you specify with -f, then call
   them with an expression you specify with -e.

(If you mix -e and -f then the expressions are evaluated in the order encountered.
Since the expression pieces are simply concatenated, please be sure to use intervening
semicolons to separate expressions.)

-s name=value: Predefines out-of-stream variable @name to have 
    Thus mlr put -s foo=97 '$column += @foo' is like
    mlr put 'begin {@foo = 97} $column += @foo'.
    The value part is subject to type-inferencing.
    May be specified more than once, e.g. -s name1=value1 -s name2=value2.
    Note: the value may be an environment variable, e.g. -s sequence=$SEQUENCE

-x (default false) Prints records for which {expression} evaluates to false, not true,
   i.e. invert the sense of the filter expression.

-q Does not include the modified record in the output stream.
   Useful for when all desired output is in begin and/or end blocks.

-S and -F: There are no-ops in Miller 6 and above, since now type-inferencing is done
   by the record-readers before filter/put is executed. Supported as no-op pass-through
   flags for backward compatibility.

-h|--help Show this message.

Parser-info options:

-w Print warnings about things like uninitialized variables.

-W Same as -w, but exit the process if there are any warnings.

-p Prints the expressions's AST (abstract syntax tree), which gives full
  transparency on the precedence and associativity rules of Miller's grammar,
  to stdout.

-d Like -p but uses a parenthesized-expression format for the AST.

-D Like -d but with output all on one line.

-E Echo DSL expression before printing parse-tree

-v Same as -E -p.

-X Exit after parsing but before stream-processing. Useful with -v/-d/-D, if you
   only want to look at parser information.
`)

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

// ----------------------------------------------------------------
func transformerPutOrFilterParseCLI(
	pargi *int,
	argc int,
	args []string,
	mainOptions *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	var dslStrings []string = make([]string, 0)
	haveDSLStringsHere := false
	echoDSLString := false
	printASTAsTree := false
	printASTMultiLine := false
	printASTSingleLine := false
	exitAfterParse := false
	doWarnings := false
	warningsAreFatal := false
	strictMode := false
	invertFilter := false
	suppressOutputRecord := false
	presets := make([]string, 0)

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
			fmt.Fprintf(os.Stderr, "%s %s: cannot load DSL expression from \"%s\": ",
				"mlr", verb, filename)
			fmt.Println(err)
			os.Exit(1)
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

		if opt == "-h" || opt == "--help" {
			transformerPutOrFilterUsage(os.Stdout, verb)
			os.Exit(0)

		} else if opt == "-f" {
			// Get a DSL string from the user-specified filename
			filename := cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			theseDSLStrings, err := lib.LoadStringsFromFileOrDir(filename, ".mlr")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s %s: cannot load DSL expression from file \"%s\": ",
					"mlr", verb, filename)
				fmt.Println(err)
				os.Exit(1)
			}
			dslStrings = append(dslStrings, theseDSLStrings...)
			haveDSLStringsHere = true

		} else if opt == "-e" {
			dslString := cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			dslStrings = append(dslStrings, dslString)
			haveDSLStringsHere = true

		} else if opt == "-s" {
			// E.g.
			//   mlr put -s sum=0
			// is like
			//   mlr put -s 'begin {@sum = 0}'
			preset := cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			presets = append(presets, preset)

		} else if opt == "-x" {
			invertFilter = true
		} else if opt == "-q" {
			suppressOutputRecord = true

		} else if opt == "-E" {
			echoDSLString = true
		} else if opt == "-p" {
			printASTAsTree = true
		} else if opt == "-v" {
			echoDSLString = true
			printASTAsTree = true
		} else if opt == "-d" {
			printASTMultiLine = true
		} else if opt == "-D" {
			printASTSingleLine = true
		} else if opt == "-X" {
			exitAfterParse = true
		} else if opt == "-w" {
			doWarnings = true
			warningsAreFatal = false
		} else if opt == "-z" {
			// TODO: perhaps doWarnings and warningsAreFatal as well.
			// But first I want to see what can be caught at runtime
			// without static analysis.
			strictMode = true
		} else if opt == "-W" {
			doWarnings = true
			warningsAreFatal = true

		} else if opt == "-S" {
			// TODO: this is a no-op in Miller 6 and above.
			// Comment this in more detail.

		} else if opt == "-F" {
			// TODO: this is a no-op in Miller 6 and above.
			// Comment this in more detail.

		} else {
			// This is inelegant. For error-proofing we advance argi already in our
			// loop (so individual if-statements don't need to). However,
			// ParseWriterOptions expects it unadvanced.
			largi := argi - 1
			if cli.FLAG_TABLE.Parse(args, argc, &largi, options) {
				// This lets mlr main and mlr put have different output formats.
				// Nothing else to handle here.
				argi = largi
			} else {
				transformerPutOrFilterUsage(os.Stderr, verb)
				os.Exit(1)
			}
		}
	}

	cli.FinalizeWriterOptions(&options.WriterOptions)

	// If they've used either of 'mlr put -f {filename}' or 'mlr put -e
	// {expression}' then that specifies their DSL expression. But if they've
	// done neither then we expect 'mlr put {expression}'.
	if !haveDSLStringsHere {
		// Get the DSL string from the command line, after the flags
		if argi >= argc {
			fmt.Fprintf(os.Stderr, "mlr %s: -f/-e requires a filename as argument.\n", verb)
			os.Exit(1)
		}
		dslString := args[argi]
		dslStrings = append(dslStrings, dslString)
		argi++
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	var dslInstanceType cst.DSLInstanceType = cst.DSLInstanceTypePut
	if verb == "filter" {
		dslInstanceType = cst.DSLInstanceTypeFilter
	}

	transformer, err := NewTransformerPut(
		dslStrings,
		dslInstanceType,
		presets,
		echoDSLString,
		printASTAsTree,
		printASTMultiLine,
		printASTSingleLine,
		exitAfterParse,
		doWarnings,
		warningsAreFatal,
		strictMode,
		invertFilter,
		suppressOutputRecord,
		options,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerPut struct {
	cstRootNode          *cst.RootNode
	runtimeState         *runtime.State
	callCount            int
	invertFilter         bool
	suppressOutputRecord bool
	executedBeginBlocks  bool
}

func NewTransformerPut(
	dslStrings []string,
	dslInstanceType cst.DSLInstanceType,
	presets []string,
	echoDSLString bool,
	printASTAsTree bool,
	printASTMultiLine bool,
	printASTSingleLine bool,
	exitAfterParse bool,
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

		func(dslString string, astNode *dsl.AST) {

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
				return nil, fmt.Errorf("missing \"=\" in preset expression \"%s\".", preset)
			}
			key := pair[0]
			svalue := pair[1]
			mvalue := mlrval.FromInferredType(svalue)
			runtimeState.Oosvars.PutCopy(key, mvalue)
		}
	}

	return &TransformerPut{
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
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
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
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			tr.executedBeginBlocks = true
		}

		tr.runtimeState.Update(inrec, &context)

		// Execute the main block on the current input record
		outrec, err := tr.cstRootNode.ExecuteMainBlock(tr.runtimeState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if !tr.suppressOutputRecord {
			filterBool, isBool := tr.runtimeState.FilterExpression.GetBoolValue()
			if !isBool {
				filterBool = false
			}
			wantToEmit := lib.BooleanXOR(filterBool, tr.invertFilter)
			if wantToEmit {
				outputRecordsAndContexts.PushBack(types.NewRecordAndContext(outrec, &context))
			}
		}

	} else {
		tr.runtimeState.Update(nil, &context)

		// If there were no input records then we never executed the
		// begin-blocks. Do so now.
		if tr.executedBeginBlocks == false {
			err := tr.cstRootNode.ExecuteBeginBlocks(tr.runtimeState)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		// Execute the end { ... } after the last input record
		err := tr.cstRootNode.ExecuteEndBlocks(tr.runtimeState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Send all registered OutputHandlerManager instances the end-of-stream
		// indicator.
		tr.cstRootNode.ProcessEndOfStream()

		outputRecordsAndContexts.PushBack(types.NewEndOfStreamMarker(&context))
	}
}
