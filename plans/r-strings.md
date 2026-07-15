# Plan: explicit r-strings for the Miller DSL

## Context

[Issue #297](https://github.com/johnkerl/miller/issues/297) ("Match regex
metacharacters: remaining step now is r-strings") is otherwise resolved in
Miller 6. The one remaining ask, per @johnkerl's final comments on the issue,
is **explicit r-strings**: `r"..."` as a DSL literal that behaves like
Python's raw strings — no backslash-escape processing — usable both directly
in regex position and, critically, **assigned to a variable that's later
used in regex position**:

```
rstar = r"\*";
$y = gsub($x, rstar, "star");
```

This is exactly the case @torbiak raised in the issue thread and that
@johnkerl agreed doesn't work today. @johnkerl's preferred design (Option 1 in
the issue) is additive: keep today's *implicit* r-string behavior for plain
string literals in regex position, and add *explicit* `r"..."` as a new,
independently-usable literal type on top.

## Why variables can't carry "raw-ness" today

Literal strings in "regex position" (2nd arg to
`sub`/`gsub`/`regextract`/`regextract_or_else`, RHS of `=~`/`!=~`) get raw
(non-unescaped) treatment via a **parse-time AST trick**, not a runtime
value property:

- `regexProtectPrePassAux` (`pkg/dsl/cst/leaves.go:230-262`, invoked from
  `pkg/dsl/cst/root.go`) walks the AST and relabels the 2nd child of those
  callsites/operators from `NodeTypeStringLiteral` to `NodeTypeRegex` —
  **only when that child is already a string-literal AST node**.
- `BuildLeafNode` (`pkg/dsl/cst/leaves.go:42-67`) then dispatches
  `NodeTypeStringLiteral` → `BuildStringLiteralNode` (calls
  `lib.UnbackslashStringLiteral`, `pkg/lib/unbackslash.go:38-97`, converting
  `\t`→TAB, `\\`→`\`, etc.) vs. `NodeTypeRegex` → `BuildRegexLiteralNode`
  (`leaves.go:270-274`, no unescaping — wraps the raw lexeme straight into
  an `Mlrval`).

A local variable or field sitting in the same argument slot is untouched by
the prepass (it's not a string-literal AST node), so its *runtime* string
value is whatever ordinary unescaping already produced when it was built —
there's nothing left for regex compilation to "undo". E.g. `star_re = "\*"`
unescapes to a single `*` byte (since `\*` isn't a recognized
general-purpose escape and the backslash is simply dropped), which is then
an invalid/misleading regex fragment on its own.

## Where regex compilation happens (and why the fix is simple)

`pkg/lib/regex.go`, `CompileMillerRegex` (lines 97-125), is the single
chokepoint for `sub`/`gsub`/`regextract`/`regextract_or_else`/
`strmatch`/`strmatchx`/`=~`/`!=~` (called from `pkg/bifs/regex.go`). It:

- strips a leading/trailing `"` (optionally with trailing `"..."i` → Go's
  `(?i)` prefix) if present, **or**
- if the string has no surrounding quotes, passes it straight to
  `regexp.Compile` unchanged (line 124, the "bare" fallback).

This bare fallback is the load-bearing fact for this design: **a
fully-unquoted raw string value (backslashes intact, no surrounding `"`)
already compiles correctly today with zero changes to `regex.go` or
`bifs/regex.go`.** If an explicit r-string literal builds a clean, unquoted
`Mlrval` at parse time, it works correctly as a regex argument whether used
directly or via a variable, with no changes to the compilation path at all.

## Existing precedents to model on

- **`bytes_literal` (`b"..."`)** — a single-letter-prefixed quoted literal
  already in the grammar (`pkg/parsing/mlr.bnf:98-99`, comment: "must
  precede non_sigil_name since 'b' is an idchar" — the same ordering
  constraint an `r` prefix needs, since `r` is also a valid identifier
  character). Its own AST/CST node type and builder,
  `BuildBytesLiteralNode` (`pkg/dsl/cst/leaves.go:342-357`), strips the
  leading `b` and surrounding quotes, then — unlike what we want for
  r-strings — still calls `UnbackslashStringLiteral`.
- **`RegexCaseInsensitive`** (`mlr.bnf:704-705`, `leaves.go:58-67`) — the
  precedent for the optional trailing `i` suffix: grammar rule
  `string_literal non_sigil_name`; CST dispatch appends `"i"` to the
  literal text if not already present and delegates to the regex-literal
  builder, which leaves the surrounding quotes on the `Mlrval` so
  `CompileMillerRegex`'s `"..."i` branch can find them later.

## Prior art (dead end — do not reuse)

Git commit `1230553eb` ("regex r-string feature", Aug 2021) added
`regex_r_string ::= 'r' '"' {...} '"'` mapping directly to `NodeTypeRegex`
— but on the pre-Miller6-restructure `go/src/parsing/...` tree, which was
deleted in the `pkg/...` rewrite and never carried forward. Its unit tests
and docs were left as TODO, unchecked, even before the tree was deleted.
Not reusable as code, but it confirms the grammar shape (`'r' '"' {...}
'"'`) is sound.

## Docs already describe the implicit-r behavior

`docs/src/reference-main-regular-expressions.md.in` (~lines 409-417):

> "...if you're familiar with r-strings in Python) all strings in regex
> position are implicit r-strings."

This is the natural place to add explicit `r"..."` documentation.

## Pinned regression test — do not touch

`test/cases/dsl-regex-matching/0016` (input
`test/input/regex-metacharacters.dkvp`) tests exactly the implicit-r
scenario from the issue (`gsub($input, "\[", "LEFT")` etc.) — added by
commits `b20a5ccd3`/`55209bfc5` ("Test case for #297"). This work must not
change its behavior.

## Design

### 1. Grammar (`pkg/parsing/mlr.bnf`)

Add an `r_string_literal` lexer rule next to `bytes_literal`, before
`non_sigil_name` (same ordering requirement, since `r` is a valid identifier
character):

```
# Raw/r-strings r"..." (must precede non_sigil_name since 'r' is an idchar)
r_string_literal ::= 'r' '"' { _string_char | _escape } '"' ;
```

Add `RStringCaseInsensitive` and `RStringLiteral` productions next to
`RegexCaseInsensitive`/`StringLiteral` (`mlr.bnf:703-709`), case-insensitive
variant declared first (same convention as the existing pair):

```
# r"a.*b" (raw, case-sensitive) or r"a.*b"i (raw, case-insensitive). Must precede RStringLiteral.
RStringCaseInsensitive ::=
  r_string_literal non_sigil_name -> { "parent": 0, "children": [0], "type": "RStringCaseInsensitive" } ;
RStringLiteral ::=
  r_string_literal -> { "parent": 0, "children": [], "type": "r_string_literal" } ;
```

Add both as alternatives in `MlrvalOrFunction` (`mlr.bnf:639-667`), next to
the existing `RegexCaseInsensitive` / `StringLiteral` / `BytesLiteral`
lines, CI variant before the plain variant:

```
  | RStringCaseInsensitive
  | RStringLiteral
```

### 2. AST node types

`pkg/dsl/ast_types.go` (near `NodeTypeBytesLiteral`/`NodeTypeRegex`,
~line 11-13):

```go
NodeTypeRStringLiteral        TNodeType = "raw string literal"
NodeTypeRStringCaseInsensitive TNodeType = "case-insensitive raw string literal"
```

`pkg/dsl/cst/ast_types.go` (near line 85 for the plain form, matching the
grammar's `"type"` string; near line 115 alongside
`NodeTypeRegexCaseInsensitive` for the CI form):

```go
NodeTypeRStringLiteral         = "r_string_literal"
NodeTypeRStringCaseInsensitive = "RStringCaseInsensitive"
```

### 3. CST builder (`pkg/dsl/cst/leaves.go`)

New node type and builder, modeled on `BuildBytesLiteralNode`'s
prefix/quote-stripping but **without** calling `UnbackslashStringLiteral`:

```go
// RStringLiteralNode is for explicit raw string literals r"..." (issue #297).
// Unlike StringLiteralNode, no backslash processing is applied: r"\*" evaluates
// to the two characters backslash-asterisk. This makes the value usable directly
// as a regex-engine pattern fragment regardless of where it later travels -- as a
// literal regex argument, or via a variable -- unlike the implicit-r-string trick
// used for plain string literals in regex position (see regexProtectPrePassAux),
// which only works at parse time and can't follow a value through a variable.
type RStringLiteralNode struct {
	literal *mlrval.Mlrval
}

func (root *RootNode) BuildRStringLiteralNode(literal string) IEvaluable {
	// The PGPG lexer produces r_string_literal token with leading 'r' in the lexeme.
	if len(literal) >= 1 && literal[0] == 'r' {
		literal = literal[1:]
	}
	// Case-insensitive form r"..."i: leave the quotes and trailing 'i' intact,
	// matching BuildRegexLiteralNode's representation, so CompileMillerRegex's
	// existing "\"...\"i" handling applies unchanged. Case-insensitivity is only
	// meaningful once compiled as a regex, so this form is not intended to double
	// as a plain string value the way the non-CI form is.
	if len(literal) >= 3 && literal[0] == '"' && strings.HasSuffix(literal, "\"i") {
		return &RStringLiteralNode{literal: mlrval.FromString(literal)}
	}
	// Plain form r"...": strip the surrounding quotes for a clean raw-string value,
	// usable both as a regex argument (via CompileMillerRegex's bare-string
	// fallback) and as an ordinary string value.
	if len(literal) >= 2 && literal[0] == '"' && literal[len(literal)-1] == '"' {
		literal = literal[1 : len(literal)-1]
	}
	return &RStringLiteralNode{literal: mlrval.FromString(literal)}
}

func (node *RStringLiteralNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.literal
}
```

Two new `BuildLeafNode` dispatch cases (~`leaves.go:42-67`), the CI one
mirroring `NodeTypeRegexCaseInsensitive`'s append-`i`-if-absent pattern:

```go
case asts.NodeType(NodeTypeRStringLiteral):
	return root.BuildRStringLiteralNode(sval), nil

case asts.NodeType(NodeTypeRStringCaseInsensitive):
	if sval == "" && astNode.Children != nil && len(astNode.Children) > 0 {
		sval = tokenLit(astNode.Children[0])
	}
	if sval != "" && !strings.HasSuffix(sval, "i") {
		sval = sval + "i"
	}
	return root.BuildRStringLiteralNode(sval), nil
```

**No change needed** to `regexProtectPrePassAux` — it only relabels nodes
that are already `NodeTypeStringLiteral`; grammar-native `r"..."` nodes
arrive pre-labeled as `NodeTypeRStringLiteral`/`NodeTypeRStringCaseInsensitive`
and pass through the prepass untouched.

### 4. No changes needed elsewhere

`pkg/lib/regex.go` and `pkg/bifs/regex.go` need no modification — see
"Where regex compilation happens" above.

### 5. Regenerate the parser

`pkg/parsing/mlr.bnf` is fed through Miller's own PGPG generator
(`github.com/johnkerl/pgpg`, not goyacc/lex) via
[`tools/build-dsl`](../tools/build-dsl) (takes a few minutes) to regenerate
`pkg/parsing/lexer/lexer.go` and `pkg/parsing/parser/parser.go`. Per
`pkg/parsing/README.md` and `README-dev.md:154`, these generated files are
committed to source control — run `tools/build-dsl` after editing
`mlr.bnf` and commit the regenerated output alongside the grammar change.

### 6. Tests

New case(s) under `test/cases/dsl-regex-matching/` (following the existing
`0016` test's structure) covering:

- (a) direct `r"\["`-style literal use as a regex argument.
- (b) the headline variable-carries-raw case:
  `rstar = r"\*"; $y = gsub($x, rstar, "star")`.
- (c) `r"..."` used as an ordinary non-regex value, to confirm it
  prints/stores raw — e.g. `r"a\tb"` is 4 literal characters (`a`, `\`,
  `t`, `b`), not a tab.
- (d) `r"..."i` case-insensitive matching, e.g. `r"abc"i" =~ "ABC"`.

Leave `test/cases/dsl-regex-matching/0016` untouched as a regression guard
for the pre-existing implicit-r behavior.

### 7. Docs

Extend `docs/src/reference-main-regular-expressions.md.in`'s existing
implicit-r-strings passage (~lines 409-417) with:

- The new explicit `r"..."` syntax.
- The variable-carries-raw example from the issue
  (`rstar = r"\*"; gsub($x, rstar, "star")`).
- A note on the CI form's quote-preserving asymmetry (documented tradeoff,
  not a bug): `r"..."i` is intended for regex position; used as a plain
  value it retains its `"..."i` wrapper, consistent with today's implicit
  `"..."i` regex-literal behavior.

Rebuild via `make -C docs/src forcebuild`.

## Verification

1. `make build`.
2. Manual checks:
   ```
   echo 'a=[' | mlr put '$a = gsub($a, r"\[", "left_square")'
   echo 'a=*' | mlr put 'rstar = r"\*"; $a = gsub($a, rstar, "STAR")'
   mlr -n put 'end { print r"\t" }'   # prints two characters, not a tab
   ```
3. `make check` (unit + regression tests), confirming
   `test/cases/dsl-regex-matching/0016` and the new r-string cases both
   pass.
4. `make lint` before pushing.
