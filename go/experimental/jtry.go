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

	mlrval, err := lib.MlrvalDecodeFromJSON(decoder)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(mlrval)
}
