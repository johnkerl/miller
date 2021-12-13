module github.com/johnkerl/miller

// The repo is 'miller' and the executable is 'mlr', going back many years and
// predating the Go port.
//
// If we had ./mlr.go then 'go build github.com/johnkerl/miller' then the
// executable would be 'miller' not 'mlr'.
//
// So we have cmd/mlr/main.go:
// * go build   github.com/johnkerl/miller/cmd/mlr
// * go install github.com/johnkerl/miller/cmd/mlr

go 1.15

require (
	github.com/goccmack/gocc v0.0.0-20210331093148-09606ea4d4d9 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/lestrrat-go/strftime v1.0.4
	github.com/mattn/go-isatty v0.0.12
	github.com/pkg/profile v1.6.0 // indirect
	golang.org/x/sys v0.0.0-20210326220804-49726bf1d181
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf
)
