package input

import (
	"encoding/json"
	"os"

	"localdeps/ordered"

	"miller/clitypes"
	"miller/lib"
)

type RecordReaderJSON struct {
}

func NewRecordReaderJSON(readerOptions *clitypes.TReaderOptions) *RecordReaderJSON {
	return &RecordReaderJSON{}
}

func (this *RecordReaderJSON) Read(
	filenames []string,
	context lib.Context,
	inrecsAndContexts chan<- *lib.RecordAndContext,
	echan chan error,
) {
	if len(filenames) == 0 { // read from stdin
		handle := os.Stdin
		this.processHandle(handle, "(stdin)", &context, inrecsAndContexts, echan)
	} else {
		for _, filename := range filenames {
			handle, err := os.Open(filename)
			if err != nil {
				echan <- err
			} else {
				this.processHandle(handle, filename, &context, inrecsAndContexts, echan)
				handle.Close()
			}
		}
	}
	inrecsAndContexts <- lib.NewRecordAndContext(
		nil, // signals end of input record stream
		&context,
	)
}

func (this *RecordReaderJSON) processHandle(
	handle *os.File,
	filename string,
	context *lib.Context,
	inrecsAndContexts chan<- *lib.RecordAndContext,
	echan chan error,
) {
	context.UpdateForStartOfFile(filename)

	jsonDecoder := json.NewDecoder(handle)

	//	// Read opening bracket
	//	t, err := jsonDecoder.Token()
	//	if err != nil {
	//		echan <- err
	//		return
	//	}
	//	fmt.Printf("%T: %v\n", t, t)

	// Ordered-map idea from:
	//   https://gitlab.com/c0b/go-ordered-json
	// found via
	//   https://github.com/golang/go/issues/27179

	for jsonDecoder.More() {

		record := lib.NewMlrmap()

		var om *ordered.OrderedMap = ordered.NewOrderedMap()
		err := jsonDecoder.Decode(om)
		if err != nil {
			echan <- err
			return
		}

		// Use an iterator func to loop over all key-value pairs.  It is OK to call Set
		// append-modify new key-value pairs, but not safe to call Delete during
		// iteration.
		iter := om.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}

			key := pair.Key // copy
			value := pair.Value
			// TODO: handle object values

			//fmt.Println("value is a ", reflect.TypeOf(value))

			// xxx make helper functions
			sval, ok := value.(string)
			if ok {
				// If it's double-quoted, leave it as a string, even if it
				// looks like something parseable as int or float.
				mval := lib.MlrvalFromString(sval)
				record.Put(&key, &mval)
			} else {
				nval, ok := value.(json.Number)
				if ok {
					// xxx look deeper into input-format-preserving operations ...
					sval = nval.String()
					mval := lib.MlrvalFromInferredType(sval)
					record.Put(&key, &mval)
				}
			}
		}

		context.UpdateForInputRecord(record)

		inrecsAndContexts <- lib.NewRecordAndContext(
			record,
			context,
		)
	}

	//	// Read closing bracket
	//	t, err = jsonDecoder.Token()
	//	if err != nil {
	//		echan <- err
	//		return
	//	}
	//	fmt.Printf("%T: %v\n", t, t)
}
