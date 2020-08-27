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

	for jsonDecoder.More() {

		// decode an array value (Message)
		var pairs map[string]interface{}

		// oh no -- pairs are *not* order-preserved. :(
		// https://github.com/golang/go/issues/27179

		err = jsonDecoder.Decode(&pairs)
		if err != nil {
			echan <- err
			return
		}

		lrec := containers.LrecAlloc()

		//fmt.Printf("%v\n", pairs)
		for key, value := range pairs {
			//fmt.Printf("-- key: %v\n", key)
			//fmt.Printf("-- value: %v\n", value)
			foo := key // copy
			// TODO:
			// * handle int values
			// * handle float values
			// * handle object values

			//fmt.Println("value is a ", reflect.TypeOf(value))

			// xxx make helper functions
			sval, ok := value.(string)
			if ok {
				lrec.Put(&foo, &sval)
			} else {
				nval, ok := value.(float64)
				if ok {
					sval = strconv.FormatFloat(nval, 'g', -1, 64)
					lrec.Put(&foo, &sval)
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
