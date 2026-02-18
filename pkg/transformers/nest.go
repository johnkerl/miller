package transformers

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameNest = "nest"

var NestSetup = TransformerSetup{
	Verb:         verbNameNest,
	UsageFunc:    transformerNestUsage,
	ParseCLIFunc: transformerNestParseCLI,
	IgnoresInput: false,
}

func transformerNestUsage(
	o *os.File,
) {
	argv0 := "mlr"
	verb := verbNameNest

	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Explodes specified field values into separate fields/records, or reverses this.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "  --explode,--implode   One is required.\n")
	fmt.Fprintf(o, "  --values,--pairs      One is required.\n")
	fmt.Fprintf(o, "  --across-records,--across-fields One is required.\n")
	fmt.Fprintf(o, "  -f {field name}       Required.\n")
	fmt.Fprintf(o, "  -r {field names}      Like -f but treat arguments as a regular expression. Match all\n")
	fmt.Fprintf(o, "                        field names and operate on each in record order. Example: `-r '^[xy]$`'.\n")
	fmt.Fprintf(o, "  --nested-fs {string}  Defaults to \";\". Field separator for nested values.\n")
	fmt.Fprintf(o, "  --nested-ps {string}  Defaults to \":\". Pair separator for nested key-value pairs.\n")
	fmt.Fprintf(o, "  --evar {string}       Shorthand for --explode --values --across-records --nested-fs {string}\n")
	fmt.Fprintf(o, "  --ivar {string}       Shorthand for --implode --values --across-records --nested-fs {string}\n")
	fmt.Fprintf(o, "Please use \"%s --usage-separator-options\" for information on specifying separators.\n",
		argv0)

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Examples:\n")

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "  %s %s --explode --values --across-records -f x\n", argv0, verb)
	fmt.Fprintf(o, "  with input record \"x=a;b;c,y=d\" produces output records\n")
	fmt.Fprintf(o, "    \"x=a,y=d\"\n")
	fmt.Fprintf(o, "    \"x=b,y=d\"\n")
	fmt.Fprintf(o, "    \"x=c,y=d\"\n")
	fmt.Fprintf(o, "  Use --implode to do the reverse.\n")

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "  %s %s --explode --values --across-fields -f x\n", argv0, verb)
	fmt.Fprintf(o, "  with input record \"x=a;b;c,y=d\" produces output records\n")
	fmt.Fprintf(o, "    \"x_1=a,x_2=b,x_3=c,y=d\"\n")
	fmt.Fprintf(o, "  Use --implode to do the reverse.\n")

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "  %s %s --explode --pairs --across-records -f x\n", argv0, verb)
	fmt.Fprintf(o, "  with input record \"x=a:1;b:2;c:3,y=d\" produces output records\n")
	fmt.Fprintf(o, "    \"a=1,y=d\"\n")
	fmt.Fprintf(o, "    \"b=2,y=d\"\n")
	fmt.Fprintf(o, "    \"c=3,y=d\"\n")

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "  %s %s --explode --pairs --across-fields -f x\n", argv0, verb)
	fmt.Fprintf(o, "  with input record \"x=a:1;b:2;c:3,y=d\" produces output records\n")
	fmt.Fprintf(o, "    \"a=1,b=2,c=3,y=d\"\n")

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Notes:\n")
	fmt.Fprintf(o, "* With --pairs, --implode doesn't make sense since the original field name has\n")
	fmt.Fprintf(o, "  been lost.\n")
	fmt.Fprintf(o, "* The combination \"--implode --values --across-records\" is non-streaming:\n")
	fmt.Fprintf(o, "  no output records are produced until all input records have been read. In\n")
	fmt.Fprintf(o, "  particular, this means it won't work in `tail -f` contexts. But all other flag\n")
	fmt.Fprintf(o, "  combinations result in streaming (`tail -f` friendly) data processing.\n")
	fmt.Fprintf(o, "  If input is coming from `tail -f`, be sure to use `--records-per-batch 1`.\n")
	fmt.Fprintf(o, "* It's up to you to ensure that the nested-fs is distinct from your data's IFS:\n")
	fmt.Fprintf(o, "  e.g. by default the former is semicolon and the latter is comma.\n")
	fmt.Fprintf(o, "See also %s reshape.\n", argv0)
}

func transformerNestParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	fieldName := ""
	doRegexes := false
	nestedFS := ";"
	nestedPS := ":"
	doExplode := true
	doPairs := true
	doAcrossFields := true
	doExplodeSpecified := false
	doPairsSpecified := false
	doAcrossFieldsSpecified := false
	evfs := ""
	ivfs := ""

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
			transformerNestUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		} else if opt == "-f" {
			s, err := cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			fieldName = s

		} else if opt == "-r" {
			doRegexes = true
			s, err := cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			fieldName = s

		} else if opt == "--explode" || opt == "-e" {
			doExplode = true
			doExplodeSpecified = true
		} else if opt == "--implode" || opt == "-i" {
			doExplode = false
			doExplodeSpecified = true

		} else if opt == "--values" || opt == "-v" {
			doPairs = false
			doPairsSpecified = true
		} else if opt == "--pairs" || opt == "-p" {
			doPairs = true
			doPairsSpecified = true

		} else if opt == "--across-fields" || opt == "-F" {
			doAcrossFields = true
			doAcrossFieldsSpecified = true
		} else if opt == "--across-records" || opt == "-R" {
			doAcrossFields = false
			doAcrossFieldsSpecified = true

		} else if opt == "--nested-fs" || opt == "-S" {
			s, err := cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			nestedFS = s
		} else if opt == "--nested-ps" || opt == "-P" {
			s, err := cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			nestedPS = s

		} else if opt == "--evar" {
			s, err := cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			evfs = s
			doExplode = true
			doExplodeSpecified = true
			doPairs = false
			doPairsSpecified = true
			doAcrossFields = false
			doAcrossFieldsSpecified = true

		} else if opt == "--ivar" {
			s, err := cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			ivfs = s
			doExplode = false
			doExplodeSpecified = true
			doPairs = false
			doPairsSpecified = true
			doAcrossFields = false
			doAcrossFieldsSpecified = true

		} else {
			transformerNestUsage(os.Stderr)
			return nil, fmt.Errorf("%s %s: option \"%s\" not recognized", "mlr", verb, opt)
		}
	}

	if evfs != "" {
		doExplode = true
		doPairs = false
		doAcrossFields = false
		nestedFS = evfs
	}
	if ivfs != "" {
		doExplode = false
		doPairs = false
		doAcrossFields = false
		nestedFS = ivfs
	}

	if fieldName == "" {
		transformerNestUsage(os.Stderr)
		return nil, fmt.Errorf("%s %s: -f or -r is required", "mlr", verb)
	}
	if !doExplodeSpecified {
		transformerNestUsage(os.Stderr)
		return nil, fmt.Errorf("%s %s: --explode or --implode is required", "mlr", verb)
	}
	if !doPairsSpecified {
		transformerNestUsage(os.Stderr)
		return nil, fmt.Errorf("%s %s: --values or --pairs is required", "mlr", verb)
	}
	if !doAcrossFieldsSpecified {
		transformerNestUsage(os.Stderr)
		return nil, fmt.Errorf("%s %s: --across-records or --across-fields is required", "mlr", verb)
	}
	if doPairs && !doExplode {
		transformerNestUsage(os.Stderr)
		return nil, fmt.Errorf("%s %s: --implode with --pairs doesn't make sense", "mlr", verb)
	}
	if doRegexes && !doExplode {
		return nil, fmt.Errorf("mlr nest: -r is only supported with --explode, not --implode")
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerNest(
		fieldName,
		doRegexes,
		nestedFS,
		nestedPS,
		doExplode,
		doPairs,
		doAcrossFields,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerNest struct {
	fieldName string
	nestedFS  string
	nestedPS  string

	doRegexes  bool
	fieldRegex *regexp.Regexp // when doRegexes, for matching field names

	// For implode across fields (when !doRegexes)
	regex *regexp.Regexp

	// For implode across records
	otherKeysToOtherValuesToBuckets *lib.OrderedMap[*lib.OrderedMap[*tNestBucket]]

	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerNest(
	fieldName string,
	doRegexes bool,
	nestedFS string,
	nestedPS string,
	doExplode bool,
	doPairs bool,
	doAcrossFields bool,
) (*TransformerNest, error) {

	tr := &TransformerNest{
		fieldName: fieldName,
		doRegexes: doRegexes,
		nestedFS:  cli.SeparatorFromArg(nestedFS), // "pipe" -> "|", etc
		nestedPS:  cli.SeparatorFromArg(nestedPS),
	}

	// For implode across fields: regex to match exploded form (e.g. x_1, x_2)
	if doRegexes {
		fieldRegex, err := lib.CompileMillerRegex(fieldName)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"%s %s: cannot compile regex [%s]\n",
				"mlr", verbNameNest, fieldName,
			)
			os.Exit(1)
		}
		tr.fieldRegex = fieldRegex
		// implode uses fieldRegex directly when doRegexes
		tr.regex = fieldRegex
	} else {
		regexString := "^" + fieldName + "_[0-9]+$"
		regex, err := lib.CompileMillerRegex(regexString)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"%s %s: cannot compile regex [%s]\n",
				"mlr", verbNameNest, regexString,
			)
			os.Exit(1)
		}
		tr.regex = regex
	}

	// For implode across records
	tr.otherKeysToOtherValuesToBuckets = lib.NewOrderedMap[*lib.OrderedMap[*tNestBucket]]()

	if doExplode {
		if doPairs {
			if doAcrossFields {
				tr.recordTransformerFunc = tr.explodePairsAcrossFields
			} else {
				tr.recordTransformerFunc = tr.explodePairsAcrossRecords
			}
		} else {
			if doAcrossFields {
				tr.recordTransformerFunc = tr.explodeValuesAcrossFields
			} else {
				tr.recordTransformerFunc = tr.explodeValuesAcrossRecords
			}
		}
	} else {
		if doPairs {
			lib.InternalCodingErrorIf(true)
			// Should have been caught in CLI-parser.
		} else {
			if doAcrossFields {
				tr.recordTransformerFunc = tr.implodeValuesAcrossFields
			} else {
				tr.recordTransformerFunc = tr.implodeValueAcrossRecords
			}
		}
	}

	return tr, nil
}

