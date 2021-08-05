<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Reference: then-chaining

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
mlr cut -x -f i,y data/big | mlr sort -n y > /dev/null

% time sh piped.sh
real 0m2.828s
user 0m3.183s
sys  0m0.137s


% cat chained.sh
mlr cut -x -f i,y then sort -n y data/big > /dev/null

% time sh chained.sh
real 0m2.082s
user 0m1.933s
sys  0m0.137s
</pre>

There are two reasons to use then-chaining: one is for performance, although I don't expect this to be a win in all cases.  Using then-chaining avoids redundant string-parsing and string-formatting at each pipeline step: instead input records are parsed once, they are fed through each pipeline stage in memory, and then output records are formatted once. On the other hand, Miller is single-threaded, while modern systems are usually multi-processor, and when streaming-data programs operate through pipes, each one can use a CPU.  Rest assured you get the same results either way.

The other reason to use then-chaining is for simplicity: you don't have re-type formatting flags (e.g. `--csv --fs tab`) at every pipeline stage.
