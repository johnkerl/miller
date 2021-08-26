<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# A note on the complexity of Miller's expression language

One of Miller's strengths is its brevity: it's much quicker -- and less
error-prone -- to type `mlr stats1 -a sum -f x,y -g a,b` than having to track
summation variables as in `awk`, or using Miller's [out-of-stream
variables](reference-dsl-variables.md#out-of-stream-variables). And the more
language features Miller's put-DSL has (for-loops, if-statements, nested
control structures, user-defined functions, etc.) then the *less* powerful it
begins to seem: because of the other programming-language features it *doesn't*
have (classes, exceptions, and so on).

When I was originally prototyping Miller in 2015, the primary decision I had
was whether to hand-code in a low-level language like C or Rust or Go, with my
own hand-rolled DSL, or whether to use a higher-level language (like Python or
Lua or Nim) and let the `put` statements be handled by the implementation
language's own `eval`: the implementation language would take the place of a
DSL. Multiple performance experiments showed me I could get better throughput
using the former, by a wide margin. So Miller is Go under the hood with a
hand-rolled DSL.

I do want to keep focusing on what Miller is good at -- concise notation, low
latency, and high throughput -- and not add too much in terms of
high-level-language features to the DSL.  That said, some sort of
customizability is a basic thing to want. As of 4.1.0 we have recursive
`for`/`while`/`if` [structures](reference-dsl-control-structures.md) on about
the same complexity level as `awk`; as of 5.0.0 we have [user-defined
functions](reference-dsl-user-defined-functions.md) and [map-valued
variables](reference-dsl-variables.md), again on about the same complexity level
as `awk` along with optional type-declaration syntax; as of Miller 6 we have
full support for [arrays](reference-main-arrays.md).  While I'm excited by these
powerful language features, I hope to keep new features focused on Miller's
sweet spot which is speed plus simplicity.

