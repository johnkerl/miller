This major release significantly expands the expressiveness of the DSL for `mlr put` and `mlr filter`. (The upcoming 5.1.0 release will add the ability to aggregate across all columns for non-DSL verbs such as `mlr stats1` and `mlr stats2`.)

**Simple but impactful features**:
* [**Line endings (CRLF vs. LF, Windows-style vs. Unix-style) are now autodetected**](http://johnkerl.org/miller-releases/miller-head/doc/file-formats.html#Autodetect_of_line_endings). For example, files (including CSV) with LF input will lead to LF output unless you specify otherwise.
* There is now an [**in-place mode**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference.html#In-place_mode) using `mlr -I`.

**Major DSL features**:
* You can now [**define your own functions and subroutines**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-dsl.html#User-defined_functions_and_subroutines): e.g. `func f(x, y) { return x**2 + y**2 }`.
* New [**local variables**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-dsl.html#Local_variables) are completely analogous to out-of-stream variables: `sum` retains its value for the duration of the expression it's defined in; `@sum` retains its value across all records in the record stream.
* Local variables, function parameters, and function return types may be defined [**untyped or typed**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-dsl.html#Type-checking) as in `x = 1` or `int x = 1`, respectively. There are also expression-inline type-assertions available. Type-checking is up to you: omit it if you want flexibility with heterogeneous data; use it if you want to help catch misspellings in your DSL code or unexpected irregularities in your input data.
* There are now four kinds of maps. Out-of-stream variables have always been scalars, maps, or multi-level maps: `@a=1`, `@b[1]=2`, `@c[1][2]=3`. The same is now true for local variables, which are new to 5.0.0. Stream records have always been single-level maps; `$*` is a map. And as of 5.0.0 there are now [**map literals**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-dsl.html#Map_literals), e.g. `{"a":1, "b":2}`, which can be defined using JSON-like syntax (with either string or integer keys) and which can be nested arbitrarily deeply.
* You can [**loop over maps**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-dsl.html#For-loops) -- `$*`, out-of-stream variables, local variables, map-literals, and map-valued function return values -- using `for (k, v in ...)` or the new `for (k in ...)` (discussed next). All flavors of map may also be used in `emit` and `dump` statements.
* [**User-defined functions**](http://johnkerl.org/miller-releases/miller-head/doc/reference-dsl.html#User-defined_functions_and_subroutines) and subroutines may take **map-valued arguments**, and may **return map values**.
* Some [**built-in functions**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-dsl.html#Built-in_functions_for_filter_and_put) now accept map-valued input: `typeof`, `length`, `depth`, `leafcount`, `haskey`. There are built-in functions producing map-valued output: `mapsum` and `mapdiff`. There are now string-to-map and map-to-string functions: `splitnv`, `splitkv`, `splitnvx`, `splitkvx`, `joink`, `joinv`, and `joinkv`.

**Minor DSL features**:
* For iterating over maps (namely, local variables, out-of-stream variables, stream records, map literals, or return values from map-valued functions) there is now a [**key-only for-loop**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-dsl.html#Key-only_for-loops) syntax: e.g. `for (k in $*) { ... }`. This is in addition to the already-existing `for (k, v in ...)` syntax.
* There are now [**triple-statement for-loops**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-dsl.html#C-style_triple-for_loops) (familiar from many other languages), e.g. `for (int i = 0; i < 10; i += 1) { ... }`.
* `mlr put` and `mlr filter` now accept multiple `-f` for script files, freely intermixable with `-e` for expressions. The suggested use case is putting user-defined functions in script files and one-liners calling them using `-e`. Example: `myfuncs.mlr` defines the function `f(...)`, then `mlr put -f myfuncs.mlr -e '$o = f($i)' myfile.dat`. More information is [**here**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-dsl.html#Expressions_from_files).
* `mlr filter` is now almost identical to `mlr put`: it can have multiple statements, it can use `begin` and/or `end` blocks, it can define and invoke functions. Its final expression must evaluate to boolean which is used as the filter criterion. More details are [**here**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-dsl.html#Overview).
* The [**min and max functions**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-dsl.html#Built-in_functions_for_filter_and_put) are now variadic: `$o = max($a, $b, $c)`.
* There is now a [**substr function**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-dsl.html#Built-in_functions_for_filter_and_put) function.
* While `ENV` has long provided read-access to environment variables on the right-hand side of assignments (as a `getenv`), it now can be at the left-hand side of assignments (as a `putenv`). This is useful for subsidiary processes created by `tee`, `emit`, `dump`, or `print` when writing to a pipe.
* Handling for the `#` in comments is now handled in the lexer, so you can now (correctly) include `#` in strings.
* Separators are now available as read-only variables in the DSL: `IPS`, `IFS`, `IRS`, `OPS`, `OFS`, `ORS`. These are particularly useful with the split and join functions: e.g. with `mlr --ifs tab ...`, the `IFS` variable within a DSL expression will evaluate to a string containing a tab character.
* Syntax errors in DSL expressions now have a little more context.
* DSL parsing and execution are a bit more transparent. There have long been `-v` and `-t` options to `mlr put` and `mlr filter`, which print the expression's abstract syntax tree and do a low-level parser trace, respectively. There are now additionally `-a` which traces stack-variable allocation and `-T` which traces statements line by line as they execute. While `-v`, `-t`, and `-a` are most useful for development of Miller, the `-T` option gives you more visibility into what your Miller scripts are doing.

**Verbs**:
* [**most-frequent**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-verbs.html#most-frequent) and [**least-frequent**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-verbs.html#least-frequent) as requested in https://github.com/johnkerl/miller/issues/110.
* [**seqgen**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-verbs.html#seqgen) makes it easy to generate data from within Miller: please also see [**here**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/cookbook2.html#Generating_random_numbers_from_various_distributions) for a usage example.
* [**unsparsify**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-verbs.html#unsparsify) makes it easy to rectangularize data where not all records have the same fields.
* [**cat -n**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-verbs.html#cat) now takes a group-by (<tt>-g</tt>) option, making it easy to number records within categories.
* [**count-distinct**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-verbs.html#count-distinct),
[**uniq**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-verbs.html#uniq),
[**most-frequent**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-verbs.html#most-frequent),
[**least-frequent**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-verbs.html#least-frequent),
[**top**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-verbs.html#top), and
[**histogram**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-verbs.html#histogram)
now take a `-o` option for specifying their output field names, as requested in https://github.com/johnkerl/miller/issues/122.
* **Median** is now a synonym for p50 in [**stats1**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-verbs.html#stats1).
* You can now start a `then` chain with an initial `then`, which is nice in backslashy/multiline-continuation contexts.
This was requested in https://github.com/johnkerl/miller/issues/130.

**I/O options**:
* The `print` statement may now be used with no arguments, which prints a newline, and a no-argument `printn` prints nothing but creates a zero-length file in redirected-output context.
* Pretty-print format now has a `--pprint --barred` option (for output only, not input). For example please see [**here**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/file-formats.html#PPRINT:_Pretty-printed_tabular).
* There are now [**keystroke-savers**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/file-formats.html#Data-conversion_keystroke-savers) of the form `--c2p` which abbreviate `--icsvlite --opprint`, and so on.
* Miller's map literals are JSON-looking but allow integer keys which JSON doesn't. The 
`--jknquoteint` and `--jvquoteall` flags for `mlr` (when using JSON output) and `mlr put` (for `dump`) provide control over double-quoting behavior.

**Documents** new since the previous release:
* [**Miller in 10 minutes**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/10-min.html) is a long-overdue addition: while Miller's detailed documentation is evident, there has been a lack of more succinct examples.
* The [**cookbook**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/cookbook.html) has likewise been expanded, and has been split out
into three parts: [**part 1**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/cookbook.html), [**part
2**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/cookbook2.html), [**part 3**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/cookbook3.html).
* A bit more background on C performance compared to other languages I experimented with, early on in the development of Miller, is [**here**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/whyc.html#C_vs._Go,_D,_Rust,_etc.;_C_is_fast).

**On-line help**:
* Help for DSL [**built-in functions**](http://johnkerl.org/miller-releases/miller-5.0.0/doc/reference-dsl.html#Built-in_functions_for_filter_and_put), DSL keywords, and verbs is accessible using `mlr -f`, `mlr -k`, and `mlr -l` respectively; name-only lists are available with `mlr -F`, `mlr -K`, and `mlr -L`.

**Bugfixes**:
* A corner-case bug causing a segmentation violation on two `sub`/`gsub` statements within a single `put`, the first one matching its pattern and the second one not matching its pattern, has been fixed.

**Backward incompatibilities**: This is Miller 5.0.0, not 4.6.0, due to the following (all relatively minor):
* The `v` variables bound in for-loops such as `for (k, v in some_multi_level_map) { ... }` can now be map-valued if the `v` specifies a non-terminal in the map.
* There are new keywords such as `var`, `int`, `float`, `num`, `str`, `bool`, `map`, `IPS`, `IFS`, `IRS`, `OPS`, `OFS`, `ORS` which can no longer be used as variable names.
* Unset of the last key in an map-valued variable's map level no longer removes the level: e.g. with `@v[1][2]=3` and `unset @v[1][2]` the `@v` variable would be empty. As of 5.0.0, `@v` has key 1 with an empty-map value.
* There is no longer type-inference on literals: `"3"+4` no longer gives 7. (That was never a good idea.)
* The `typeof` function used to say things like `MT_STRING`; now it says things like `string`.
