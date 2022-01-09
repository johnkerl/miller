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
# DSL operators

## Detailed listing

Operators are listed on the [DSL built-in functions page](reference-dsl-builtin-functions.md).

## Operator precedence

Operators are listed in order of decreasing precedence, highest first.

| Operators                     | Associativity |
|-------------------------------|---------------|
| `()` `{}` `[]`                | left to right |
| `**`                          | right to left |
| `!` `~` unary`+` unary`-` `&` | right to left |
| binary`*` `/` `//` `%`        | left to right |
| `.`                           | left to right |
| binary`+` binary`-`           | left to right |
| `<<` `>>` `>>>`               | left to right |
| `&`                           | left to right |
| `^`                           | left to right |
| `|`                           | left to right |
| `<` `<=` `>` `>=`             | left to right |
| `==` `!=` `=~` `!=~` `<=>`    | left to right |
| `???`                         | left to right |
| `??`                          | left to right |
| `&&`                          | left to right |
| `^^`                          | left to right |
| `||`                          | left to right |
| `? :`                         | right to left |
| `=`                           |  N/A for Miller (there is no $a=$b=$c) |

## Operator and function semantics

* Functions are often pass-throughs straight to the system-standard Go libraries.

* The [`min`](reference-dsl-builtin-functions.md#min) and [`max`](reference-dsl-builtin-functions.md#max) functions are different from other multi-argument functions which return null if any of their inputs are null: for [`min`](reference-dsl-builtin-functions.md#min) and [`max`](reference-dsl-builtin-functions.md#max), by contrast, if one argument is absent-null, the other is returned. Empty-null loses min or max against numeric or boolean; empty-null is less than any other string.

* Symmetrically with respect to the bitwise OR, XOR, and AND operators
[`|`](reference-dsl-builtin-functions.md#bitwise-or),
[`&`](reference-dsl-builtin-functions.md#bitwise-and), and
[`^`](reference-dsl-builtin-functions.md#bitwise-xor), Miller has logical operators
[`||`](reference-dsl-builtin-functions.md#logical-or),
[`&&`](reference-dsl-builtin-functions.md#logical-and), and
[`^^`](reference-dsl-builtin-functions.md#logical-xor).

* The exponentiation operator [`**`](reference-dsl-builtin-functions.md#exponentiation) is familiar from many languages, except that an integer raised to an int power is int, not float.

* The regex-match and regex-not-match operators [`=~`](reference-dsl-builtin-functions.md#regmatch) and [`!=~`](reference-dsl-builtin-functions.md#regnotmatch) are similar to those in Ruby and Perl.

## The double-purpose dot operator

The main use for the `.` operator is for string concatenation: `"abc" . "def"` is `"abc.def"`.

However, in Miller 6 it has optional use for map traversal. Example:

<pre class="pre-highlight-in-pair">
<b>cat data/server-log.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "hostname": "localhost",
  "pid": 12345,
  "req": {
    "id": 6789,
    "method": "GET",
    "path": "api/check",
    "host": "foo.bar",
    "headers": {
      "host": "bar.baz",
      "user-agent": "browser"
    }
  },
  "res": {
    "status_code": 200,
    "header": {
      "content-type": "text",
      "content-encoding": "plain"
    }
  }
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --json --from data/server-log.json put -q '</b>
<b>  print $req["headers"]["host"];</b>
<b>  print $req.headers.host;</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
bar.baz
bar.baz
[
]
</pre>

This also works on the left-hand sides of assignment statements:

<pre class="pre-highlight-in-pair">
<b>mlr --json --from data/server-log.json put '</b>
<b>  $req.headers.host = "UPDATED";</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "hostname": "localhost",
  "pid": 12345,
  "req": {
    "id": 6789,
    "method": "GET",
    "path": "api/check",
    "host": "foo.bar",
    "headers": {
      "host": "UPDATED",
      "user-agent": "browser"
    }
  },
  "res": {
    "status_code": 200,
    "header": {
      "content-type": "text",
      "content-encoding": "plain"
    }
  }
}
]
</pre>

A few caveats:

* This is why `.` has higher precedece than `+` in the table above -- in Miller 5 and below, where `.` was only used for concatenation, it had the same precedence as `+`. So you can now do this:

<pre class="pre-highlight-in-pair">
<b>mlr --json --from data/server-log.json put -q '</b>
<b>  print $req.id + $res.status_code</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
6989
[
]
</pre>

* However (awkwardly), if you want to use `.` for map-traversal as well as string-concatenation in the same statement, you'll need to insert parentheses, as the default associativity is left-to-right:

<pre class="pre-highlight-in-pair">
<b>mlr --json --from data/server-log.json put -q '</b>
<b>  print $req.method . " -- " . $req.path</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
(error)
[
]
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --json --from data/server-log.json put -q '</b>
<b>  print ($req.method) . " -- " . ($req.path)</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
GET -- api/check
[
]
</pre>
