<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flags</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verbs</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Functions</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="../release-docs/index.html">Release docs</a>
</span>
</div>
# Building from source

Please also see [Installation](installing-miller.md) for information about pre-built executables.

You will need to first install Go version 1.15 or higher: please see [https://go.dev](https://go.dev).

## Miller license

Two-clause BSD license [https://github.com/johnkerl/miller/blob/master/LICENSE.txt](https://github.com/johnkerl/miller/blob/master/LICENSE.txt).

## From release tarball

* Obtain `mlr-i.j.k.tar.gz` from [https://github.com/johnkerl/miller/tags](https://github.com/johnkerl/miller/tags), replacing `i.j.k` with the desired release, e.g. `6.1.0`.
* `tar zxvf mlr-i.j.k.tar.gz`
* `cd mlr-i.j.k`
* `cd go`
* `make` creates the `./mlr` (or `.\mlr.exe` on Windows) executable
    * Without `make`: `go build github.com/johnkerl/miller/cmd/mlr`
* `make check` runs tests
    * Without `make`: `go test github.com/johnkerl/miller/internal/pkg/...` and `mlr regtest`
* `make install` installs the `mlr` executable and the `mlr` manpage
    * Without make: `go install github.com/johnkerl/miller/cmd/mlr` will install to _GOPATH_`/bin/mlr`

## From git clone

* `git clone https://github.com/johnkerl/miller`
* `make`/`go build github.com/johnkerl/miller/cmd/mlr` as above

## In case of problems

If you have any build errors, feel free to open an issue with "New Issue" at [https://github.com/johnkerl/miller/issues](https://github.com/johnkerl/miller/issues).

## Dependencies

### Required external dependencies

These are necessary to produce the `mlr` executable.

* Go version 1.15 or higher: please see [https://go.dev](https://go.dev)
* Others packaged within `go.mod` and `go.sum` which you don't need to deal with manually -- the Go build process handles them for us

### Optional external dependencies

This documentation pageset is built using [https://www.mkdocs.org/](MkDocs). Please see [https://github.com/johnkerl/miller/blob/main/docs/README.md](https://github.com/johnkerl/miller/blob/main/docs/README.md) for details.
