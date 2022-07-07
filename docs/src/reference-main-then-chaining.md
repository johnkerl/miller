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
# Then-chaining

In accord with the [Unix philosophy](http://en.wikipedia.org/wiki/Unix_philosophy), you can pipe data into or out of Miller. For example:

<pre class="pre-highlight-non-pair">
<b>mlr cut --complement -f os_version *.dat | mlr sort -f hostname,uptime</b>
</pre>

You can, if you like, instead simply chain commands together using the `then` keyword:

<pre class="pre-highlight-non-pair">
<b>mlr cut --complement -f os_version then sort -f hostname,uptime *.dat</b>
</pre>

(You can precede the very first verb with `then`, if you like, for symmetry.)

Here's a performance comparison:

<pre class="pre-non-highlight-non-pair">
% cat piped.sh
mlr cut -x -f i,y data/big | mlr sort -n y &gt; /dev/null

% time sh piped.sh
real    0m2.321s
user    0m4.878s
sys     0m1.564s

% cat chained.sh
mlr cut -x -f i,y then sort -n y data/big &gt; /dev/null

% time sh chained.sh
real    0m2.070s
user    0m2.738s
sys     0m1.259s
</pre>

There are two reasons to use then-chaining: one is for performance, although I don't expect this to be a win in all cases.  Using then-chaining avoids redundant string-parsing and string-formatting at each pipeline step: instead input records are parsed once, they are fed through each pipeline stage in memory, and then output records are formatted once.

The other reason to use then-chaining is for simplicity: you don't have re-type formatting flags (e.g. `--csv --fs tab`) at every pipeline stage.

As of Miller 6.3.0, `+` is an alias for `then`.
