package mappers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/lib"
	"miller/mapping"
)

// ----------------------------------------------------------------
var CutSetup = mapping.MapperSetup{
	Verb:         "cut",
	ParseCLIFunc: mapperCutParseCLI,
	IgnoresInput: false,
}

func mapperCutParseCLI(
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

	pFieldNames := flagSet.String(
		"f",
		"",
		"Comma-separated field names for cut, e.g. a,b,c",
	)

	pDoArgOrder := flagSet.Bool(
		"o",
		false,
		`Retain fields in the order specified here in the argument list.
Default is to retain them in the order found in the input data.`,
	)

	pDoComplement := flagSet.Bool(
		"x",
		false,
		"Exclude, rather than include, field names specified by -f.\n",
	)

	pDoComplementLong := flagSet.Bool(
		"complement",
		false,
		"Synonymous with -x",
	)

	//	ap_define_true_flag(pstate, "-r",           &do_regexes);
	//	fmt.Fprintf(o, "-r               Treat field names as regular expressions. \"ab\", \"a.*b\" will\n");
	//	fmt.Fprintf(o, "                 match any field name containing the substring \"ab\" or matching\n");
	//	fmt.Fprintf(o, "                 \"a.*b\", respectively; anchors of the form \"^ab$\", \"^a.*b$\" may\n");
	//	fmt.Fprintf(o, "                 be used. The -o flag is ignored when -r is present.\n");

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		mapperCutUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentioally requested
		return nil
	}

	if *pFieldNames == "" {
		flagSet.Usage()
		os.Exit(1)
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	fieldNameList := strings.Split(*pFieldNames, ",")
	mapper, _ := NewMapperCut(
		fieldNameList,
		*pDoArgOrder,
		*pDoComplement,
		*pDoComplementLong,
	)

	*pargi = argi
	return mapper
}

func mapperCutUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Passes through input records with specified fields included/excluded.\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})

	fmt.Fprintf(o, "Examples:\n")
	fmt.Fprintf(o, "  %s %s -f hostname,status\n", argv0, verb)
	fmt.Fprintf(o, "  %s %s -x -f hostname,status\n", argv0, verb)
	//	fmt.Fprintf(o, "  %s %s -r -f '^status$,sda[0-9]'\n", argv0, verb);
	//	fmt.Fprintf(o, "  %s %s -r -f '^status$,\"sda[0-9]\"'\n", argv0, verb);
	//	fmt.Fprintf(o, "  %s %s -r -f '^status$,\"sda[0-9]\"i' (this is case-insensitive)\n", argv0, verb);

}

// ----------------------------------------------------------------
type MapperCut struct {
	fieldNameList    []string
	fieldNameSet     map[string]bool
	doArgOrder       bool
	doComplement     bool
	doComplementLong bool
}

func NewMapperCut(
	fieldNameList []string,
	doArgOrder bool,
	doComplement bool,
	doComplementLong bool,
) (*MapperCut, error) {

	// xxx make/find-reuse util func
	fieldNameSet := make(map[string]bool)
	for _, fieldName := range fieldNameList {
		fieldNameSet[fieldName] = true
	}

	return &MapperCut{
		fieldNameList:    fieldNameList,
		fieldNameSet:     fieldNameSet,
		doArgOrder:       doArgOrder,
		doComplement:     doComplement,
		doComplementLong: doComplementLong,
	}, nil
}

//	if (!do_regexes) {
//		pstate->pfield_name_list   = pfield_name_list;
//		slls_reverse(pstate->pfield_name_list);
//		pstate->pfield_name_set    = hss_from_slls(pfield_name_list);
//		pstate->nregex             = 0;
//		pstate->regexes            = NULL;
//		pmapper->pprocess_func     = mapper_cut_process_no_regexes;
//	} else {
//		pstate->pfield_name_list   = NULL;
//		pstate->pfield_name_set    = NULL;
//		pstate->nregex = pfield_name_list->length;
//		pstate->regexes = mlr_malloc_or_die(pstate->nregex * sizeof(regex_t));
//		int i = 0;
//		for (sllse_t* pe = pfield_name_list->phead; pe != NULL; pe = pe->pnext, i++) {
//			// Let them type in a.*b if they want, or "a.*b", or "a.*b"i.
//			// Strip off the leading " and trailing " or "i.
//			regcomp_or_die_quoted(&pstate->regexes[i], pe->value, REG_NOSUB);
//		}
//		slls_free(pfield_name_list);
//		pmapper->pprocess_func = mapper_cut_process_with_regexes;
//	}

