module github.com/johnkerl/miller/v6

// The repo is 'miller' and the executable is 'mlr', going back many years and
// predating the Go port.
//
// If we had ./mlr.go then 'go build github.com/johnkerl/miller' then the
// executable would be 'miller' not 'mlr'.
//
// So we have cmd/mlr/main.go:
// * go build   github.com/johnkerl/miller/v6/cmd/mlr
// * go install github.com/johnkerl/miller/v6/cmd/mlr

// go get github.com/johnkerl/lumin@v1.0.0
// Local development:
// replace github.com/johnkerl/lumin => /Users/kerl/git/johnkerl/lumin

go 1.24.0

require (
	github.com/facette/natsort v0.0.0-20181210072756-2cd4dd1e2dcb
	github.com/johnkerl/lumin v1.0.0
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/klauspost/compress v1.18.2
	github.com/kshedden/statmodel v0.0.0-20210519035403-ee97d3e48df1
	github.com/lestrrat-go/strftime v1.1.1
	github.com/mattn/go-isatty v0.0.20
	github.com/nine-lives-later/go-windows-terminal-sequences v1.0.4
	github.com/pkg/profile v1.7.0
	github.com/stretchr/testify v1.11.1
	golang.org/x/sys v0.38.0
	golang.org/x/term v0.36.0
	golang.org/x/text v0.31.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/fgprof v0.9.3 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/google/pprof v0.0.0-20211214055906-6f57359322fd // indirect
	github.com/kshedden/dstream v0.0.0-20190512025041-c4c410631beb // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/tools v0.38.0 // indirect
	gonum.org/v1/gonum v0.16.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
