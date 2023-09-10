// ================================================================
// Experiments for type-inference performance optimization
// ================================================================

package main

import (
	"fmt"
	"os"

	"github.com/johnkerl/miller/pkg/scan"
)

func main() {
	for _, arg := range os.Args[1:] {
		scanType := scan.FindScanType(arg)
		fmt.Printf("%-10s -> %s\n", arg, scan.TypeNames[scanType])
	}
}