// ----------------------------------------------------------------
func (this *MapperCut) Map(
	inrecAndContext *lib.LrecAndContext,
	outrecsAndContexts chan<- *lib.LrecAndContext,
) {
	if !this.doComplement {
		if !this.doArgOrder {
			this.includeWithInputOrder(inrecAndContext, outrecsAndContexts)
		} else {
			this.includeWithArgOrder(inrecAndContext, outrecsAndContexts)
		}
	} else {
		this.exclude(inrecAndContext, outrecsAndContexts)
	}
}

// ----------------------------------------------------------------
// mlr cut -f a,b,c
func (this *MapperCut) includeWithInputOrder(
	inrecAndContext *lib.LrecAndContext,
	outrecsAndContexts chan<- *lib.LrecAndContext,
) {
	inrec := inrecAndContext.Lrec
	if inrec != nil { // not end of record stream
		outrec := lib.NewLrec()
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			fieldName := *pe.Key
			_, wanted := this.fieldNameSet[fieldName]
			if wanted {
				outrec.Put(&fieldName, pe.Value) // inrec will be GC'ed
			}
		}
		outrecAndContext := lib.NewLrecAndContext(outrec, &inrecAndContext.Context)
		outrecsAndContexts <- outrecAndContext
	} else {
		outrecsAndContexts <- inrecAndContext
	}
}

// ----------------------------------------------------------------
// mlr cut -o -f a,b,c
func (this *MapperCut) includeWithArgOrder(
	inrecAndContext *lib.LrecAndContext,
	outrecsAndContexts chan<- *lib.LrecAndContext,
) {
	inrec := inrecAndContext.Lrec
	if inrec != nil { // not end of record stream
		outrec := lib.NewLrec()
		for _, fieldName := range this.fieldNameList {
			value := inrec.Get(&fieldName)
			if value != nil {
				outrec.Put(&fieldName, value)
			}
		}
		outrecAndContext := lib.NewLrecAndContext(outrec, &inrecAndContext.Context)
		outrecsAndContexts <- outrecAndContext
	} else {
		outrecsAndContexts <- inrecAndContext
	}
}

// ----------------------------------------------------------------
// mlr cut -x -f a,b,c
func (this *MapperCut) exclude(
	inrecAndContext *lib.LrecAndContext,
	outrecsAndContexts chan<- *lib.LrecAndContext,
) {
	inrec := inrecAndContext.Lrec
	if inrec != nil { // not end of record stream
		for _, fieldName := range this.fieldNameList {
			if inrec.Has(&fieldName) {
				inrec.Remove(&fieldName)
			}
		}
	}
	outrecsAndContexts <- inrecAndContext
}

//// ----------------------------------------------------------------
//static sllv_t* mapper_cut_process_with_regexes(lrec_t* pinrec, context_t* pctx, void* pvstate) {
//	if (pinrec != NULL) {
//		mapper_cut_state_t* pstate = (mapper_cut_state_t*)pvstate;
//		// Loop over the record and free the fields to be discarded, being
//		// careful about the fact that we're modifying what we're looping over.
//		for (lrece_t* pe = pinrec->phead; pe != NULL; /* next in loop */) {
//			int matches_any = FALSE;
//			for (int i = 0; i < pstate->nregex; i++) {
//				if (regmatch_or_die(&pstate->regexes[i], pe->key, 0, NULL)) {
//					matches_any = TRUE;
//					break;
//				}
//			}
//			if (matches_any ^ pstate->do_complement) {
//				pe = pe->pnext;
//			} else {
//				lrece_t* pf = pe->pnext;
//				lrec_remove(pinrec, pe->key);
//				pe = pf;
//			}
//		}
//		return sllv_single(pinrec);
//	}
//	else {
//		return sllv_single(NULL);
//	}
//}
