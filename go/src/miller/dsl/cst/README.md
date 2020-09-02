See [go/src/miller/dsl/README.md](https://github.com/johnkerl/miller/blob/master/go/src/miller/dsl/README.md) for more information about Miller's use of abstract syntax trees (ASTs) and concrete syntax trees (CSTs) within the Miller `put`/`filter` domain-specific language (DSL).

## Files

* `types.go` is a starting point for seeing datatypes involved in the concrete syntax tree.
  * `IExecutable` is the interface for executable nodes, such as assignment statements, or statement blocks (if-bodies, etc.).
  * `IEvaluable` is the interface for evaluable expressions (e.g. right-hand sides of assignment statements).
* `root.go` contains the top-level logic for building a CST from an AST at parse time (`cst.Root.Build`), as well as executing the CST on a per-record basis (`cst.Root.Execute`). See also the [`put` mapper](https://github.com/johnkerl/miller/blob/master/go/src/miller/mappers/put.go).

## Notes

Go is a strongly typed language, but the AST is polymorphic. This results in if/else or switch statemens as an AST is walked.

Also, when we modify code, there can be changes in the [BNF grammar](https://github.com/johnkerl/miller/blob/master/go/src/miller/parsing/mlr.bnf) not yet reflected in the [AST](https://github.com/johnkerl/miller/blob/master/go/src/miller/dsl/ast.go). Likewise, there can be AST changes not yet reflected here. (Example: you are partway through adding a new binary operator to the grammar.)

As a result, throughout the code, there are error checks which may seem redundant but which are in place to make incremental development more pleasant and robust.

During CST build from an AST, one starts from the AST root and walks down through the nodes of the AST. Within a caller method, there is an if/else or switch statement on the AST node type. (Example: is this a leaf node, like the string literal `"abcd"`, int literal `3`, field-name `$x`? Or a binary operator like `+`, or function call like `cos`?).

Different builder methods are invoked for leaves, operators, etc. There is also, redundantly, a precondition assertion within each builder method: the leaf-builder method checks to make sure it's given an AST leaf node to build from; the operator-builder method checks to make sure it's given an AST operator node to build from; etc. The latter return Go `error` in case there is something new in the caller which has not yet been implemented in the callee.

This is all done to make development more happy: when you see things like `CST build: AST unary operator node unhandled.` you can check the code here and see what you need to do next to continue development.
