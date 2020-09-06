package main

import (
	"fmt"
	"os"

	"miller/lib"
)

// ----------------------------------------------------------------
func main() {
	for _, arg := range(os.Args[1:]) {
		mlrval := lib.MlrvalFromPending()
		err := mlrval.UnmarshalJSON([]byte(arg))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(mlrval)
		}
	}
}
