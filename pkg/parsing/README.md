This directory contains a single source file, `mlr.bnf`, which is the lexical/semantic grammar file
for the Miller `put`/`filter` DSL (domain-specific language) using the [PGPG
framework](https://github.com/johnkerl/pgpg). (In a classical Lex/Yacc framework, there would be
separate `mlr.l` and `mlr.y` files; using PGPG, there is a single `mlr.bnf` file.)

All subdirectories of `pkg/parsing/` are autogen code created by PGPG's processing of `mlr.bnf`.
They are nonetheless committed to source control, since running PGPG takes a bit longer than the `go
build` does, and the BNF file doesn't often change. (_BNF_ is for _Backus-Naur Form_ which is the
phrasing of the grammar file that PGPG support.)

Run [tools/build-dsl](../../tools/build-dsl) from the repo root to regenerate the lexer and parser.
