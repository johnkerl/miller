Logic for parsing the Miller command line.

* `internal/pkg/climain` is the flag-parsing logic for supporting Miller's command-line interface. When you type something like `mlr --icsv --ojson put '$sum = $a + $b' then filter '$sum > 1000' myfile.csv`, it's the CLI parser which makes it possible for Miller to construct a CSV record-reader, a transformer chain of `put` then `filter`, and a JSON record-writer.
* `internal/pkg/cli` contains datatypes for the CLI-parser, which was split out to avoid a Go package-import cycle.
* I don't use the Go [`flag`](https://golang.org/pkg/flag/) package. The `flag` package is quite fine; Miller's command-line processing is multi-purpose between serving CLI needs per se as well as for manpage/docfile generation, and I found it simplest to roll my own command-line handling here. More importantly, some Miller verbs such as ``sort`` take flags more than once -- ``mlr sort -f field1 -n field2 -f field3`` -- which is not supported by the `flag` package.
