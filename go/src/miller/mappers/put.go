package mappers

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"miller/clitypes"
	"miller/dsl"
	"miller/dsl/cst"
	"miller/lib"
	"miller/mapping"
	"miller/parsing/lexer"
	"miller/parsing/parser"
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
	pVerbose := flagSet.Bool(
		"v",
		false,
		`Prints the expressions's AST (abstract syntax tree), which gives
    full transparency on the precedence and associativity rules of
    Miller's grammar, to stdout.`,
	)
	pExpressionFileName := flagSet.String(
		"f",
		"",
		`File containing a DSL expression.`,
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

	mapper, err := NewMapperPut(dslString, *pVerbose)
	if err != nil {
		// xxx make sure better parse-fail info is printed by the DSL parser
		fmt.Fprintf(os.Stderr, "%s %s: cannot parse DSL expression.\n",
			args[0], verb)
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
	astRoot *dsl.AST
	cstRoot *cst.Root
}

func NewMapperPut(
	dslString string,
	verbose bool,
) (*MapperPut, error) {
	astRoot, err := NewASTFromString(dslString)
	if err != nil {
		return nil, err
	}
	if verbose {
		fmt.Println("DSL EXPRESSION:")
		fmt.Println(dslString)
		fmt.Println("RAW AST:")
		astRoot.Print()
		fmt.Println()
	}
	cstRoot, err := cst.Build(astRoot)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	return &MapperPut{
		astRoot: astRoot,
		cstRoot: cstRoot,
	}, nil
}

// xxx note (package cycle) why not a dsl.AST constructor :(
// xxx maybe split out dsl into two package ... and/or put the astRoot.go into miller/parsing -- ?
//   depends on TBD split-out of AST and CST ...
func NewASTFromString(dslString string) (*dsl.AST, error) {
	theLexer := lexer.NewLexer([]byte(dslString))
	theParser := parser.NewParser()
	interfaceAST, err := theParser.Parse(theLexer)
	if err != nil {
		return nil, err
	}
	astRoot := interfaceAST.(*dsl.AST)
	return astRoot, nil
}

func (this *MapperPut) Map(
	inrecAndContext *lib.RecordAndContext,
	outputChannel chan<- *lib.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	context := inrecAndContext.Context
	if inrec != nil {

		cstState := cst.NewState(inrec, &context)
		outrec, err := this.cstRoot.Execute(cstState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		outputChannel <- lib.NewRecordAndContext(
			outrec,
			&context,
		)
	} else {
		outputChannel <- lib.NewRecordAndContext(
			nil, // signals end of input record stream
			&context,
		)

	}
}
