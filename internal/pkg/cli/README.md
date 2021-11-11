Datatypes for parsing the Miller command line, and the flags table.

* `internal/pkg/climain` is the flag-parsing logic for supporting Miller's command-line interface. When you type something like `mlr --icsv --ojson put '$sum = $a + $b' then filter '$sum > 1000' myfile.csv`, it's the CLI parser which makes it possible for Miller to construct a CSV record-reader, a transformer chain of `put` then `filter`, and a JSON record-writer.
* `internal/pkg/cli` contains datatypes and the flags table for the CLI-parser, which was split out to avoid a Go package-import cycle.
