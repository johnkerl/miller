// ================================================================
// Experiments for type-inference performance optimization
// ================================================================

package main

import (
	"fmt"
	"os"

	"github.com/johnkerl/miller/internal/pkg/scan"
)

// const (
//     scanTypeString     ScanType = 0
//     scanTypeDecimalInt          = 1
//     scanTypeOctalInt            = 2
//     scanTypeHexInt              = 3
//     scanTypeBinaryInt           = 4
//     scanTypeMaybeFloat          = 5
//     scanTypeBool                = 6
// )

func main() {
	// TODO:
	// func ParseInt(s string, base int, bitSize int) (int64, error)
	// func ParseUint(s string, base int, bitSize int) (uint64, error)

	for _, arg := range os.Args[1:] {
		scanType := scan.FindScanType(arg)
		fmt.Printf("%-10s -> %s\n", arg, scan.TypeNames[scanType])
	}
}
