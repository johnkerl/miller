# Status of the Go port

* This is not necessarily a full Go port of Miller. At the moment, it's a little spot for some experimentation. Things are very rough and very iterative and very incomplete. I don't commit to finishing a Go port but I very much hope to.
* One reason Miller exists is to be a useful tool for myself and others; another is it's fun to write. At bare minimum, I'll re-teach myself some Go.
* In all likelihood, though, this will turn into a full port which will someday become Miller 6.0.
* I hope to retain backward compatibility at the command-line level as much as possible.
* Benefits of porting to Go:
  * The lack of a streaming (record-by-record) JSON reader in the C implementation ([issue 99](https://github.com/johnkerl/miller/issues/99)) is immediately solved in the Go implementation.
  * The quoted-DKVP feature from [issue 266](https://github.com/johnkerl/miller/issues/266) will be easily addressed.
  * String/number-formatting issues in [issue 211](https://github.com/johnkerl/miller/issues/211), [issue 178](https://github.com/johnkerl/miller/issues/178), [issue 151](https://github.com/johnkerl/miller/issues/151), and [issue 259](https://github.com/johnkerl/miller/issues/259) will be fixed during the Go port.
  * I think some DST/timezone issues such as [issue 359](https://github.com/johnkerl/miller/issues/359) will be easier to fix using the Go datetime library than using the C datetime library
  * The code will be easier to read and, I hope, easier for others to contribute to.
* In the meantime I will still keep fixing bugs, doing some features, etc. in C on Miller 5.x -- in the near term, support for Miller's C implementation continues as before.

# Directory structure

Information here is for the benefit of anyone reading/using the Miller Go code. To use Miller, you don't need to know any of this if you don't want to. :)

Miller is a multi-format record-stream processor, where a **record** is a
sequence of key-value pairs. The basic **stream** operation is:

* **read** records in some specified file format;
* **map** the input records to output records in some user-specified way, using a **chain** of **verbs** (sort, filter, cut, put, etc.);
* **write** the records in some specified file format.

## Directory-structure overview

So, in broad overview, the key packages are:

* `src/miller/stream`   -- connect input -> mapping -> output via Go channels
* `src/miller/input`    -- read input records
* `src/miller/mapping`  -- map input records to output records
* `src/miller/output`   -- write output records
* The rest are details to support this

## Directory-structure details

### Dependencies

* Miller dependencies are all in the Go standard library, except a couple local ones:
  * `src/localdeps/ordered`
    * Insertion-ordered (order-preserving) maps from [gitlab.com/c0b/go-ordered-json](https://gitlab.com/c0b/go-ordered-json):
    * If you have a JSON data record `{"x":3,"y":4,"z":5}` then the keys `x,y,z` should stay that way. This package makes that happen.
  * `src/github.com/goccmack`
    * GOCC lexer/parser code-generator from [github.com/goccmack/gocc](https://github.com/goccmack/gocc):
    * This package defines the grammar for Miller's domain-specific language (DSL) for the Miller `put` and `filter` verbs. And, GOCC is a joy to use. :)

I didn't put GOCC into `src/localdeps` since `go get github.com/goccmack/gocc` uses this directory path, and is nice enough to also create `bin/gocc` for me -- so I thought I would just let it continue to do that. :)

### Miller per se

* Main entry point in `mlr.go`; everything else in `src/miller`
* `src/miller/lib`
  * `Mlrval` which includes string/int/float/boolean/void/absent/error types. These are used for record values, as well as expression/variable values in the Miller `put`/`filter` DSL. See also below for more details.
* `src/miller/containers`
  * `Lrec` is the sequence of key-value pairs which represents a Miller record. The key-lookup mechanism is optimized for Miller read/write usage patterns -- please see `lrec.go` for more details.
  * `context` supports AWK-like variables such as `FILENAME`, `NF`, `NR`, and so on.
* `src/miller/cli` is the flag-parsing logic for supporting Miller's command-line interface. When you type something like `mlr --icsv --ojson put '$sum = $a + $b' then filter '$sum > 1000' myfile.csv`, it's the CLI parser which makes it possible for Miller to construct a CSV record-reader, a mapper chain of `put` then `filter`, and a JSON record-writer.
* `src/miller/clitypes` contains datatypes for the CLI-parser, which was split out to avoid a Go package-import cycle.

* `src/miller/stream` is as above -- it uses Go channels to pipe together file-reads, to record-reading/parsing, to a chain of record-mappers, to record-writing/formatting, to terminal standard output.
* `src/miller/input` is as above -- one record-reader type per supported input file format, and a factory method.
* `src/miller/output` is as above -- one record-writer type per supported output file format, and a factory method.
* `src/miller/mapping` contains the abstract record-mapper interface datatype, as well as the Go-channel chaining mechanism for piping one mapper into the next.
* `src/miller/mappers` is all the concreate record-mappers such as `cat`, `tac`, `sort`, `put`, and so on. I put it here, not in `mapping`, so all files in `mappers` would be of the same type.
* `src/miller/parsing` contains a single source file, `mlr.bnf`, which is the lexical/semantic grammar file for the Miller `put`/`filter` DSL using the GOCC framework. All subdirectories of `src/miller/parsing/` are autogen code created by GOCC's processing of `mlr.bnf`.
* `src/miller/dsl` contains `ast.go` which is the abstract syntax tree datatype shared between GOCC and Miller. I didn't use a `src/miller/dsl/ast` naming convention, although that would have been nice, in order to avoid a Go package-dependency cycle.
* `src/miller/dsl/cst` is the concrete syntax tree, constructed from an AST produced by GOCC. The CST is what is actually executed on every input record when you do things like `$z = $x * 0.3 * $y`. Please see the `README.md` in `src/miller/dsl/cst` for more information.


### More about mlrvals

`Mlrval` is the datatype of record values, as well as expression/variable values in the Miller `put`/`filter` DSL. It includes string/int/float/boolean/void/absent/error types, not unlike PHP's `zval`.

* Miller's `absent` type is like Javascript's `undefined` -- it's for times when there is no such key, as in a DSL expression `$out = $foo` when the input record is `$x=3,y=4` -- there is no `$foo` so `$foo` has `absent` type.
* Miller's `void` type is like Javascript's `null` -- it's for times when there is a key with no value, as in `$out = $x` when the input record is `$x=,$y=4`.
* Miller's `error` type is for things like doing type-uncoerced addition of strings.
