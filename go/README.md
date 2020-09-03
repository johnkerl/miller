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

# Efficiency of the Go port

As I wrote [here](http://johnkerl.org/miller/doc/whyc.html) back in 2015 I couldn't get Rust or Go (or any other language I tried) to do some test-case processing as quickly as C, so I stuck with C. 

Either Go has improved since 2015, or I'm a better Go programmer than I used to be, or both -- but as of 2020 I can get Go-Miller to process data about as quickly as C-Miller. 

Note: in some sense Go-Miller is *less* efficient but in a way that doesn't significantly affect wall time. Namely, doing `mlr cat` on a million-record data file on my bargain-value MacBook Pro, the C version takes about 2.5 seconds and the Go version takes about 3 seconds.  But using `htop` (while processing an even larger file, to see the steady-state resource consumption) shows that the C version -- which is purely single-threaded -- is taking 100% CPU, while the Go version -- which uses concurrency and channels and `MAXPROCS=4`, with reader/mapper/writer each on their own CPU -- is taking about 240% CPU. So Go-Miller is taking up quite a bit more CPU, but does more work in parallel -- to finish the job in about the same amount of time. 

Even commodity hardware has multiple CPUs these days -- and the Go code is *much* easier to read than the C code -- so I'll call this a net win for Go.

# Source-code goals

Donald Knuth famously said: *Programs are meant to be read by humans and only incidentally for computers to execute.*

During the coding of Miller, I've been guided by the following:

* *Miller should be pleasant to read.*
  * If you want to fix a bug, you should be able to quickly and confidently find out where an how.
  * If you want to learn something about Go channels, or lexing/parsing in Go -- especially if you don't already know much about them -- the comments should help you learn what you want to.
  * If you're the kind of person who reads other people's code for fun, well, the code should be fun, as well as readable.
  * `README.md` files throughout the directory tree are intended to give you a sense of what is where, what to read first and and what doesn't need reading right away, and so on -- so you spend a minimum of time being confused or frustrated.
  * Names of files, variables, functions, etc. should be fully spelled out (e.g. `NewEvaluableLeafNode`), except for a small number of most-used names where a longer name would cause unnecessary line-wraps (e.g. `Mlrval` instead of `MillerValue` since this appears very very often).
  * Code should not be too clever. This includes some reasonable amounts of code duplication from time to time, to keep things inline, rather than lasagna code.
  * Things should be transparent.  For example, `mlr -n put -v '$y = 3 + 0.1 * $x'` shows you the abstract syntax tree derived from the DSL expression.
  * Comments should be robust with respect to reasonably anticipated changes. For example, one package should cross-link to another in its comments, but I try to avoid mentioning specific filenames too much in the comments and README files since these may change over time. I make an exception for stable points such as `mlr.go`, `mlr.bnf`, `stream.go`, etc.
* *Miller should be pleasant to write.*
  * It should be quick to find out if you've made a mistake -- hence the `reg_test/run` regression script.
  * It should be quick to find out what to do next as you iteratively develop -- see for example [cst/README.md](https://github.com/johnkerl/miller/blob/master/go/src/miller/dsl/cst/README.md).
* *The language should be an asset, not a liability.*
  * One of the reasons I chose Go is that (personally anyway) I find it to be reasonably efficient, well-supported with standard libraries, straightforward to read, and fun to write.  I hope you enjoy it as much as I have.

# Directory structure

Information here is for the benefit of anyone reading/using the Miller Go code. To use the Miller tool at the command line, you don't need to know any of this if you don't want to. :)

## Directory-structure overview

Miller is a multi-format record-stream processor, where a **record** is a
sequence of key-value pairs. The basic **stream** operation is:

* **read** records in some specified file format;
* **map** the input records to output records in some user-specified way, using a **chain** of **mappers** (also sometimes called **verbs**) -- sort, filter, cut, put, etc.;
* **write** the records in some specified file format.

So, in broad overview, the key packages are:

* `src/miller/stream`   -- connect input -> mapping -> output via Go channels
* `src/miller/input`    -- read input records
* `src/miller/mapping`  -- map input records to output records
* `src/miller/output`   -- write output records
* The rest are details to support this.

## Directory-structure details

### Dependencies

* Miller dependencies are all in the Go standard library, except a couple local ones:
  * `src/localdeps/ordered`
    * Insertion-ordered (order-preserving) maps from [gitlab.com/c0b/go-ordered-json](https://gitlab.com/c0b/go-ordered-json):
    * If you have a JSON data record `{"x":3,"y":4,"z":5}` then the keys `x,y,z` should stay that way. This package makes that happen.
  * `src/github.com/goccmack`
    * GOCC lexer/parser code-generator from [github.com/goccmack/gocc](https://github.com/goccmack/gocc):
    * This package defines the grammar for Miller's domain-specific language (DSL) for the Miller `put` and `filter` verbs. And, GOCC is a joy to use. :)
    * I didn't put GOCC into `src/localdeps` since `go get github.com/goccmack/gocc` uses this directory path, and is nice enough to also create `bin/gocc` for me -- so I thought I would just let it continue to do that. :)
* I kept these locally so I could source-control them with Miller and guarantee their stability. They are used on the terms of the open-source licenses within their respective directories.

### Miller per se

* Main entry point is `mlr.go`; everything else in `src/miller`.
* `src/miller/lib`:
  * Implementation of the `Mlrval` datatype which includes string/int/float/boolean/void/absent/error types. These are used for record values, as well as expression/variable values in the Miller `put`/`filter` DSL. See also below for more details.
* `src/miller/containers`:
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

### Nil-lrec conventions

Through out the code, records are passed by reference (as are most things, for
that matter, to reduce unnecessary data copies). In particular, records can be
nil through the reader/mapper/writer sequence.

* Record-readers produce a nil lrec-pointer to signify end of input stream.
* Each mapper takes an lrec-pointer as input and produces a sequence of zero or more lrec-pointers.
  * Many mappers, such as `cat`, `cut`, `rename`, etc. produce one output record per input record.
  * The `filter` mapper produces one or zero output records per input record depending on whether the record passed the filter.
  * The `nothing` mapper produces zero output records.
  * The `sort` and `tac` mappers are *non-streaming* -- they produce zero output records per input record, and instead retain each input record in a list. Then, when the nil-lrec end-of-stream marker is received, they sort/reverse the records and emit them, then they emit the nil-lrec end-of-stream marker.
  * Many mappers such as `stats1` and `count` also retain input records, then produce output once there is no more input to them.
* A null lrec-pointer at end of stream is passed to lrec writers so that they may produce final output.
  * Most writers produce their output one record at a time.
  * The pretty-print writer produces no output until end of stream, since it needs to compute the max width down each column.

### More about mlrvals

`Mlrval` is the datatype of record values, as well as expression/variable values in the Miller `put`/`filter` DSL. It includes string/int/float/boolean/void/absent/error types, not unlike PHP's `zval`.

* Miller's `absent` type is like Javascript's `undefined` -- it's for times when there is no such key, as in a DSL expression `$out = $foo` when the input record is `$x=3,y=4` -- there is no `$foo` so `$foo` has `absent` type. Nothing is written to the `$out` field in this case. See also [here](http://johnkerl.org/miller/doc/reference.html#Null_data:_empty_and_absent) for more information.
* Miller's `void` type is like Javascript's `null` -- it's for times when there is a key with no value, as in `$out = $x` when the input record is `$x=,$y=4`. This is an overlap with `string` type, since a void value looks like an empty string. I've gone back and forth on this (including when I was writing the C implementation) -- whether to retain `void` as a distinct type from empty-string, or not. I ended up keeping it as it made the `Mlrval` logic easier to understand.
* Miller's `error` type is for things like doing type-uncoerced addition of strings. Data-dependent errors are intended to result in `(error)`-valued output, rather than crashing Miller. See also [here](http://johnkerl.org/miller/doc/reference.html#Data_types) for more information.
* Miller's number handling makes auto-overflow from int to float transparent, while preserving the possibility of 64-bit bitwise arithmetic.
  * This is different from JavaScript, which has only double-precision floats and thus no support for 64-bit numbers (note however that there is now [`BigInt`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/BigInt)).
  * This is also different from C and Go, wherein casts are necessary -- without which int arithmetic overflows.
  * See also [here](http://johnkerl.org/miller/doc/reference.html#Arithmetic) for the semantics of Miller arithmetic, which the `Mlrval` class implements.
