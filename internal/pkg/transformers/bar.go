package transformers

import (
	"bytes"
	"container/list"
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/types"
)

const barDefaultFillString = "*"
const barDefaultOOBString = "#"
const barDefaultBlankString = "."
const barDefaultLo = 0.0
const barDefaultHi = 100.0
const barDefaultWidth = 40

// ----------------------------------------------------------------
const verbNameBar = "bar"

var BarSetup = TransformerSetup{
	Verb:         verbNameBar,
	UsageFunc:    transformerBarUsage,
	ParseCLIFunc: transformerBarParseCLI,
	IgnoresInput: false,
}

func transformerBarUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameBar)
	fmt.Fprintf(o, "Replaces a numeric field with a number of asterisks, allowing for cheesy\n")
	fmt.Fprintf(o, "bar plots. These align best with --opprint or --oxtab output format.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f   {a,b,c}      Field names to convert to bars.\n")
	fmt.Fprintf(o, "--lo {lo}         Lower-limit value for min-width bar: default '%f'.\n", barDefaultLo)
	fmt.Fprintf(o, "--hi {hi}         Upper-limit value for max-width bar: default '%f'.\n", barDefaultHi)
	fmt.Fprintf(o, "-w   {n}          Bar-field width: default '%d'.\n", barDefaultWidth)
	fmt.Fprintf(o, "--auto            Automatically computes limits, ignoring --lo and --hi.\n")
	fmt.Fprintf(o, "                  Holds all records in memory before producing any output.\n")
	fmt.Fprintf(o, "-c   {character}  Fill character: default '%s'.\n", barDefaultFillString)
	fmt.Fprintf(o, "-x   {character}  Out-of-bounds character: default '%s'.\n", barDefaultOOBString)
	fmt.Fprintf(o, "-b   {character}  Blank character: default '%s'.\n", barDefaultBlankString)
	fmt.Fprintf(o, "Nominally the fill, out-of-bounds, and blank characters will be strings of length 1.\n")
	fmt.Fprintf(o, "However you can make them all longer if you so desire.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerBarParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	var fieldNames []string = nil
	lo := barDefaultLo
	hi := barDefaultHi
	width := barDefaultWidth
	doAuto := false
	fillString := barDefaultFillString
	oobString := barDefaultOOBString
	blankString := barDefaultBlankString

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerBarUsage(os.Stdout, true, 0)

		} else if opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--lo" {
			lo = cli.VerbGetFloatArgOrDie(verb, opt, args, &argi, argc)
		} else if opt == "-w" {
			width = cli.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)
		} else if opt == "--hi" {
			hi = cli.VerbGetFloatArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-c" {
			fillString = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
		} else if opt == "-x" {
			oobString = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
		} else if opt == "-b" {
			blankString = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--auto" {
			doAuto = true

		} else {
			transformerBarUsage(os.Stderr, true, 1)
		}
	}

	if fieldNames == nil {
		transformerBarUsage(os.Stderr, true, 1)
	}

	transformer, err := NewTransformerBar(
		fieldNames,
		lo,
		hi,
		width,
		doAuto,
		fillString,
		oobString,
		blankString,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerBar struct {
	fieldNames         []string
	lo                 float64
	hi                 float64
	width              int
	fillString         string
	oobString          string
	blankString        string
	bars               []string
	recordsForAutoMode *list.List

	recordTransformerFunc RecordTransformerFunc
}

// ----------------------------------------------------------------
func NewTransformerBar(
	fieldNames []string,
	lo float64,
	hi float64,
	width int,
	doAuto bool,
	fillString string,
	oobString string,
	blankString string,
) (*TransformerBar, error) {

	tr := &TransformerBar{
		fieldNames:  fieldNames,
		lo:          lo,
		hi:          hi,
		width:       width,
		fillString:  fillString,
		oobString:   oobString,
		blankString: blankString,
	}

	tr.bars = make([]string, width+1)
	for i := 0; i <= tr.width; i++ {
		var bar = ""
		if i == 0 {
			bar = tr.oobString + strings.Repeat(tr.blankString, width-1)
		} else if i < width {
			bar = strings.Repeat(tr.fillString, i) + strings.Repeat(tr.blankString, width-i)
		} else {
			bar = strings.Repeat(tr.fillString, width-1) + tr.oobString
		}

		tr.bars[i] = bar
	}

	if doAuto {
		tr.recordTransformerFunc = tr.processAuto
		tr.recordsForAutoMode = list.New()
	} else {
		tr.recordTransformerFunc = tr.processNoAuto
		tr.recordsForAutoMode = nil
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerBar) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerBar) simpleBar(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	outputChannel <- inrecAndContext
}

// ----------------------------------------------------------------
func (tr *TransformerBar) processNoAuto(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for _, fieldName := range tr.fieldNames {
			mvalue := inrec.Get(fieldName)
			if mvalue == nil {
				continue
			}
			floatValue, ok := mvalue.GetNumericToFloatValue()
			if !ok {
				continue
			}
			idx := int(float64(tr.width) * (floatValue - tr.lo) / (tr.hi - tr.lo))
			if idx < 0 {
				idx = 0
			}
			if idx > tr.width {
				idx = tr.width
			}
			inrec.PutReference(fieldName, types.MlrvalFromString(tr.bars[idx]))
		}

		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // emit end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerBar) processAuto(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		tr.recordsForAutoMode.PushBack(inrecAndContext.Copy())
		return
	}

	// Else, end of stream

	// Loop over field names to be barred
	for _, fieldName := range tr.fieldNames {
		lo := 0.0
		hi := 0.0

		// The first pass computes lo and hi from the data
		onFirst := true
		for e := tr.recordsForAutoMode.Front(); e != nil; e = e.Next() {
			recordAndContexts := e.Value.(*types.RecordAndContext)
			record := recordAndContexts.Record
			mvalue := record.Get(fieldName)
			if mvalue == nil {
				continue
			}
			floatValue, ok := mvalue.GetNumericToFloatValue()
			if !ok {
				continue
			}

			if onFirst || floatValue < lo {
				lo = floatValue
			}
			if onFirst || floatValue > hi {
				hi = floatValue
			}
			onFirst = false
		}

		// The second pass applies the bars. There is some redundant computation
		// which could be hoisted out of the loop for performance ... but this
		// verb computes data solely for visual inspection and I take the
		// nominal use case to be tens or hundreds of records. So, optimization
		// isn't worth the effort here.

		slo := fmt.Sprintf("%g", lo)
		shi := fmt.Sprintf("%g", hi)

		for e := tr.recordsForAutoMode.Front(); e != nil; e = e.Next() {
			recordAndContext := e.Value.(*types.RecordAndContext)
			record := recordAndContext.Record
			mvalue := record.Get(fieldName)
			if mvalue == nil {
				continue
			}
			floatValue, ok := mvalue.GetNumericToFloatValue()
			if !ok {
				continue
			}

			idx := int((float64(tr.width) * (floatValue - lo) / (hi - lo)))
			if idx < 0 {
				idx = 0
			}
			if idx > tr.width {
				idx = tr.width
			}

			var buffer bytes.Buffer // faster than fmt.Print() separately
			buffer.WriteString("[")
			buffer.WriteString(slo)
			buffer.WriteString("]")
			buffer.WriteString(tr.bars[idx])
			buffer.WriteString("[")
			buffer.WriteString(shi)
			buffer.WriteString("]")
			record.PutReference(fieldName, types.MlrvalFromString(buffer.String()))
		}
	}

	for e := tr.recordsForAutoMode.Front(); e != nil; e = e.Next() {
		recordAndContext := e.Value.(*types.RecordAndContext)
		outputChannel <- recordAndContext
	}

	outputChannel <- inrecAndContext // Emit the end-of-stream marker
}
