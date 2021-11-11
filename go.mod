module mlr
// 'module github.com/johnkerl/miller' would be more standard, but it has the
// fatal flaw that 'go build' would produce a file named 'miller', not 'mlr' --
// and this naming goes back many years for Miller with executable named 'mlr',
// predating the Go port, across many platforms.

go 1.15

require (
	github.com/goccmack/gocc v0.0.0-20210331093148-09606ea4d4d9 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/lestrrat-go/strftime v1.0.4
	github.com/mattn/go-isatty v0.0.12
	golang.org/x/sys v0.0.0-20210326220804-49726bf1d181
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf
)
