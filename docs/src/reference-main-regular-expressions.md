<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flags</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verbs</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Functions</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="../release-docs/index.html">Release docs</a>
</span>
</div>
# Regular expressions

Miller lets you use regular expressions (of the [types accepted by Go](https://pkg.go.dev/regexp)) in the following contexts:

* In `mlr filter` with `=~` or `!=~`, e.g. `mlr filter '$url =~ "http.*com"'`

* In `mlr put` with `regextract`, e.g. `mlr put '$output = regextract($input, "[a-z][a-z][0-9][0-9]")`

* In `mlr put` with `sub` or `gsub`, e.g. `mlr put '$url = sub($url, "http.*com", "")'`

* In `mlr having-fields`, e.g. `mlr having-fields --any-matching '^sda[0-9]'`

* In `mlr cut`, e.g. `mlr cut -r -f '^status$,^sda[0-9]'`

* In `mlr rename`, e.g. `mlr rename -r '^(sda[0-9]).*$,dev/\1'`

* In `mlr grep`, e.g. `mlr --csv grep 00188555487 myfiles*.csv`

Points demonstrated by the above examples:

* There are no implicit start-of-string or end-of-string anchors; please use `^` and/or `$` explicitly.

* Miller regexes are wrapped with double quotes rather than slashes.

* The `i` after the ending double quote indicates a case-insensitive regex.

* Capture groups are wrapped with `(...)` rather than `\(...\)`; use `\(` and `\)` to match against parentheses.

Example:

<pre class="pre-highlight-in-pair">
<b>cat data/regex-in-data.dat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
name=jane,regex=^j.*e$
name=bill,regex=^b[ou]ll$
name=bull,regex=^b[ou]ll$
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr filter '$name =~ $regex' data/regex-in-data.dat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
name=jane,regex=^j.*e$
name=bull,regex=^b[ou]ll$
</pre>

## Regex captures for the `=~` operator

Regex captures of the form `\0` through `\9` are supported as follows:

* Captures have in-function context for `sub` and `gsub`. For example, the first `\1,\2` pair belong to the first `sub` and the second `\1,\2` pair belong to the second `sub`:

<pre class="pre-highlight-non-pair">
<b>mlr put '$b = sub($a, "(..)_(...)", "\2-\1"); $c = sub($a, "(..)_(.)(..)", ":\1:\2:\3")'</b>
</pre>

* Captures endure for the entirety of a `put` for the `=~` and `!=~` operators. For example, here the `\1,\2` are set by the `=~` operator and are used by both subsequent assignment statements:

<pre class="pre-highlight-non-pair">
<b>mlr put '$a =~ "(..)_(....); $b = "left_\1"; $c = "right_\2"'</b>
</pre>

* Each user-defined function has its own frame for captures. For example:

<pre class="pre-highlight-non-pair">
<b>mlr -n put '</b>
<b>func f() {</b>
<b>    if ("456 defg" =~ "([0-9]+) ([a-z]+)") {</b>
<b>        print "INNER: \1 \2";</b>
<b>    }</b>
<b>}</b>
<b>end {</b>
<b>    if ("123 abc" =~ "([0-9]+) ([a-z]+)") {</b>
<b>        print "OUTER PRE:  \1 \2";</b>
<b>        f();</b>
<b>        print "OUTER POST: \1 \2";</b>
<b>    }</b>
<b>}'</b>
</pre>

* The captures are not retained across multiple puts. For example, here the `\1,\2` won't be expanded from the regex capture:

<pre class="pre-highlight-non-pair">
<b>mlr put '$a =~ "(..)_(....)' then {... something else ...} then put '$b = "left_\1"; $c = "right_\2"'</b>
</pre>

* Up to nine matches are supported: `\1` through `\9`, while `\0` is the entire match string; `\15` is treated as `\1` followed by an unrelated `5`.

## Resetting captures

If you use `(...)` in your regular expression, then up to 9 matches are supported for the `=~`
operator, and an arbitrary number of matches are supported for the `match` DSL function.

* Before any match is done, `"\1"` etc. in a string evaluate to themselves.
* After a successful match is done, `"\1"` etc. in a string evaluate to the matched substring.
* After an unsuccessful match is done, `"\1"` etc. in a string evaluate to the empty string.
* You can match against `null` to reset to the original state.

<pre class="pre-highlight-in-pair">
<b>mlr repl</b>
</pre>
<pre class="pre-non-highlight-in-pair">

[mlr] "\1:\2"
"\1:\2"

[mlr] "abc" =~ "..."
true

[mlr] "\1:\2"
":"

[mlr] "abc" =~ "(.).(.)"
true

[mlr] "\1:\2"
"a:c"

[mlr] "abc" =~ "(.)x(.)"
false

[mlr] "\1:\2"
":"

[mlr] "abc" =~ null

[mlr] "\1:\2"
"\1:\2"
</pre>

## The `strmatch` and `strmatchx` DSL functions

The `=~` and `!=~` operators have been in Miller for a long time, and they will continue to be
supported.  They do, however, have some deficiencies. As of Miller 6.11 and beyond, the `strmatch`
and `strmatchx` provide more robust ways to do capturing.

First, some examples.

The `strmatch` function only returns a boolean result, and it doesn't set `\0..\9`:

<pre class="pre-highlight-in-pair">
<b>mlr repl</b>
</pre>
<pre class="pre-non-highlight-in-pair">

[mlr] strmatch("abc", "....")
false

[mlr] strmatch("abc", "...")
true

[mlr] strmatch("abc", "(.).(.)")
true

[mlr] strmatch("[ab:3458]", "([a-z]+):([0-9]+)")
true
</pre>

The `strmatchx` function also doesn't set `\0..\9`, but returns a map-valued result:

<pre class="pre-highlight-in-pair">
<b>mlr repl</b>
</pre>
<pre class="pre-non-highlight-in-pair">

[mlr] strmatchx("abc", "....")
{
  "matched": false
}

[mlr] strmatchx("abc", "...")
{
  "matched": true,
  "full_capture": "abc",
  "full_start": 1,
  "full_end": 3
}

[mlr] strmatchx("abc", "(.).(.)")
{
  "matched": true,
  "full_capture": "abc",
  "full_start": 1,
  "full_end": 3,
  "captures": ["a", "c"],
  "starts": [1, 3],
  "ends": [1, 3]
}

[mlr] "[ab:3458]" =~ "([a-z]+):([0-9]+)"
true

[mlr] "\1"
"ab"

[mlr] "\2"
"3458"

[mlr] strmatchx("[ab:3458]", "([a-z]+):([0-9]+)")
{
  "matched": true,
  "full_capture": "ab:3458",
  "full_start": 2,
  "full_end": 8,
  "captures": ["ab", "3458"],
  "starts": [2, 5],
  "ends": [3, 8]
}
</pre>

Notes:

* When there is no match, the result from `strmatchx` only has the `"matched":false` key/value pair.
* When there is a match with no captures, the result from `strmatchx` has the `"matched":true` key/value pair,
  as well as `full_capture` (taking the place of `\0` set by `=~`), and `full_start` and `full_end`
  which `=~` does not offer.
* When there is a match with no captures, the result from `strmatchx` also has the `captures` array
  whose slots 1, 2, 3, ... are the same as would have been set by `=~` via `\1, \2, \3, ...`.
  However, `strmatchx` offers an arbitrary number of captures, not just `\1..\9`.
  Additionally, the `starts` and `ends` arrays are indices into the input string.
* Since you hold the return value from `strmatchx`, you can operate on it as you wish --- instead of
  relying on the (function-scoped) globals `\0..\9`.
* The price paid is that using `strmatchx` does indeed tend to take more keystrokes than `=~`.

## More information

Regular expressions are those supported by the [Go regexp package](https://pkg.go.dev/regexp), which in turn are of type [RE2](https://github.com/google/re2/wiki/Syntax) except for `\C`:

<pre class="pre-highlight-in-pair">
<b>go doc regexp/syntax</b>
</pre>
<pre class="pre-non-highlight-in-pair">
package syntax // import "regexp/syntax"

Package syntax parses regular expressions into parse trees and compiles parse
trees into programs. Most clients of regular expressions will use the facilities
of package regexp (such as Compile and Match) instead of this package.

# Syntax

The regular expression syntax understood by this package when parsing with
the Perl flag is as follows. Parts of the syntax can be disabled by passing
alternate flags to Parse.

Single characters:

    .              any character, possibly including newline (flag s=true)
    [xyz]          character class
    [^xyz]         negated character class
    \d             Perl character class
    \D             negated Perl character class
    [[:alpha:]]    ASCII character class
    [[:^alpha:]]   negated ASCII character class
    \pN            Unicode character class (one-letter name)
    \p{Greek}      Unicode character class
    \PN            negated Unicode character class (one-letter name)
    \P{Greek}      negated Unicode character class

Composites:

    xy             x followed by y
    x|y            x or y (prefer x)

Repetitions:

    x*             zero or more x, prefer more
    x+             one or more x, prefer more
    x?             zero or one x, prefer one
    x{n,m}         n or n+1 or ... or m x, prefer more
    x{n,}          n or more x, prefer more
    x{n}           exactly n x
    x*?            zero or more x, prefer fewer
    x+?            one or more x, prefer fewer
    x??            zero or one x, prefer zero
    x{n,m}?        n or n+1 or ... or m x, prefer fewer
    x{n,}?         n or more x, prefer fewer
    x{n}?          exactly n x

Implementation restriction: The counting forms x{n,m}, x{n,}, and x{n} reject
forms that create a minimum or maximum repetition count above 1000. Unlimited
repetitions are not subject to this restriction.

Grouping:

    (re)           numbered capturing group (submatch)
    (?P<name>re)   named & numbered capturing group (submatch)
    (?:re)         non-capturing group
    (?flags)       set flags within current group; non-capturing
    (?flags:re)    set flags during re; non-capturing

    Flag syntax is xyz (set) or -xyz (clear) or xy-z (set xy, clear z). The flags are:

    i              case-insensitive (default false)
    m              multi-line mode: ^ and $ match begin/end line in addition to begin/end text (default false)
    s              let . match \n (default false)
    U              ungreedy: swap meaning of x* and x*?, x+ and x+?, etc (default false)

Empty strings:

    ^              at beginning of text or line (flag m=true)
    $              at end of text (like \z not \Z) or line (flag m=true)
    \A             at beginning of text
    \b             at ASCII word boundary (\w on one side and \W, \A, or \z on the other)
    \B             not at ASCII word boundary
    \z             at end of text

Escape sequences:

    \a             bell (== \007)
    \f             form feed (== \014)
    \t             horizontal tab (== \011)
    \n             newline (== \012)
    \r             carriage return (== \015)
    \v             vertical tab character (== \013)
    \*             literal *, for any punctuation character *
    \123           octal character code (up to three digits)
    \x7F           hex character code (exactly two digits)
    \x{10FFFF}     hex character code
    \Q...\E        literal text ... even if ... has punctuation

Character class elements:

    x              single character
    A-Z            character range (inclusive)
    \d             Perl character class
    [:foo:]        ASCII character class foo
    \p{Foo}        Unicode character class Foo
    \pF            Unicode character class F (one-letter name)

Named character classes as character class elements:

    [\d]           digits (== \d)
    [^\d]          not digits (== \D)
    [\D]           not digits (== \D)
    [^\D]          not not digits (== \d)
    [[:name:]]     named ASCII class inside character class (== [:name:])
    [^[:name:]]    named ASCII class inside negated character class (== [:^name:])
    [\p{Name}]     named Unicode property inside character class (== \p{Name})
    [^\p{Name}]    named Unicode property inside negated character class (== \P{Name})

Perl character classes (all ASCII-only):

    \d             digits (== [0-9])
    \D             not digits (== [^0-9])
    \s             whitespace (== [\t\n\f\r ])
    \S             not whitespace (== [^\t\n\f\r ])
    \w             word characters (== [0-9A-Za-z_])
    \W             not word characters (== [^0-9A-Za-z_])

ASCII character classes:

    [[:alnum:]]    alphanumeric (== [0-9A-Za-z])
    [[:alpha:]]    alphabetic (== [A-Za-z])
    [[:ascii:]]    ASCII (== [\x00-\x7F])
    [[:blank:]]    blank (== [\t ])
    [[:cntrl:]]    control (== [\x00-\x1F\x7F])
    [[:digit:]]    digits (== [0-9])
    [[:graph:]]    graphical (== [!-~] == [A-Za-z0-9!"#$%&'()*+,\-./:;<=>?@[\\\]^_`{|}~])
    [[:lower:]]    lower case (== [a-z])
    [[:print:]]    printable (== [ -~] == [ [:graph:]])
    [[:punct:]]    punctuation (== [!-/:-@[-`{-~])
    [[:space:]]    whitespace (== [\t\n\v\f\r ])
    [[:upper:]]    upper case (== [A-Z])
    [[:word:]]     word characters (== [0-9A-Za-z_])
    [[:xdigit:]]   hex digit (== [0-9A-Fa-f])

Unicode character classes are those in unicode.Categories and unicode.Scripts.

func IsWordChar(r rune) bool
type EmptyOp uint8
    const EmptyBeginLine EmptyOp = 1 << iota ...
    func EmptyOpContext(r1, r2 rune) EmptyOp
type Error struct{ ... }
type ErrorCode string
    const ErrInternalError ErrorCode = "regexp/syntax: internal error" ...
type Flags uint16
    const FoldCase Flags = 1 << iota ...
type Inst struct{ ... }
type InstOp uint8
    const InstAlt InstOp = iota ...
type Op uint8
    const OpNoMatch Op = 1 + iota ...
type Prog struct{ ... }
    func Compile(re *Regexp) (*Prog, error)
type Regexp struct{ ... }
    func Parse(s string, flags Flags) (*Regexp, error)
</pre>

One caveat: for strings in "regex position" -- e.g. the second argument to
[`sub`](reference-dsl-builtin-functions.md#sub) or
[`gsub`](reference-dsl-builtin-functions.md#gsub), or after `=~` -- `"\t"`
means a backslash and a `t` -- which is the right thing -- whereas for strings
in "non-regex position", e.g. anywhere else, `"\t"` becomes the tab character.
This is to say (if you're familiar with r-strings in Python) all strings in
regex position are implicit r-strings.  Generally this is the right thing and
should cause little confusion. Note however that this means `"\t"."\t"` in the
second argument to `sub` isn't the same as `"\t\t"`.
