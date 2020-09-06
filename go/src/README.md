* Miller dependencies are all in the Go standard library, except a local one:
  * `src/github.com/goccmack`
    * GOCC lexer/parser code-generator from [github.com/goccmack/gocc](https://github.com/goccmack/gocc):
    * This package defines the grammar for Miller's domain-specific language (DSL) for the Miller `put` and `filter` verbs. And, GOCC is a joy to use. :)
    * Note on the path: `go get github.com/goccmack/gocc` uses this directory path, and is nice enough to also create `bin/gocc` for me -- so I thought I would just let it continue to do that by using that local path. :)
* I kept this locally so I could source-control it along with Miller and guarantee its stability. It is used on the terms of its open-source license.
