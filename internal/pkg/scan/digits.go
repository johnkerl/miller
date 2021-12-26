package scan

// TODO: comment re context

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
