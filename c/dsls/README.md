# Miller domain-specific languages

These exist for Miller's `put` and `filter` functions. The grammars are not at
all profound: just parsing 101 as familiar from an introducy compilers course.
I use `lex` and `lemon` rather than `lex` and `yacc`: I find Lemon far more
transparent.

Concrete syntax trees (CSTs) are embodied in the `lex`/`lemon` files. Abstract
syntax trees (ASTs) are in the Miller `containers` directory.
