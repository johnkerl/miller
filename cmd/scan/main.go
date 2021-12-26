// ================================================================
// Experiments for type-inference performance optimization
// ================================================================

package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"

	"github.com/pkg/profile" // for trace.out
)

type tScanType int

const (
	scanTypeString     tScanType = 0
	scanTypeDecimalInt           = 1
	scanTypeOctalInt             = 2
	scanTypeHexInt               = 3
	scanTypeBinaryInt            = 4
	scanTypeMaybeFloat           = 5
	scanTypeBool                 = 6
)

var scanTypeNames = []string{
	"string",
	"decint",
	"octint",
	"hexint",
	"binint",
	"float?",
	"bool",
}

// 00000000: 00 01 02 03  04 05 06 07  08 09 0a 0b  0c 0d 0e 0f |................|
// 00000010: 10 11 12 13  14 15 16 17  18 19 1a 1b  1c 1d 1e 1f |................|
// 00000020: 20 21 22 23  24 25 26 27  28 29 2a 2b  2c 2d 2e 2f | !"#$%&'()*+,-./|
// 00000030: 30 31 32 33  34 35 36 37  38 39 3a 3b  3c 3d 3e 3f |0123456789:;<=>?|
// 00000040: 40 41 42 43  44 45 46 47  48 49 4a 4b  4c 4d 4e 4f |@ABCDEFGHIJKLMNO|
// 00000050: 50 51 52 53  54 55 56 57  58 59 5a 5b  5c 5d 5e 5f |PQRSTUVWXYZ[\]^_|
// 00000060: 60 61 62 63  64 65 66 67  68 69 6a 6b  6c 6d 6e 6f |`abcdefghijklmno|
// 00000070: 70 71 72 73  74 75 76 77  78 79 7a 7b  7c 7d 7e 7f |pqrstuvwxyz{|}~.|

var isDecimalDigitTable = []bool{
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 00-0f
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 10-1f
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 20-2f
	true, true, true, true, true, true, true, true, true, true, false, false, false, false, false, false, // 30-3f
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 40-4f
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 50-5f
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 60-6f
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 70-7f
}

var isHexDigitTable = []bool{
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 00-0f
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 10-1f
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 20-2f
	true, true, true, true, true, true, true, true, true, true, false, false, false, false, false, false, // 30-3f
	false, true, true, true, true, true, true, false, false, false, false, false, false, false, false, false, // 40-4f
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 50-5f
	false, true, true, true, true, true, true, false, false, false, false, false, false, false, false, false, // 60-6f
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 70-7f
}

// Possible character in floats include '.', 0-9, [eE], [-+] -- the latter two for things like 1.2e-8.
// Miller intentionally does not accept 'inf' or 'NaN' as float numbers in file-input data.
var isFloatDigitTable = []bool{
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 00-0f
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 10-1f
	false, false, false, false, false, false, false, false, false, false, false, true, false, true, true, false, // 20-2f
	true, true, true, true, true, true, true, true, true, true, false, false, false, false, false, false, // 30-3f
	false, false, false, false, false, true, false, false, false, false, false, false, false, false, false, false, // 40-4f
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 50-5f
	false, false, false, false, false, true, false, false, false, false, false, false, false, false, false, false, // 60-6f
	false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, // 70-7f
}

func isDecimalDigit(c byte) bool {
	if c < 128 { // byte is unsigned in Go
		return isDecimalDigitTable[c]
	} else {
		return false
	}
}

func isHexDigit(c byte) bool {
	if c < 128 { // byte is unsigned in Go
		return isHexDigitTable[c]
	} else {
		return false
	}
}

func isFloatDigit(c byte) bool {
	if c < 128 { // byte is unsigned in Go
		return isFloatDigitTable[c]
	} else {
		return false
	}
}

// ----------------------------------------------------------------

// TODO: UT the type-names LUT
// TODO: inout tabls & CLI access & UT access & bench access

func findScanType(input []byte) tScanType {
	if len(input) == 0 {
		return scanTypeString
	}

	i0 := input[0]
	if i0 == '-' {
		return findScanTypePositiveNumberOrString(input[1:])
	}
	if i0 >= '0' && i0 <= '9' {
		return findScanTypePositiveNumberOrString(input)
	}
	if i0 == '.' {
		return findScanTypePositiveDecimalOrFloatOrString(input)
	}

	sinput := string(input)
	if sinput == "true" || sinput == "false" {
		return scanTypeBool
	}

	return scanTypeString
}

// TODO: type up the why
//  o make a grammar for numbers & case-through
//    k len 0
//    - len 1
//    k has leading minus; strip & rest
//    - 0x, 0b, 0[0-9]
//    - decimal: leading minus; [0-9]+
//    - octal:   leading minus; 0[0-7]+
//    - hex:     leading minus; 0[xX][0-9a-fA-F]+
//    - float:   leadinug minus; [0-9] or '.'
//
//  o float literals:
//    123 123.  123.4 .234
//    1e2 1e-2 1.2e3 1.e3 1.2e-3 1.e-3
//    .2e3 .2e-3 1.e-3
//
//    ?- [0-9]+
//    ?- [0-9]+ '.' [0-9]*
//    ?- [0-9]* '.' [0-9]+
//    ?- [0-9]+            [eE] ?- [0-9]+
//    ?- [0-9]+ '.' [0-9]* [eE] ?- [0-9]+
//    ?- [0-9]* '.' [0-9]+ [eE] ?- [0-9]+

