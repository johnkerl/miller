Logic for parsing the Miller command line.

* `src/miller/cli` is the flag-parsing logic for supporting Miller's command-line interface. When you type something like `mlr --icsv --ojson put '$sum = $a + $b' then filter '$sum > 1000' myfile.csv`, it's the CLI parser which makes it possible for Miller to construct a CSV record-reader, a transformer chain of `put` then `filter`, and a JSON record-writer.
* `src/miller/cliutil` contains datatypes for the CLI-parser, which was split out to avoid a Go package-import cycle.
* I don't use the Go [`flag`](https://golang.org/pkg/flag/) package here, although I do use it within the transformers' subcommand flag-handling. The `flag` package is quite fine; Miller's command-line processing is multi-purpose between serving CLI needs per se as well as for manpage/docfile generation, and I found it simplest to roll my own command-line handling here.
