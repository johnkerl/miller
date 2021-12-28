# Scope

This note is for a developer point of view. For a user point of view, please see [https://miller.readthedocs.io/en/latest/new-in-miller-6](https://miller.readthedocs.io/en/latest/new-in-miller-6).

# Quickstart for developers

See `makefile` in the repo base directory.

# Continuous integration

The Go implementation is auto-built using GitHub Actions: see [.github/workflows/go.yml](.github/workflows/go.yml). This works splendidly on Linux, MacOS, and Windows.

# Benefits of porting to Go

* The lack of a streaming (record-by-record) JSON reader in the C implementation ([issue 99](https://github.com/johnkerl/miller/issues/99)) is immediately solved in the Go implementation.
* In the C implementation, arrays were not supported in the DSL; in the Go implementation they are.
* Flattening nested map structures to output records was clumsy. Now, Miller will be a JSON-to-JSON processor, if your inputs and outputs are both JSON; JSON input and output will be idiomatic.
* The quoted-DKVP feature from [issue 266](https://github.com/johnkerl/miller/issues/266) will be easily addressed.
* String/number-formatting issues in [issue 211](https://github.com/johnkerl/miller/issues/211), [issue 178](https://github.com/johnkerl/miller/issues/178), [issue 151](https://github.com/johnkerl/miller/issues/151), and [issue 259](https://github.com/johnkerl/miller/issues/259) will be fixed during the Go port.
* I think some DST/timezone issues such as [issue 359](https://github.com/johnkerl/miller/issues/359) will be easier to fix using the Go datetime library than using the C datetime library
* The code will be easier to read and, I hope, easier for others to contribute to. What this means is it should be quicker and easier to add new features to Miller -- after the development-time cost of the port itself is paid, of course.

# Why Go

* As noted above, multiple Miller issues will benefit from stronger library support.
* Channels/goroutines are an excellent for Miller's reader/mapper/mapper/mapper/writer record-stream architecture.
* Since I did timing experiments in 2015, I found Go to be faster than it was then.
* In terms of CPU-cycle-count, Go is a bit slower than C (it does more things, like bounds-checking arrays and so on) -- but by leveraging concurrency over a couple processors, I find that it's competitive in terms of wall-time.
* Go is an up-and-coming language, with good reason -- it's mature, stable, with few of C's weaknesses and many of C's strengths.
* The source code will be easier to read/maintain/write, by myself and others.

# Efficiency of the Go port

As I wrote [here](https://johnkerl.org//miller-docs-by-release/1.0.0/performance.html) back in 2015 I couldn't get Rust or Go (or any other language I tried) to do some test-case processing as quickly as C, so I stuck with C.

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
  * Things should be transparent.  For example, the `-v` in `mlr -n put -v '$y = 3 + 0.1 * $x'` shows you the abstract syntax tree derived from the DSL expression.
  * Comments should be robust with respect to reasonably anticipated changes. For example, one package should cross-link to another in its comments, but I try to avoid mentioning specific filenames too much in the comments and README files since these may change over time. I make an exception for stable points such as [cmd/mlr/main.go](./cmd/mlr/main.go), [mlr.bnf](./internal/pkg/parsing/mlr.bnf), [stream.go](./internal/pkg/stream/stream.go), etc.
* *Miller should be pleasant to write.*
  * It should be quick to answer the question *Did I just break anything?* -- hence `mlr regtest` functionality.
  * It should be quick to find out what to do next as you iteratively develop -- see for example [cst/README.md](./internal/pkg/dsl/cst/README.md).
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

* [internal/pkg/stream](./internal/pkg/stream) -- connect input -> transforms -> output via Go channels
* [internal/pkg/input](./internal/pkg/input) -- read input records
* [internal/pkg/transformers](./internal/pkg/transformers) -- transform input records to output records
* [internal/pkg/output](./internal/pkg/output) -- write output records
* The rest are details to support this.

## Directory-structure details

### Dependencies

* Miller dependencies are all in the Go standard library, except two:
  * GOCC lexer/parser code-generator from [github.com/goccmack/gocc](https://github.com/goccmack/gocc):
    * This package defines the grammar for Miller's domain-specific language (DSL) for the Miller `put` and `filter` verbs. And, GOCC is a joy to use. :)
    * It is used on the terms of its open-source license.
  * [golang.org/x/term](https://pkg.go.dev/golang.org/x/term):
    * Just a one-line Miller callsite for is-a-terminal checking for the [Miller REPL](./internal/pkg/auxents/repl/README.md).
    * It is used on the terms of its open-source license.
* See also [./go.mod](go.mod). Setup:
  * `go get github.com/goccmack/gocc`
  * `go get golang.org/x/term`

### Miller per se

* The main entry point is [cmd/mlr/main.go](./cmd/mlr/main.go); everything else in [internal/pkg](./internal/pkg).
* [internal/pkg/entrypoint](./internal/pkg/entrypoint): All the usual contents of `main()` are here, for ease of testing.
* [internal/pkg/platform](./internal/pkg/platform): Platform-dependent code, which as of early 2021 is the command-line parser. Handling single quotes and double quotes is different on Windows unless particular care is taken, which is what this package does.
* [internal/pkg/lib](./internal/pkg/lib):
  * Implementation of the [`Mlrval`](./internal/pkg/types/mlrval.go) datatype which includes string/int/float/boolean/void/absent/error types. These are used for record values, as well as expression/variable values in the Miller `put`/`filter` DSL. See also below for more details.
  * [`Mlrmap`](./internal/pkg/types/mlrmap.go) is the sequence of key-value pairs which represents a Miller record. The key-lookup mechanism is optimized for Miller read/write usage patterns -- please see [mlrmap.go](./internal/pkg/types/mlrmap.go) for more details.
  * [`context`](./internal/pkg/types/context.go) supports AWK-like variables such as `FILENAME`, `NF`, `NR`, and so on.
* [internal/pkg/cli](./internal/pkg/cli) is the flag-parsing logic for supporting Miller's command-line interface. When you type something like `mlr --icsv --ojson put '$sum = $a + $b' then filter '$sum > 1000' myfile.csv`, it's the CLI parser which makes it possible for Miller to construct a CSV record-reader, a transformer-chain of `put` then `filter`, and a JSON record-writer.
* [internal/pkg/climain](./internal/pkg/climain) contains a layer which invokes `internal/pkg/cli`, which was split out to avoid a Go package-import cycle.
* [internal/pkg/stream](./internal/pkg/stream) is as above -- it uses Go channels to pipe together file-reads, to record-reading/parsing, to a chain of record-transformers, to record-writing/formatting, to terminal standard output.
* [internal/pkg/input](./internal/pkg/input) is as above -- one record-reader type per supported input file format, and a factory method.
* [internal/pkg/output](./internal/pkg/output) is as above -- one record-writer type per supported output file format, and a factory method.
* [internal/pkg/transformers](./internal/pkg/transformers) contains the abstract record-transformer interface datatype, as well as the Go-channel chaining mechanism for piping one transformer into the next. It also contains all the concrete record-transformers such as `cat`, `tac`, `sort`, `put`, and so on.
* [internal/pkg/parsing](./internal/pkg/parsing) contains a single source file, `mlr.bnf`, which is the lexical/semantic grammar file for the Miller `put`/`filter` DSL using the GOCC framework. All subdirectories of `internal/pkg/parsing/` are autogen code created by GOCC's processing of `mlr.bnf`. If you need to edit `mlr.bnf`, please use [tools/build-dsl](./tools/build-dsl) to autogenerate Go code from it (using the GOCC tool). (This takes several minutes to run.)
* [internal/pkg/dsl](./internal/pkg/dsl) contains [`ast_types.go`](internal/pkg/dsl/ast_types.go) which is the abstract syntax tree datatype shared between GOCC and Miller. I didn't use a `internal/pkg/dsl/ast` naming convention, although that would have been nice, in order to avoid a Go package-dependency cycle.
* [internal/pkg/dsl/cst](./internal/pkg/dsl/cst) is the concrete syntax tree, constructed from an AST produced by GOCC. The CST is what is actually executed on every input record when you do things like `$z = $x * 0.3 * $y`. Please see the [internal/pkg/dsl/cst/README.md](./internal/pkg/dsl/cst/README.md) for more information.

## Nil-record conventions

Through out the code, records are passed by reference (as are most things, for
that matter, to reduce unnecessary data copies). In particular, records can be
nil through the reader/transformer/writer sequence.

* Record-readers produce an end-of-stream marker (within the `RecordAndContext` struct) to signify end of input stream.
* Each transformer takes a record-pointer as input and produces a sequence of zero or more record-pointers.
  * Many transformers, such as `cat`, `cut`, `rename`, etc. produce one output record per input record.
  * The `filter` transformer produces one or zero output records per input record depending on whether the record passed the filter.
  * The `nothing` transformer produces zero output records.
  * The `sort` and `tac` transformers are *non-streaming* -- they produce zero output records per input record, and instead retain each input record in a list. Then, when the end-of-stream marker is received, they sort/reverse the records and emit them, then they emit the end-of-stream marker.
  * Many transformers such as `stats1` and `count` also retain input records, then produce output once there is no more input to them.
* An end-of-stream marker is passed to record-writers so that they may produce final output.
  * Most writers produce their output one record at a time.
  * The pretty-print writer produces no output until end of stream (or schema change), since it needs to compute the max width down each column.

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

[`Mlrval`](./internal/pkg/types/mlrval.go) is the datatype of record values, as well as expression/variable values in the Miller `put`/`filter` DSL. It includes string/int/float/boolean/void/absent/error types, not unlike PHP's `zval`.

* Miller's `absent` type is like Javascript's `undefined` -- it's for times when there is no such key, as in a DSL expression `$out = $foo` when the input record is `$x=3,y=4` -- there is no `$foo` so `$foo` has `absent` type. Nothing is written to the `$out` field in this case. See also [here](https://miller.readthedocs.io/en/latest/reference-main-null-data) for more information.
* Miller's `void` type is like Javascript's `null` -- it's for times when there is a key with no value, as in `$out = $x` when the input record is `$x=,$y=4`. This is an overlap with `string` type, since a void value looks like an empty string. I've gone back and forth on this (including when I was writing the C implementation) -- whether to retain `void` as a distinct type from empty-string, or not. I ended up keeping it as it made the `Mlrval` logic easier to understand.
* Miller's `error` type is for things like doing type-uncoerced addition of strings. Data-dependent errors are intended to result in `(error)`-valued output, rather than crashing Miller. See also [here](https://miller.readthedocs.io/en/latest/reference-main-data-types) for more information.
* Miller's number handling makes auto-overflow from int to float transparent, while preserving the possibility of 64-bit bitwise arithmetic.
  * This is different from JavaScript, which has only double-precision floats and thus no support for 64-bit numbers (note however that there is now [`BigInt`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/BigInt)).
  * This is also different from C and Go, wherein casts are necessary -- without which int arithmetic overflows.
  * See also [here](https://miller.readthedocs.io/en/latest/reference-main-arithmetic) for the semantics of Miller arithmetic, which the [`Mlrval`](./internal/pkg/types/mlrval.go) class implements.

## Performance optimizations

Key performance-related PRs for the Go port include:

* [https://github.com/johnkerl/miller/pull/424](#424): In C, lots of little mallocs are fine. In Go, not so much. This is not the garbage-collection penalty -- it's the penalty of _allocating_ -- lots of `duffcopy` and `madvise` appearing in the flame graphs. The idea was to reduce data-copies in the DSL.
* [https://github.com/johnkerl/miller/pull/765](#765): In C, prints to `stdout` are buffered a line at a time if the output is to the terminal, or a few KB at a time if not (i.e. file or pipe). Note the cost is how often the process does a `write` system call with associated overhead of context-switching into the kernel and back out. The C behavior is the right thing to do. In the Go port, very early on writes were all unbuffered -- several per record. Then buffering was soon switched to per-record, which was an improvement. But as of #765, the buffering is done at the library level, and it's done C-style -- much less frequently when output is not to a terminal.
* [https://github.com/johnkerl/miller/pull/774](#774): For CSV-lite and DKVP, this avoids using regexes to split strings when `strings.Split` will do.
* [https://github.com/johnkerl/miller/pull/779](#779): The basic idea of the Miller Go port was that the record-reader writes a record at a time over a channel to the first verb; the first verb writes records one at a time to the second verb, and so on; the last verb writes records one at a time to the record-writer. This is very simple, but for large files, the Go runtime scheduler overhead is too large -- data are chopped up into too many pieces. On #779 records are written 500 (or fewer) per batch, and all the channels from record-reader, to verbs, to record-writer are on record-batches. This lets Miller spend more time doing its job and less time yielding to the goroutine scheduler.
* [https://github.com/johnkerl/miller/pull/787](#786): In the C version, all values were strings until operated on specifically (expliclitly) by a verb. In the Go port, initially, all values were type-inferred on read, with types retained throughout the processing chain. This was an incredibly elegant and empowering design decision -- central to the Go port, in fact -- but it came with the cost that _all_ fields were being scanned as float/int even if they weren't used in the processing chain. On #786, fields are left as raw strings with type "pending", only just-in-time inferred to string/int/float only when used within the processing chain.
* [https://github.com/johnkerl/miller/pull/787](#787): This removed an unnecessary data copy in the `mlrval.String()` method. Originally this method had non-pointer receiver to conform with the `fmt.Stringer` interface. Hoewver, that's a false economy: `fmt.Println(someMlrval)` is a corner case, and stream processing is the primary concern. Implementing this as a pointer-receiver method was a performance improvement.
* [https://github.com/johnkerl/miller/pull/809](#809): This reduced the number of passes through fields for just-in-time type-inference. For example, for `$y = $x + 1`, each record's `$x` field's raw string (if not already accessed in the processing chain) needs to be checked to see if it's int (like `123`), float (like `123.4` or `1.2e3`), or string (anything else). Previously, succinct calls to built-in Go library functions were used. That was easy to code, but made too many expensive calls that were avoidable by lighter peeking of strings. In particular, an is-octal regex was being invoked unnecessarily on every field type-infer operation.

See also [./README-profiling.md](./README-profiling.md) and [https://miller.readthedocs.io/en/latest/new-in-miller-6/#performance-benchmarks](https://miller.readthedocs.io/en/latest/new-in-miller-6/#performance-benchmarks).

In summary:

* #765, #774, and #787 were low-hanging fruit.
* #424 was a bit more involved, and reveals that memory allocation -- not just GC -- needs to be handled more mindfully in Go than in C.
* #779 was a bit more involved, and reveals that Go's elegant goroutine/channel processing model comes with the caveat that channelized data should not be organized in many, small pieces.
* #809 was also bit more involved, and reveals that library functions are convenient, but profiling and analysis can sometimes reveal an opportunity for an impact, custom solution.
* #786 was a massive refactor involving about 10KLOC -- in hindsight it would have been best to do this work at the start of the Go port, not at the end.

## Software-testing methodology

See [./test/README.md](./test/README.md).

## Godoc

As of September 2021, `godoc` support is minimal: package-level synopses exist;
most `func`/`const`/etc content lacks `godoc`-style comments.

To view doc material, you can:

* `go get golang.org/x/tools/cmd/godoc`
* `cd go`
* `godoc -http=:6060 -goroot .`
* Browse to `http://localhost:6060`
* Note: control-C and restart the server, then reload in the browser, to pick up edits to source files

## Source-code indexing

Please see https://sourcegraph.com/github.com/johnkerl/miller
