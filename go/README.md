# Status of the Go port

* This will be a full Go port of [Miller](https://miller.readthedocs.io/). Things are currently rough and iterative and incomplete. I don't have a firm timeline but I suspect it will take a few more months of late-evening/spare-time work.
* The released Go port will become Miller 6.0. As noted below, this will be a win both at the source-code level, and for users of Miller.
* I hope to retain backward compatibility at the command-line level as much as possible.
* In the meantime I will still keep fixing bugs, doing some features, etc. in C on Miller 5.x -- in the near term, support for Miller's C implementation continues as before.

# Port-completion criteria

* `reg-test/run` completes -- either completing/fixing the C/Go source-code discrepancies, or accepting the changes as backward incomptabilities
* Double-checking all Miller issues ever, in case I fixed/implemented something but didn't have reg-test coverage
* All `TODO`/`xxx` comments in Go, BNF source code, and case-files are resolved
* Release notes including Go-only features, and C/Go backward-incompatibilities
* Docs updated at https://miller.readthedocs.io/ (source-controlled at [../docs](../docs/))
* Equivalent of `./configure`, whatever that turns out to be

# Trying out the Go port

* Caveat: *lots* of things present in the C implementation are currently missing in the Go implementation. So if something doesn't work, it's almost certainly because it doesn't work *yet*.
* That said, if anyone is interested in playing around with it and giving early feedback, I'll be happy for it.
* Building:
  * Clone the Miller repo
  * `cd go`
  * `./build` should create `mlr`, and print the two lines `Compile OK` and `Test OK`. If it doesn't do this on your platform, please [file an issue](https://github.com/johnkerl/miller/issues).
* Platforms tried so far:
  * macOS with Go 1.14, and Linux Mint with Go 1.10
  * Windows I have not tried at all
* On-line help:
  * `mlr --help` advertises some things the Go implementation doesn't actually do yet.
  * `mlr --help-all-verbs` correctly lists verbs which do things in the Go implementation.
* See also https://github.com/johnkerl/miller/issues/372

# Benefits of porting to Go

* The [lack of a streaming (record-by-record) JSON reader](http://johnkerl.org/miller/doc/file-formats.html#JSON_non-streaming) in the C implementation ([issue 99](https://github.com/johnkerl/miller/issues/99)) is immediately solved in the Go implementation.
* In the C implementation, [arrays were not supported in the DSL](http://johnkerl.org/miller/doc/file-formats.html#Arrays); in the Go implementation they are.
* [Flattening nested map structures to output records](http://johnkerl.org/miller/doc/file-formats.html#Formatting_JSON_options) was clumsy. Now, Miller will be a JSON-to-JSON processor, if your inputs and outputs are both JSON; JSON input and output will be idiomatic.
* The quoted-DKVP feature from [issue 266](https://github.com/johnkerl/miller/issues/266) will be easily addressed.
* String/number-formatting issues in [issue 211](https://github.com/johnkerl/miller/issues/211), [issue 178](https://github.com/johnkerl/miller/issues/178), [issue 151](https://github.com/johnkerl/miller/issues/151), and [issue 259](https://github.com/johnkerl/miller/issues/259) will be fixed during the Go port.
* I think some DST/timezone issues such as [issue 359](https://github.com/johnkerl/miller/issues/359) will be easier to fix using the Go datetime library than using the C datetime library
* The code will be easier to read and, I hope, easier for others to contribute to. What this means is it should be quicker and easier to add new features to Miller -- after the development-time cost of the port itself is paid, of course.

# Things which may change

Please see https://github.com/johnkerl/miller/issues/372.

# Efficiency of the Go port

As I wrote [here](http://johnkerl.org/miller/doc/whyc.html) back in 2015 I couldn't get Rust or Go (or any other language I tried) to do some test-case processing as quickly as C, so I stuck with C.

Either Go has improved since 2015, or I'm a better Go programmer than I used to be, or both -- but as of 2020 I can get Go-Miller to process data about as quickly as C-Miller.

Note: in some sense Go-Miller is *less* efficient but in a way that doesn't significantly affect wall time. Namely, doing `mlr cat` on a million-record data file on my bargain-value MacBook Pro, the C version takes about 2.5 seconds and the Go version takes about 3 seconds. So in terms of wall time -- which is what we care most about, how long we have to wait -- it's about the same.

A way to look a little deeper at resource usage is to run `htop`, while processing a 10x larger file, so it'll take 25 or 30 seconds rather than 2.5 or 3. This way we can look at the steady-state resource consumption. I found that the C version -- which is purely single-threaded -- is taking 100% CPU. And the Go version, which uses concurrency and channels and `MAXPROCS=4`, with reader/transformer/writer each on their own CPU, is taking about 240% CPU. So Go-Miller is taking up not just a little more CPU, but a lot more -- yet, it does more work in parallel, and finishes the job in about the same amount of time.

Even commodity hardware has multiple CPUs these days -- and the Go code is *much* easier to read, extend, and improve than the C code -- so I'll call this a net win for Miller.

# Developer information

## Source-code goals

Donald Knuth famously said: *Programs are meant to be read by humans and only incidentally for computers to execute.*

During the coding of Miller, I've been guided by the following:

* *Miller should be pleasant to read.*
  * If you want to fix a bug, you should be able to quickly and confidently find out where and how.
  * If you want to learn something about Go channels, or lexing/parsing in Go -- especially if you don't already know much about them -- the comments should help you learn what you want to.
  * If you're the kind of person who reads other people's code for fun, well, the code should be fun, as well as readable.
  * `README.md` files throughout the directory tree are intended to give you a sense of what is where, what to read first and and what doesn't need reading right away, and so on -- so you spend a minimum of time being confused or frustrated.
  * Names of files, variables, functions, etc. should be fully spelled out (e.g. `NewEvaluableLeafNode`), except for a small number of most-used names where a longer name would cause unnecessary line-wraps (e.g. `Mlrval` instead of `MillerValue` since this appears very very often).
  * Code should not be too clever. This includes some reasonable amounts of code duplication from time to time, to keep things inline, rather than lasagna code.
  * Things should be transparent.  For example, `mlr -n put -v '$y = 3 + 0.1 * $x'` shows you the abstract syntax tree derived from the DSL expression.
  * Comments should be robust with respect to reasonably anticipated changes. For example, one package should cross-link to another in its comments, but I try to avoid mentioning specific filenames too much in the comments and README files since these may change over time. I make an exception for stable points such as [mlr.go](./mlr.go), [mlr.bnf](./src/miller/parsing/mlr.bnf), [stream.go](./src/miller/stream/stream.go), etc.
* *Miller should be pleasant to write.*
  * It should be quick to answer the question *Did I just break anything?* -- hence the `build` and `reg_test/run` regression scripts.
  * It should be quick to find out what to do next as you iteratively develop -- see for example [cst/README.md](https://github.com/johnkerl/miller/blob/master/go/src/miller/dsl/cst/README.md).
* *The language should be an asset, not a liability.*
  * One of the reasons I chose Go is that (personally anyway) I find it to be reasonably efficient, well-supported with standard libraries, straightforward, and fun.  I hope you enjoy it as much as I have.

## Directory structure

Information here is for the benefit of anyone reading/using the Miller Go code. To use the Miller tool at the command line, you don't need to know any of this if you don't want to. :)

## Directory-structure overview

Miller is a multi-format record-stream processor, where a **record** is a
sequence of key-value pairs. The basic **stream** operation is:

* **read** records in some specified file format;
* **transform** the input records to output records in some user-specified way, using a **chain** of **transformers** (also sometimes called **verbs**) -- sort, filter, cut, put, etc.;
* **write** the records in some specified file format.

So, in broad overview, the key packages are:

* [src/miller/stream](./src/miller/stream) -- connect input -> transforms -> output via Go channels
* [src/miller/input](./src/miller/input) -- read input records
* [src/miller/transforming](./src/miller/transforming) -- transform input records to output records
* [src/miller/output](./src/miller/output) -- write output records
* The rest are details to support this.

## Directory-structure details

### Dependencies

* Miller dependencies are all in the Go standard library, except a local one:
  * `src/github.com/goccmack`
    * GOCC lexer/parser code-generator from [github.com/goccmack/gocc](https://github.com/goccmack/gocc):
    * This package defines the grammar for Miller's domain-specific language (DSL) for the Miller `put` and `filter` verbs. And, GOCC is a joy to use. :)
    * Note on the path: `go get github.com/goccmack/gocc` uses this directory path, and is nice enough to also create `bin/gocc` for me -- so I thought I would just let it continue to do that by using that local path. :)
* I kept this locally so I could source-control it along with Miller and guarantee its stability. It is used on the terms of its open-source license.

### Miller per se

* The main entry point is [mlr.go](./mlr.go); everything else in [src/miller](./src/miller).
* [src/miller/lib](./src/miller/lib):
  * Implementation of the [`Mlrval`](./src/miller/types/mlrval.go) datatype which includes string/int/float/boolean/void/absent/error types. These are used for record values, as well as expression/variable values in the Miller `put`/`filter` DSL. See also below for more details.
  * [`Mlrmap`](./src/miller/types/mlrmap.go) is the sequence of key-value pairs which represents a Miller record. The key-lookup mechanism is optimized for Miller read/write usage patterns -- please see [mlrmap.go](./src/miller/types/mlrmap.go) for more details.
  * [`context`](./src/miller/types/context.go) supports AWK-like variables such as `FILENAME`, `NF`, `NR`, and so on.
* [src/miller/cli](./src/miller/cli) is the flag-parsing logic for supporting Miller's command-line interface. When you type something like `mlr --icsv --ojson put '$sum = $a + $b' then filter '$sum > 1000' myfile.csv`, it's the CLI parser which makes it possible for Miller to construct a CSV record-reader, a transformer-chain of `put` then `filter`, and a JSON record-writer.
* [src/miller/clitypes](./src/miller/clitypes) contains datatypes for the CLI-parser, which was split out to avoid a Go package-import cycle.
* [src/miller/stream](./src/miller/stream) is as above -- it uses Go channels to pipe together file-reads, to record-reading/parsing, to a chain of record-transformers, to record-writing/formatting, to terminal standard output.
* [src/miller/input](./src/miller/input) is as above -- one record-reader type per supported input file format, and a factory method.
* [src/miller/output](./src/miller/output) is as above -- one record-writer type per supported output file format, and a factory method.
* [src/miller/transforming](./src/miller/transforming) contains the abstract record-transformer interface datatype, as well as the Go-channel chaining mechanism for piping one transformer into the next.
* [src/miller/transformers](./src/miller/transformers) is all the concrete record-transformers such as `cat`, `tac`, `sort`, `put`, and so on. I put it here, not in `transforming`, so all files in `transformers` would be of the same type.
* [src/miller/parsing](./src/miller/parsing) contains a single source file, `mlr.bnf`, which is the lexical/semantic grammar file for the Miller `put`/`filter` DSL using the GOCC framework. All subdirectories of `src/miller/parsing/` are autogen code created by GOCC's processing of `mlr.bnf`.
* [src/miller/dsl](./src/miller/dsl) contains [`ast_types.go`](src/miller/dsl/ast_types.go) which is the abstract syntax tree datatype shared between GOCC and Miller. I didn't use a `src/miller/dsl/ast` naming convention, although that would have been nice, in order to avoid a Go package-dependency cycle.
* [src/miller/dsl/cst](./src/miller/dsl/cst) is the concrete syntax tree, constructed from an AST produced by GOCC. The CST is what is actually executed on every input record when you do things like `$z = $x * 0.3 * $y`. Please see the [src/miller/dsl/cst/README.md](./src/miller/dsl/cst/README.md) for more information.

## Nil-record conventions

Through out the code, records are passed by reference (as are most things, for
that matter, to reduce unnecessary data copies). In particular, records can be
nil through the reader/transformer/writer sequence.

* Record-readers produce a nil record-pointer to signify end of input stream.
* Each transformer takes a record-pointer as input and produces a sequence of zero or more record-pointers.
  * Many transformers, such as `cat`, `cut`, `rename`, etc. produce one output record per input record.
  * The `filter` transformer produces one or zero output records per input record depending on whether the record passed the filter.
  * The `nothing` transformer produces zero output records.
  * The `sort` and `tac` transformers are *non-streaming* -- they produce zero output records per input record, and instead retain each input record in a list. Then, when the nil-record end-of-stream marker is received, they sort/reverse the records and emit them, then they emit the nil-record end-of-stream marker.
  * Many transformers such as `stats1` and `count` also retain input records, then produce output once there is no more input to them.
* A null record-pointer at end of stream is passed to record-writers so that they may produce final output.
  * Most writers produce their output one record at a time.
  * The pretty-print writer produces no output until end of stream, since it needs to compute the max width down each column.

## Memory management

* Go has garbage collection which immediately simplifies the coding compared to the C port.
* Pointers are used freely for record-processing: record-readers allocate pointed records; pointed records are passed on Go channels from record-readers to record-transformers to record-writers.
  * Any transformer which passes an input record through is fine -- be it unmodifed as in `mlr cat` or modified as in `mlr cut`.
  * If a transformer drops a record (`mlr filter` in false cases, for example, or `mlr nothing`) it will be GCed.
  * One caveat is any transformer which produces multiples, e.g. `mlr repeat` -- this needs to explicitly copy records instead of producing multiple pointers to the same record.
* Right-hand-sides of DSL expressions all pass around pointers to records and Mlrvals.
  * Lvalue expressions return pointed `*types.Mlrmap` so they can be assigned to; rvalue expressions return non-pointed `types.Mlrval` but these are very shallow copies -- the int/string/etc types are copied but maps/arrays are passed by reference in the rvalue expression-evaluators.
* Copy-on-write is done on map/array put -- for example, in the assignment phase of a DSL statement, where an rvalue is assigned to an lvalue.

## More about mlrvals

[`Mlrval`](./src/miller/types/mlrval.go) is the datatype of record values, as well as expression/variable values in the Miller `put`/`filter` DSL. It includes string/int/float/boolean/void/absent/error types, not unlike PHP's `zval`.

* Miller's `absent` type is like Javascript's `undefined` -- it's for times when there is no such key, as in a DSL expression `$out = $foo` when the input record is `$x=3,y=4` -- there is no `$foo` so `$foo` has `absent` type. Nothing is written to the `$out` field in this case. See also [here](http://johnkerl.org/miller/doc/reference.html#Null_data:_empty_and_absent) for more information.
* Miller's `void` type is like Javascript's `null` -- it's for times when there is a key with no value, as in `$out = $x` when the input record is `$x=,$y=4`. This is an overlap with `string` type, since a void value looks like an empty string. I've gone back and forth on this (including when I was writing the C implementation) -- whether to retain `void` as a distinct type from empty-string, or not. I ended up keeping it as it made the `Mlrval` logic easier to understand.
* Miller's `error` type is for things like doing type-uncoerced addition of strings. Data-dependent errors are intended to result in `(error)`-valued output, rather than crashing Miller. See also [here](http://johnkerl.org/miller/doc/reference.html#Data_types) for more information.
* Miller's number handling makes auto-overflow from int to float transparent, while preserving the possibility of 64-bit bitwise arithmetic.
  * This is different from JavaScript, which has only double-precision floats and thus no support for 64-bit numbers (note however that there is now [`BigInt`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/BigInt)).
  * This is also different from C and Go, wherein casts are necessary -- without which int arithmetic overflows.
  * See also [here](http://johnkerl.org/miller/doc/reference.html#Arithmetic) for the semantics of Miller arithmetic, which the [`Mlrval`](./src/miller/types/mlrval.go) class implements.
 
## Software-testing methodology

See [./reg-test/README.md](./reg-test/README.md).
