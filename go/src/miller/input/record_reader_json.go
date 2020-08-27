package input

import (
	// System:
	"encoding/json"
	//"fmt"
	"os"
	//"reflect"
	"strconv"

	// Miller:
	"miller/containers"

	// Local dependencies:
	"deps/ordered"
)

type RecordReaderJSON struct {
}

func NewRecordReaderJSON() *RecordReaderJSON {
	return &RecordReaderJSON{}
}

func (this *RecordReaderJSON) Read(
	filenames []string,
	inrecs chan<- *containers.Lrec,
	echan chan error,
) {
	// TODO: loop over filenames
	// TODO: handle empty filenames array as read-from-stdin
	filename := filenames[0]

	handle, err := os.Open(filename)
	if err != nil {
		echan <- err
		return
	}

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

		lrec := containers.LrecAlloc()

		var om *ordered.OrderedMap = ordered.NewOrderedMap()
		err = jsonDecoder.Decode(om)
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
				lrec.Put(&key, &sval)
			} else {
				nval, ok := value.(float64)
				if ok {
					sval = strconv.FormatFloat(nval, 'g', -1, 64)
					lrec.Put(&key, &sval)
				}
			}
		}
		inrecs <- lrec
	}

	//	// Read closing bracket
	//	t, err = jsonDecoder.Token()
	//	if err != nil {
	//		echan <- err
	//		return
	//	}
	//	fmt.Printf("%T: %v\n", t, t)

	inrecs <- nil // signals end of input record stream
}
