<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verb list</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Function list</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="https://github.com/johnkerl/miller" target="_blank">Repository ↗</a>
</span>
</div>
# Regular expressions

Miller lets you use regular expressions (of type POSIX.2) in the following contexts:

* In `mlr filter` with `=~` or `!=~`, e.g. `mlr filter '$url =~ "http.*com"'`

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

## Regex captures

Regex captures of the form `\0` through `\9` are supported as

* Captures have in-function context for `sub` and `gsub`. For example, the first `\1,\2` pair belong to the first `sub` and the second `\1,\2` pair belong to the second `sub`:

<pre class="pre-highlight-non-pair">
<b>mlr put '$b = sub($a, "(..)_(...)", "\2-\1"); $c = sub($a, "(..)_(.)(..)", ":\1:\2:\3")'</b>
</pre>

* Captures endure for the entirety of a `put` for the `=~` and `!=~` operators. For example, here the `\1,\2` are set by the `=~` operator and are used by both subsequent assignment statements:

<pre class="pre-highlight-non-pair">
<b>mlr put '$a =~ "(..)_(....); $b = "left_\1"; $c = "right_\2"'</b>
</pre>

* The captures are not retained across multiple puts. For example, here the `\1,\2` won't be expanded from the regex capture:

<pre class="pre-highlight-non-pair">
<b>mlr put '$a =~ "(..)_(....)' then {... something else ...} then put '$b = "left_\1"; $c = "right_\2"'</b>
</pre>

* Up to nine matches are supported: `\1` through `\9`, while `\0` is the entire match string; `\15` is treated as `\1` followed by an unrelated `5`.

