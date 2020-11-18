package output

import (
	"fmt"
	"os"

	"miller/clitypes"
	"miller/types"
)

// ostream *os.File in constructors/factory
type RecordWriterJSON struct {
	onFirst bool
}

func NewRecordWriterJSON(writerOptions *clitypes.TWriterOptions) *RecordWriterJSON {
	return &RecordWriterJSON{
		onFirst: true,
	}
}

func (this *RecordWriterJSON) Write(
	outrec *types.Mlrmap,
) {
	// End of record stream
	if outrec == nil {
		return
	}

	// TODO: --jlistwrap using onFirst and outrec == nil

	bytes, err := outrec.MarshalJSON()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Stdout.Write(bytes)
}
