# Supported indexable lvalues

* Direct/indirect field name like `$x` or `$["x"]`
* Direct/indirect oosvar like `@x` or `@["x"]`
* Local variable like `x`
* Full srec `$*`
* Full oosvar `@*`

# Supported indexing

Each level by int or string:

* `$x[1]` or `$x["a"]`
* `@x[1]` or `@x["a"]`
* `x[1]` or `x["a"]`
* `$*[1]` or `$*["a"]`
* `@*[1]` (not supported) or `@*["a"]`

Multiple levels:

* Each can be further indexed, e.g. `$x[1]["a"][3]`

Auto-expand:

* `x[1][2][3] = 4` should not auto-expand (or should it?)
  * Accesses should be all in-bounds for each at-level array
  * I don't want to absent-fill the intervenings if someone sets `x[3]` when `x` is empty ... or maybe I should?
* `x["a"]["b"]["c"] = 4` should auto-expand
  * Create new maps at each level if necessary, unless they're already something else -- like `x["a"]` is already int/array/etc.

# Indexed types

* `$x` is a `Mlrval`
* `@x` is a `Mlrval`
* `x` is a `Mlrval
* `$*` is a `Mlrmap`
* `@*` is a `Mlrmap`

# Implementation

* `*Mlrval` needs a `PutIndexed` which takes `indices []*Mlrval` and `rvalue *Mlrval`.
* `*Mlrmap` needs a `PutIndexed` which takes `indices []*Mlrval` and `rvalue *Mlrval`.
