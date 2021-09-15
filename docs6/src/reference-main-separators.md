<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flag list</a>
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

Miller lets you use the same separators for input and output, or, to change
them between input and output, if you wish to transform your data in that way.

Miller uses the names `IRS` and `ORS` for the input and output record
separators, `IFS` and `OFS` for the input and output field separators, and
`IPS` and `OPS` for input and output pair separators.

For example:

<pre class="pre-highlight-in-pair">
<b>mlr --ifs , --ofs ';' --ips = --ops : cut -o -f c,a,b data/a.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
c:3;a:1;b:2
c:6;a:4;b:5
</pre>

If your data has non-default separators and you don't want to change those
between input and output, you can use `--rs`, `--fs`, and `--ps`. Setting `--fs
:` is the same as setting `--ifs : --ofs :`, but with fewer keystrokes.

<pre class="pre-highlight-in-pair">
<b>mlr --fs ';' --ps : cut -o -f c,a,b data/modsep.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
c:3;a:1;b:2
c:6;a:4;b:5
</pre>

## Multi-character separators

The separators default to single characters, but can be multiple characters if you like:

<pre class="pre-highlight-in-pair">
<b>mlr --ifs ';' --ips : --ofs ';;;' --ops := cut -o -f c,a,b data/modsep.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
c:=3;;;a:=1;;;b:=2
c:=6;;;a:=4;;;b:=5
</pre>

While the separators can be multiple characters, [regular
expressions](reference-main-regular-expressions.md) (which Miller supports in
many ways) are not (as of mid-2021) supported by Miller. So, in the above
example, you can say the field-separator is one semicolon, or three, but two or
four won't be recognized using `--ifs ';;;'`.

To fill this need, in the absence of full regular-expression support, Miller
has a `--repifs` option for input. This means, for example, using `--ifs
' ' --repifs` you can have the field separator be one _or more_ spaces. (Mixes
of spaces and tabs, however, won't be recognized as a separator.)

The `--repifs` flag means that multiple successive occurrences of the field
separator count as one.  For example, in CSV data we often signify nulls by
empty strings, e.g. `2,9,,,,,6,5,4`. On the other hand, if the field separator
is a space, it might be more natural to parse `2 4    5` the same as `2 4 5`:
`--repifs --ifs ' '` lets this happen.  In fact, the `--ipprint` option above
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

(TODO: FIXME)

<pre class="pre-highlight-in-pair">
<b>mlr --ifs ' ' --repifs --inidx --oxtab cat  data/extra-spaces.txt</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1 oh
2 
3 
4 
5 say
6 
7 
8 can
9 you

1 see
2 
3 
4 by
5 
6 
7 
8 the
9 dawn's

1 early
2 light
3 what
4 so
</pre>

## Command-line flags

Given the above, we now have seen the following flags:

<pre class="pre-non-highlight-non-pair">
--rs --irs --ors
--fs --ifs --ofs --repifs
--ps --ips --ops
</pre>

Also note that you can use names for certain characters: e.g. `--fs space` is
the same as `--fs ' '`.  A full list is: `colon`, `comma`, `equals`, `newline`,
`pipe`, `semicolon`, `slash`, `space`, `tab`.

## DSL built-in variables

Miller exposes for you read-only [built-in variables](reference-dsl-variables.md#built-in-variables) with
names `IRS`, `ORS`, `IFS`, `OFS`, `IPS`, and `OPS`. Unlike in AWK, you can't set these in begin-blocks --
their values indicate what you set at the command line -- so their use is limited.

<pre class="pre-highlight-in-pair">
<b>mlr --ifs , --ofs ';' --ips = --ops : --from data/a.dkvp put '$d = ">>>" . IFS . "|||" . OFS . "<<<"'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a:1;b:2;c:3;d:>>>,|||;<<<
a:4;b:5;c:6;d:>>>,|||;<<<
</pre>

## Which separators apply to which file formats

* CSV/TSV/ASV/USV/etc.:
    * Record separator is newline (Linux/BSDs/MacOS) or carriage-return-newline (Windows); programmable in Miller 5 and below; TODO to support for the Miller 6 release.
    * If field separator is tab, we have TSV; see more examples (ASV, USV, etc.) at in the [CSV section](file-formats.md#csvtsvasvusvetc).
    * No pair separator.
* JSON: ignores all separator flags from the command line.
* PPRINT
    * Record separator is newline (Linux/BSDs/MacOS) or carriage-return-newline (Windows); programmable in Miller 5 and below; TODO to support for the Miller 6 release.
    * TODO: write up
    * TODO: write up
* Markdown tabular: ignores all separator flags from the command line.
* XTAB
    * TODO: write up
    * TODO: write up
    * TODO: write up
* DKVP: lets you specify record, field, and pair separators.
* NIDX
    * TODO: write up
    * TODO: write up
    * No pair separator.
