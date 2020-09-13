package main

import (
	"fmt"
	"os"

	"encoding/json"

	"miller/lib"
)

// ----------------------------------------------------------------
func main() {
	decoder := json.NewDecoder(os.Stdin)

	for {
		mlrval, eof, err := types.MlrvalDecodeFromJSON(decoder)
		if eof {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(mlrval)
	}
}
