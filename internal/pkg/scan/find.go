package scan

import (
)

// TODO: comment re context

//  o grammar for numbers & case-through
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

func FindScanType(sinput string) ScanType {
	input := []byte(sinput)

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
		if len(input) == 1 {
			return scanTypeString
		} else {
			return findScanTypePositiveDecimalOrFloatOrString(input)
		}
	}

	return scanTypeString
}

// Convenience function for unit test
func findScanTypeName(sinput string) string {
	return TypeNames[FindScanType(sinput)]
}

func findScanTypePositiveNumberOrString(input []byte) ScanType {
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
				if len(input) == 2 {
					return scanTypeString
				} else {
					return findScanTypePositiveHexOrString(input[2:])
				}
			}
			if i1 == 'o' || i1 == 'O' {
				if len(input) == 2 {
					return scanTypeString
				} else {
					return findScanTypePositiveOctalOrString(input[2:])
				}
			}
			if i1 == 'b' || i1 == 'B' {
				if len(input) == 2 {
					return scanTypeString
				} else {
					return findScanTypePositiveBinaryOrString(input[2:])
				}
			}

			allOctal := true
			allDecimal := true
			for _, c := range input[1:] {
				if !isOctalDigit(c) {
					allOctal = false
				}
				if !isDecimalDigit(c) {
					allDecimal = false
					break
				}
			}
			if allOctal {
				return scanTypeLeadingZeroOctalInt
			}
			if allDecimal {
				return scanTypeLeadingZeroDecimalInt
			}
			// else fall through
		}

		return findScanTypePositiveDecimalOrFloatOrString(input)
	}

	return scanTypeString
}

func findScanTypePositiveFloatOrString(input []byte) ScanType {
	for _, c := range []byte(input) {
		if !isFloatDigit(c) {
			return scanTypeString
		}
	}
	return scanTypeMaybeFloat
}

func findScanTypePositiveDecimalOrFloatOrString(input []byte) ScanType {
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

// Leading 0o has already been stripped
func findScanTypePositiveOctalOrString(input []byte) ScanType {
	for _, c := range []byte(input) {
		if !isOctalDigit(c) {
			return scanTypeString
		}
	}
	return scanTypeOctalInt
}

// Leading 0x has already been stripped
func findScanTypePositiveHexOrString(input []byte) ScanType {
	for _, c := range []byte(input) {
		if !isHexDigit(c) {
			return scanTypeString
		}
	}
	return scanTypeHexInt
}

// Leading 0b has already been stripped
func findScanTypePositiveBinaryOrString(input []byte) ScanType {
	for _, c := range []byte(input) {
		if c < '0' || c > '1' {
			return scanTypeString
		}
	}
	return scanTypeBinaryInt
}