func (tr *TransformerNest) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

// getMatchingFieldNames returns field names matching tr.fieldRegex in record order.
// When !tr.doRegexes, returns [tr.fieldName] if present, else [].
func (tr *TransformerNest) getMatchingFieldNames(inrec *mlrval.Mlrmap) []string {
	if !tr.doRegexes {
		if inrec.Get(tr.fieldName) != nil {
			return []string{tr.fieldName}
		}
		return nil
	}
	var names []string
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		if tr.fieldRegex.MatchString(pe.Key) {
			names = append(names, pe.Key)
		}
	}
	return names
}

func (tr *TransformerNest) explodeValuesAcrossFields(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {

		inrec := inrecAndContext.Record
		fieldNames := tr.getMatchingFieldNames(inrec)
		if len(fieldNames) == 0 {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
			return
		}

		for _, fieldName := range fieldNames {
			originalEntry := inrec.GetEntry(fieldName)
			if originalEntry == nil {
				continue
			}

			recordEntry := originalEntry
			mvalue := originalEntry.Value
			svalue := mvalue.String()

			// Not lib.SplitString so 'x=' will map to 'x_1=', rather than no field at all
			pieces := strings.Split(svalue, tr.nestedFS)
			i := 1
			for _, piece := range pieces {
				key := fieldName + "_" + strconv.Itoa(i)
				value := mlrval.FromString(piece)
				recordEntry = inrec.PutReferenceAfter(recordEntry, key, value)
				i++
			}

			inrec.Unlink(originalEntry)
		}
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // emit end-of-stream marker
	}
}

func (tr *TransformerNest) explodeValuesAcrossRecords(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		fieldNames := tr.getMatchingFieldNames(inrec)
		if len(fieldNames) == 0 {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
			return
		}
		fieldName := fieldNames[0]

		mvalue := inrec.Get(fieldName)
		if mvalue == nil {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
			return
		}
		svalue := mvalue.String()

		// Not lib.SplitString so 'x=' will map to 'x=', rather than no field at all
		pieces := strings.SplitSeq(svalue, tr.nestedFS)
		for piece := range pieces {
			outrec := inrec.Copy()
			outrec.PutReference(fieldName, mlrval.FromString(piece))
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(outrec, &inrecAndContext.Context))
		}

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // emit end-of-stream marker
	}
}

func (tr *TransformerNest) explodePairsAcrossFields(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {

		inrec := inrecAndContext.Record
		fieldNames := tr.getMatchingFieldNames(inrec)
		if len(fieldNames) == 0 {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
			return
		}

		for _, fieldName := range fieldNames {
			originalEntry := inrec.GetEntry(fieldName)
			if originalEntry == nil {
				continue
			}

			mvalue := originalEntry.Value
			svalue := mvalue.String()

			recordEntry := originalEntry
			pieces := lib.SplitString(svalue, tr.nestedFS)
			for _, piece := range pieces {
				pair := strings.SplitN(piece, tr.nestedPS, 2)
				if len(pair) == 2 { // there is a pair
					recordEntry = inrec.PutReferenceAfter(
						recordEntry,
						pair[0],
						mlrval.FromString(pair[1]),
					)
				} else { // there is not a pair
					recordEntry = inrec.PutReferenceAfter(
						recordEntry,
						fieldName,
						mlrval.FromString(piece),
					)
				}
			}

			inrec.Unlink(originalEntry)
		}
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // emit end-of-stream marker
	}
}

