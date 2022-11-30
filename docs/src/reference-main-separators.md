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
# Separators

## Record, field, and pair separators

Miller has record separators, field separators, and pair separators. For
example, given the following [DKVP](file-formats.md#dkvp-key-value-pairs)
records:

<pre class="pre-highlight-in-pair">
<b>cat data/a.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=1,b=2,c=3
a=4,b=5,c=6
</pre>

* the **record separator** is newline -- it separates records from one another;
* the **field separator** is `,` -- it separates fields (key-value pairs) from one another;
* and the **pair separator** is `=` -- it separates the key from the value within each key-value pair.

These are the default values, which you can override with flags such as `--ips`
and `--ops` (below).

Not all [file formats](file-formats.md) have all three of these: for example,
CSV does not have a pair separator, since keys are on the header line and
values are on each data line.

Also, separators are not programmable for all file formats.  For example, in
[JSON objects](file-formats.md#json), the pair separator is `:` and the
field-separator is `,` -- we write `{"a":1,"b":2,"c":3}` -- but these aren't
modifiable.  If you do `mlr --json --ips : --ips '=' cat myfile.json` then you
don't get `{"a"=1,"b"=2,"c"=3}`.  This is because the pair-separator `:` is
part of the JSON specification.

## Input and output separators

Miller lets you use the same separators for input and output (e.g. CSV input,
CSV output), or, to change them between input and output (e.g. CSV input, JSON
output), if you wish to transform your data in that way.

Miller uses the names `IRS` and `ORS` for the input and output record
separators, `IFS` and `OFS` for the input and output field separators, and
`IPS` and `OPS` for input and output pair separators.

For example:

<pre class="pre-highlight-in-pair">
<b>cat data/a.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=1,b=2,c=3
a=4,b=5,c=6
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ifs , --ofs ';' --ips = --ops : cut -o -f c,a,b data/a.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
c:3;a:1;b:2
c:6;a:4;b:5
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv head -n 2 example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
yellow,triangle,true,1,11,43.6498,9.8870
red,square,true,2,15,79.2778,0.0130
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv --ofs pipe head -n 2 example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color|shape|flag|k|index|quantity|rate
yellow|triangle|true|1|11|43.6498|9.8870
red|square|true|2|15|79.2778|0.0130
</pre>

If your data has non-default separators and you don't want to change those
between input and output, you can use `--rs`, `--fs`, and `--ps`. Setting `--fs
:` is the same as setting `--ifs : --ofs :`, but with fewer keystrokes.

<pre class="pre-highlight-in-pair">
<b>cat data/modsep.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a:1;b:2;c:3
a:4;b:5;c:6
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --fs ';' --ps : cut -o -f c,a,b data/modsep.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
c:3;a:1;b:2
c:6;a:4;b:5
</pre>

## Multi-character separators

All separators can be multi-character, except for file formats which don't
allow parameterization (see below). And for CSV (CSV-lite doesn't have these
restrictions), IRS must be `\n` and IFS must be a single character.

<pre class="pre-highlight-in-pair">
<b>mlr --ifs ';' --ips : --ofs ';;;' --ops := cut -o -f c,a,b data/modsep.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
c:=3;;;a:=1;;;b:=2
c:=6;;;a:=4;;;b:=5
</pre>

If your data has field separators which are one or more consecutive spaces, you
can use `--ifs space --repifs`.
More generally, the `--repifs` flag means that multiple successive occurrences of the field
separator count as one.  For example, in CSV data we often signify nulls by
empty strings, e.g. `2,9,,,,,6,5,4`. On the other hand, if the field separator
is a space, it might be more natural to parse `2 4    5` the same as `2 4 5`:
`--repifs --ifs ' '` lets this happen.  In fact, the `--ipprint` option
is internally implemented in terms of `--repifs`.

For example:

<pre class="pre-highlight-in-pair">
<b>cat data/extra-spaces.txt</b>
</pre>
<pre class="pre-non-highlight-in-pair">
oh    say   can you
see   by    the dawn's
early light what so
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ifs ' ' --repifs --inidx --oxtab cat  data/extra-spaces.txt</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1 oh
2 say
3 can
4 you

1 see
2 by
3 the
4 dawn's

1 early
2 light
3 what
4 so
</pre>

## Regular-expression separators

`IFS` and `IPS` can be regular expressions: use `--ifs-regex` or `--ips-regex` in place of
`--ifs` or `--ips`, respectively.

You can also use either `--ifs space --repifs` or `--ifs-regex '()+'`. (But that gets a little tedious,
so there are aliases listed below.) Note however that `--ifs space --repifs` is about 3x faster than
`--ifs-regex '( )+'` -- regular expressions are powerful, but slower.

## Aliases

Many things we'd like to write as separators need to be escaped from the shell
-- e.g. `--ifs ';'` or `--ofs '|'`, and so on. You can use the following if you like:

<pre class="pre-highlight-in-pair">
<b>mlr help list-separator-aliases</b>
</pre>
<pre class="pre-non-highlight-in-pair">
ascii_esc  = "\x1b"
ascii_etx  = "\x04"
ascii_fs   = "\x1c"
ascii_gs   = "\x1d"
ascii_null = "\x01"
ascii_rs   = "\x1e"
ascii_soh  = "\x02"
ascii_stx  = "\x03"
ascii_us   = "\x1f"
asv_fs     = "\x1f"
asv_rs     = "\x1e"
colon      = ":"
comma      = ","
cr         = "\r"
crcr       = "\r\r"
crlf       = "\r\n"
crlfcrlf   = "\r\n\r\n"
equals     = "="
lf         = "\n"
lflf       = "\n\n"
newline    = "\n"
pipe       = "|"
semicolon  = ";"
slash      = "/"
space      = " "
tab        = "\t"
usv_fs     = "\xe2\x90\x9f"
usv_rs     = "\xe2\x90\x9e"
</pre>

And for `--ifs-regex` and `--ips-regex`:

<pre class="pre-highlight-in-pair">
<b>mlr help list-separator-regex-aliases</b>
</pre>
<pre class="pre-non-highlight-in-pair">
spaces     = "( )+"
tabs       = "(\t)+"
whitespace = "([ \t])+"
</pre>

Note that `spaces`, `tabs`, and `whitespace` already are regexes so you
shouldn't use `--repifs` with them. (In fact, the `--repifs` flag is ignored
when `--ifs-regex` is provided.)

## Command-line flags

Given the above, we now have seen the following flags:

<pre class="pre-non-highlight-non-pair">
--rs --irs --ors
--fs --ifs --ofs --repifs --ifs-regex
--ps --ips --ops --ips-regex
</pre>

See also the [separator-flags section](reference-main-flag-list.md#separator-flags).

## DSL built-in variables

Miller exposes for you read-only [built-in variables](reference-dsl-variables.md#built-in-variables) with
names `IRS`, `ORS`, `IFS`, `OFS`, `IPS`, and `OPS`. Unlike in AWK, you can't set these in begin-blocks --
their values indicate what you specified at the command line -- so their use is limited.

<pre class="pre-highlight-in-pair">
<b>mlr --ifs , --ofs ';' --ips = --ops : --from data/a.dkvp put '$d = ">>>" . IFS . "|||" . OFS . "<<<"'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a:1;b:2;c:3;d:>>>,|||;<<<
a:4;b:5;c:6;d:>>>,|||;<<<
</pre>

## Which separators apply to which file formats

Notes:

* CSV IRS and ORS must be newline, and CSV IFS must be a single character. (CSV-lite does not have these restrictions.)
* TSV IRS and ORS must be newline, and TSV IFS must be a tab. (TSV-lite does not have these restrictions.)
* See the [CSV section](file-formats.md#csvtsvasvusvetc) for information about ASV and USV.
* JSON: ignores all separator flags from the command line.
* Headerless CSV overlaps quite a bit with NIDX format using comma for IFS. See also the page on [CSV with and without headers](csv-with-and-without-headers.md).
* For XTAB, the record separator is a repetition of the field separator. For example, if one record has `x=1,y=2` and the next has `x=3,y=4`, and OFS is newline, then output lines are `x 1`, then `y 2`, then an extra newline, then `x 3`, then `y 4`. This means: to customize XTAB, set `OFS` rather than `ORS`.

|            | **RS**  | **FS**  | **PS**   |
|------------|---------|---------|----------|
| [**CSV**](file-formats.md#csvtsvasvusvetc)    | Always `\n`; not alterable * | Default `,`; must be single-character    | None     |
| [**TSV**](file-formats.md#csvtsvasvusvetc)    | Always `\n`; not alterable * |  Default `\t`; must be single-character   | None     |
| [**CSV-lite**](file-formats.md#csvtsvasvusvetc)    | Default `\n` *   | Default `,`    | None     |
| [**TSV-lite**](file-formats.md#csvtsvasvusvetc)    | Default `\n` *  |  Default `\t`   | None     |
| [**JSON**](file-formats.md#json)   | N/A; records are between `{` and `}` | Always `,`; not alterable    | Always `:`; not alterable |
| [**DKVP**](file-formats.md#dkvp-key-value-pairs)   | Default `\n`    | Default `,`    | Default `=` |
| [**NIDX**](file-formats.md#nidx-index-numbered-toolkit-style)   | Default `\n`    | Default space    | None     |
| [**XTAB**](file-formats.md#xtab-vertical-tabular)   | Not used; records are separated by an extra FS    | `\n` *    | Default: space with repeats  |
| [**PPRINT**](file-formats.md#pprint-pretty-printed-tabular) | Default `\n` *    | Space with repeats    | None     |
| [**Markdown**](file-formats.md#markdown-tabular) | Always `\n`; not alterable * | One or more spaces, then `|`, then one or more spaces; not alterable | None     |

\* or `\r\n` on Windows
