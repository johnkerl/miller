# Purpose of the REPL

The Miller read-evaluate-print loop is an interactive counterpart to record-processing using the put/filter domain-specific language.

Using put and filter, you can do the following:

* Specify input format (e.g. `--icsv`), output format (e.g. `--ojson`), etc. using command-line flags.
* Specify filenames on the command line.
* Define `begin {...}` blocks which are executed before the first record is read.
* Define `end {...}` blocks which are executed after the last record is read.
* Define user-defined functions/subroutines using `func` and `subr`.
* Specify statements to be executed on each record -- which are anything outside of `begin`/`end`/`func`/`subr`.
* Example:
  `mlr --icsv --ojson put 'begin {print "HELLO"} $z = $x + $y; end {print "GOODBYE"}`

Using the REPL, by contrast, you get interactive control over those same steps:

* Specify input format (e.g. `--icsv`), output format (e.g. `--ojson`), etc. using command-line flags.
* Specify filenames either on the command line or via `:open` at the Miller REPL.
* Read records one at a time using `:read`.
* Skip ahead using statements `:skip 10` or `:skip until NR == 100` or `:skip until $status_code != 200`.
* Similarly, but processing records rather than skipping past them, using `:process` rather than `:skip`.
* Define `begin {...}` blocks; invoke them at will using `:begin`.
* Define `end {...}` blocks; invoke them at will using `:end`.
* Define user-defined functions/subroutines using `func`/`subr`; call them from other statements.
* Interactively specify statements to be executed on the current record.
* Load any of the above from Miller-script files using `:load`.
* Furthermore, any DSL statements other than `begin`/`end`/`func`/`subr` loaded using `:load` -- or from _multiline input mode_ which is where you type `<` on a line by itself, enter the code, then type `>` on a line by itself -- will be remembered and can be invoked on a given record using `:main`.  In multiline mode and load-from-file, semicolons are required between statements; otherwise they are not needed.

At this REPL prompt you can enter any Miller DSL expression.  REPL-only statements (non-DSL statements) start with `:`, such as `:help` or `:quit`.  Type `:help` to see more about your options.

No command-line-history-editing feature is built in but `rlwrap mlr repl` is a delight. You may need `brew install rlwrap`, `sudo apt-get install rlwrap`, etc. depending on your platform.

The input "record" by default is the empty map but you can do things like `$x=3`, or `unset $y`, or `$* = {"x": 3, "y": 4}` to populate it. Or, `:open foo.dat` followed by `:read` to populate it from a data file.

Non-assignment expressions, such as `7` or `true`, operate as filter conditions in the put DSL: they can be used to specify whether a record will or won't be included in the output-record stream.  But here in the REPL, they are simply printed to the terminal, e.g. if you type `1+2`, you will see `3`.

Examples:

```
$ mlr repl
Miller v6.0.0-dev
Type ':help' for on-line help; ':quit' to quit.
[mlr] 

```

```
[mlr] 1+2
3
```

```
[mlr] x=3  # These are local variables
[mlr] y=4
[mlr] x+y
7
```

```
[mlr] <
func f(a,b) {
  return a**b
}
>
[mlr] f(7,5)
16807
```

```
[mlr] :open foo.dat
[mlr] :read
[mlr] :context
FILENAME="foo.dat",FILENUM=1,NR=1,FNR=1
[mlr] $*
{
  "a": "eks",
  "b": "wye",
  "i": 4,
  "x": 0.38139939387114097,
  "y": 0.13418874328430463
}
[mlr] f($x,$i)
0.021160211005187134
[mlr] $z = f($x, $i)
[mlr] $*
{
  "a": "eks",
  "b": "wye",
  "i": 4,
  "x": 0.38139939387114097,
  "y": 0.13418874328430463,
  "z": 0.021160211005187134
}
```

# Implementation of the REPL

This is a small modification around the [CST](../../dsl/cst/) and the [`put` verb](../../transformers/put_or_filter.go).. Most of the keystroking here is for online help and command-line parsing.

One subtlety is that non-assignment expressions like `NR < 10` are filter statements within `put` -- they can be used to control whether or not a given record is included in the outut stream. Here, in the REPL, these expressions are simply printed to the terminal. And for `:skip until ...` or `:process until ...`, they're used as the exit condition to break out of reading input records.

# File structure

* [types.go](./types.go) -- data types including the `Repl` class
* [entry.go](./entry.go) -- shell command-line entry point to the Miller repl command line. E.g. handles `mlr repl --json` which is typed at the shell prompt, and starts a command-line session at the Miller REPL prompt.
* [session.go](./session.go) -- constructs a `Repl` object and ingests command lines, dispatching them either to the DSL (e.g. `$z = $x + $y`) or to the non-DSL verb handler (e.g. `:open foo.dat` or `help`).
* [prompt.go](./prompt.go) -- Handling for default and customized banners/prompts for the Miller REPL.
* [dsl.go](./dsl.go) -- Handler for taking DSL statements typed in interactively by the user, parsing them to an AST, building a CST from the AST, and executing the CST.
* [ast.go](./ast.go) -- Interface between the REPL and the DSL-to-AST parser.
* [verbs.go](./verbs.go) -- Handlers for non-DSL statements like `:open foo.dat` or `:help`.
