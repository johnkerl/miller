package mappers

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

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

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)
	pExpressionFileName := flagSet.String(
		"f",
		"",
		`File containing a DSL expression.`,
	)
	pVerbose := flagSet.Bool(
		"v",
		false,
		`Prints the expressions's AST (abstract syntax tree), which gives
    full transparency on the precedence and associativity rules of
    Miller's grammar, to stdout.`,
	)
	pSuppressOutputRecord := flagSet.Bool(
		"q",
		false,
		`Does not include the modified record in the output stream.
    Useful for when all desired output is in begin and/or end blocks.`,
	)
	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		mapperPutUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentioally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	dslString := ""
	if *pExpressionFileName != "" {
		// Get the DSL string from the user-specified filename
		data, err := ioutil.ReadFile(*pExpressionFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s: cannot load DSL expression: ", args[0], verb)
			fmt.Println(err)
			return nil
		}
		dslString = string(data)
	} else {
		// Get the DSL string from the command line, after the flags
		if argi >= argc {
			flagSet.Usage()
			os.Exit(1)
		}
		dslString = args[argi]
		argi += 1
	}

	mapper, err := NewMapperPut(dslString, *pVerbose, *pSuppressOutputRecord)
	if err != nil {
		// Error message already printed out
		os.Exit(1)
	}

	*pargi = argi
	return mapper
}

func mapperPutUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {DSL expression}\n", argv0, verb)
	fmt.Fprintf(o, "TODO: put detailed on-line help here.\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type MapperPut struct {
	astRootNode          *dsl.AST
	cstRootNode          *cst.RootNode
	cstState             *cst.State
	callCount            int64
	suppressOutputRecord bool
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
		}

		this.cstState.Update(inrec, &context)

		// Execute the main block on the current input record
		outrec, err := this.cstRootNode.ExecuteMainBlock(this.cstState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if !this.suppressOutputRecord {
			outputChannel <- types.NewRecordAndContext(
				outrec,
				&context,
			)
		}
	} else {

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