func findScanTypePositiveNumberOrString(input []byte) tScanType {
	if len(input) == 0 {
		return scanTypeString
	}
	i0 := input[0]

	if i0 == '.' {
		return findScanTypePositiveFloatOrString(input)
	}

	if isDecimalDigit(i0) {
		if len(input) == 1 {
			return scanTypeDecimalInt
		}
		if i0 == '0' {
			i1 := input[1]
			if i1 == 'x' || i1 == 'X' {
				return findScanTypePositiveHexOrString(input[2:])
			}
			if i1 == 'b' || i1 == 'B' {
				return findScanTypePositiveBinaryOrString(input[2:])
			}
		}

		// TODO: nope, could be float too
		return findScanTypePositiveDecimalOrFloatOrString(input)
	}

	return scanTypeString
}

func findScanTypePositiveFloatOrString(input []byte) tScanType {
	for _, c := range []byte(input) {
		if !isFloatDigit(c) {
			return scanTypeString
		}
	}
	return scanTypeMaybeFloat
}

func findScanTypePositiveDecimalOrFloatOrString(input []byte) tScanType {
	maybeInt := true
	for _, c := range []byte(input) {
		// All float digits are decimal-int digits so if the current character
		// is not a float digit, this can't be either a float or a decimal int.
		// Example: "1x2"
		if !isFloatDigit(c) {
			return scanTypeString
		}

		// Examples: "1e2" or "1x2".
		if !isDecimalDigit(c) {
			maybeInt = false
		}
	}
	if maybeInt {
		return scanTypeDecimalInt
	} else {
		return scanTypeMaybeFloat
	}
}

// Leading 0x has already been stripped
func findScanTypePositiveHexOrString(input []byte) tScanType {
	for _, c := range []byte(input) {
		if !isHexDigit(c) {
			return scanTypeString
		}
	}
	return scanTypeHexInt
}

// Leading 0b has already been stripped
func findScanTypePositiveBinaryOrString(input []byte) tScanType {
	for _, c := range []byte(input) {
		if c < '0' || c > '1' {
			return scanTypeString
		}
	}
	return scanTypeBinaryInt
}

// ----------------------------------------------------------------
func scanMain() {

	//	var c byte
	//
	//	fmt.Printf("dec: ")
	//	for c = 0x20; c <= 0x6f; c++ {
	//		if isDecimalDigit(c) {
	//			fmt.Printf("%c", c)
	//		}
	//	}
	//	fmt.Println()
	//
	//	fmt.Printf("hex: ")
	//	for c = 0x20; c <= 0x6f; c++ {
	//		if isHexDigit(c) {
	//			fmt.Printf("%c", c)
	//		}
	//	}
	//	fmt.Println()
	//
	//	fmt.Printf("float: ")
	//	for c = 0x20; c <= 0x6f; c++ {
	//		if isFloatDigit(c) {
	//			fmt.Printf("%c", c)
	//		}
	//	}
	//	fmt.Println()

	// TODO:
	// func ParseInt(s string, base int, bitSize int) (int64, error)
	// func ParseUint(s string, base int, bitSize int) (uint64, error)

	for _, arg := range os.Args[1:] {
		scanType := findScanType([]byte(arg))
		fmt.Printf("%-10s -> %s\n", arg, scanTypeNames[scanType])
	}
}

// ----------------------------------------------------------------
func main() {

	// Respect env $GOMAXPROCS, if provided, else set default.
	haveSetGoMaxProcs := false
	goMaxProcsString := os.Getenv("GOMAXPROCS")
	if goMaxProcsString != "" {
		goMaxProcs, err := strconv.Atoi(goMaxProcsString)
		if err != nil {
			runtime.GOMAXPROCS(goMaxProcs)
			haveSetGoMaxProcs = true
		}
	}
	if !haveSetGoMaxProcs {
		// As of Go 1.16 this is the default anyway. For 1.15 and below we need
		// to explicitly set this.
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	debug.SetGCPercent(500) // Empirical: See README-profiling.md

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// CPU profiling
	//
	// We do this here, not in the command-line parser, since
	// pprof.StopCPUProfile() needs to be called at the very end of everything.
	// Putting this pprof logic into a go func running in parallel with main,
	// and properly stopping the profile only when main ends via chan-sync,
	// results in a zero-length pprof file.
	//
	// Please see README-profiling.md for more information.

	if len(os.Args) >= 3 && os.Args[1] == "--cpuprofile" {
		profFilename := os.Args[2]
		handle, err := os.Create(profFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, os.Args[0], ": ", "Could not start CPU profile: ", err)
			return
		}
		defer handle.Close()

		if err := pprof.StartCPUProfile(handle); err != nil {
			fmt.Fprintln(os.Stderr, os.Args[0], ": ", "Could not start CPU profile: ", err)
			return
		}
		defer pprof.StopCPUProfile()

		fmt.Fprintf(os.Stderr, "CPU profile started.\n")
		defer fmt.Fprintf(os.Stderr, "CPU profile finished.\ngo tool pprof -http=:8080 %s\n", profFilename)
	}

	if len(os.Args) >= 3 && os.Args[1] == "--traceprofile" {
		defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()
		defer fmt.Fprintf(os.Stderr, "go tool trace trace.out\n")
	}

	scanMain()
}
