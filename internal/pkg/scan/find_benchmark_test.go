package scan

import (
	"testing"
)

// go test -run=nonesuch -bench=. github.com/johnkerl/miller/internal/pkg/scan/...

func BenchmarkFromNormalCases(b *testing.B) {

	data := []string{
		"yellow", "triangle", "true", "1", "11", "43.6498", "9.8870",
		"red", "square", "true", "2", "15", "79.2778", "0.0130",
		"red", "circle", "true", "3", "16", "13.8103", "2.9010",
		"red", "square", "false", "4", "48", "77.5542", "7.4670",
		"purple", "triangle", "false", "5", "51", "81.2290", "8.5910",
		"red", "square", "false", "6", "64", "77.1991", "9.5310",
		"purple", "triangle", "false", "7", "65", "80.1405", "5.8240",
		"yellow", "circle", "true", "8", "73", "63.9785", "4.2370",
		"yellow", "circle", "true", "9", "87", "63.5058", "8.3350",
		"purple", "square", "false", "10", "91", "72.3735", "8.2430",
	}
	ndata := len(data)

	for i := 0; i < b.N; i++ {
		_ = FindScanType(data[i%ndata])
	}
}

func BenchmarkFromAbnormalCases(b *testing.B) {

	data := []string{
		"", "-",
		"abc", "-abc",
		"0", "-0",
		"1", "-1",
		"2", "-2",
		"123", "-123",
		"1.", "-1.",
		".2", "-.2",
		".", "-.",
		"1.2", "-1.2",
		"1.2.3", "-1.2.3",
		"1e2e3", "-1e2e3",
		"12e-2", "-12e-2",
		"1e2x3", "-1e2x3",
		"0x", "-0x",
		"0x0", "-0x0",
		"0xcafe", "-0xcafe",
		"0xcape", "-0xcape",
		"0o", "-0o",
		"0o0", "-0o0",
		"0o1234", "-0o1234",
		"0b", "-0b",
		"0b0", "-0b0",
		"0b1011", "-0b1011",
		"0b1021", "-0b1021",
		"true", "true",
		"false", "false",
		"True", "True",
		"False", "False",
	}
	ndata := len(data)

	for i := 0; i < b.N; i++ {
		_ = FindScanType(data[i%ndata])
	}
}
