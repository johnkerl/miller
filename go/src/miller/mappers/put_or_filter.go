package mappers

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"miller/clitypes"
	"miller/dsl"
	"miller/dsl/cst"
	"miller/mapping"
	"miller/parsing/lexer"
	"miller/parsing/parser"
	"miller/types"
)

// ----------------------------------------------------------------
var PutSetup = mapping.MapperSetup{
	Verb:         "put",
	ParseCLIFunc: mapperPutParseCLI,
	IgnoresInput: false,
}

// TODO:
// * rename this file to put_or_filter.go
// * check other things from the C impl
var FilterSetup = mapping.MapperSetup{
	Verb:         "filter",
	ParseCLIFunc: mapperPutParseCLI,
	IgnoresInput: false,
}

func mapperPutParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) mapping.IRecordMapper {

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	dslString := ""
	verbose := false
	suppressOutputRecord := false
	needExpressionArg := true

	// Parse local flags.
	//
	// Unlike other mappers, we can't use flagSet here. The syntax of 'mlr put'
	// and 'mlr filter' is they need to be able to take -f and/or -e more than
	// once, and Go flags can't handle that.

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			mapperPutUsage(os.Stdout, 0, errorHandling, args[0], verb)
			return nil // help intentionally requested

		} else if args[argi] == "-f" {
			checkArgCountPut(args, argi, argc, 2)

			// Get a DSL string from the user-specified filename
			data, err := ioutil.ReadFile(args[argi+1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s %s: cannot load DSL expression: ", args[0], verb)
				fmt.Println(err)
				return nil
			}
			if dslString != "" {
				dslString += ";\n"
			}
			dslString += string(data)

			needExpressionArg = false
			argi += 2

		} else if args[argi] == "-e" {
			checkArgCountPut(args, argi, argc, 2)

			// Get a DSL string from the user-specified filename
			if dslString != "" {
				dslString += ";\n"
			}
			dslString += args[argi+1]

			needExpressionArg = false
			argi += 2

		} else if args[argi] == "-v" {
			verbose = true
			argi++
		} else if args[argi] == "-q" {
			suppressOutputRecord = true
			argi++
		} else {
			mapperPutUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
			os.Exit(1)
		}
	}

	// If they've used either of 'mlr put -f {filename}' or 'mlr put -e
	// {expression}' then that specifies their DSL expression. But if they've
	// done neither then we expect 'mlr put {expression}'.
	if needExpressionArg {
		// Get the DSL string from the command line, after the flags
		if argi >= argc {
			mapperPutUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
			os.Exit(1)
		}
		dslString = args[argi]
		argi += 1
	}

	mapper, err := NewMapperPut(dslString, verbose, suppressOutputRecord)
	if err != nil {
		// Error message already printed out
		os.Exit(1)
	}

	*pargi = argi
	return mapper
}

// For flags with values, e.g. ["-n" "10"], while we're looking at the "-n"
// this let us see if the "10" slot exists.
func checkArgCountPut(args []string, argi int, argc int, n int) {
	if (argc - argi) < n {
		fmt.Fprintf(os.Stderr, "%s: option \"%s\" missing argument(s).\n", args[0], args[argi])
		mapperPutUsage(os.Stderr, 1, flag.ExitOnError, os.Args[0], "sort")
		os.Exit(1)
	}
}

func mapperPutUsage(
	o *os.File,
	exitCode int,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	argv0 string,
	verb string,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {DSL expression}\n", argv0, verb)
	fmt.Fprintf(o, "TODO: put detailed on-line help here.\n")
	fmt.Fprintf(o,
		` -f {file name} File containing a DSL expression.
 -q (default false) Does not include the modified record in the output stream.
    Useful for when all desired output is in begin and/or end blocks.
 -v (default false) Prints the expressions's AST (abstract syntax tree), which gives
    full transparency on the precedence and associativity rules of
    Miller's grammar, to stdout.
`)
}

// ----------------------------------------------------------------
type MapperPut struct {
	astRootNode          *dsl.AST
	cstRootNode          *cst.RootNode
	cstState             *cst.State
	callCount            int64
	suppressOutputRecord bool
	executedBeginBlocks  bool
}

func NewMapperPut(
	dslString string,
	verbose bool,
	suppressOutputRecord bool,
) (*MapperPut, error) {
	astRootNode, err := BuildASTFromString(dslString)
	if err != nil {
		// Leave this out until we get better control over the error-messaging.
		// At present it's overly parser-internal, and confusing. :(
		// fmt.Fprintln(os.Stderr, err)
		fmt.Fprintf(os.Stderr, "%s: cannot parse DSL expression.\n",
			os.Args[0])
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	if verbose {
		fmt.Println("DSL EXPRESSION:")
		fmt.Println(dslString)
		fmt.Println("RAW AST:")
		astRootNode.Print()
		fmt.Println()
	}
	cstRootNode, err := cst.Build(astRootNode)
	cstState := cst.NewEmptyState()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	return &MapperPut{
		astRootNode:          astRootNode,
		cstRootNode:          cstRootNode,
		cstState:             cstState,
		callCount:            0,
		suppressOutputRecord: suppressOutputRecord,
		executedBeginBlocks:  false,
	}, nil
}

// xxx note (package cycle) why not a dsl.AST constructor :(
// xxx maybe split out dsl into two packages ... and/or put the ast.go into miller/parsing -- ?
//   depends on TBD split-out of AST and CST ...
func BuildASTFromString(dslString string) (*dsl.AST, error) {

	theLexer := lexer.NewLexer([]byte(dslString))
	theParser := parser.NewParser()
	interfaceAST, err := theParser.Parse(theLexer)
	if err != nil {
		return nil, err
	}
	astRootNode := interfaceAST.(*dsl.AST)
	return astRootNode, nil
}

func (this *MapperPut) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.cstState.OutputChannel = outputChannel

	inrec := inrecAndContext.Record
	context := inrecAndContext.Context
	if inrec != nil {

		// Execute the begin { ... } before the first input record
		this.callCount++
		if this.callCount == 1 {
			err := this.cstRootNode.ExecuteBeginBlocks(this.cstState)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			this.executedBeginBlocks = true
		}

		this.cstState.Update(inrec, &context)

		// Execute the main block on the current input record
		outrec, err := this.cstRootNode.ExecuteMainBlock(this.cstState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if !this.suppressOutputRecord {
			if this.cstState.FilterResult == true {
				outputChannel <- types.NewRecordAndContext(
					outrec,
					&context,
				)
			}
		}
	} else {
		this.cstState.Update(nil, &context)

		// If there were no input records then we never executed the
		// begin-blocks. Do so now.
		if this.executedBeginBlocks == false {
			err := this.cstRootNode.ExecuteBeginBlocks(this.cstState)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		// Execute the end { ... } after the last input record
		err := this.cstRootNode.ExecuteEndBlocks(this.cstState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		outputChannel <- types.NewRecordAndContext(
			nil, // signals end of input record stream
			&context,
		)

	}
}
