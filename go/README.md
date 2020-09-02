# Status of the Go port

* This is not necessarily a full Go port of Miller. At the moment, it's a little spot for some experimentation. Things are very rough and very iterative and very incomplete. I don't commit to finishing a Go port but I very much hope to.
* One reason Miller exists is to be a useful tool for myself and others; another is it's fun to write. At bare minimum, I'll re-teach myself some Go.
* In all likelihood, though, this will turn into a full port which will someday become Miller 6.0.
* I hope to retain backward compatibility at the command-line level as much as possible.
* Benefits of porting to Go:
  * The lack of a streaming (record-by-record) JSON reader in the C implementation ([issue 99](https://github.com/johnkerl/miller/issues/99)) is immediately solved in the Go implementation.
  * The quoted-DKVP feature from [issue 266](https://github.com/johnkerl/miller/issues/266) will be easily addressed.
  * String/number-formatting issues in [issue 211](https://github.com/johnkerl/miller/issues/211), [issue 178](https://github.com/johnkerl/miller/issues/178), [issue 151](https://github.com/johnkerl/miller/issues/151), and [issue 259](https://github.com/johnkerl/miller/issues/259) will be fixed during the Go port.
  * I think some DST/timezone issues such as [issue 359](https://github.com/johnkerl/miller/issues/359) will be easier to fix using the Go datetime library than using the C datetime library
  * The code will be easier to read and, I hope, easier for others to contribute to.
* In the meantime I will still keep fixing bugs, doing some features, etc. in C on Miller 5.x -- in the near term, support for Miller's C implementation continues as before.

# Directory structure

Miller is a multi-format record-stream processor, where a *record* is a
sequence of key-value pairs. The basic *stream* operation is:

* *read* records in some specified file format;
* *map* the input records to output records in some user-specified way, using a *chain* of *verbs* (sort, filter, cut, put, etc.);
* *write* the records in some specified file format.

## Directory-structure overview

So, in broad overview, the key packages are:

```
src/miller/stream   -- connect input -> mapping -> output via Go channels
src/miller/input    -- read input records
src/miller/mapping  -- map input records to output records
src/miller/output   -- write output records
```

## Directory-structure details

### Dependencies

* Miller dependencies are all in the Go standard library, except a couple local ones:
  * Insertion-ordered (order-preserving) maps from [gitlab.com/c0b/go-ordered-json](https://gitlab.com/c0b/go-ordered-json):
    * If you have a JSON data record `{"x":3,"y":4,"z":5}` then the keys `x,y,z` should stay that way. This package makes that happen.
  * GOCC lexer/parser code-generator from [github.com/goccmack/gocc](https://github.com/goccmack/gocc):
    * This package defines the grammar for Miller's domain-specific language (DSL) for the Miller `put` and `filter` verbs. And, GOCC is a joy to use. :)


```
src/localdeps/ordered
src/github.com/goccmack
src/miller
```

### More

```
mlr.go
src/miller/lib
src/miller/containers

src/miller/cli
src/miller/clitypes
src/miller/stream
src/miller/input
src/miller/mapping
src/miller/output

src/miller/mappers
src/miller/parsing
src/miller/parsing/token
src/miller/parsing/util
src/miller/parsing/lexer
src/miller/parsing/parser
src/miller/parsing/errors
src/miller/dsl
src/miller/dsl/cst
```
