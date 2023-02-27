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

// go get github.com/johnkerl/lumin@v1.0.0
// Local development:
// replace github.com/johnkerl/lumin => /Users/kerl/git/johnkerl/lumin

go 1.15

require (
	github.com/facette/natsort v0.0.0-20181210072756-2cd4dd1e2dcb
	github.com/johnkerl/lumin v1.0.0
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/lestrrat-go/strftime v1.0.6
	github.com/mattn/go-isatty v0.0.17
	github.com/nine-lives-later/go-windows-terminal-sequences v1.0.4
	github.com/pkg/profile v1.7.0
	github.com/stretchr/testify v1.8.2
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf
)