func (tr *TransformerNest) explodePairsAcrossRecords(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		fieldNames := tr.getMatchingFieldNames(inrec)
		if len(fieldNames) == 0 {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
			return
		}
		fieldName := fieldNames[0]

		mvalue := inrec.Get(fieldName)
		if mvalue == nil {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
			return
		}

		svalue := mvalue.String()
		pieces := lib.SplitString(svalue, tr.nestedFS)
		for _, piece := range pieces {
			outrec := inrec.Copy()

			originalEntry := outrec.GetEntry(fieldName)

			// Put the new field where the old one was -- unless there's already a field with the new
			// name, in which case replace its value.
			pair := strings.SplitN(piece, tr.nestedPS, 2)
			if len(pair) == 2 { // there is a pair
				outrec.PutReferenceAfter(originalEntry, pair[0], mlrval.FromString(pair[1]))
			} else { // there is not a pair
				outrec.PutReferenceAfter(originalEntry, fieldName, mlrval.FromString(piece))
			}

			outrec.Unlink(originalEntry)
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(outrec, &inrecAndContext.Context))
		}

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // emit end-of-stream marker
	}
}

func (tr *TransformerNest) implodeValuesAcrossFields(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		var previousEntry *mlrval.MlrmapEntry = nil
		fieldCount := 0
		var buffer bytes.Buffer
		for pe := inrec.Head; pe != nil; /* increment in loop */ {
			if tr.regex.MatchString(pe.Key) {
				if fieldCount > 0 {
					buffer.WriteString(tr.nestedFS)
				}
				buffer.WriteString(pe.Value.String())
				fieldCount++

				// Keep the location so we can implode in-place.
				if previousEntry == nil {
					previousEntry = pe.Prev
				}
				pnext := pe.Next
				inrec.Unlink(pe)
				pe = pnext
			} else {
				pe = pe.Next
			}
		}

		if fieldCount > 0 {
			newValue := mlrval.FromString(buffer.String())
			if previousEntry == nil { // No record before the unlinked one, i.e. list-head.
				inrec.PrependReference(tr.fieldName, newValue)
			} else {
				inrec.PutReferenceAfter(previousEntry, tr.fieldName, newValue)
			}
		}

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // emit end-of-stream marker
	}
}

func (tr *TransformerNest) implodeValueAcrossRecords(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		originalEntry := inrec.GetEntry(tr.fieldName)
		if originalEntry == nil {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
			return
		}

		fieldValueCopy := originalEntry.Value.Copy()

		// Don't unset tr.fieldName in the record, so we can implode in-place at the end.
		otherKeysJoined := inrec.GetKeysJoinedExcept(originalEntry)
		var otherValuesToBuckets *lib.OrderedMap[*tNestBucket] = nil
		iOtherValuesToBuckets := tr.otherKeysToOtherValuesToBuckets.Get(otherKeysJoined)
		if iOtherValuesToBuckets == nil {
			otherValuesToBuckets = lib.NewOrderedMap[*tNestBucket]()
			tr.otherKeysToOtherValuesToBuckets.Put(otherKeysJoined, otherValuesToBuckets)
		} else {
			otherValuesToBuckets = iOtherValuesToBuckets
		}

		otherValuesJoined := inrec.GetValuesJoinedExcept(originalEntry)
		var bucket *tNestBucket = nil
		bucket = otherValuesToBuckets.Get(otherValuesJoined)
		if bucket == nil {
			bucket = newNestBucket(inrec)
			otherValuesToBuckets.Put(otherValuesJoined, bucket)
		}

		pair := mlrval.NewMlrmapAsRecord()
		pair.PutReference(tr.fieldName, fieldValueCopy)
		bucket.pairs = append(bucket.pairs, pair)

	} else { // end of input stream

		for pe := tr.otherKeysToOtherValuesToBuckets.Head; pe != nil; pe = pe.Next {
			otherValuesToBuckets := pe.Value
			for pf := otherValuesToBuckets.Head; pf != nil; pf = pf.Next {
				var buffer bytes.Buffer
				bucket := pf.Value
				outrec := bucket.representative
				bucket.representative = nil // ownership transfer

				i := 0
				for _, pr := range bucket.pairs {
					if i > 0 {
						buffer.WriteString(tr.nestedFS)
					}
					i++
					buffer.WriteString(pr.Head.Value.String())
				}

				// tr.fieldName was already present so we'll overwrite it in-place here.
				outrec.PutReference(tr.fieldName, mlrval.FromString(buffer.String()))
				*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(outrec, &inrecAndContext.Context))
			}
		}

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // emit end-of-stream marker
	}
}

type tNestBucket struct {
	representative *mlrval.Mlrmap
	pairs          []*mlrval.Mlrmap
}

func newNestBucket(representative *mlrval.Mlrmap) *tNestBucket {
	return &tNestBucket{
		representative: representative,
		pairs:          []*mlrval.Mlrmap{},
	}
}
