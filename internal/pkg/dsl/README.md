Parsing a Miller DSL (domain-specific language) expression goes through three representations:

* Source code which is a string of characters.
* Abstract syntax tree (AST)
* Concrete syntax tree (AST)

The job of the GOCC parser is to turn the DSL string into an AST.

The job of the CST builder is to turn the AST into a CST.

The job of the `put` and `filter` transformers is to execute the CST statements on each input record.

# Source-code representation

For example, the part between the single quotes in

`mlr put '$v = $i + $x * 4 + 100.7 * $y' myfile.dat`

# AST representation

Use `put -v` to display the AST:

```
mlr -n put -v '$v = $i + $x * 4 + 100.7 * $y'
RAW AST:
* StatementBlock
    * SrecDirectAssignment "=" "="
        * DirectFieldName "md_token_field_name" "v"
        * Operator "+" "+"
            * Operator "+" "+"
                * DirectFieldName "md_token_field_name" "i"
                * Operator "*" "*"
                    * DirectFieldName "md_token_field_name" "x"
                    * IntLiteral "md_token_int_literal" "4"
            * Operator "*" "*"
                * FloatLiteral "md_token_float_literal" "100.7"
                * DirectFieldName "md_token_field_name" "y"
```

Note the following about the AST:

* Parentheses, commas, semicolons, line endings, whitespace are all stripped away
* Variable names and literal values remain as leaf nodes of the AST
* Operators like `=` `+` `-` `*` `/` `**`, function names, and so on remain as non-leaf nodes of the AST
* Operator precedence is clear from the tree structure

Operator-precedence examples:

```
$ mlr -n put -v '$x = 1 + 2 * 3'
RAW AST:
* StatementBlock
    * SrecDirectAssignment "=" "="
        * DirectFieldName "md_token_field_name" "x"
        * Operator "+" "+"
            * IntLiteral "md_token_int_literal" "1"
            * Operator "*" "*"
                * IntLiteral "md_token_int_literal" "2"
                * IntLiteral "md_token_int_literal" "3"
```

```
$ mlr -n put -v '$x = 1 * 2 + 3'
RAW AST:
* StatementBlock
    * SrecDirectAssignment "=" "="
        * DirectFieldName "md_token_field_name" "x"
        * Operator "+" "+"
            * Operator "*" "*"
                * IntLiteral "md_token_int_literal" "1"
                * IntLiteral "md_token_int_literal" "2"
            * IntLiteral "md_token_int_literal" "3"
```

```
$ mlr -n put -v '$x = 1 * (2 + 3)'
RAW AST:
* StatementBlock
    * SrecDirectAssignment "=" "="
        * DirectFieldName "md_token_field_name" "x"
        * Operator "*" "*"
            * IntLiteral "md_token_int_literal" "1"
            * Operator "+" "+"
                * IntLiteral "md_token_int_literal" "2"
                * IntLiteral "md_token_int_literal" "3"
```

# CST representation

There's no `-v` display for the CST, but it's simply a reshaping of the AST
with pre-processed setup of function pointers to handle each type of statement
on a per-record basis.

The if/else and/or switch statements to decide what to do with each AST node
are done at CST-build time, so they don't need to be re-done when the syntax
tree is executed once on every data record.

# Source directories/files

* The AST logic is in `./ast*.go`.  I didn't use a `internal/pkg/dsl/ast` naming convention, although that would have been nice, in order to avoid a Go package-dependency cycle.
* The CST logic is in [`./cst`](./cst). Please see [cst/README.md](./cst/README.md) for more information.
